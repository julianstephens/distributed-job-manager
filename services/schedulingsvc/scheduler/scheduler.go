package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/auth0client"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/queue"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
	"github.com/rabbitmq/amqp091-go"
)

type Scheduler struct {
	baseURL    string
	conf       *models.Config
	api        *JobAPI
	logger     *graylogger.GrayLogger
	ScheduleCh *amqp091.Channel
}

func NewScheduler(config *models.Config, logger *graylogger.GrayLogger) (*Scheduler, error) {
	conn, err := queue.GetConnection(config.Rabbit.Username, config.Rabbit.Password, "", config)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Scheduler{
		conf: config,
		api: NewJobAPI(
			auth0client.NewAuth0Client(config.Schedule.Auth0ClientID, config.Schedule.Auth0ClientSecret),
			logger,
			config.JobAPIEndpoint,
		),
		logger:     logger,
		ScheduleCh: ch,
	}, nil
}

func (s *Scheduler) Run(ctx context.Context) {
	s.pollTable(ctx, s.ScheduleCh)

	tick := time.NewTicker(30 * time.Second)

	go s.runner(ctx, tick, s.ScheduleCh)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	tick.Stop()
}

func (s *Scheduler) pollTable(ctx context.Context, ch *amqp091.Channel) error {
	startTime := time.Now().Add(time.Second * 60)
	endTime := startTime.Add(time.Second * 60)

	s.logger.Info(fmt.Sprintf("looking for scheduled jobs between %s and %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)), nil)
	queuedSchedules, err := s.api.GetSchedules(&startTime, &endTime)
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("found %d scheduled jobs", len(*queuedSchedules)), nil)

	for _, sched := range *queuedSchedules {
		job, err := s.api.GetJob(sched.JobID)
		if err != nil {
			s.logger.Error(fmt.Sprintf("failed to get job %s for scheduling", sched.JobID), &err)
			return err
		}

		if job.Status != models.JobStatusPending {
			return nil
		}

		jobJson, err := json.Marshal(job)
		if err != nil {
			s.logger.Error("failed to marshal scheduled job to json", &err)
			return err
		}

		if err := ch.PublishWithContext(ctx, s.conf.Rabbit.Name, "", false, false, amqp091.Publishing{ContentType: "application/json", Body: jobJson}); err != nil {
			s.logger.Error(fmt.Sprintf("failed to publish job %s to queue", job.JobID), &err)
			return err
		}
		s.logger.Info(fmt.Sprintf("sent job %s to queue", job.JobID), nil)

		if _, err = s.api.UpdateJob(job.JobID, &models.JobUpdateRequest{Status: utils.StringPtr(models.JobStatusScheduled)}); err != nil {
			return err
		}
		s.logger.Info(fmt.Sprintf("updated job %s status to scheduled", job.JobID), nil)
	}

	return nil
}

func (s *Scheduler) runner(ctx context.Context, tick *time.Ticker, ch *amqp091.Channel) {
	for range tick.C {
		s.pollTable(ctx, ch)
	}
}
