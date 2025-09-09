# Cloud Storage

module "documents_bucket" {
  source = "../../modules/cloudstorage"

  bucket_name = "${var.environment}-${var.service}-documents-bucket"
  region      = var.region
  environment = var.environment
  service     = var.service
}

# Cloud Run

module "backend_admin_cloudrun" {
  source = "../../modules/cloudrun"

  cloudrun_service_name = "${var.environment}-${var.service}-backend-admin"
  region                = var.region
  env_vars              = local.common_env_vars
  vpc_name              = var.vpc_id
  subnet_name           = var.subnet_id
  enable_public_access  = true
  environment           = var.environment
  service               = var.service
}

module "backend_agent_cloudrun" {
  source = "../../modules/cloudrun"

  cloudrun_service_name = "${var.environment}-${var.service}-backend-agent"
  region                = var.region
  env_vars              = local.common_env_vars
  vpc_name              = var.vpc_id
  subnet_name           = var.subnet_id
  enable_public_access  = true
  environment           = var.environment
  service               = var.service
}

module "backend_vector_cloudrun" {
  source = "../../modules/cloudrun"

  cloudrun_service_name = "${var.environment}-${var.service}-backend-vector"
  region                = var.region
  env_vars              = local.common_env_vars
  vpc_name              = var.vpc_id
  subnet_name           = var.subnet_id
  enable_public_access  = true
  environment           = var.environment
  service               = var.service
}

# Cloud Build

module "backend_admin_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-backend-admin-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for backend admin development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/backend-admin.yaml"
  included_files                = ["backend/**", "infrastructure/deployments/cloudbuild/${var.environment}/backend-admin.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^develop$"
}

# Secret Manager
module "secret_manage_vector_db_password" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-vector-db-password"
  region      = var.region
}

module "secret_manage_vector_db_username" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-vector-db-username"
  region      = var.region
}

module "secret_manage_app_db_password" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-app-db-password"
  region      = var.region
}

module "secret_manage_app_db_username" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-app-db-username"
  region      = var.region
}
