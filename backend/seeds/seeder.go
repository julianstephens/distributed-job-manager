package seeds

import (
	"log"
	"reflect"

	"github.com/julianstephens/distributed-task-scheduler/backend/internal/config"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
	"gorm.io/gorm"
)

type Seed struct {
	db    *gorm.DB
	conf  *model.Config
	count int
}

func Execute(db *gorm.DB, seedRowCount int, seedMethodNames ...string) {
	s := Seed{db: db, conf: config.GetConfig(), count: seedRowCount}

	seedType := reflect.TypeOf(s)

	if len(seedMethodNames) == 0 {
		logger.Infof("Running all seeders...")

		for i := range seedType.NumMethod() {
			method := seedType.Method(i)
			seed(s, method.Name)
		}
	} else {
		for _, method := range seedMethodNames {
			seed(s, method)
		}
	}
}

func seed(s Seed, seedMethodName string) {

	m := reflect.ValueOf(s).MethodByName(seedMethodName)
	if !m.IsValid() {
		log.Fatal("No method called ", seedMethodName)
	}

	logger.Infof("Seeding %s...", seedMethodName)
	m.Call(nil)

	logger.Infof("Seed %s succeeded", seedMethodName)
}
