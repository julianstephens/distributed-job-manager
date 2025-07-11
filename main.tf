locals {
  app_name = "dts"
}



resource "aws_dynamodb_table" "tasks" {
  name         = "${local.app_name}-tasks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
  attribute {
    name = "id"
    type = "S"
  }
  point_in_time_recovery {
    enabled = true
  }
}
