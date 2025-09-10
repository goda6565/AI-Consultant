output "cloudrun_service_url" {
  description = "URL of the Cloud Run service"
  value       = google_cloud_run_v2_service.cloudrun_service.uri
}