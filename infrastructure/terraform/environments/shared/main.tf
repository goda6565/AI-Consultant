module "github_wip" {
  source = "../../modules/github-wip"

  project_id             = var.project_id
  github_organization    = "goda6565"
  github_repository_name = "AI-Consultant"
  backend_buckets        = [var.backend_bucket_name]
}