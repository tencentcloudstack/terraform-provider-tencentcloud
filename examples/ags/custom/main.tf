# Local candidate only; images/specifications must come from a verified fixture.
terraform {
  required_version = ">= 1.17.0"
  required_providers {
    tencentcloud = {
      source  = "tencentcloudstack/tencentcloud"
      version = "= 0.0.1-local"
    }
  }
}

provider "tencentcloud" {
  region  = var.region
  profile = var.profile
}

module "environment" {
  source          = "../modules/sandbox-environment"
  tool_name       = var.tool_name
  tool_type       = "custom"
  role_arn        = var.role_arn
  default_timeout = "5m"
  network         = { network_mode = "PUBLIC" }
  custom_configuration = {
    image               = var.image
    image_registry_type = var.image_registry_type
    env                 = { TF_CASE = "v1" }
    resources = {
      cpu    = var.cpu
      memory = var.memory
    }
  }
  tags = { scenario = "custom" }
}

variable "region" {
  type    = string
  default = "ap-guangzhou"
}
variable "profile" {
  type    = string
  default = "default"
}
variable "tool_name" {
  type = string
}
variable "image" {
  type        = string
  description = "Fixed digest of an existing, AGS-compatible runtime image."
}
variable "role_arn" {
  type        = string
  description = "Existing AGS data-access role for registry access. This example creates no role or IAM grants."
}
variable "image_registry_type" {
  type = string
}
variable "cpu" {
  type = string
}
variable "memory" {
  type = string
}
output "sandbox_tool_id" {
  value = module.environment.tool_id
}
