variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default     = "meetup-devops"
}

variable "region" {
  description = "Region where resources are deployed"
  type        = string
  default     = "europe-west1"
}

variable "zones" {
  description = "Zones where resources are deployed"
  type        = list(string)
  default     = [
    "europe-west1-b",
    "europe-west1-c",
    "europe-west1-d",
  ]
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "meetup-keda"
}

variable "network_subnets" {
  description = "Subnets to create"
  type        = object({
    subnet_ip_cidr_range     = string
    pods_ip_cidr_range = string
    services_ip_cidr_range = string
  })
  default = {
    subnet_ip_cidr_range   = "10.48.64.0/18"
    pods_ip_cidr_range = "10.52.0.0/14"
    services_ip_cidr_range = "10.50.0.0/15"
  }
}