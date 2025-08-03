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

type ScheduleController struct {
	Controller
	repo *repository.ScheduleRepository
}

func NewScheduleController(db *store.DBSession, conf *models.Config, log *graylogger.GrayLogger) *ScheduleController {
	return &ScheduleController{
		Controller: Controller{
			DB:     db,
			Config: conf,
			Logger: log,
		},
		repo: repository.NewScheduleRepository(db, log),
	}
}

func (s *ScheduleController) GetSchedules(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	schedules, err := s.repo.GetSchedules(queryParams)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, *schedules, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

func (s *ScheduleController) GetSchedule(c *gin.Context) {
	id := httputil.GetId(c)

	schedule, err := s.repo.GetSchedule(id)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, *schedule, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

func (s *ScheduleController) CreateSchedule(c *gin.Context) {
	var req models.JobSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	schedule, err := s.repo.CreateSchedule(req)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, *schedule, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

func (s *ScheduleController) UpdateSchedule(c *gin.Context) {
	id := httputil.GetId(c)

	var req models.JobScheduleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	schedule, err := s.repo.UpdateSchedule(id, req)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, *schedule, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Patch})
}

func (s *ScheduleController) DeleteSchedule(c *gin.Context) {
	id := httputil.GetId(c)

	if err := s.repo.DeleteSchedule(id); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, id, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}
