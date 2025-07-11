package models

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusInProgress
	StatusCompleted
	StatusFailed
)

type Task struct {
	ID          string     `faker:"uuid_digit" json:"id" dynamodbav:"id"`
	Title       string     `faker:"word" json:"title" dynamodbav:"title"`
	Description string     `faker:"sentence" json:"description" dynamodbav:"description"`
	Status      TaskStatus `faker:"oneof: 0" json:"status" dynamodbav:"status"`
	CreatedAt   int64      `faker:"unix_time" json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt   int64      `faker:"unix_time" json:"updatedAt" dynamodbav:"updatedAt"`
}
