package table

import (
	"github.com/julianstephens/distributed-task-scheduler/backend/pkg/model"
)

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusInProgress
	StatusCompleted
	StatusFailed
	StatusCancelled
)

type TaskRecurrence int

const (
	Once TaskRecurrence = iota
	Daily
	Weekly
	Monthly
)

type Task struct {
	model.Base
	Title         string         `faker:"word" json:"title"  binding:"required"`
	Description   string         `faker:"sentence" json:"description"`
	Status        TaskStatus     `faker:"oneof: 0" gorm:"index" json:"status"`
	Recurrence    TaskRecurrence `faker:"oneof: 0" json:"recurrence"`
	ScheduledTime int64          `faker:"unix_time" json:"scheduledTime" binding:"required"`
	CreatedAt     int64          `faker:"-" json:"createdAt" gorm:"autoUpdateTime" binding:"-"`
	UpdatedAt     int64          `faker:"-" json:"updatedAt" gorm:"autoCreateTime" binding:"-"`
}
