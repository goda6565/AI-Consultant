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
  databases        = ["dev_app_db", "dev_vector_db", "staging_app_db", "staging_vector_db"]
}

module "cloudbuild-connection" {
  source = "../../modules/cloudbuild-connection"

  project_id                 = var.project_id
  github_organization        = "goda6565"
  github_repository_name     = "AI-Consultant"
  github_app_installation_id = "58532167"
  enable_roles = [
    "roles/artifactregistry.writer",
    "roles/cloudbuild.workerPoolUser",
    "roles/run.admin",
    "roles/secretmanager.secretAccessor",
    "roles/iam.serviceAccountUser",
    "roles/logging.logWriter",
  ]
}

module "cloudbuild-pool" {
  source = "../../modules/cloudbuild-pool"

  name   = "${var.environment}-${var.service}-cloudbuild-pool"
  vpc_id = module.network.vpc_id
}

# secret manager
module "secret_manager_identity_platform_api_key" {
  source = "../../modules/secret-manager"

  secret_name = "${var.environment}-${var.service}-identity-platform-api-key"
  region      = var.region
}

module "secret_manager_identity_platform_auth_domain" {
  source = "../../modules/secret-manager"

  secret_name = "${var.environment}-${var.service}-identity-platform-auth-domain"
  region      = var.region
}

module "secret_manager_redis_url" {
  source = "../../modules/secret-manager"

  secret_name = "${var.environment}-${var.service}-redis-url"
  region      = var.region
}