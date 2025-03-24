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

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket name.
* `domain_certificate` - (Required, List, ForceNew) The certificate of specified doamin.

The `certificate` object of `domain_certificate` supports the following:

* `cert_type` - (Required, String) Certificate type.
* `custom_cert` - (Required, List) Custom certificate.

The `custom_cert` object of `certificate` supports the following:

* `cert` - (Required, String) Public key of certificate.
* `private_key` - (Required, String) Private key of certificate.
* `cert_id` - (Optional, String) ID of certificate.

The `domain_certificate` object supports the following:

* `certificate` - (Required, List) Certificate info.
* `domain` - (Required, String) The name of domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



