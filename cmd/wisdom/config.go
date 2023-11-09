package main

import (
	"flag"

	zaplib "github.com/farawaygg/go-stdlib/zap"
)

var configFile = flag.String("config", "app/config.yaml", "config file to load")

type Config struct {
	Listen struct {
		GRPC    string `validate:"nonzero"`
		Metrics string `validate:"nonzero"`
	}
	LogLevel     zaplib.LogLevel
	DBConnstring string `validate:"nonzero"`
}
