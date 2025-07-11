package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/julianstephens/distributed-task-scheduler/internal/config"
	"github.com/julianstephens/distributed-task-scheduler/internal/router"
	"github.com/julianstephens/distributed-task-scheduler/pkg/database"
	"github.com/julianstephens/distributed-task-scheduler/pkg/httputil"
	"github.com/julianstephens/distributed-task-scheduler/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/seeds"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
		case "seed":
			db, err := database.GetDB()
			if err != nil {
				log.Fatalf("unable to init dynamodb client, %v", err)
			}
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
