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

variable "service" {
  type        = string
  description = "The name of the service"
  default     = "ai-consultant"
}

variable "environment" {
  type        = string
  description = "The name of the environment"
  default     = "dev"
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC"
  default     = "shared-ai-consultant-vpc"
}

variable "subnet_id" {
  type        = string
  description = "The ID of the subnet"
  default     = "shared-ai-consultant-subnet"
}