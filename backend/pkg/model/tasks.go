package model

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
	ID            string         `faker:"-" json:"id" dynamodbav:"id" binding:"-"`
	Title         string         `faker:"word" json:"title" dynamodbav:"title" binding:"required"`
	Description   string         `faker:"sentence" json:"description" dynamodbav:"description"`
	Status        TaskStatus     `faker:"oneof: 0" json:"status" dynamodbav:"status" binding:"required"`
	Recurrence    TaskRecurrence `faker:"oneof: 0" json:"recurrence" dynamodbav:"recurrence" binding:"required"`
	ScheduledTime int64          `faker:"unix_time" json:"scheduledTime" dynamodbav:"scheduledTime" binding:"required"`
	CreatedAt     int64          `faker:"-" json:"createdAt" dynamodbav:"createdAt" binding:"-"`
	UpdatedAt     int64          `faker:"-" json:"updatedAt" dynamodbav:"updatedAt" binding:"-"`
	Version       string         `faker:"-" json:"version" dynamodbav:"version" binding:"-"`
}
