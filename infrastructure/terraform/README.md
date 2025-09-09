# Terraform

This directory manages the infrastructure for AI Consultant.

## What's Not Managed
- Google Cloud Project
- Google Cloud API’s Enablement (Should be enabled manually)
- Terraform backend bucket

## Backend Bucket
Now, we use a single backend bucket for all environments.
- `ai-consultant-tf-backend`

## Naming Convention

#### 1. Google Cloud Resource ID/Name
- **Format**: `{environment}-{service}-{resource_type}`
- **Case**: kebab-case
- **Example**:
  - Service Account: `prod-ai-consultant-github-actions`
  - Storage Bucket: `dev-ai-consultant-documents-bucket`

#### 2. Environment Prefix
- `dev`: Development environment
- `stg`: Staging environment
- `prod`: Production environment
- `shared`: Shared environment

### Tag・Label Naming Conventions

#### Required Labels
```hcl
labels = {
  environment = "prod"           # dev/stg/prod/shared
  service     = "ai-consultant"  # service name
}
```