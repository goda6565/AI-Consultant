variable "project_id" {
  type        = string
  description = "The ID of the project to deploy to"
}

variable "github_organization" {
  type        = string
  description = "The name of the GitHub organization"
}

variable "github_repository_name" {
  type        = string
  description = "The name of the GitHub repository"
}

variable "backend_buckets" {
  type        = list(string)
  description = "The list of backend buckets"
}