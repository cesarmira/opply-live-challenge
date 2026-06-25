variable "aws_region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "opencode_api_key" {
  description = "API key for opencode.ai / DeepSeek"
  type        = string
  sensitive   = true
}

variable "opencode_base_url" {
  description = "Base URL for the LLM API"
  type        = string
}
