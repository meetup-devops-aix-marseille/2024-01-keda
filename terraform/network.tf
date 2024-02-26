# Create a VPC network
module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "7.3"

  project_id   = var.project_id
  network_name = var.environment
  routing_mode = "GLOBAL"

  subnets = [
    {
      subnet_name   = "${var.environment}-gke"
      subnet_ip     = var.network_subnets.subnet_ip_cidr_range
      subnet_region = var.region
      description = "GKE node subnet"
    },
  ]

  secondary_ranges = {
    "${var.environment}-gke" = [
      {
        range_name    = "${var.environment}-gke-pods"
        ip_cidr_range = var.network_subnets.pods_ip_cidr_range
      },
      {
        range_name    = "${var.environment}-gke-services"
        ip_cidr_range = var.network_subnets.services_ip_cidr_range
      },
    ],
  }
}