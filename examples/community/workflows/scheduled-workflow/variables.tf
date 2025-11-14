variable "n8n_api_url" {
  description = "N8N API URL"
  type        = string
  default     = "http://localhost:5678"
}

variable "n8n_api_key" {
  description = "N8N API Key"
  type        = string
  sensitive   = true
}
