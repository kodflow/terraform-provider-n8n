variable "n8n_api_key" {
  description = "n8n API key for authentication"
  type        = string
  sensitive   = true
}

variable "n8n_base_url" {
  description = "Base URL of your n8n instance"
  type        = string
}

variable "run_id" {
  description = "Unique run identifier for cattle-style resource naming"
  type        = string
  default     = "local"
}

variable "name_suffix" {
  description = "Name suffix for resources (v1 = initial, v2 = renamed)"
  type        = string
  default     = "v1"
}

variable "project_id" {
  description = "Project ID for E2E test isolation"
  type        = string
  default     = ""
}
