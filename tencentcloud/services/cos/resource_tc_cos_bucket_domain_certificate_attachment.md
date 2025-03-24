Provides a resource to attach/detach the corresponding certificate for the domain name in specified cos bucket.

~> **NOTE:** The current resource does not support cdc.

Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket      = "private-bucket-${local.app_id}"
  acl         = "private"
  force_clean = true
}

resource "tencentcloud_cos_bucket_domain_certificate_attachment" "example" {
  bucket = tencentcloud_cos_bucket.example.id
  domain_certificate {
    domain = "www.example.com"
    certificate {
      cert_type = "CustomCert"
      custom_cert {
        cert_id     = "Mbx45wts"
        cert        = "-----BEGIN CERTIFICATE-----"
        private_key = "-----BEGIN RSA PRIVATE_KEY-----"
      }
    }
  }
}
```