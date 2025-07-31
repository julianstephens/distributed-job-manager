package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/utils"
	"github.com/oklog/ulid/v2"
	"github.com/scylladb/gocqlx/v3/qb"
)

// GetJobs godoc
// @Summary Get all jobs
// @Description retrieves all jobs
// @Tags jobs
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs [get]
func (base *Controller) GetJobs(c *gin.Context) {
	userId := httputil.GetId(c)

	logger.Infof("getting jobs for user %s", userId)

	var jobs []models.Job
	stmt, names := qb.Select(models.Jobs.Name()).Where(qb.Eq("user_id")).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(userId).SelectRelease(&jobs); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, jobs, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// GetJob godoc
// @Summary Get a specific jobs
// @Description retrieves a single job
// @Tags jobs
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs/:id [get]
func (base *Controller) GetJob(c *gin.Context) {
	jobId, ok := c.Params.Get("id")
	if !ok {
		httputil.NewError(c, http.StatusBadRequest, errors.New("no job id provided"))
		return
	}

	userId := httputil.GetId(c)

	logger.Infof("getting jobs for user %s", userId)

	var job models.Job
	stmt, names := qb.Select(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id")).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(jobId, userId).Get(&job); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, job, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// CreateJob godoc
// @Summary Create a job
// @Description creates a new job
// @Tags jobs
// @Security ApiKey
// @Success 201 {object} httputil.HTTPResponse[models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs [post]
func (base *Controller) CreateJob(c *gin.Context) {
	var job models.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	userId := httputil.GetId(c)

	now := time.Now().UTC()

	job.CreatedAt = now
	job.UpdatedAt = now
	job.JobID = ulid.Make().String()
	job.UserID = userId
	job.RetryCount = 0
	job.Status = models.JobStatusPending

	parser := &utils.Parser{}
	if err := parser.Parse(job.Payload); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	supportedLanguages := utils.GetSupportedLanguages()
	for _, block := range parser.Result {
		if supportedLanguages[block.Language] == "" {
			httputil.NewError(c, http.StatusBadRequest, fmt.Errorf("%s is not a supported code language", block.Language))
			return
		}
	}
	job.Payload = parser.SanitizedInput

	if err := base.DB.Client.Query(models.Jobs.Insert()).BindStruct(job).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	jobSchedule := models.JobSchedule{
		JobID:       job.JobID,
		NextRunTime: job.ExecutionTime,
	}

	if err := base.DB.Client.Query(models.JobSchedules.Insert()).BindStruct(jobSchedule).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, job, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

// UpdateJob godoc
// @Summary Update a job
// @Description updates a new job
// @Tags jobs
// @Security ApiKey
// @Success 201 {object} httputil.HTTPResponse[models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs/:id [patch]
func (base *Controller) UpdateJob(c *gin.Context) {
	userId := httputil.GetId(c)

	jobId := c.Param("id")
	if jobId == "" {
		httputil.NewError(c, http.StatusBadRequest, errors.New("no job id provided"))
		return
	}

	var jobUpdate models.JobUpdateRequest
	if err := c.ShouldBindJSON(&jobUpdate); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if jobUpdate.Payload != nil {
		parser := &utils.Parser{}
		if err := parser.Parse(*jobUpdate.Payload); err != nil {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		}
		supportedLanguages := utils.GetSupportedLanguages()
		for _, block := range parser.Result {
			if supportedLanguages[block.Language] == "" {
				httputil.NewError(c, http.StatusBadRequest, fmt.Errorf("%s is not a supported code language", block.Language))
				return
			}
		}
		jobUpdate.Payload = &parser.SanitizedInput
	}

	var job models.Job
	stmt, names := qb.Select(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id")).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(jobId, userId).Get(&job); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, fmt.Errorf("unable to get existing job %s", jobId))
		return
	}

	err := copier.Copy(&job, &jobUpdate)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, errors.New("Something went wrong"))
		return
	}

	stmt, names = qb.Delete(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id")).ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(job.JobID, job.UserID).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, errors.New("unable to update job"))
		return
	}

	if err := base.DB.Client.Query(models.Jobs.Insert()).BindStruct(job).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, errors.New("unable to update job"))
		return
	}

	if jobUpdate.ExecutionTime != nil {
		jobSchedule := models.JobSchedule{
			JobID:       job.JobID,
			NextRunTime: job.ExecutionTime,
		}

		if err := base.DB.Client.Query(models.JobSchedules.Update()).BindStruct(jobSchedule).ExecRelease(); err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
	}

	httputil.NewResponse(c, job, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Patch})
}

// DeleteJob godoc
// @Summary Delete a job
// @Description removes an existing job
// @Tags jobs
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[string]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs/:id [delete]
func (base *Controller) DeleteJob(c *gin.Context) {
	userId := httputil.GetId(c)

	jobId := c.Param("id")
	if jobId == "" {
		httputil.NewError(c, http.StatusBadRequest, errors.New("no job id provided"))
		return
	}

	var job models.Job
	stmt, names := qb.Select(models.Jobs.Name()).Where(qb.Eq("job_id"), qb.Eq("user_id")).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(jobId, userId).Get(&job); err != nil {
		httputil.NewError(c, http.StatusNotFound, fmt.Errorf("job %s not found", jobId))
		return
	}

	stmt, names = qb.Delete(models.Jobs.Name()).Where(qb.Eq("user_id"), qb.Eq("job_id")).Where(qb.EqNamed("job_id", jobId)).ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(userId, jobId).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, fmt.Errorf("unable to delete job %s", jobId))
		return
	}

	httputil.NewResponse(c, jobId, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}
