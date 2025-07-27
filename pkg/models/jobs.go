package models

import "time"

type Job struct {
	JobID         string    `binding:"-" json:"job_id"`
	UserID        string    `json:"user_id"`
	JobName       string    `json:"job_name"`
	Frequency     string    `json:"frequency"`
	Status        string    `json:"status"`
	Payload       string    `json:"payload"`
	RetryCount    int       `json:"retry_count"`
	MaxRetries    int       `json:"max_retries"`
	ExecutionTime time.Time `json:"execution_time"`
}

type JobSchedule struct {
	JobID       string    `json:"job_id"`
	NextRunTime time.Time `json:"next_run_time"`
	LastRunTime time.Time `json:"last_run_time"`
}

func (j *Job) GetJobFrequencyIntervalSeconds() int {
	switch j.Frequency {
	case "hourly":
		return int(time.Hour.Seconds())
	case "daily":
		return int(time.Hour.Seconds() * 24)
	}
	return -1
}

// type JobExecution struct {
// 	ExecutionID  string
// 	JobID        string
// 	WorkerID     string
// 	StartTime    time.Time
// 	EndTime      time.Time
// 	Status       string
// 	ErrorMessage string
// }

// type WorkerNode struct {
// 	WorkerID      string
// 	IPAddress     string
// 	Status        string
// 	LastHeartbeat string
// 	Capacity      int
// 	CurrentLoad   int
// }
