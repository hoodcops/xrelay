package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	development = "development"
	production  = "production"
)

var env = struct {
	Port                      int    `envconfig:"PORT" required:"true"`
	Environment               string `envconfig:"ENVIRONMENT" default:"development"`
	BrokerDSN                 string `envconfig:"BROKER_DSN" required:"true"`
	VerificationWorkers       int    `envconfig:"VERIFICATION_WORKERS" default:"5"`
	InvitationWorkers         int    `envconfig:"INVITATION_WORKERS" default:"5"`
	AlertWorkers              int    `envconfig:"ALERT_WORKERS" default:"10"`
	City                      string `envconfig:"CITY" required:"true"`
	Locale                    string `envconfig:"LOCALE" default:"en"`
	TwilioVerificationAPIHost string `envconfig:"TWILIO_VERIFICATION_API_HOST" required:"true"`
	TwilioVerificationAPIKey  string `envconfig:"TWILIO_VERIFICATION_API_KEY" required:"true"`
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

	conn, err := amqp.Dial(env.BrokerDSN)
	if err != nil {
		logger.Fatal("failed connecting to rabbitmq", zap.Error(err))
	}

	_ = conn
}
