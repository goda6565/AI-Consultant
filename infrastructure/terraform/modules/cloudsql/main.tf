resource "google_sql_database_instance" "instance" {
  name                = var.name
  region              = var.region
  database_version    = var.database_version
  deletion_protection = false

  settings {
    tier = var.tier
    edition = "ENTERPRISE"

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }
  }
}

resource "google_sql_database" "blog_service_database" {
  for_each = toset(var.databases)

  name     = each.value
  instance = google_sql_database_instance.instance.name
}