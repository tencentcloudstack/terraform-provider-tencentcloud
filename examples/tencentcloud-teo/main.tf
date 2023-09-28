terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
      # 通过version指定版本
      # version = "1.79.10"
    }
  }
}

# TEO
provider "tencentcloud" {
  region = "ap-guangzhou"

  secret_id  = ""
  secret_key = ""

}

# Domain
provider "tencentcloud" {
  alias      = "tfdomain"
  region     = "ap-guangzhou"
  secret_id  = ""
  secret_key = ""
}


variable "zone_name" {
  default = "tf-teo.com"
}

# cname
resource "tencentcloud_teo_zone" "zone" {
  area            = "overseas"
  alias_zone_name = "tftest"
  paused          = false
  plan_id         = "edgeone-2kfv1h391n6w"
  tags = {
    "createdBy" = "terraform"
  }
  type      = "partial"
  zone_name = var.zone_name
}

resource "tencentcloud_dnspod_record" "demo" {
  provider = tencentcloud.tfdomain

  domain      = var.zone_name
  record_type = tencentcloud_teo_zone.zone.ownership_verification.0.dns_verification.0.record_type
  record_line = "默认"
  value       = tencentcloud_teo_zone.zone.ownership_verification.0.dns_verification.0.record_value
  sub_domain  = tencentcloud_teo_zone.zone.ownership_verification.0.dns_verification.0.subdomain

}

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [tencentcloud_dnspod_record.demo]
}

variable "sub_domain" {
  default = "aaa"
}

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id     = tencentcloud_teo_zone.zone.id
  domain_name = "${var.sub_domain}.${var.zone_name}"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  depends_on = [tencentcloud_teo_ownership_verify.ownership_verify]
}

resource "tencentcloud_dnspod_record" "acceleration_domain_record" {
  provider = tencentcloud.tfdomain

  domain      = var.zone_name
  record_type = "CNAME"
  record_line = "默认"
  value       = "${tencentcloud_teo_acceleration_domain.acceleration_domain.domain_name}.eo.dnse0.com."
  sub_domain  = var.sub_domain

}

resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = tencentcloud_teo_acceleration_domain.acceleration_domain.domain_name
  mode    = "eofreecert"
  zone_id = tencentcloud_teo_zone.zone.id

  depends_on = [tencentcloud_dnspod_record.acceleration_domain_record]
}