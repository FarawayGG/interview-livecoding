package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
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
	DBConnString string `yaml:"DBConnString" validate:"required"`
}

func mustLoadConfig() Config {
	configFile, err := os.Open(*configFile)
	if err != nil {
		panic(fmt.Errorf("os.Open: %w", err))
	}
	defer configFile.Close()

	var config Config
	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(&config); err != nil {
		panic(fmt.Errorf("yamlDecoder.Decode: %w", err))
	}

	if err := validator.New().Struct(config); err != nil {
		panic(fmt.Errorf("validator.New().Struct: %w", err))
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
