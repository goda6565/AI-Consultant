output "github_connection_id" {
  description = "ID of the GitHub connection"
  value       = google_cloudbuildv2_connection.github_connection.id
}

output "cloudbuild_service_account_id" {
  description = "ID of the Cloud Build service account"
  value       = google_service_account.cloudbuild_service_account.id
}

output "github_repository_id" {
  description = "ID of the GitHub repository"
  value       = google_cloudbuildv2_repository.github_repository.id
}