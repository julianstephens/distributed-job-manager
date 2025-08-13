package models

import "time"

type Job struct {
	JobID          string    `binding:"-" json:"job_id"`
	UserID         string    `binding:"-" json:"user_id"`
	JobName        string    `json:"job_name"`
	JobDescription string    `json:"job_description"`
	JobMetadata    string    `binding:"-" json:"job_metadata"`
	Frequency      string    `json:"frequency"`
	Status         string    `binding:"-" json:"status"`
	Payload        string    `json:"payload"`
	RetryCount     int       `binding:"-" json:"retry_count"`
	MaxRetries     int       `json:"max_retries"`
	ExecutionTime  time.Time `json:"execution_time"`
	CreatedAt      time.Time `binding:"-" json:"created_at"`
	UpdatedAt      time.Time `binding:"-" json:"updated_at"`
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

const (
	JobStatusReady      = "ready"
	JobStatusPending    = "pending"
	JobStatusScheduled  = "scheduled"
	JobStatusInProgress = "in-progress"
	JobStatusCompleted  = "completed"
	JobStatusCancelled  = "cancelled"
	JobStatusFailed     = "failed"
)

const (
	JobFrequencyOnce    = "one-time"
	JobFrequencyDaily   = "daily"
	JobFrequencyWeekly  = "weekly"
	JobFrequencyMonthly = "monthly"
)

type JobUpdateRequest struct {
	JobName        *string    `json:"job_name"`
	JobDescription *string    `json:"job_description"`
	Frequency      *string    `json:"frequency"`
	Status         *string    `json:"status"`
	Payload        *string    `json:"payload"`
	MaxRetries     *int       `json:"max_retries"`
	ExecutionTime  *time.Time `json:"execution_time"`
}

type JobSchedule struct {
	JobID       string    `json:"job_id"`
	NextRunTime time.Time `json:"next_run_time"`
	LastRunTime time.Time `json:"last_run_time"`
}

type JobScheduleUpdateRequest struct {
	NextRunTime *time.Time `json:"next_run_time"`
	LastRunTime *time.Time `json:"last_run_time"`
}

type JobExecution struct {
	ExecutionID  string    `binding:"-" json:"execution_id"`
	JobID        string    `json:"job_id"`
	WorkerID     string    `json:"worker_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `binding:"-" json:"status"`
	Output       string    `json:"output"`
	ErrorMessage string    `json:"error_message"`
}

type JobExecutionUpdateRequest struct {
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Status       *string    `json:"status"`
	Output       *string    `json:"output"`
	ErrorMessage *string    `json:"error_message"`
}

// type WorkerNode struct {
// 	WorkerID      string
// 	IPAddress     string
// 	Status        string
// 	LastHeartbeat string
// 	Capacity      int
// 	CurrentLoad   int
// }
