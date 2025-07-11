package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-task-scheduler/docs"
	"github.com/julianstephens/distributed-task-scheduler/internal/controller"
	"github.com/julianstephens/distributed-task-scheduler/pkg/database"
)

const BasePath = "/api/v1"

func Setup() *gin.Engine {
	r := gin.New()

	// f, err := os.OpenFile("ls.access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("failed to create access log file: %v", err)
	// } else {
	// 	gin.DefaultWriter = io.MultiWriter(f)
	// }

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = BasePath

	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	api := controller.Controller{DB: db}

	baseGroup := r.Group(BasePath)

	taskGroup := baseGroup.Group("/tasks")
	{
		taskGroup.GET("/", api.GetTasks)
		taskGroup.GET("/:id", api.GetTask)
	}

	return r
}
