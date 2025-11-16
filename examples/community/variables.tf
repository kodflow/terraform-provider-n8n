variable "n8n_base_url" {
  description = "n8n Base URL"
  type        = string
}

variable "n8n_api_key" {
  description = "n8n API Key"
  type        = string
  sensitive   = true
}
