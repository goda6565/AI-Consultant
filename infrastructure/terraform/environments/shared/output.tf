output "cloudbuild_service_account_id" {
  description = "ID of the Cloud Build service account"
  value       = module.cloudbuild-connection.cloudbuild_service_account_id
}

output "github_repository_id" {
  description = "ID of the GitHub repository"
  value       = module.cloudbuild-connection.github_repository_id
}
