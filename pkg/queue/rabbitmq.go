package queue

import (
	"fmt"
	"sync"

	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn *amqp091.Connection
	RabbitCh   *amqp091.Channel
	once       sync.Once
)

func GetConnection(conf *models.Config) (*amqp091.Connection, *amqp091.Channel, error) {
	var err error
	once.Do(func() {
		err = setup(conf)
	})
	if err != nil {
		return nil, nil, err
	}
	return RabbitConn, RabbitCh, nil
}

func CloseConnection(conn *amqp091.Connection, ch *amqp091.Channel) {
	defer conn.Close()
	defer ch.Close()
}

func setup(conf *models.Config) error {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/graylog", "grayloguser", "0v3rth3r3", conf.Rabbit.Host, conf.Rabbit.Port))
	if err != nil {
		logger.Fatalf("unable to get rabbitmq connection: %v", err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("unable to get queue channel: %v", err)
		return err
	}

	RabbitConn = conn
	RabbitCh = ch
	return nil
}
