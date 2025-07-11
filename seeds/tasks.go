package seeds

import (
	"context"

	"github.com/go-faker/faker/v4"
	"github.com/guregu/dynamo/v2"
	"github.com/julianstephens/distributed-task-scheduler/pkg/model"
)

func (s Seed) TaskSeed() {
	table := s.db.Table(s.conf.TaskTableName)

	for range s.count {
		var task model.Task

		_ = faker.FakeData(&task)
		task.Version = s.conf.TaskTableVersion

		if err := table.Put(dynamo.AWSEncoding(task)).Run(context.Background()); err != nil {
			panic(err)
		}
	}
}
