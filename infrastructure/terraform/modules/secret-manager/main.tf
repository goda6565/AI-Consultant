resource "google_secret_manager_secret" "secret" {
  secret_id = var.secret_name

  replication {
    user_managed {
      replicas {
        location = var.region
      }
    }
  }
}