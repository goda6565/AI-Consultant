variable "cloudrun_service_name" {
  description = "Name of the Cloud Run service"
  type        = string
}

variable "region" {
  description = "Google Cloud Region"
  type        = string
}

variable "env_vars" {
  description = "List of environment variables for the Cloud Run service"
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

variable "max_instance_count" {
  description = "Maximum number of instances for the Cloud Run service"
  type        = number
  default     = 1
}

variable "min_instance_count" {
  description = "Minimum number of instances for the Cloud Run service"
  type        = number
  default     = 0
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "service" {
  description = "Service"
  type        = string
}

variable "enable_public_access" {
  description = "Enable public access to the Cloud Run service (allows unauthenticated requests)"
  type        = bool
  default     = false
}

variable "ingress" {
  description = "Ingress configuration for the Cloud Run service"
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}