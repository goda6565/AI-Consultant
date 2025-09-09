locals {
  common_env_vars = [
    {
      name  = "ENV"
      value = "development"
    },
    {
      name  = "LISTEN_ADDRESS"
      value = "0.0.0.0:8080"
    },
    {
      name  = "SHUTDOWN_TIMEOUT"
      value = "10s"
    },
    {
      name  = "VECTOR_DB_HOST"
      value = "${var.vector_db_host}"
    },
    {
      name  = "VECTOR_DB_PORT"
      value = "5432"
    },
    {
      name  = "VECTOR_DB_NAME"
      value = "${var.environment}-vector-db"
    },
    {
      name  = "VECTOR_DB_SSL_MODE"
      value = "allow"
    },
    {
      name  = "APP_DB_HOST"
      value = "${var.app_db_host}"
    },
    {
      name  = "APP_DB_PORT"
      value = "5432"
    },
    {
      name  = "APP_DB_NAME"
      value = "${var.environment}-app-db"
    },
    {
      name  = "APP_DB_SSL_MODE"
      value = "allow"
    },
    {
      name  = "GOOGLE_CLOUD_PROJECT_ID"
      value = "165573575904"
    },
    {
      name  = "CLOUD_STORAGE_BUCKET_NAME"
      value = "${var.environment}-${var.service}-documents-bucket"
    },
    {
      name  = "DOCUMENT_AI_PROCESSOR_ID"
      value = "fdb37644f200783e"
    },
    {
      name  = "DOCUMENT_AI_LOCATION"
      value = "us"
    },
    {
      name  = "VERTEX_AI_LOCATION"
      value = "global"
    }
  ]
}