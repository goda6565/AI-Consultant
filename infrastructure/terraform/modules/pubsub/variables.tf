variable "topic_name" {
  description = "Name of the Pub/Sub topic"
  type        = string
}

variable "subscription_name" {
  description = "Name of the Pub/Sub subscription"
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

variable "push_endpoint" {
  description = "Push endpoint URL for Cloud Run service"
  type        = string
}

variable "retry_minimum_backoff" {
  description = "Minimum retry backoff duration (e.g., '1s', '10s')"
  type        = string
  default     = "10s"
}

variable "retry_maximum_backoff" {
  description = "Maximum retry backoff duration (e.g., '10s', '600s')"
  type        = string
  default     = "600s"
}

variable "message_retention_duration" {
  description = "Message retention duration (minimum 10m, maximum 744h)"
  type        = string
  default     = "10m"
}
