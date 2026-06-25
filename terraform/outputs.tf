output "api_url" {
  description = "Public HTTPS URL of the ingredient-suggestion API"
  value       = aws_lambda_function_url.api.function_url
}
