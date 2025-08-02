package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/repository"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
)

type JobController struct {
	Controller
	repo *repository.JobRepository
}

func NewJobController(db *store.DBSession, config *models.Config, logger *graylogger.GrayLogger) *JobController {
	return &JobController{
		Controller: Controller{
			DB:     db,
			Config: config,
			Logger: logger,
		},
		repo: repository.NewJobRepository(db, logger),
	}
}

// GetJobs godoc
// @Summary Get all jobs
// @Description retrieves all jobs
// @Tags jobs
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs [get]
func (j *JobController) GetJobs(c *gin.Context) {
	userId := httputil.GetUserId(c)
	isAdmin := c.GetBool("isAdmin")

	jobs, err := j.repo.GetJobs(userId, isAdmin)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, *jobs, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// GetJob godoc
// @Summary Get a specific jobs
// @Description retrieves a single job
// @Tags jobs
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs/:id [get]
func (j *JobController) GetJob(c *gin.Context) {
	userId := httputil.GetUserId(c)
	jobId := httputil.GetId(c)
	isAdmin := c.GetBool("isAdmin")

	job, err := j.repo.GetJob(jobId, userId, isAdmin)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, fmt.Errorf("unable to get job %s", jobId))
		return
	}

	httputil.NewResponse(c, *job, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// CreateJob godoc
// @Summary Create a job
// @Description creates a new job
// @Tags jobs
// @Security ApiKey
// @Success 201 {object} httputil.HTTPResponse[models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs [post]
func (j *JobController) CreateJob(c *gin.Context) {
	userId := httputil.GetUserId(c)

	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	res, err := j.repo.CreateJob(job, userId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, fmt.Errorf("unable to create job: %w", err))
		return
	}

	httputil.NewResponse(c, *res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

// UpdateJob godoc
// @Summary Update a job
// @Description updates a new job
// @Tags jobs
// @Security ApiKey
// @Success 201 {object} httputil.HTTPResponse[models.Job]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs/:id [patch]
func (j *JobController) UpdateJob(c *gin.Context) {
	userId := httputil.GetUserId(c)
	jobId := httputil.GetId(c)

	var jobUpdate models.JobUpdateRequest
	if err := c.ShouldBindJSON(&jobUpdate); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	job, err := j.repo.UpdateJob(jobUpdate, jobId, userId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, fmt.Errorf("unable to update job %s: %w", jobId, err))
		return
	}

	httputil.NewResponse(c, *job, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Patch})
}

// DeleteJob godoc
// @Summary Delete a job
// @Description removes an existing job
// @Tags jobs
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[string]
// @Failure 500 {object} httputil.HTTPError
// @Router /jobs/:id [delete]
func (j *JobController) DeleteJob(c *gin.Context) {
	jobId := httputil.GetId(c)

	if err := j.repo.DeleteJob(jobId); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, jobId, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}
