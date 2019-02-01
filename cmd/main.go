package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		logger.Fatal("failed binding to port", zap.Int("port", env.Port))
	}
	defer listener.Close()

	// routes := api.InitRoutes(logger, conn)
	server := http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           nil,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	connsClosed := make(chan struct{})
	go func() {
		defer close(connsClosed)

		recv := <-sigs
		logger.Info("received signal, shutting down", zap.Any("signal", recv.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("failed shutting down server", zap.Error(err))
		}
	}()

	url := fmt.Sprintf("http://%s", listener.Addr())
	logger.Info("server listening on ", zap.String("url", url))

	if err = server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			logger.Fatal("failed starting server", zap.Error(err))
		}
	}

	<-connsClosed
	logger.Info("server shutdown successfully")
}
