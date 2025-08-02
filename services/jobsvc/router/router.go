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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key"},
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

	privateGroup := r.Group(BasePath, middleware.Guard())

	jobGroup := privateGroup.Group("/jobs")
	{
		jobGroup.GET("/", api.GetJobs)
		jobGroup.GET("/:id", api.GetJob)
		jobGroup.POST("", api.CreateJob)
		jobGroup.PATCH("/:id", api.UpdateJob)
		jobGroup.DELETE("/:id", api.DeleteJob)
	}

	return r
}
