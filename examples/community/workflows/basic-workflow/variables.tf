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

variable "run_id" {
  description = "Unique run identifier for cattle-style resource naming"
  type        = string
  default     = "local"
}

variable "project_id" {
  description = "Project ID for E2E test isolation"
  type        = string
  default     = ""
}
