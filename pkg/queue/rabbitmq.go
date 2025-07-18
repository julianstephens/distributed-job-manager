package queue

import (
	"fmt"
	"sync"

	"github.com/julianstephens/distributed-task-scheduler/pkg/models"
	"github.com/rabbitmq/amqp091-go"
)

var (
	RabbitConn *amqp091.Connection
	once       sync.Once
)

func GetConnection(conf *models.Config) (*amqp091.Connection, error) {
	var err error
	once.Do(func() {
		err = setup(conf)
	})
	if err != nil {
		return nil, err
	}
	return RabbitConn, nil
}

func setup(conf *models.Config) error {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.Rabbit.Username, conf.Rabbit.Password, conf.Rabbit.Host, conf.Rabbit.Port))
	if err != nil {
		return err
	}

	RabbitConn = conn
	return nil
}
