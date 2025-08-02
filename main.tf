locals {
  app_name = "djm"
}

resource "aws_cloudwatch_log_group" "jobsvc" {
  name_prefix = "${local.app_name}-jobsvc"
}


resource "aws_cloudwatch_log_group" "schedulingsvc" {
  name_prefix = "${local.app_name}-schedsvc"
}
