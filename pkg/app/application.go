package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hoodcops/xrelay/pkg/config"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// App represents application dependencies
type App struct {
	brokerConn *amqp.Connection
	config     *config.Config
	logger     *zap.Logger
}

// NewApp returns a new App
func NewApp(brokerConn *amqp.Connection, config *config.Config, logger *zap.Logger) *App {
	return &App{
		brokerConn: brokerConn,
		config:     config,
		logger:     logger,
	}
}

// Run boots the service
func (s *App) Run(ctx context.Context) {
	s.logger.Info("starting service")

	srv := &http.Server{
		Handler: nil,
		Addr:    fmt.Sprintf(":%d", s.config.Port),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("http server error occured", zap.Error(err))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	recv := <-sigs

	s.logger.Info("signal received, shutting down server", zap.String("signal", recv.String()))

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Warn("error shutting down server", zap.Error(err))
	}

	s.logger.Info("service shutdown successfully")
}
