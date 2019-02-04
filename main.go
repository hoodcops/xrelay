package main

import (
	"context"
	"log"

	"github.com/hoodcops/xrelay/pkg/app"
	"github.com/hoodcops/xrelay/pkg/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	development = "development"
	production  = "production"
)

var confs config.Config

func init() {
	err := envconfig.Process("", &confs)
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
	logger, err := initLogger(confs.Environment)
	if err != nil {
		log.Fatalf("failed initializing logger : %s", err)
	}

	if confs.Environment == development {
		logger.Info("successfully loaded env vars", zap.Any("configuration", confs))
	}

	conn, err := amqp.Dial(confs.BrokerDSN)
	if err != nil {
		logger.Fatal("failed connecting to rabbitmq", zap.Error(err))
	}
	logger.Info("connected to broker successfully")

	app := app.NewApp(conn, &confs, logger)
	app.Run(context.Background())

}
