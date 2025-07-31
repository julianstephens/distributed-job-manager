package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/store"
	"github.com/julianstephens/distributed-job-manager/services/jobsvc/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			DJM Job API
//	@version		0.1.0
//	@description	REST API for managing job submission

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes http
// @securityDefinitions.apikey ApiKey
// @in header
// @name X-API-KEY
// @description User-specific API key
func main() {
	conf := config.GetConfig()

	db, err := store.GetDB(conf.Cassandra.Keyspace)
	if err != nil {
		logger.Fatalf("unable to get cassandra store connection: %v", err)
		return
	}

	r := router.Setup(conf, db)
	r.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.NoRoute(func(c *gin.Context) {
		httputil.NewError(c, http.StatusNotFound, fmt.Errorf("endpoint not found"))
	})

	logger.Infof("DJM Job SVC starting at %s:%s", conf.JobService.Host, conf.JobService.Port)
	logger.Fatalf("%v", r.Run(conf.JobService.Host+":"+conf.JobService.Port))
}
