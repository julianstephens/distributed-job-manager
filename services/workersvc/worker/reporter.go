package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/auth0client"
	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

type Reporter struct {
	conf    *models.Config
	token   *string
	client  *auth0client.Auth0Client
	baseURL string
}

func NewReporter() *Reporter {
	conf := config.GetConfig()
	r := &Reporter{
		conf:   conf,
		token:  nil,
		client: auth0client.NewAuth0Client(conf.Auth0Worker.ClientId, conf.Auth0Worker.ClientSecret),
	}
	r.baseURL = fmt.Sprintf("%s/%s", r.conf.JobAPIEndpoint, "executions")
	return r
}

func (r *Reporter) RegisterExecution(jobId string) (*models.JobExecution, error) {
	exec := models.JobExecution{
		JobID:    jobId,
		WorkerID: r.conf.WorkerID,
		Status:   models.JobStatusScheduled,
	}

	data, err := r.makeRequest("POST", nil, &exec)
	if err != nil {
		return nil, err
	}

	var jobExec models.JobExecution
	if err = json.Unmarshal(*data, &jobExec); err != nil {
		return nil, err
	}

	return &jobExec, nil
}

func (r *Reporter) StartExecution(executionId string) (*models.JobExecution, error) {
	start := time.Now().UTC()
	status := models.JobStatusInProgress

	update := models.JobExecutionUpdateRequest{
		StartTime: &start,
		Status:    &status,
	}

	data, err := r.makeRequest("PATCH", &executionId, update)
	if err != nil {
		return nil, err
	}

	var jobExec models.JobExecution
	if err = json.Unmarshal(*data, &jobExec); err != nil {
		return nil, err
	}

	return &jobExec, nil
}

func (r *Reporter) CompleteExecution(executionId string, response RunnerResponse) (*models.JobExecution, error) {
	end := time.Now().UTC()
	update := models.JobExecutionUpdateRequest{
		EndTime: &end,
	}

	data, err := r.makeRequest("PATCH", &executionId, update)
	if err != nil {
		return nil, err
	}

	var jobExec models.JobExecution
	if err = json.Unmarshal(*data, &jobExec); err != nil {
		return nil, err
	}

	return &jobExec, nil
}

func (r *Reporter) makeRequest(method string, path *string, body any) (*[]byte, error) {
	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(
		method,
		utils.If(path != nil, fmt.Sprintf("%s/%s", r.baseURL, *path), r.baseURL),
		utils.If(body != nil, bytes.NewReader(jsonBody), nil),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *r.token))

	res, err := r.client.Request(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
