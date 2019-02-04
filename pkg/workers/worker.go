package workers

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Worker represents any object that listens for message on
// a queue and acts on them
type Worker interface {
	Run()
}

// VerificationWorker consumes messages on the verification queue
// and sends verification codes to msisds to verify their authenticity
type VerificationWorker struct {
	conn   *amqp.Connection
	logger *zap.Logger
}

// NewVerificationWorker returns a new verification worker
func NewVerificationWorker(conn *amqp.Connection, logger *zap.Logger) *VerificationWorker {
	return &VerificationWorker{
		conn:   conn,
		logger: logger,
	}
}

// Run ...
func (vw VerificationWorker) Run() {
	channel, err := vw.conn.Channel()
	if err != nil {
		vw.logger.Fatal("failed to open a channel", zap.Error(err))
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"msisdn.verification", // quene name
		true,                  // durability
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no wait
		nil,                   //arguments
	)
	if err != nil {
		vw.logger.Fatal("failed declaring msisdn verification queue", zap.Error(err))
	}

	messages, err := channel.Consume(
		queue.Name, // queue name
		"",         // consumer
		true,       // auto-acknowledge
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		vw.logger.Fatal("failed consuming messages on queue", zap.Error(err))
	}

	go func() {
		vw.logger.Info("wating for messages")
		for msg := range messages {
			vw.logger.Debug("message received on msisdn.verification queue", zap.Any("body", string(msg.Body)))
		}

	}()
}
