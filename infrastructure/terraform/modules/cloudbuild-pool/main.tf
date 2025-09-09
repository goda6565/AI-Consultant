resource "google_cloudbuild_worker_pool" "pool" {
  name     = var.name
  location = "us-west2"
  worker_config {
    disk_size_gb   = 100
    machine_type   = "e2-medium"
    no_external_ip = false
  }
  network_config {
    peered_network          = var.vpc_id
    peered_network_ip_range = "/29"
  }
}