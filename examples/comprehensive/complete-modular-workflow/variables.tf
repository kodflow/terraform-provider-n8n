# Variables for complete modular workflow example
variable "n8n_base_url" {
  description = "Base URL of the n8n instance"
  type        = string
  default     = "http://localhost:5678"
}

variable "n8n_api_key" {
  description = "API key for n8n authentication"
  type        = string
  sensitive   = true
}

variable "run_id" {
  description = "Unique run identifier for CI/CD"
  type        = string
  default     = "local"
}
