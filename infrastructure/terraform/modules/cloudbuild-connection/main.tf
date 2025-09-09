resource "google_cloudbuildv2_connection" "github_connection" {
  location = "us-west2"
  name     = "github-connection"

  github_config {
    app_installation_id = var.github_app_installation_id
    authorizer_credential {
      oauth_token_secret_version = data.google_secret_manager_secret_version.github_oauth_token.name
    }
  }
}

resource "google_cloudbuildv2_repository" "github_repository" {
  name              = var.github_repository_name
  parent_connection = google_cloudbuildv2_connection.github_connection.id
  remote_uri        = "https://github.com/${var.github_organization}/${var.github_repository_name}.git"
}

resource "google_project_iam_member" "cloudbuild_roles" {
  for_each = toset(var.enable_roles)

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.cloudbuild_service_account.email}"
}

data "google_secret_manager_secret_version" "github_oauth_token" {
  project = var.project_id
  secret  = "host-github-oauthtoken-0339fd"
  version = "latest"
}

resource "google_service_account" "cloudbuild_service_account" {
  project      = var.project_id
  account_id   = "cloudbuild-deployer"
  display_name = "Cloud Build Deployer"
  description  = "Service account for Cloud Build to deploy to Cloud Run"
}