package seeds

import (
	"context"

	"github.com/go-faker/faker/v4"
	"github.com/guregu/dynamo/v2"
	"github.com/julianstephens/distributed-task-scheduler/internal/models"
)

func (s Seed) TaskSeed() {
	table := s.db.Table("dts-tasks")

	for range s.count {
		var task models.Task

		_ = faker.FakeData(&task)

		if err := table.Put(dynamo.AWSEncoding(task)).Run(context.Background()); err != nil {
			panic(err)
		}
	}
}
