package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/farawaygg/wisdom/internal/psql"
	"github.com/farawaygg/wisdom/internal/services/private"
	"github.com/farawaygg/wisdom/internal/storage/postgres"
	"github.com/farawaygg/wisdom/internal/wisdom"
	privAPI "github.com/farawaygg/wisdom/pkg/wisdom"
)

func main() {
	flag.Parse()

	cfg := mustLoadConfig()

	zl, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(errors.WithMessage(err, "zapcore.ParseLevel"))
	}
	zcfg := zap.NewProductionConfig()
	zcfg.Level = zap.NewAtomicLevelAt(zl)
	zlog, err := zcfg.Build()
	if err != nil {
		panic(errors.WithMessage(err, "zcfg.Build"))
	}
	zap.ReplaceGlobals(zlog)
	backgroundCtx := ctxzap.ToContext(context.Background(), zlog)

	db, err := psql.Open(cfg.DBConnstring)
	if err != nil {
		zlog.Panic("openSQL", zap.Error(err))
	}
	pgStorage, err := postgres.New(sqlx.NewDb(db, "postgres"))
	if err != nil {
		zlog.Panic("postgres.New", zap.Error(err))
	}

	wisdomsRepo := wisdom.New(pgStorage)

	privSrv := private.New(wisdomsRepo)

	eg, ctx := errgroup.WithContext(backgroundCtx)
	eg.Go(func() error {
		grpcserv := grpc.NewServer(grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(zlog, grpc_zap.WithLevels(codeToLevels)),
			grpc_validator.UnaryServerInterceptor(),
		),
			grpc.ChainStreamInterceptor(
				grpc_zap.StreamServerInterceptor(zlog, grpc_zap.WithLevels(codeToLevels)),
				grpc_validator.StreamServerInterceptor(),
			))
		privAPI.RegisterWisdomSvcServer(grpcserv, privSrv)

		zlog.Info("grpc starting", zap.String("addr", cfg.Listen.GRPC))
		return listenAndServeContext(ctx, cfg.Listen.GRPC, grpcserv)
	})
	eg.Go(func() error {
		return listenSignalContext(ctx)
	})

	if err := eg.Wait(); err != nil {
		// Do not panic or fatal
		// in order to continue graceful shutdown.
		zlog.Error(err.Error())
	}
}

func listenAndServeContext(ctx context.Context, addr string, srv *grpc.Server) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()
	return srv.Serve(lis)
}

func listenSignalContext(ctx context.Context, sig ...os.Signal) error {
	ch := make(chan os.Signal, 1)

	ss := make([]os.Signal, 0, 2+len(sig))
	ss = append(ss, syscall.SIGTERM, syscall.SIGINT)
	ss = append(ss, sig...)
	signal.Notify(ch, ss...)

	select {
	case <-ctx.Done():
		return nil
	case sig := <-ch:
		return fmt.Errorf("cancelled by signal (%v)", sig)
	}
}
