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
  region = var.region
}

# CLS delivery currently requires a custom Tool in the service. This draft
# composition validates locally but must be adapted to a compatible custom
# image before applying with log_configuration enabled.
# Every identifier is an existing, caller-owned dependency. The example
# neither creates these dependencies nor modifies their access policies.
module "environment" {
  source          = "../modules/sandbox-environment"
  tool_name       = var.tool_name
  tool_type       = "code-interpreter"
  description     = "Existing dependency composition"
  default_timeout = "5m"
  role_arn        = var.role_arn
  network = {
    network_mode = "VPC"
    vpc_config = {
      subnet_ids         = var.subnet_ids
      security_group_ids = var.security_group_ids
    }
  }
  storage_mounts = [
    {
      name       = "cos-data"
      mount_path = "/mnt/cos"
      read_only  = true
      cos = {
        endpoint    = var.cos_endpoint
        bucket_name = var.cos_bucket
        bucket_path = "/terraform-fixture"
      }
    },
    {
      name       = "cfs-data"
      mount_path = "/mnt/cfs"
      read_only  = true
      cfs = {
        file_system_id = var.cfs_file_system_id
        path           = "/terraform-fixture"
      }
    },
  ]
  log_configuration = {
    topic_id = var.cls_topic_id
    files    = ["/logs/terraform-e2e.log"]
  }
}

variable "region" {
  type    = string
  default = "ap-guangzhou"
}
variable "tool_name" { type = string }
variable "role_arn" { type = string }
variable "subnet_ids" { type = set(string) }
variable "security_group_ids" { type = set(string) }
variable "cos_endpoint" { type = string }
variable "cos_bucket" { type = string }
variable "cfs_file_system_id" { type = string }
variable "cls_topic_id" { type = string }

output "sandbox_tool_id" { value = module.environment.tool_id }
