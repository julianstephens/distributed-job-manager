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

func (j *Job) GetJobFrequencyIntervalSeconds() int {
	switch j.Frequency {
	case "hourly":
		return int(time.Hour.Seconds())
	case "daily":
		return int(time.Hour.Seconds() * 24)
	}
	return -1
}

type JobSchedule struct {
	JobID       string    `json:"job_id"`
	NextRunTime time.Time `json:"next_run_time"`
	LastRunTime time.Time `json:"last_run_time"`
}

type JobExecution struct {
	ExecutionID  string    `json:"execution_id"`
	JobID        string    `json:"job_id"`
	WorkerID     string    `json:"worker_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
}

// type WorkerNode struct {
// 	WorkerID      string
// 	IPAddress     string
// 	Status        string
// 	LastHeartbeat string
// 	Capacity      int
// 	CurrentLoad   int
// }
