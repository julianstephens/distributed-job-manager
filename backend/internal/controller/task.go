package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/service"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/httputil"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// GetTasks godoc
// @Summary Get all tasks
// @Description retrieves all tasks
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]model.Task]
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks [get]
func (base *Controller) GetTasks(c *gin.Context) {
	logger.Infof("retrieving all tasks")

	tasks, err := service.GetAll[model.Task](base.DB, base.Config.TaskTableName)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, tasks, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// GetTask godoc
// @Summary Get a specific task
// @Description retrieves a task by its ID
// @Param task_id path string true "Task ID"
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.Task]
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks/:id [get]
func (base *Controller) GetTask(c *gin.Context) {
	id := c.Param("id")
	logger.Infof("retrieving task %s", id)

	task, err := service.FindById[model.Task](base.DB, id, base.Config.TaskTableName)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, task, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// PutTask godoc
// @Summary Put a task
// @Description creates or updates a task
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.Task]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks [put]
func (base *Controller) PutTask(c *gin.Context) {
	logger.Infof("putting task")

	var task model.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		httputil.HandleFieldError(c, err)
		return
	}

	if task.ID == "" {
		id, err := gonanoid.New()
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		task.ID = id
		task.CreatedAt = time.Now().Unix()
		task.Version = base.Config.TaskTableVersion
	}

	task.UpdatedAt = time.Now().Unix()

	newTask, err := service.Put[model.Task](base.DB, task, base.Config.TaskTableName)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, newTask, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Put})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description deletes a specific task
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.Task]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks/:id [delete]
func (base *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	logger.Infof("deleting task %s", id)

	err := service.Delete[model.Task](base.DB, id, base.Config.TaskTableName)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, "Record deleted", httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}
