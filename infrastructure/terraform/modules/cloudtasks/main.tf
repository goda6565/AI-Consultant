resource "google_cloud_tasks_queue" "document_sync" {
  name     = var.name
  location = var.location

  rate_limits {
    max_dispatches_per_second = var.max_dispatches_per_second
    max_concurrent_dispatches = var.max_concurrent_dispatches
  }

  retry_config {
    max_attempts = var.max_attempts
    min_backoff  = var.min_backoff
    max_backoff  = var.max_backoff
  }
}

resource "google_service_account" "cloudtasks_service_account" {
  account_id   = "${var.environment}-${var.service}-tasks-sa"
  display_name = "Service Account for ${var.name} Cloud Tasks"
  description  = "Service account used for Cloud Tasks to invoke Cloud Run"
}

resource "google_project_iam_member" "cloudrun_invoker" {
  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.cloudtasks_service_account.email}"
}
