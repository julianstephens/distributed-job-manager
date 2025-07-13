locals {
  app_name = "dts"
}

resource "aws_ssm_parameter" "api_key" {
  name  = "api_key"
  type = "String"
  value = var.api_key
}
