variable "n8n_base_url" {
  description = "N8N instance URL"
  type        = string
}

variable "n8n_api_key" {
  description = "N8N API key"
  type        = string
  sensitive   = true
}
