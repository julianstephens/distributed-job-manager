package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
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
	user_id := c.Query("user_id")
	if user_id == "" {
		httputil.NewError(c, http.StatusBadRequest, errors.New("user_id query param is required"))
		return
	}

	var jobs []models.Job
	q := base.DB.Client.Query(models.Jobs.SelectAll()).BindMap(qb.M{"user_id": user_id})
	if err := q.SelectRelease(&jobs); err != nil {
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
	job_id, ok := c.Params.Get("id")
	if !ok {
		httputil.NewError(c, http.StatusBadRequest, errors.New("no job id provided"))
		return
	}

	var job models.Job
	stmt, _ := qb.Select(models.Jobs.Name()).Where(qb.EqNamed("job_id", job_id)).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, []string{"job_id"}).Bind(job_id).Get(&job); err != nil {
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
	job.JobID = ulid.Make().String()

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

	q := base.DB.Client.Query(models.Jobs.Insert()).BindStruct(job)
	if err := q.ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	jobSchedule := models.JobSchedule{
		JobID:       job.JobID,
		NextRunTime: job.ExecutionTime,
	}

	q = base.DB.Client.Query(models.JobSchedules.Insert()).BindStruct(jobSchedule)
	if err := q.ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, job, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
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
	job_id, ok := c.Params.Get("id")
	if !ok {
		httputil.NewError(c, http.StatusBadRequest, errors.New("no job id provided"))
		return
	}

	var job models.Job
	stmt, _ := qb.Select(models.Jobs.Name()).Where(qb.EqNamed("job_id", job_id)).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, []string{"job_id"}).Bind(job_id).Get(&job); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	stmt, _ = qb.Delete(models.Jobs.Name()).Where(qb.EqNamed("user_id", job.UserID)).Where(qb.EqNamed("job_id", job_id)).ToCql()
	if err := base.DB.Client.Query(stmt, []string{"user_id", "job_id"}).Bind(job.UserID, job_id).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, job_id, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}
