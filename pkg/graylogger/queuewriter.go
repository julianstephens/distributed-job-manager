package graylogger

import (
	"encoding/json"

	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
)

type QueueWriter struct {
	config    *models.Config
	conn      *amqp091.Connection
	ch        *amqp091.Channel
	queueName string
}

func NewQueueWriter(config *models.Config, queueName string) (*QueueWriter, error) {
	conn, err := queue.GetConnection(config.Rabbit.LoggingUsername, config.Rabbit.LoggingPassword, "graylog", config)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &QueueWriter{
		config:    config,
		conn:      conn,
		ch:        ch,
		queueName: queueName,
	}, nil
}

// Write sends the log message to the configured RabbitMQ queue
func (w *QueueWriter) Write(p []byte) (n int, err error) {
	logStr := string(p)

	log := formatGELFLog(w.queueName, logStr, LogLevelInfo, nil, nil)
	body, err := json.Marshal(log)
	if err != nil {
		return 0, err
	}

	err = w.ch.Publish(
		"logs",      // exchange
		w.queueName, // routing key (queue name)
		false,       // mandatory
		false,       // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		logger.Errorf("failed to publish log message to RabbitMQ: %v", err)
		return 0, err
	}
	return len(p), nil
}

// Close closes the AMQP channel and connection
func (w *QueueWriter) Close() {
	w.ch.Close()
	w.conn.Close()
}
