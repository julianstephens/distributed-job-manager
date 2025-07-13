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
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/aws/ddb"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
)

const BasePath = "/api/v1"

func Setup(conf *model.Config) *gin.Engine {
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = BasePath

	db, err := ddb.GetDB()
	if err != nil {
		panic(err)
	}

	api := controller.Controller{DB: db, Config: conf}

	fmt.Printf("%v", conf)

	baseGroup := r.Group(BasePath, middleware.AuthGuard())

	taskGroup := baseGroup.Group("/tasks")
	{
		taskGroup.GET("", api.GetTasks)
		taskGroup.GET("/:task_id", api.GetTask)
		taskGroup.PUT("/", api.PutTask)
		taskGroup.DELETE("/:task_id", api.DeleteTask)
	}

	return r
}
