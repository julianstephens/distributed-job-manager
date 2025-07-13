package router

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-task-scheduler/backend/docs"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/controller"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/middleware"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/database"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
)

const BasePath = "/api/v1"

func Setup(conf *model.Config) *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = BasePath

	db := database.GetDB()

	api := controller.Controller{DB: db, Config: conf}

	fmt.Printf("%v", conf)

	baseGroup := r.Group(BasePath, middleware.AuthGuard())

	taskGroup := baseGroup.Group("/tasks")
	{
		taskGroup.GET("", api.GetTasks)
		taskGroup.GET("/:task_id", api.GetTask)
		taskGroup.POST("", api.CreateTask)
		taskGroup.DELETE("/:task_id", api.DeleteTask)
	}

	return r
}
