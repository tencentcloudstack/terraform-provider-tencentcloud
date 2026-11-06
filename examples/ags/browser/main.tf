# Local candidate only; replace 0.0.1-local with an accepted published version later.
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
  tool_type       = "browser"
  description     = "Terraform browser scenario"
  default_timeout = "5m"
  network         = { network_mode = "PUBLIC" }
  tags            = { scenario = "browser" }
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
