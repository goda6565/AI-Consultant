resource "google_pubsub_topic" "topic" {
  name = var.topic_name

  labels = {
    environment = var.environment
    service     = var.service
  }
}

resource "google_pubsub_subscription" "subscription" {
  name  = var.subscription_name
  topic = google_pubsub_topic.topic.name

  push_config {
    push_endpoint = var.push_endpoint

    oidc_token {
      service_account_email = google_service_account.pubsub_push_service_account.email
    }
  }


  labels = {
    environment = var.environment
    service     = var.service
  }
}


resource "google_service_account" "pubsub_push_service_account" {
  account_id   = "${var.environment}-${var.service}-pubsub-sa"
  display_name = "Service Account for ${var.topic_name} Pub/Sub Push"
  description  = "Service account used for Pub/Sub push delivery to Cloud Run"
}

resource "google_project_iam_member" "cloudrun_invoker" {
  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.pubsub_push_service_account.email}"
}

resource "google_project_iam_member" "pubsub_token_creator" {
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

data "google_project" "project" {
  project_id = var.project_id
}
