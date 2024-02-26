resource "google_artifact_registry_repository" "docker" {
  project       = var.project_id
  location      = var.region
  repository_id = "images"
  description   = "Docker images repository"
  format        = "DOCKER"
}
