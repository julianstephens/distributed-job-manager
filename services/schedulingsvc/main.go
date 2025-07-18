package main

import (
	"github.com/julianstephens/distributed-task-scheduler/pkg/config"
	"github.com/julianstephens/distributed-task-scheduler/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/pkg/queue"
)

func main() {
	conf := config.GetConfig()

	conn, err := queue.GetConnection(conf)
	if err != nil {
		logger.Fatalf("unable to get rabbitmq connection: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("unable to get queue channel: %v", err)
		return
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		logger.Fatalf("failed to set QoS: %v", err)
		return
	}

}
