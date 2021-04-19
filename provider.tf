terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
#      version = "1.53.0"
    }
  }
}

provider "tencentcloud" {
  region = "ap-guangzhou"
}

