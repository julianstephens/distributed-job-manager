locals {
  app_name = "djm"
}

resource "aws_ssm_parameter" "api_key" {
  name  = "api_key"
  type = "String"
  value = var.api_key
}


resource "aws_dynamodb_table" "jobs" {
  name         = "${local.app_name}-jobs"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "user_id"
  range_key     = "job_id"
  attribute {
    name = "job_id"
    type = "S"
  }
  attribute {
    name = "user_id"
    type = "S"
  }
  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_dynamodb_table" "job_exec" {
  name         = "${local.app_name}-job-executions"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "execution_id"
  attribute {
    name = "execution_id"
    type = "S"
  }
  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_dynamodb_table" "job_sched" {
  name         = "${local.app_name}-job-schedules"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "job_id"
  attribute {
    name = "job_id"
    type = "S"
  }
  point_in_time_recovery {
    enabled = true
  }
}

resource "aws_dynamodb_table" "workers" {
  name         = "${local.app_name}-workers"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "worker_id"
  attribute {
    name = "worker_id"
    type = "S"
  }
  point_in_time_recovery {
    enabled = true
  }
}
