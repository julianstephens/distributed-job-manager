package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-task-scheduler/internal/models"
	"github.com/julianstephens/distributed-task-scheduler/internal/service"
	"github.com/julianstephens/distributed-task-scheduler/pkg/httputil"
	"github.com/julianstephens/distributed-task-scheduler/pkg/logger"
)

// GetTasks godoc
// @Summary Get all tasks
// @Description retrieves all tasks
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]models.Task]
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks [get]
func (base *Controller) GetTasks(c *gin.Context) {
	logger.Infof("retrieving all tasks")

	tasks, err := service.GetAll[models.Task](base.DB, "dts-tasks")
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	for _, t := range *tasks {
		logger.Infof("%v", t)
	}

	httputil.NewResponse(c, tasks, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// GetTask godoc
// @Summary Get a specific task
// @Description retrieves a task by its ID
// @Param task_id path string true "Task ID"
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[models.Task]
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks/{task_id} [get]
func (base *Controller) GetTask(c *gin.Context) {
	id := c.Param("id")
	logger.Infof("retrieving task %s", id)

	task, err := service.FindById[models.Task](base.DB, id, "dts-tasks")
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, task, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}
