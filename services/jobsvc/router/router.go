package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-job-manager/pkg/controller"
	"github.com/julianstephens/distributed-job-manager/pkg/graylogger"
	"github.com/julianstephens/distributed-job-manager/pkg/middleware"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	docs "github.com/julianstephens/distributed-job-manager/services/jobsvc/docs"
)

const BasePath = "/api/v1"

func Setup(conf *models.Config, db *store.DBSession, log *graylogger.GrayLogger) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = BasePath

	api := controller.Controller{
		DB:     db,
		Config: conf,
		Logger: log,
	}

	baseGroup := r.Group(BasePath, middleware.Guard())

	jobAPI := controller.NewJobController(db, conf, log)
	jobGroup := baseGroup.Group("/jobs", middleware.RequireScopes("read:jobs", "write:jobs"))
	{
		jobGroup.GET("/", jobAPI.GetJobs)
		jobGroup.GET("/:id", jobAPI.GetJob)
		jobGroup.POST("", jobAPI.CreateJob)
		jobGroup.PATCH("/:id", jobAPI.UpdateJob)
		jobGroup.DELETE("/:id", jobAPI.DeleteJob)
	}

	executionAPI := controller.NewExecutionController(db, conf, log)
	executionGroup := baseGroup.Group("/executions", middleware.RequireScopes("read:executions", "write:executions"))
	{
		executionGroup.POST("/", executionAPI.CreateExecution)
		executionGroup.PATCH("/:id", executionAPI.UpdateExecution)
	}

	scheduleGroup := baseGroup.Group("/schedules", middleware.RequireScopes("read:schedules", "write:schedules"))
	{
		scheduleGroup.GET("/", api.GetSchedules)
		scheduleGroup.GET("/:id", api.GetSchedule)
		scheduleGroup.POST("", api.CreateSchedule)
		scheduleGroup.PATCH("/:id", api.UpdateSchedule)
		scheduleGroup.DELETE("/:id", api.DeleteSchedule)
	}

	return r
}
