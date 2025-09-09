resource "google_storage_bucket" "storage_bucket" {
  name          = var.bucket_name
  location      = var.region
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  labels = {
    environment = var.environment
    service     = var.service
  }
}