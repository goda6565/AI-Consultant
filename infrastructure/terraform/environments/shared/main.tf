module "github_wip" {
  source = "../../modules/github-wip"

  project_id             = var.project_id
  github_organization    = "goda6565"
  github_repository_name = "AI-Consultant"
  backend_buckets        = [var.backend_bucket_name]
}

module "artifact_registry" {
  source = "../../modules/artifact-registry"

  region        = var.region
  repository_id = "${var.environment}-${var.service}-artifact-registry"
  description   = "Docker Repository for AI Consultant"
  format        = "DOCKER"
  environment   = var.environment
  service       = var.service
}

module "network" {
  source = "../../modules/network"

  region      = var.region
  vpc_name    = "${var.environment}-${var.service}-vpc"
  subnet_name = "${var.environment}-${var.service}-subnet"
  environment = var.environment
  service     = var.service
}

# Now all applications are running on the same database
module "cloudsql" {
  source = "../../modules/cloudsql"

  name             = "${var.environment}-${var.service}-database"
  region           = var.region
  database_version = "POSTGRES_17"
  tier             = "db-f1-micro"
  vpc_id           = module.network.vpc_id
  databases        = ["dev_app_db", "dev_vector_db"]
}
