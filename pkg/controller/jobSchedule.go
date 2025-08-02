package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
)

func (base *Controller) GetSchedules(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	var schedules []models.JobSchedule
	q, err := base.getFilteredQuery(queryParams, models.JobSchedules.Name(), models.JobSchedules.Metadata().Columns)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err := q.SelectRelease(&schedules); err != nil {
		base.Logger.Error("unable to get job schedules from db", &err)
		httputil.NewError(c, http.StatusInternalServerError, errors.New("unable to get job schedules"))
		return
	}

	httputil.NewResponse(c, schedules, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

func (base *Controller) GetSchedule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		httputil.NewError(c, http.StatusBadRequest, ErrMissingJobId)
		return
	}

	var schedule models.JobSchedule
	if err := base.DB.Client.Query(models.JobSchedules.Select()).Bind(id).Get(&schedule); err != nil {
		base.Logger.Error(fmt.Sprintf("unable to get job schedule %s from db", id), &err)
		httputil.NewError(c, http.StatusInternalServerError, errors.New("unable to get job schedule"))
		return
	}

	httputil.NewResponse(c, schedule, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

func (base *Controller) CreateSchedule(c *gin.Context) {
}

func (base *Controller) UpdateSchedule(c *gin.Context) {
}

func (base *Controller) DeleteSchedule(c *gin.Context) {
}
