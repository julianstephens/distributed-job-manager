package models

import "github.com/scylladb/gocqlx/v3/table"

var (
	Jobs = table.New(table.Metadata{
		Name: "jobs",
		Columns: []string{
			"job_id",
			"user_id",
			"job_name",
			"job_description",
			"job_metadata",
			"frequency",
			"status",
			"payload",
			"retry_count",
			"max_retries",
			"execution_time",
			"created_at",
			"updated_at",
		},
		PartKey: []string{
			"user_id",
		},
		SortKey: []string{
			"job_id",
			"status",
		},
	})

	JobSchedules = table.New(table.Metadata{
		Name: "job_schedules",
		Columns: []string{
			"job_id",
			"next_run_time",
			"last_run_time",
		},
		PartKey: []string{
			"job_id",
		},
		SortKey: []string{
			"next_run_time",
		},
	})

	JobExecutions = table.New(table.Metadata{
		Name: "job_executions",
		Columns: []string{
			"execution_id",
			"job_id",
			"worker_id",
			"start_time",
			"end_time",
			"status",
			"output",
			"error_message",
		},
		PartKey: []string{
			"job_id",
		},
		SortKey: []string{
			"worker_id",
			"status",
		},
	})
)
