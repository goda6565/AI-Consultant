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
  env_vars = concat(
    local.common_env_vars,
    [
      {
        name  = "SYNC_TARGET_URL"
        value = "${module.backend_vector_cloudrun.cloudrun_service_url}/webhook"
      }
    ]
  )
  vpc_name             = var.vpc_id
  subnet_name          = var.subnet_id
  enable_public_access = true
  environment          = var.environment
  service              = var.service
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
  ingress               = "INGRESS_TRAFFIC_INTERNAL_ONLY"
  environment           = var.environment
  service               = var.service
}

module "frontend_cloudrun" {
  source = "../../modules/cloudrun"

  cloudrun_service_name = "${var.environment}-${var.service}-frontend"
  region                = var.region
  env_vars = [
    {
      name  = "ENV"
      value = "staging"
    }
  ]
  vpc_name             = var.vpc_id
  subnet_name          = var.subnet_id
  enable_public_access = true
  environment          = var.environment
  service              = var.service
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
  branch_pattern                = "^main$"
}

module "backend_vector_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-backend-vector-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for backend vector development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/backend-vector.yaml"
  included_files                = ["backend/**", "infrastructure/deployments/cloudbuild/${var.environment}/backend-vector.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^main$"
}

module "backend_agent_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-backend-agent-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for backend agent development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/backend-agent.yaml"
  included_files                = ["backend/**", "infrastructure/deployments/cloudbuild/${var.environment}/backend-agent.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^main$"
}

module "migrate_app_db_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-migrate-app-db-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for migrate app db development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/migrate-app-db.yaml"
  included_files                = ["infrastructure/schemas/migrations/app/**", "infrastructure/deployments/cloudbuild/${var.environment}/migrate-app-db.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^main$"
}

module "migrate_vector_db_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-migrate-vector-db-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for migrate vector db development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/migrate-vector-db.yaml"
  included_files                = ["infrastructure/schemas/migrations/vector/**", "infrastructure/deployments/cloudbuild/${var.environment}/migrate-vector-db.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^main$"
}

module "frontend_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-frontend-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for frontend development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/frontend.yaml"
  included_files                = ["frontend/**", "infrastructure/deployments/cloudbuild/${var.environment}/frontend.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^main$"
}

module "backend_proposal_job_cloudbuild_trigger" {
  source = "../../modules/cloudbuild-trigger-branch"

  trigger_name                  = "${var.environment}-${var.service}-backend-proposal-job-cloudbuild-trigger"
  trigger_description           = "Cloud Build trigger for backend proposal job development"
  file_name                     = "infrastructure/deployments/cloudbuild/${var.environment}/backend-proposal-job.yaml"
  included_files                = ["backend/**", "infrastructure/deployments/cloudbuild/${var.environment}/backend-proposal-job.yaml"]
  cloudbuild_service_account_id = var.cloudbuild_service_account_id
  github_repository_id          = var.github_repository_id
  branch_pattern                = "^main$"
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

module "secret_manage_admin_api_url" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-admin-api-url"
  region      = var.region
}

module "secret_manage_agent_api_url" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-agent-api-url"
  region      = var.region
}

module "secret_manage_custom_search_api_key" {
  source      = "../../modules/secret-manager"
  secret_name = "${var.environment}-${var.service}-custom-search-api-key"
  region      = var.region
}

# Cloud Tasks
module "document_processing_cloudtasks" {
  source = "../../modules/cloudtasks"

  name                      = "${var.environment}-${var.service}-document-processing"
  location                  = var.region
  project_id                = var.project_id
  environment               = var.environment
  service                   = var.service
  max_dispatches_per_second = 5
  max_concurrent_dispatches = 2
  max_attempts              = 10
  min_backoff               = "1s"
  max_backoff               = "10s"
}

# Cloud Run Job

module "backend_proposal_job" {
  source = "../../modules/cloudrun-job"

  cloudrun_job_name = "${var.environment}-${var.service}-backend-proposal-job"
  region            = var.region
  env_vars          = local.common_env_vars
  vpc_name          = var.vpc_id
  subnet_name       = var.subnet_id
  runners           = ["serviceAccount:${module.backend_agent_cloudrun.cloudrun_service_account_email}"]
  environment       = var.environment
  service           = var.service
  cpu               = "1"
  memory            = "4Gi"
}