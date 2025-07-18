CREATE TABLE IF NOT EXISTS job_schedules (
  job_id text,
  next_run_time timestamp,
  last_run_time timestamp,
  PRIMARY KEY (job_id, next_run_time)
);
