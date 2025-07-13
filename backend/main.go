package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/config"
	"github.com/julianstephens/distributed-task-scheduler/backend/internal/router"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/database"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/httputil"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model/table"
	"github.com/julianstephens/distributed-task-scheduler/backend/seeds"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var models = []any{
	&table.Task{},
}

//	@title			DTS API
//	@version		0.1.0
//	@description	REST API for managing task scheduling

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes http
// @securityDefinitions.apikey ApiKey
// @in header
// @name X-API-KEY
// @description User-specific API key
func main() {
	godotenv.Load()
	handleArgs()
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()
	logger.Infof(strings.Join(args, ","))
	if len(args) >= 1 {
		switch args[0] {
		case "migrate":
			fmt.Println("HERE")
			stmts, err := gormschema.New("postgres").Load(models...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
				os.Exit(1)
			}
			_, err = io.WriteString(os.Stdout, stmts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write planned schema: %v\n", err)
			}
		case "seed":
			db := database.GetDB()
			masterSeedCount := 10
			seeds.Execute(db, masterSeedCount, args[1:]...)
			os.Exit(0)
		case "start":
			conf := config.GetConfig()
			r := router.Setup(conf)
			r.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

			r.NoRoute(func(c *gin.Context) {
				httputil.NewError(c, http.StatusNotFound, fmt.Errorf("resource not found"))
			})

			logger.Infof("DTS Server starting at %s:%s", conf.Host, conf.Port)
			logger.Fatalf("%v", r.Run(conf.Host+":"+conf.Port))
		}
	}
}
