package main

import (
	"context"
	"flag"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/farawaygg/wisdom/internal/services/private"
	"github.com/farawaygg/wisdom/internal/storage/postgres"
	"github.com/farawaygg/wisdom/internal/wisdom"
	privAPI "github.com/farawaygg/wisdom/pkg/wisdom"

	"github.com/farawaygg/configurator"
	"github.com/farawaygg/go-stdlib/graceful"
	"github.com/farawaygg/go-stdlib/grpc"
	"github.com/farawaygg/go-stdlib/psql"
	zaplib "github.com/farawaygg/go-stdlib/zap"
)

func main() {
	flag.Parse()

	c, err := configurator.New(configurator.WithFile(*configFile))
	if err != nil {
		panic(err)
	}

	var cfg Config
	c.MustLoad(&cfg)

	zlog := zaplib.MustNew(cfg.LogLevel.Level())
	backgroundCtx := ctxzap.ToContext(context.Background(), zlog)

	db, err := psql.Open(cfg.DBConnstring)
	if err != nil {
		zlog.Panic("psql.Open", zap.Error(err))
	}
	pgStorage, err := postgres.New(sqlx.NewDb(db, "postgres"))
	if err != nil {
		zlog.Panic("postgres.New", zap.Error(err))
	}

	wisdomsRepo := wisdom.New(pgStorage)

	privSrv := private.New(wisdomsRepo)

	eg, ctx := errgroup.WithContext(backgroundCtx)
	eg.Go(func() error {
		var opts grpc.ServerOptions
		opts.WithLogger(zlog)

		grpcserv := opts.NewServer()
		privAPI.RegisterWisdomSvcServer(grpcserv, privSrv)

		zlog.Info("grpc starting", zap.String("addr", cfg.Listen.GRPC))
		return grpc.ListenAndServeContext(ctx, cfg.Listen.GRPC, grpcserv)
	})
	eg.Go(func() error {
		return graceful.ListenSignalContext(ctx)
	})

	if err := eg.Wait(); err != nil {
		// Do not panic or fatal
		// in order to continue graceful shutdown.
		zlog.Error(err.Error())
	}
}
