variable "n8n_api_url" {
  description = "n8n API URL"
  type        = string
}

variable "n8n_api_key" {
  description = "n8n API Key"
  type        = string
  sensitive   = true
}
