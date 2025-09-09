locals {
  github_oidc_workload_identity_user = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_actions_pool.name}/attribute.repository/${var.github_organization}/${var.github_repository_name}"
}
