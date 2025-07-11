package seeds

import (
	"log"
	"reflect"

	"github.com/guregu/dynamo/v2"
	"github.com/julianstephens/distributed-task-scheduler/pkg/logger"
)

type Seed struct {
	db    *dynamo.DB
	count int
}

func Execute(db *dynamo.DB, seedRowCount int, seedMethodNames ...string) {
	s := Seed{db: db, count: seedRowCount}

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
