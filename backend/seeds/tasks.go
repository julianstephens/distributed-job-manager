package seeds

import (
	"github.com/go-faker/faker/v4"
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model/table"
)

func (s Seed) TaskSeed() {
	for range s.count {
		var task table.Task

		_ = faker.FakeData(&task)

		if err := s.db.Create(&task).Error; err != nil {
			panic(err)
		}
	}
}
