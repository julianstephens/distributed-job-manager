package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	schedulesURL string
	jobsURL      string
}

type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func NewJobAPI(client *auth0client.Auth0Client, logger *graylogger.GrayLogger, baseURL string) *JobAPI {
	return &JobAPI{
		client:       client,
		schedulesURL: fmt.Sprintf("%s/schedules", baseURL),
		jobsURL:      fmt.Sprintf("%s/jobs", baseURL),
		logger:       logger,
	}
}

func (api *JobAPI) GetSchedules(startTime *time.Time, endTime *time.Time) (schedules *[]models.JobSchedule, err error) {
	params := url.Values{}

	if startTime != nil {
		params.Add("next_run_time[gte]", startTime.Format(time.RFC3339))
	}
	if endTime != nil {
		params.Add("next_run_time[lt]", endTime.Format(time.RFC3339))
	}

	queryStr := params.Encode()

	req, err := http.NewRequest("GET", utils.If(queryStr == "", api.schedulesURL, fmt.Sprintf("%s?%s", api.schedulesURL, queryStr)), nil)
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
		err = fmt.Errorf("failed to fetch scheduled jobs: %s", res.Status)
		return
	}

	var apiResponse APIResponse[[]models.JobSchedule]
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		api.logger.Error("failed to unmarshal scheduled jobs", &err)
		return
	}

	schedules = &apiResponse.Data

	return
}

func (api *JobAPI) GetJob(id string) (job *models.Job, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", api.jobsURL, id), nil)
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
		err = fmt.Errorf("failed to fetch job: %s", res.Status)
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

func (api *JobAPI) UpdateJob(id string, updates *models.JobUpdateRequest) (updatedJob *models.Job, err error) {
	updatesJSON, err := json.Marshal(updates)
	if err != nil {
		api.logger.Error("failed to marshal job updates", &err)
		err = errors.New("failed to marshal job updates")
		return
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", api.jobsURL, id), strings.NewReader(string(updatesJSON)))
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

	updatedJob = &apiResponse.Data

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
	api.logger.ErrorWithData(fmt.Sprintf("failed to fetch scheduled jobs: %s", res.Status), &err, &data)
}
