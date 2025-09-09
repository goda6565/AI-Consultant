terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.34"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 6.34"
    }
  }
  required_version = ">= 1.12.0"
  backend "gcs" {
    bucket = "ai-consultant-tf-backend"
    prefix = "shared"
  }
}