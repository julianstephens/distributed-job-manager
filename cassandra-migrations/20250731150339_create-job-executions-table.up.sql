CREATE TABLE IF NOT EXISTS job_executions (
  execution_id text,
  job_id text,
  worker_id text,
  start_time timestamp,
  end_time timestamp,
  status text,
  output text,
  error_message text,
  PRIMARY KEY (job_id, worker_id, status)
);
