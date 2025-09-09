variable "trigger_name" {
  description = "Name of the trigger"
  type        = string
}

variable "trigger_description" {
  description = "Description of the trigger"
  type        = string
}

variable "file_name" {
  description = "Name of the file"
  type        = string
}

variable "included_files" {
  description = "Included files"
  type        = list(string)
}

variable "cloudbuild_service_account_id" {
  description = "ID of the Cloud Build service account"
  type        = string
}

variable "github_repository_id" {
  description = "ID of the GitHub repository"
  type        = string
}

variable "branch_pattern" {
  description = "Branch pattern"
  type        = string
}