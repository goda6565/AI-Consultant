output "cloudrun_service_url" {
  description = "URL of the Cloud Run service"
  value       = google_cloud_run_v2_service.cloudrun_service.uri
}

output "cloudrun_service_account_email" {
  description = "Email of the Cloud Run service account"
  value       = google_cloud_run_v2_service.cloudrun_service.template[0].service_account
}

