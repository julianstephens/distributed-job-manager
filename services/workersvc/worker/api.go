package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/julianstephens/distributed-job-manager/pkg/auth0client"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
)

type JobAPI struct {
	client       *auth0client.Auth0Client
	logger       *graylogger.GrayLogger
	executionURL string
	jobURL       string
}

type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func NewJobAPI(client *auth0client.Auth0Client, logger *graylogger.GrayLogger, baseURL string) *JobAPI {
	return &JobAPI{
		client:       client,
		executionURL: fmt.Sprintf("%s/executions", baseURL),
		jobURL:       fmt.Sprintf("%s/jobs", baseURL),
		logger:       logger,
	}
}

func (api *JobAPI) CreateExecution(data models.JobExecution) (execution *models.JobExecution, err error) {
	req, err := http.NewRequest("POST", api.executionURL, strings.NewReader(string(utils.MustMarshalJson(data))))
	if err != nil {
		return
	}

	res, err := api.client.Request(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusCreated {
		api.logApiError(*req, *res, body)
		err = fmt.Errorf("failed to create job execution: %s", res.Status)
		return
	}

	var apiResponse APIResponse[models.JobExecution]
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		api.logger.Error("failed to unmarshal job execution", &err)
		return
	}

	execution = &apiResponse.Data

	return
}

func (api *JobAPI) UpdateExecution(id string, data models.JobExecutionUpdateRequest) (execution *models.JobExecution, err error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", api.executionURL, id), strings.NewReader(string(utils.MustMarshalJson(data))))
	if err != nil {
		return
	}

	res, err := api.client.Request(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		api.logApiError(*req, *res, body)
		err = fmt.Errorf("failed to update job execution: %s", res.Status)
		return
	}

	var apiResponse APIResponse[models.JobExecution]
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		api.logger.Error("failed to unmarshal job execution", &err)
		return
	}

	jobUpdates := models.JobUpdateRequest{}
	if data.Status != nil {
		switch *data.Status {
		case models.JobStatusFailed:
			jobUpdates.Status = utils.StringPtr(models.JobStatusReady)
		case models.JobStatusCancelled:
			jobUpdates.Status = utils.StringPtr(models.JobStatusReady)
		case models.JobStatusInProgress:
			jobUpdates.Status = utils.StringPtr(models.JobStatusInProgress)
		}
	}

	if jobUpdates.Status != nil {
		if _, err = api.UpdateJob(apiResponse.Data.JobID, jobUpdates); err != nil {
			err = fmt.Errorf("failed to update job status: %w", err)
			return
		}
	}

	execution = &apiResponse.Data

	return
}

func (api *JobAPI) UpdateJob(id string, data models.JobUpdateRequest) (job *models.Job, err error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", api.jobURL, id), strings.NewReader(string(utils.MustMarshalJson(data))))
	if err != nil {
		return
	}

	res, err := api.client.Request(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		api.logApiError(*req, *res, body)
		err = fmt.Errorf("failed to update job: %s", res.Status)
		return
	}

	var apiResponse APIResponse[models.Job]
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		api.logger.Error("failed to unmarshal job", &err)
		return
	}

	job = &apiResponse.Data

	return
}

func (api *JobAPI) logApiError(req http.Request, res http.Response, body []byte) {
	data := map[string]any{
		"status_code": res.StatusCode,
		"status":      res.Status,
		"method":      req.Method,
		"url":         req.URL.String(),
		"timestamp":   time.Now().Format(time.RFC3339),
		"error":       "failed to fetch scheduled jobs",
	}

	err := errors.New(string(body))
	api.logger.ErrorWithData(fmt.Sprintf("error on job execution request: %s", res.Status), &err, &data)
}
