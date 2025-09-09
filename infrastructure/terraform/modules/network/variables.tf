variable "region" {
  description = "The region to deploy the network"
  type        = string
}

variable "vpc_name" {
  description = "The name of the VPC network"
  type        = string
}

variable "subnet_name" {
  description = "The name of the subnet"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/24"
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