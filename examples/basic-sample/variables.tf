variable "n8n_api_key" {
  description = "n8n API key for authentication"
  type        = string
  sensitive   = true
}

variable "n8n_base_url" {
  description = "Base URL of your n8n instance"
  type        = string
  default     = "https://n8n.example.com"
}

variable "api_key" {
  description = "API key for Mocky.io (can be any value for testing)"
  type        = string
  default     = "test-api-key-12345"
  sensitive   = true
}
