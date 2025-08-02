package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/scylladb/gocqlx/v3/qb"
)

func (base *Controller) CreateExecution(c *gin.Context) {
	var req models.JobExecution
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	req.ExecutionID = uuid.New().String()
	req.Status = models.JobStatusScheduled

	if err := base.DB.Client.Query(models.JobExecutions.Insert()).BindStruct(&req).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, errors.New("unable to create job execution"))
		return
	}

	httputil.NewResponse(c, req, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

func (base *Controller) UpdateExecution(c *gin.Context) {
	id := c.Param("id")

	var req models.JobExecutionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	var jobExec models.JobExecution
	stmt, names := qb.Select(models.JobExecutions.Name()).Where(qb.Eq("execution_id")).AllowFiltering().ToCql()
	if err := base.DB.Client.Query(stmt, names).Bind(id).SelectRelease(&jobExec); err != nil {
		httputil.NewError(c, http.StatusNotFound, fmt.Errorf("job execution %s not found", id))
		return
	}

	copier.Copy(&jobExec, &req)

	if err := base.DB.Client.Query(models.JobExecutions.Update()).BindStruct(&req).ExecRelease(); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, errors.New("unable to create job execution"))
		return
	}

	httputil.NewResponse(c, req, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Patch})
}
