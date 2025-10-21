---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_domain_certificate_attachment"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_domain_certificate_attachment"
description: |-
  Provides a resource to attach/detach the corresponding certificate for the domain name in specified cos bucket.
---

# tencentcloud_cos_bucket_domain_certificate_attachment

Provides a resource to attach/detach the corresponding certificate for the domain name in specified cos bucket.

~> **NOTE:** The current resource does not support cdc.

## Example Usage

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
        cert_id = "JG65alUy"
        cert    = <<-EOF
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

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket name.
* `domain_certificate` - (Required, List, ForceNew) The certificate of specified doamin.

The `certificate` object of `domain_certificate` supports the following:

* `cert_type` - (Required, String, ForceNew) Certificate type.
* `custom_cert` - (Required, List, ForceNew) Custom certificate.

The `custom_cert` object of `certificate` supports the following:

* `cert` - (Required, String, ForceNew) Public key of certificate.
* `private_key` - (Required, String, ForceNew) Private key of certificate.
* `cert_id` - (Optional, String, ForceNew) ID of certificate.

The `domain_certificate` object supports the following:

* `certificate` - (Required, List, ForceNew) Certificate info.
* `domain` - (Required, String, ForceNew) The name of domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



