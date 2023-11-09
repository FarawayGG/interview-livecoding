package main

import (
	"flag"
	"os"

	"github.com/go-playground/validator/v10"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v3"
)

var configFile = flag.String("config", "app/config.yaml", "config file to load")

type Config struct {
	Listen struct {
		GRPC string `yaml:"GRPC" validate:"required"`
	} `yaml:"Listen" validate:"required"`
	LogLevel     string `yaml:"LogLevel" validate:"required"`
	DBConnstring string `yaml:"DBConnstring" validate:"required"`
}

func mustLoadConfig() Config {
	configFile, err := os.Open(*configFile)
	if err != nil {
		panic(errors.WithMessage(err, "os.Open"))
	}
	defer configFile.Close()

	var config Config
	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(&config); err != nil {
		panic(errors.WithMessage(err, "yamlDecoder.Decode"))
	}

	if err := validator.New().Struct(config); err != nil {
		panic(errors.WithMessage(err, "validator.New().Struct"))
	}

	return config
}

func codeToLevels(code codes.Code) zapcore.Level {
	switch code {
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated, codes.FailedPrecondition:
		return zapcore.InfoLevel
	case codes.OK:
		return zapcore.DebugLevel
	}

	return grpc_zap.DefaultCodeToLevel(code)
}
