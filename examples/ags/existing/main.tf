# This example reads only; no import is needed and no resource is managed.
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
data "tencentcloud_ags_sandbox_tool" "existing" {
  tool_id = var.tool_id
}
variable "region" {
  type    = string
  default = "ap-guangzhou"
}
variable "profile" {
  type    = string
  default = "default"
}
variable "tool_id" {
  type = string
}
output "sandbox_tool_id" {
  value = data.tencentcloud_ags_sandbox_tool.existing.id
}
