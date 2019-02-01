package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

const (
	development = "development"
	production  = "production"
)

var env = struct {
	Port        int    `envconfig:"PORT" required:"true"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	BrokerDSN   string `envconfig:"BROKER_DSN" required:"true"`
}{}

func init() {
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("failed loading env vars : %s", err)
	}
}

func initLogger(environment string) (*zap.Logger, error) {
	if environment == production {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}

func main() {
	logger, err := initLogger(env.Environment)
	if err != nil {
		log.Fatalf("failed initializing logger : %s", err)
	}

	if env.Environment == development {
		logger.Info("successfully loaded env vars", zap.Any("configuration", env))
	}
}
