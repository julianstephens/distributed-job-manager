package main

import (
	"encoding/json"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/services/workersvc/worker"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conf := config.GetConfig()
	conn, ch, err := queue.GetConnection(conf)
	defer queue.CloseConnection(conn, ch)

	if err != nil {
		logger.Fatalf("unable to get queue connection: %v", err)
		return
	}

	msgs, err := ch.Consume(
		conf.Rabbit.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatalf("unable to register queue consumer: %v", err)
		return
	}

	pool := worker.NewSandboxPool(conf.SandboxCount)
	pool.ScheduleCleanup()
	logger.Infof("initialized sandbox pool. %d sandboxes available", pool.AvailableCount())

	runner := worker.NewRunner()

	var forever chan struct{}
	go processJobs(runner, pool, msgs)
	<-forever
}

func processJobs(runner *worker.Runner, pool *worker.SandboxPool, msgs <-chan amqp091.Delivery) {
	var job models.Job
	for d := range msgs {
		if err := json.Unmarshal(d.Body, &job); err != nil {
			panic(err)
		}

		logger.Infof("%d sandboxes available", pool.AvailableCount())
		req, err := runner.NewRequest(job)
		if err != nil {
			panic(err)
		}

		box, err := pool.Reserve(job.UserID)
		if err != nil {
			panic(err)
		}
		logger.Infof("user %s reserved sandbox %d", job.UserID, box.ID)

		req.BoxID = box.ID

		// runner.RunCode(*req)
		time.Sleep(time.Second * 10)

		pool.Release(job.UserID)
		logger.Infof("user %s released sandbox %d", job.UserID, box.ID)
		logger.Infof("%d sandboxes available", pool.AvailableCount())
	}
}
