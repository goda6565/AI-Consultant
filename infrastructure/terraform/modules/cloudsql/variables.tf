variable "name" {
  description = "Name of the Cloud SQL instance"
  type        = string
}

variable "region" {
  description = "Google Cloud Region"
  type        = string
}

variable "database_version" {
  description = "Database version"
  type        = string
}

variable "tier" {
  description = "Database tier (machine type)"
  type        = string
  default     = "db-f1-micro"
}

variable "vpc_id" {
  description = "VPC network name for private IP"
  type        = string
}

variable "databases" {
  description = "List of database names to create"
  type        = list(string)
}