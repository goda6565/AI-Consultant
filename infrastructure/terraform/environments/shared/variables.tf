variable "project_id" {
  type        = string
  description = "The project ID"
  default     = "ai-consultant-471305"
}

variable "region" {
  type        = string
  description = "The default region"
  default     = "asia-northeast1"
}

variable "backend_bucket_name" {
  type        = string
  description = "The name of the backend bucket"
  default     = "ai-consultant-tf-backend"
}

variable "service" {
  type        = string
  description = "The name of the service"
  default     = "ai-consultant"
}

variable "environment" {
  type        = string
  description = "The name of the environment"
  default     = "shared"
}