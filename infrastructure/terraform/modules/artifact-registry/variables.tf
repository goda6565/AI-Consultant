variable "region" {
  description = "Google Cloud Region"
  type        = string
}

variable "repository_id" {
  description = "Artifact Registry Repository ID"
  type        = string
}

variable "description" {
  description = "Artifact Registry Repository Description"
  type        = string
}

variable "format" {
  description = "Artifact Registry Repository Format"
  type        = string
  default     = "DOCKER"
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "service" {
  description = "Service"
  type        = string
}