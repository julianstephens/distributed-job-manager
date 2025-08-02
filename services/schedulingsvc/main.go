package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	"github.com/rabbitmq/amqp091-go"
	"github.com/scylladb/gocqlx/v3/qb"
)

// TODO:
// 1. switch to api requests instead of direct db interactions
// 2. clean up main func

func main() {
	conf := config.GetConfig()
	db, err := store.GetDB(conf.Cassandra.Keyspace)
	if err != nil {
		logger.Fatalf("unable to get db connection: %v", err)
		return
	}

	log, err := graylogger.NewLogger("schedsvc")
	if err != nil {
		logger.Fatalf("unable to init logger: %v", err)
		return
	}
	defer queue.CloseConnection(log.LogConn, log.LogCh)

	log.Info("scheduling service initialized", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pollTable(ctx, log.LogCh, conf, db, log)

	tick := time.NewTicker(30 * time.Second)

	go scheduler(ctx, tick, log.LogCh, conf, db, log)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	tick.Stop()
}

func scheduler(ctx context.Context, tick *time.Ticker, ch *amqp091.Channel, conf *models.Config, db *store.DBSession, log *graylogger.GrayLogger) {
	for _ = range tick.C {
		pollTable(ctx, ch, conf, db, log)
	}
}

func pollTable(ctx context.Context, ch *amqp091.Channel, conf *models.Config, db *store.DBSession, log *graylogger.GrayLogger) error {
	log.Info("scanning table...", nil)
	startTime := time.Now()
	endTime := time.Now().Add(time.Second * 60)

	var queuedSchedules []models.JobSchedule
	stmt, names := qb.Select(models.JobSchedules.Name()).Where(qb.GtOrEq("next_run_time"), qb.LtOrEq("next_run_time")).AllowFiltering().ToCql()
	if err := db.Client.Query(stmt, names).Bind(startTime, endTime).SelectRelease(&queuedSchedules); err != nil {
		log.Error("unable to scan job schedules table", &err)
		return err
	}

	log.Info(fmt.Sprintf("found %d scheduled jobs", len(queuedSchedules)), nil)

	for _, sched := range queuedSchedules {
		var job models.Job
		stmt, _ := qb.Select(models.Jobs.Name()).Where(qb.EqNamed("job_id", sched.JobID)).AllowFiltering().ToCql()
		if err := db.Client.Query(stmt, []string{"job_id"}).Bind(sched.JobID).Get(&job); err != nil {
			log.Error("unable to select scheduled job from db", &err)
			return err
		}

		jobJson, err := json.Marshal(job)
		if err != nil {
			log.Error("failed to marshal scheduled job to json", &err)
			return err
		}

		if err := ch.PublishWithContext(ctx, conf.Rabbit.Name, "", false, false, amqp091.Publishing{ContentType: "application/json", Body: jobJson}); err != nil {
			return err
		}
		log.Info(fmt.Sprintf("sent job %s to queue", job.JobID), nil)

		stmt, names := qb.Delete(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id"), qb.Eq("status")).ToCql()
		q := db.Client.Query(stmt, names).Bind(job.JobID, job.UserID, job.Status)
		if err := q.ExecRelease(); err != nil {
			return err
		}

		updatedJob := job
		updatedJob.Status = models.JobStatusScheduled
		q = db.Client.Query(models.Jobs.Insert()).BindStruct(updatedJob)
		if err := q.ExecRelease(); err != nil {
			return err
		}
	}

	return nil
}
