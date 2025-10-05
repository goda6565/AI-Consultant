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
  default     = "staging"
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

variable "cloudbuild_service_account_id" {
  type        = string
  description = "The ID of the Cloud Build service account"
  default     = "projects/ai-consultant-471305/serviceAccounts/cloudbuild-deployer@ai-consultant-471305.iam.gserviceaccount.com"
}

variable "github_repository_id" {
  type        = string
  description = "The ID of the GitHub repository"
  default     = "projects/ai-consultant-471305/locations/us-west2/connections/github-connection/repositories/AI-Consultant"
}

variable "vector_db_host" {
  type        = string
  description = "The host of the vector database"
  default     = "10.18.208.3"
}

variable "app_db_host" {
  type        = string
  description = "The host of the app database"
  default     = "10.18.208.3"
}