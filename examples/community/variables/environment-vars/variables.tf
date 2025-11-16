variable "n8n_base_url" {
  description = "N8N Base URL"
  type        = string
  default     = "http://localhost:5678"
}

variable "n8n_api_key" {
  description = "N8N API Key"
  type        = string
  sensitive   = true
}

variable "environment" {
  description = "Environment name (dev, staging, production)"
  type        = string
  default     = "development"
}
