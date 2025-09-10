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
