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
	cache      = make(map[string]*amqp091.Connection)
	cacheMutex sync.RWMutex
)

func GetConnection(username string, password string, conf *models.Config) (*amqp091.Connection, error) {
	cacheMutex.RLock()
	if val, ok := cache[username]; ok {
		cacheMutex.RUnlock()
		return val, nil
	}
	cacheMutex.RUnlock()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if val, ok := cache[username]; ok {
		return val, nil
	}

	conn, err := setup(username, password, conf)
	if err != nil {
		return nil, err
	}

	cache[username] = conn
	return conn, err
}

func CloseConnection(username string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if conn, ok := cache[username]; ok {
		if err := conn.Close(); err != nil {
			logger.Errorf("error closing connection for %s: %v", username, err)
		}
		delete(cache, username)
		logger.Infof("closed connection for %s", username)
	} else {
		logger.Warnf("no connection found for %s to close", username)
	}
}

func setup(username string, password string, conf *models.Config) (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/graylog", username, password, conf.Rabbit.Host, conf.Rabbit.Port))
	if err != nil {
		logger.Fatalf("unable to get rabbitmq connection: %v", err)
		return nil, err
	}

	return conn, nil
}
