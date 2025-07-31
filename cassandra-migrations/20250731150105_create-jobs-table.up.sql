CREATE TABLE IF NOT EXISTS jobs(
  job_id text,
  user_id text,
  job_name text,
  job_description text,
  job_metadata text,
  frequency text,
  status text,
  payload text,
  retry_count int,
  max_retries int,
  execution_time timestamp,
  created_at timestamp,
  updated_at timestamp,
  PRIMARY KEY (user_id, job_id, status)
)
