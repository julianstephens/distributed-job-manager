package worker

import (
	"errors"
	"fmt"

	"github.com/julianstephens/distributed-job-manager/pkg/auth0client"
	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

type Reporter struct {
	conf    *models.Config
	token   *string
	client  *auth0client.Auth0Client
	baseURL string
	log     *graylogger.GrayLogger
	api     *JobAPI
}

func NewReporter(log *graylogger.GrayLogger) *Reporter {
	conf := config.GetConfig()
	r := &Reporter{
		conf:   conf,
		token:  nil,
		log:    log,
		client: auth0client.NewAuth0Client(conf.Auth0Worker.ClientId, conf.Auth0Worker.ClientSecret),
	}
	r.baseURL = r.conf.JobAPIEndpoint
	r.api = NewJobAPI(r.client, log, r.baseURL)
	return r
}

func (r *Reporter) RegisterExecution(jobId string) (*models.JobExecution, error) {
	exec := models.JobExecution{
		JobID:    jobId,
		WorkerID: r.conf.WorkerID,
		Status:   models.JobStatusScheduled,
	}

	data, err := r.api.CreateExecution(exec)
	if err != nil {
		return nil, err
	}

	if data == nil {
		msg := "no data returned from execution registration"
		r.log.Error(msg, nil)
		return nil, errors.New(msg)
	}

	return data, nil
}

func (r *Reporter) StartExecution(executionId string) (*models.JobExecution, error) {
	status := models.JobStatusInProgress

	update := models.JobExecutionUpdateRequest{
		Status: &status,
	}

	data, err := r.api.UpdateExecution(executionId, update)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no data returned from execution start")
	}

	return data, nil
}

func (r *Reporter) CompleteExecution(executionId string, response RunnerResponse) (*models.JobExecution, error) {
	r.log.Info(fmt.Sprintf("completing execution %s with status %s", executionId, utils.If(response.Error == nil, "completed", "failed")), nil)

	update := models.JobExecutionUpdateRequest{
		StartTime:    &response.StartTime,
		EndTime:      &response.EndTime,
		Status:       utils.StringPtr(utils.If(response.Error == nil, models.JobStatusCompleted, models.JobStatusFailed)),
		ErrorMessage: response.Error,
		Output:       response.Output,
	}

	data, err := r.api.UpdateExecution(executionId, update)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no data returned from execution completion")
	}

	return data, nil
}
