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
