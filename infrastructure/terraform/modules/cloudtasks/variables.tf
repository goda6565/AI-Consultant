variable "name" {
  description = "Name of the Cloud Tasks queue"
  type        = string
}

variable "location" {
  description = "Location for the Cloud Tasks queue"
  type        = string
}

variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "service" {
  description = "Service name"
  type        = string
}

variable "max_dispatches_per_second" {
  description = "Maximum dispatches per second"
  type        = number
  default     = 1
}

variable "max_concurrent_dispatches" {
  description = "Maximum concurrent dispatches"
  type        = number
  default     = 1
}

variable "max_attempts" {
  description = "Maximum retry attempts"
  type        = number
  default     = 5
}

variable "min_backoff" {
  description = "Minimum backoff duration"
  type        = string
  default     = "10s"
}

variable "max_backoff" {
  description = "Maximum backoff duration"
  type        = string
  default     = "300s"
}
