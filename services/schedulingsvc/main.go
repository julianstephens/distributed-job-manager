package main

import (
	"context"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/services/schedulingsvc/scheduler"
)

func main() {
	conf := config.GetConfig()

	log, err := graylogger.NewLogger("schedsvc")
	if err != nil {
		logger.Fatalf("unable to init logger: %v", err)
		return
	}
	defer queue.CloseConnection(conf.Rabbit.LoggingUsername)
	defer log.LogCh.Close()

	sched, err := scheduler.NewScheduler(conf, log)
	if err != nil {
		logger.Fatalf("unable to create scheduler: %v", err)
		return
	}
	defer queue.CloseConnection(conf.Rabbit.Username)
	defer sched.ScheduleCh.Close()

	log.Info("scheduling service initialized", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sched.Run(ctx)
}
