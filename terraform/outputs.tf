output "api_url" {
  description = "Public HTTPS URL of the ingredient-suggestion API"
  value       = aws_apigatewayv2_stage.default.invoke_url
}
