resource "google_artifact_registry_repository" "artifact_registry_repository" {
  location      = var.region
  repository_id = var.repository_id
  description   = var.description
  format        = var.format
  labels = {
    environment = var.environment
    service     = var.service
  }
}