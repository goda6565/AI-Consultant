variable "cloudrun_job_name" {
  description = "Name of the Cloud Run Job"
  type        = string
}

variable "region" {
  description = "Google Cloud Region"
  type        = string
}

variable "env_vars" {
  description = "List of environment variables for the Cloud Run job"
  type = list(object({
    name  = string
    value = optional(string)
    secret_key_ref = optional(object({
      secret  = string
      version = optional(string, "latest")
    }))
  }))
  default = []
}

variable "vpc_name" {
  description = "Name of the VPC network"
  type        = string
  default     = null
}

variable "subnet_name" {
  description = "Name of the subnet"
  type        = string
  default     = null
}

variable "service_account_email" {
  description = "Service Account email used by the job"
  type        = string
  default     = null
}

variable "max_retries" {
  description = "Maximum number of retries for job tasks"
  type        = number
  default     = 3
}

variable "timeout" {
  description = "Execution timeout for the job (e.g., 3600s)"
  type        = string
  default     = "3600s"
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "service" {
  description = "Service"
  type        = string
}

variable "runners" {
  description = "Members who can run the Cloud Run Job (roles/run.jobRunner)"
  type        = list(string)
  default     = []
}

variable "memory" {
  description = "Memory for the Cloud Run Job"
  type        = string
}

variable "cpu" {
  description = "CPU for the Cloud Run Job"
  type        = string
}