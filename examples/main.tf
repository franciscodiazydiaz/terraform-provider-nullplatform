locals {
  dimensions = {
    "environment" = lower(var.environment),
    "country"     = "arg"
  }
}

resource "nullplatform_scope" "test" {
  scope_name          = "${var.environment}-terraform-test-00"
  null_application_id = var.null_application_id

  lambda_function_name            = "test-00"
  lambda_current_function_version = "2"
  lambda_function_role            = "arn:aws:iam::300001300842:role/LambdaRole"
  lambda_function_main_alias      = upper(var.environment)
  lambda_function_warm_alias      = "WARM"

  capabilities_serverless_memory       = 512
  capabilities_serverless_handler_name = "thehandler"
  capabilities_serverless_runtime_id   = "java11"
  log_group_name                       = "/aws/lambda/test-00"

  dimensions = local.dimensions
}

resource "nullplatform_parameter" "param0" {
  nrn      = var.null_application_nrn
  name     = "LOG_LEVEL"
  variable = "LOG_LEVEL"
}

resource "nullplatform_parameter_value" "param0_value0" {
  parameter_id = nullplatform_parameter.param0.id
  nrn          = nullplatform_scope.test.nrn
  value        = "DEBUG"
}

resource "nullplatform_parameter_value" "param0_value1" {
  parameter_id = nullplatform_parameter.param0.id
  nrn          = var.null_application_nrn
  value        = "INFO"
  dimensions   = local.dimensions
}
