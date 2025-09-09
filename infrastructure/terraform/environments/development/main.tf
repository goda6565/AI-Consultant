module "documents_bucket" {
  source = "../../modules/cloudstorage"

  bucket_name = "${var.environment}-${var.service}-documents-bucket"
  region      = var.region
  environment = var.environment
  service     = var.service
}

module "backend_admin_cloudrun" {
  source = "../../modules/cloudrun"

  cloudrun_service_name = "${var.environment}-${var.service}-backend-admin"
  region                = var.region
  env_vars = [
    {
      name  = "ENV"
      value = "development"
    }
  ]
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
  env_vars = [
    {
      name  = "ENV"
      value = "development"
    }
  ]
  vpc_name             = var.vpc_id
  subnet_name          = var.subnet_id
  enable_public_access = true
  environment          = var.environment
  service              = var.service
}

module "backend_vector_cloudrun" {
  source = "../../modules/cloudrun"

  cloudrun_service_name = "${var.environment}-${var.service}-backend-vector"
  region                = var.region
  env_vars = [
    {
      name  = "ENV"
      value = "development"
    }
  ]
  vpc_name             = var.vpc_id
  subnet_name          = var.subnet_id
  enable_public_access = true
  environment          = var.environment
  service              = var.service
}