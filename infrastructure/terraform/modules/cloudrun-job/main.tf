resource "google_cloud_run_v2_job" "cloudrun_job" {
  name                = var.cloudrun_job_name
  location            = var.region
  deletion_protection = false

  template {
    template {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/job:latest"

        resources {
          limits = {
            cpu    = var.cpu
            memory = var.memory
          }
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

      service_account = var.service_account_email
      max_retries     = var.max_retries
      timeout         = var.timeout
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].template[0].containers[0].image,
    ]
  }

  labels = {
    environment = var.environment
    service     = var.service
  }
}

resource "google_cloud_run_v2_job_iam_member" "runner" {
  for_each = toset(var.runners)
  name     = google_cloud_run_v2_job.cloudrun_job.name
  location = google_cloud_run_v2_job.cloudrun_job.location
  project  = google_cloud_run_v2_job.cloudrun_job.project
  role     = "roles/run.developer"
  member   = each.value
}


