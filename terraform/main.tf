terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  profile = "cesarmira"
  region  = var.aws_region
}

locals {
  function_name = "opply-ingredient-suggest"
  zip_path      = "${path.module}/../dist/function.zip"
}

resource "aws_lambda_function" "api" {
  function_name    = local.function_name
  role             = aws_iam_role.lambda_exec.arn
  filename         = local.zip_path
  source_code_hash = filebase64sha256(local.zip_path)

  runtime = "provided.al2023"
  handler = "server"

  # Lambda Web Adapter layer proxies Lambda events to the HTTP server on :8080
  layers = [
    "arn:aws:lambda:${var.aws_region}:753240598075:layer:LambdaAdapterLayerX86:24"
  ]

  environment {
    variables = {
      OPENCODE_API_KEY        = var.opencode_api_key
      OPENCODE_BASE_URL       = var.opencode_base_url
      PORT                    = "8080"
      AWS_LAMBDA_EXEC_WRAPPER = "/opt/bootstrap"
      READINESS_CHECK_PATH    = "/healthz"
    }
  }

  timeout     = 30
  memory_size = 256
}

resource "aws_lambda_function_url" "api" {
  function_name      = aws_lambda_function.api.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "public_url" {
  statement_id           = "FunctionURLAllowPublicAccess"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.api.function_name
  principal              = "*"
  function_url_auth_type = "NONE"
}
