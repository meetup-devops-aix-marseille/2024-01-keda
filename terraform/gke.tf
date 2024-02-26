# google_client_config and kubernetes provider must be explicitly specified like the following.
data "google_client_config" "default" {}

# provider "kubernetes" {
#   host                   = "https://${module.gke.endpoint}"
#   token                  = data.google_client_config.default.access_token
#   cluster_ca_certificate = base64decode(module.gke.ca_certificate)
# }

# Define the gke service account name via a local variable to reuse it in the node pool definition to avoir the following error:
# local.node_pools is a map of map of string, known only after apply
# Issue https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/issues/991
locals {
  cluster_sa_name = "gke-${var.environment}-sa"
}

module "service_account" {
  source = "terraform-google-modules/service-accounts/google"
  version = "~> 4.2"

  project_id = var.project_id
  names = [
    local.cluster_sa_name,
  ]
  project_roles = [
    "${var.project_id}=>roles/artifactregistry.reader",
  ]
}

module "gke" {
  source                     = "terraform-google-modules/kubernetes-engine/google"
  version = "~> 27.0"

  project_id                 = var.project_id
  name                       = var.environment
  region                     = var.region
  zones                      = var.zones
  network                    = module.vpc.network_name
  subnetwork                 = module.vpc.subnets_names[0]
  ip_range_pods              = "${var.environment}-gke-pods"
  ip_range_services          = "${var.environment}-gke-services"
  http_load_balancing        = false
  network_policy             = false
  horizontal_pod_autoscaling = true
  filestore_csi_driver       = false
  gce_pd_csi_driver          = true
  release_channel            = "RAPID"
  identity_namespace         = "enabled"

  node_pools = [
    {
      name                      = "main-node-pool"
      machine_type              = "e2-standard-4"
      node_locations            = var.zones[0]
      min_count                 = 1
      max_count                 = 1
      local_ssd_count           = 0
      spot                      = false
      disk_size_gb              = 60
      disk_type                 = "pd-standard"
      image_type                = "COS_CONTAINERD"
      enable_gcfs               = false
      enable_gvnic              = false
      auto_repair               = true
      auto_upgrade              = true
      service_account           = format("%s@%s.iam.gserviceaccount.com", local.cluster_sa_name, var.project_id)
      preemptible               = false
      initial_node_count        = 1
    },
    {
      name                      = "cpu-node-pool"
      machine_type              = "e2-standard-4"
      node_locations            = var.zones[0]
      min_count                 = 0
      max_count                 = 10
      local_ssd_count           = 0
      spot                      = true
      disk_size_gb              = 60
      disk_type                 = "pd-standard"
      image_type                = "COS_CONTAINERD"
      enable_gcfs               = false
      enable_gvnic              = false
      auto_repair               = true
      auto_upgrade              = true
      service_account           = format("%s@%s.iam.gserviceaccount.com", local.cluster_sa_name, var.project_id)
      preemptible               = false
      initial_node_count        = 1
    },
  ]

  node_pools_oauth_scopes = {
    all = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/trace.append",
    ]
  }

  node_pools_labels = {
    cpu-node-pool = {
      worker = "meetup"
    }
  }

  node_pools_taints = {
    all = []

    cpu-node-pool = [
      {
        key    = "worker"
        value  = "meetup"
        effect = "NO_SCHEDULE"
      },
    ]
  }

  node_pools_tags = {
    all = [var.environment]
  }
}