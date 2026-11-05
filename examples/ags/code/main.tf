# Candidate validation only. 0.0.1-local is the isolated local build installed by the
# E2E harness, not a claimed Registry release. Pin the actual accepted release
# before publishing this example as a production tutorial.
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
  tool_type       = "code-interpreter"
  description     = "Terraform code scenario"
  default_timeout = "5m"
  network         = { network_mode = "PUBLIC" }
  tags            = { scenario = "code" }
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

output "sandbox_tool_id" {
  value = module.environment.tool_id
}
