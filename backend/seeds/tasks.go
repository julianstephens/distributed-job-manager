package seeds

import (
	"context"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/guregu/dynamo/v2"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (s Seed) TaskSeed() {
	table := s.db.Table(s.conf.TaskTableName)

	for range s.count {
		var task model.Task

		_ = faker.FakeData(&task)
		id, err := gonanoid.New()
		if err != nil {
			panic(err)
		}
		task.ID = id
		task.CreatedAt = time.Now().Unix()
		task.UpdatedAt = time.Now().Unix()
		task.Version = s.conf.TaskTableVersion

		if err := table.Put(dynamo.AWSEncoding(task)).Run(context.Background()); err != nil {
			panic(err)
		}
	}
}
