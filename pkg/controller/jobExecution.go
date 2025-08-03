package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/repository"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
)

type ExecutionController struct {
	Controller
	repo *repository.ExecutionRepository
}

func NewExecutionController(db *store.DBSession, config *models.Config, logger *graylogger.GrayLogger) *ExecutionController {
	return &ExecutionController{
		Controller: Controller{
			DB:     db,
			Config: config,
			Logger: logger,
		},
		repo: repository.NewExecutionRepository(db, logger),
	}
}

func (e *ExecutionController) CreateExecution(c *gin.Context) {
	var req models.JobExecution
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	exec, err := e.repo.CreateExecution(req)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
	}

	httputil.NewResponse(c, *exec, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

func (e *ExecutionController) UpdateExecution(c *gin.Context) {
	id := httputil.GetId(c)

	var req models.JobExecutionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	exec, err := e.repo.UpdateExecution(req, id)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, *exec, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Patch})
}
