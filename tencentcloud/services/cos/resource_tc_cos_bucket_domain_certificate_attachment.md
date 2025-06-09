Provides a resource to attach/detach the corresponding certificate for the domain name in specified cos bucket.

~> **NOTE:** The current resource does not support cdc.

Example Usage

```hcl
variable "custom_origin_domain" {
  default = "tf.example.com"
}

data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket      = "private-bucket-${local.app_id}"
  acl         = "private"
  force_clean = true

  origin_domain_rules {
    domain = var.custom_origin_domain
    status = "ENABLED"
    type   = "REST"
  }
}

resource "tencentcloud_cos_bucket_domain_certificate_attachment" "example" {
  bucket = tencentcloud_cos_bucket.example.id
  domain_certificate {
    domain = var.custom_origin_domain
    certificate {
      cert_type = "CustomCert"
      custom_cert {
        cert_id     = "JG65alUy"
        cert        = <<-EOF
-----BEGIN CERTIFICATE-----
MIIGQjCCBSqgAwIBAgIQfTllN2vZr7vcoGF3ZTHwxjANBgkqhkiG9w0BAQsFADBA
...
...
...
9YSJrdvskqI3v/3SkVezzNiWQMuMTg==
-----END CERTIFICATE-----
EOF

        private_key = <<-EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAsmwAXXVh6N4fd281K0671jYBrSV2v/5+TCeewsNx6ys3kC8o
...
...
...
MgbOv6byAafSQWU+5+KFfK3Nj7eezx6yfQQM0Kxl4ZPm1w3Fb6gIFBc=
-----END RSA PRIVATE KEY-----
EOF
      }
    }
  }
}
```