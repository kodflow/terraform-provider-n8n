variable "n8n_api_key" {
  description = "N8N API key for authentication"
  type        = string
  sensitive   = true
}

variable "n8n_base_url" {
  description = "Base URL of the N8N instance"
  type        = string
}

variable "run_id" {
  description = "Unique run identifier for cattle-style resource naming"
  type        = string
  default     = "local"
}

variable "timestamp" {
  description = "Timestamp for guaranteed uniqueness (format: RUN_NUMBER-RUN_ATTEMPT)"
  type        = string
  default     = "0"
}
