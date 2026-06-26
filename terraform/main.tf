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

# API Gateway HTTP API — no Block Public Access restrictions
resource "aws_apigatewayv2_api" "api" {
  name          = local.function_name
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda" {
  api_id                 = aws_apigatewayv2_api.api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.api.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*"
}
