# プロジェクト情報を取得
data "google_project" "project" {}

# Create a workload identity pool for GitHub Actions
resource "google_iam_workload_identity_pool" "github_actions_pool" {
  project                   = var.project_id
  workload_identity_pool_id = "github-actions-wip-pool"
  display_name              = "GitHub Actions WIP"
  description               = "GitHub Actions Workload Identity Pool for ${var.github_repository_name}"
  disabled                  = false
}

# Set up an OIDC provider for GitHub Actions
resource "google_iam_workload_identity_pool_provider" "github_actions_pool_provider" {
  project                            = var.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_actions_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-actions-pool-provider"
  display_name                       = "GitHub Actions WIP Provider"
  description                        = "GitHub Actions Workload Identity Pool Provider for ${var.github_repository_name}"
  disabled                           = false
  attribute_condition                = "assertion.repository == \"${var.github_organization}/${var.github_repository_name}\""
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }
  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

# Get the basic role for the storage object user
data "google_iam_role" "roles_storage_object_user" {
  name = "roles/storage.objectUser"
}

# Create a custom storage object user role for the Terraform backend
resource "google_project_iam_custom_role" "custom_storage_object_user" {
  role_id     = "customStorageObjectUser"
  title       = "customStorageObjectUser"
  description = "Custom role for storage object user"
  permissions = setsubtract(
    concat([
      "storage.buckets.getIamPolicy",
      ],
      data.google_iam_role.roles_storage_object_user.included_permissions,
    ),
    [
      "resourcemanager.projects.list",
    ],
  )
  project = var.project_id
}

# Grant access to the Terraform backend bucket
resource "google_storage_bucket_iam_member" "terraform_executor_storage_object_user" {
  for_each = toset(var.backend_buckets)
  bucket   = each.value
  role     = google_project_iam_custom_role.custom_storage_object_user.id
  member   = "serviceAccount:${google_service_account.terraform_executor.email}"
}

# Create Terraform executor service account and grant owner role to it.
resource "google_service_account" "terraform_executor" {
  project      = var.project_id
  account_id   = "terraform-executor"
  display_name = "terraform-executor"
}

resource "google_project_iam_member" "terraform_executor" {
  project = var.project_id
  role    = "roles/owner"
  member  = "serviceAccount:${google_service_account.terraform_executor.email}"
}

resource "google_service_account_iam_member" "terraform_workload_identity_user" {
  service_account_id = google_service_account.terraform_executor.name
  role               = "roles/iam.workloadIdentityUser"
  member             = local.github_oidc_workload_identity_user
}

# Enable services
resource "google_project_service" "service" {
  provider = google-beta
  for_each = toset([
    "iamcredentials.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
  ])

  service                               = each.value
  disable_on_destroy                    = false # Do not delete the service
  check_if_service_has_usage_on_destroy = true  # Check if the service is being used
}