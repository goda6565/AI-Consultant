resource "google_cloud_run_v2_service" "cloudrun_service" {
  name                = var.cloudrun_service_name
  location            = var.region
  deletion_protection = false
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "nginx:latest" # This is a dummy image to avoid terraform plan error. Cloud Build will override this.
      ports {
        container_port = 8080 # This is a dummy port to avoid terraform plan error. Cloud Build will override this.
      }

      dynamic "env" {
        for_each = var.env_vars
        content {
          name  = env.value.name
          value = env.value.value

          dynamic "value_source" {
            for_each = env.value.secret_key_ref != null ? [env.value.secret_key_ref] : []
            content {
              secret_key_ref {
                secret  = value_source.value.secret
                version = value_source.value.version
              }
            }
          }
        }
      }
    }

    dynamic "vpc_access" {
      for_each = var.vpc_name != null && var.subnet_name != null ? [1] : []
      content {
        network_interfaces {
          network    = var.vpc_name
          subnetwork = var.subnet_name
        }
      }
    }

    scaling {
      max_instance_count = var.max_instance_count
      min_instance_count = var.min_instance_count
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
      template[0].containers[0].ports[0].container_port,
    ]
  }

  labels = {
    environment = var.environment
    service     = var.service
  }
}

resource "google_cloud_run_v2_service_iam_member" "public_access" {
  count    = var.enable_public_access ? 1 : 0
  name     = google_cloud_run_v2_service.cloudrun_service.name
  location = google_cloud_run_v2_service.cloudrun_service.location
  project  = google_cloud_run_v2_service.cloudrun_service.project
  role     = "roles/run.invoker"
  member   = "allUsers"
}