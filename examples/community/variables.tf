# Community Edition Variables

variable "n8n_base_url" {
  description = "n8n instance URL"
  type        = string
}

variable "n8n_api_key" {
  description = "n8n API key"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "Project ID for E2E test isolation"
  type        = string
  default     = ""
}
