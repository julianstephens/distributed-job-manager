package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	"github.com/rabbitmq/amqp091-go"
	"github.com/scylladb/gocqlx/v3/qb"
)

func main() {
	conf := config.GetConfig()
	db, err := store.GetDB(conf.Cassandra.Keyspace)
	if err != nil {
		logger.Fatalf("unable to get db connection: %v", err)
	}

	conn, ch, err := queue.GetConnection(conf)
	if err != nil {
		return
	}
	defer queue.CloseConnection(conn, ch)

	err = ch.Qos(1, 0, false)
	if err != nil {
		logger.Fatalf("failed to set QoS: %v", err)
		return
	}

	err = ch.ExchangeDeclare(
		conf.Rabbit.Name,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatalf("unable to declare exchange: %v", err)
		return
	}

	err = ch.QueueBind(
		conf.Rabbit.Name, // queue name
		"",               // routing key
		conf.Rabbit.Name, // exchange
		false,
		nil,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pollTable(ctx, ch, conf, db)

	tick := time.NewTicker(30 * time.Second)

	go scheduler(ctx, tick, ch, conf, db)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	tick.Stop()
}

func scheduler(ctx context.Context, tick *time.Ticker, ch *amqp091.Channel, conf *models.Config, db *store.DBSession) {
	for _ = range tick.C {
		pollTable(ctx, ch, conf, db)
	}
}

func pollTable(ctx context.Context, ch *amqp091.Channel, conf *models.Config, db *store.DBSession) error {
	var queuedSchedules []models.JobSchedule
	q := db.Client.Query(models.JobSchedules.SelectAll())
	if err := q.SelectRelease(&queuedSchedules); err != nil {
		return err
	}

	logger.Infof("found %d scheduled jobs", len(queuedSchedules))

	for _, sched := range queuedSchedules {
		var job models.Job
		stmt, _ := qb.Select(models.Jobs.Name()).Where(qb.EqNamed("job_id", sched.JobID)).AllowFiltering().ToCql()
		if err := db.Client.Query(stmt, []string{"job_id"}).Bind(sched.JobID).Get(&job); err != nil {
			return err
		}

		jobJson, err := json.Marshal(job)
		if err != nil {
			return err
		}

		if err := ch.PublishWithContext(ctx, conf.Rabbit.Name, "", false, false, amqp091.Publishing{ContentType: "application/json", Body: jobJson}); err != nil {
			return err
		}
		logger.Infof("sent job %s to queue", job.JobID)

		stmt, names := qb.Delete(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id"), qb.Eq("status")).ToCql()
		q := db.Client.Query(stmt, names).Bind(job.JobID, job.UserID, job.Status)
		if err := q.ExecRelease(); err != nil {
			return err
		}

		updatedJob := job
		updatedJob.Status = "scheduled"
		q = db.Client.Query(models.Jobs.Insert()).BindStruct(updatedJob)
		if err := q.ExecRelease(); err != nil {
			return err
		}
	}

	return nil
}
