package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/service"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/httputil"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model/table"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// GetTasks godoc
// @Summary Get all tasks
// @Description retrieves all tasks
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]table.Task]
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks [get]
func (base *Controller) GetTasks(c *gin.Context) {
	logger.Infof("retrieving all tasks")

	tasks, err := service.GetAll[table.Task](base.DB)
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
// @Success 200 {object} httputil.HTTPResponse[table.Task]
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks/{task_id} [get]
func (base *Controller) GetTask(c *gin.Context) {
	id := c.Param("task_id")
	logger.Infof("retrieving task %s", id)

	task, err := service.FindById[table.Task](base.DB, id)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, task, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// CreateTask godoc
// @Summary Create a task
// @Description creates a new task
// @Param data body table.Task true "Task data"
// @Tags tasks
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[table.Task]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks [post]
func (base *Controller) CreateTask(c *gin.Context) {
	logger.Infof("creating task")

	var task table.Task

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
	}

	task.UpdatedAt = time.Now().Unix()

	newTask, err := service.Create[table.Task](base.DB, task)
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
// @Param task_id path string true "Task ID"
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[table.Task]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /tasks/{task_id} [delete]
func (base *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("task_id")
	logger.Infof("deleting task %s", id)

	_, err := service.Delete(base.DB, id, table.Task{})
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, "Record deleted", httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}
