package main

import (
	"encoding/json"
	"fmt"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
	"github.com/julianstephens/distributed-job-manager/services/workersvc/worker"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conf := config.GetConfig()

	log, err := graylogger.NewLogger("worker")
	if err != nil {
		logger.Fatalf("unable to create logger: %v", err)
		return
	}
	defer queue.CloseConnection(conf.Rabbit.LoggingUsername)
	defer log.LogCh.Close()

	conn, err := queue.GetConnection(conf.Rabbit.Username, conf.Rabbit.Password, "", conf)
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
	log.Info(fmt.Sprintf("sandbox pool created with %d sandboxes", conf.SandboxCount), nil)

	runner := worker.NewRunner(log)
	reporter := worker.NewReporter(log)

	forever := make(chan bool)
	go processJobs(runner, pool, reporter, log, msgs)
	<-forever
}

func processJobs(runner *worker.Runner, pool *worker.SandboxPool, reporter *worker.Reporter, log *graylogger.GrayLogger, msgs <-chan amqp091.Delivery) {
	log.Info("worker started, waiting for jobs...", nil)

	var job models.Job
	for d := range msgs {
		if err := json.Unmarshal(d.Body, &job); err != nil {
			log.Error("failed to unmarshal job", &err)
			return
		}

		log.Info(fmt.Sprintf("worker received job %s for user %s", job.JobID, job.UserID), nil)

		jobExec, err := reporter.RegisterExecution(job.JobID)
		if err != nil {
			log.Error(fmt.Sprintf("failed to register job execution for job %s", job.JobID), &err)
			return
		}
		log.Info(fmt.Sprintf("registered job execution %s for job %s", jobExec.ExecutionID, jobExec.JobID), nil)

		req, err := runner.NewRequest(job, jobExec.ExecutionID)
		if err != nil {
			log.Error(fmt.Sprintf("failed to create request for job %s", job.JobID), &err)
			return
		}
		data, _ := json.Marshal(req)
		log.Info(fmt.Sprintf("created request for job %s with execution ID %s", job.JobID, req.ExecutionID), utils.StringPtr(string(data)))

		box, err := pool.Reserve(job.UserID)
		if err != nil {
			log.Error(fmt.Sprintf("failed to reserve sandbox for user %s", job.UserID), &err)
			return
		}
		log.Info(fmt.Sprintf("reserved sandbox %d for user %s", box.ID, job.UserID), nil)

		req.BoxID = box.ID

		if err := runner.RunCode(*req, reporter); err != nil {
			log.Error(fmt.Sprintf("failed to run code for job %s in sandbox %d", job.JobID, box.ID), &err)
			return
		}

		pool.Release(job.UserID)
		log.Info(fmt.Sprintf("released sandbox %d for user %s", box.ID, job.UserID), nil)
	}
}
