variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "github_organization" {
  description = "GitHub organization"
  type        = string
}

variable "github_repository_name" {
  description = "GitHub repository name"
  type        = string
}

variable "github_app_installation_id" {
  description = "GitHub App Installation ID"
  type        = string
}

variable "enable_roles" {
  description = "List of Google Cloud IAM roles to enable"
  type        = list(string)
}