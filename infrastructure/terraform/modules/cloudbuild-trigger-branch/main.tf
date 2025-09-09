resource "google_cloudbuild_trigger" "github_trigger" {
  name            = var.trigger_name
  description     = var.trigger_description
  location        = "us-west2"
  service_account = var.cloudbuild_service_account_id

  repository_event_config {

    repository = var.github_repository_id

    push {
      branch = var.branch_pattern
    }
  }

  filename           = var.file_name
  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"

  included_files = var.included_files
} 