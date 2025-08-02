package main

import (
	"encoding/json"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/services/workersvc/worker"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conf := config.GetConfig()

	conn, err := queue.GetConnection(conf.Rabbit.Username, conf.Rabbit.Password, conf)
	if err != nil {
		logger.Fatalf("unable to get queue connection: %v", err)
		return
	}
	defer queue.CloseConnection(conf.Rabbit.Username)

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("unable to get queue channel: %v", err)
		return
	}
	defer ch.Close()

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
	reporter := worker.NewReporter()

	go processJobs(runner, pool, reporter, msgs)
	select {}
}

func processJobs(runner *worker.Runner, pool *worker.SandboxPool, reporter *worker.Reporter, msgs <-chan amqp091.Delivery) {
	var job models.Job
	for d := range msgs {
		if err := json.Unmarshal(d.Body, &job); err != nil {
			panic(err)
		}

		logger.Infof("got job %s. starting processing...", job.JobID)

		jobExec, err := reporter.RegisterExecution(job.JobID)
		if err != nil {
			panic(err)
		}

		logger.Infof("%d sandboxes available", pool.AvailableCount())
		req, err := runner.NewRequest(job, jobExec.ExecutionID)
		if err != nil {
			panic(err)
		}

		box, err := pool.Reserve(job.UserID)
		if err != nil {
			panic(err)
		}
		logger.Infof("user %s reserved sandbox %d", job.UserID, box.ID)

		req.BoxID = box.ID

		if err := runner.RunCode(*req, reporter); err != nil {
			panic(err)
		}

		pool.Release(job.UserID)
		logger.Infof("user %s released sandbox %d", job.UserID, box.ID)
		logger.Infof("%d sandboxes available", pool.AvailableCount())
	}
}
