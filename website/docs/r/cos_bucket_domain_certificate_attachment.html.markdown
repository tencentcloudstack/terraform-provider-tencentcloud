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

### Use cert_id

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
      }
    }
  }
}
```

### Use cert and key

```hcl
resource "tencentcloud_cos_bucket_domain_certificate_attachment" "example" {
  bucket = tencentcloud_cos_bucket.example.id
  domain_certificate {
    domain = var.custom_origin_domain
    certificate {
      cert_type = "CustomCert"
      custom_cert {
        cert = <<-EOF
-----BEGIN CERTIFICATE-----
MIIG1DCCBLygAwIBAgIQDpfXbVCbQpEy5NNNSXxeeDANBgkqhkiG9w0BAQsFADBb
***
***
***
ynZ7SbC03yR+gKZQDeTXrNP1kk5Qhe7jSXgw+nhbspe0q/M1ZcNCz+sPxeOwdCcC
gJE=
-----END CERTIFICATE-----
EOF

        private_key = <<-EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAlnWPIMF4BnVyezE7KCoL+7Y1OpJ8V76g1Q9EvwWRbHus8xSM
***
***
***
Z8SK8+vMkRO9T9PBsZVMYmtQ0EtOLFtElep59iI3Mb3SdRyu+sCPmw==
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

* `cert_id` - (Optional, String, ForceNew) ID of certificate.
* `cert` - (Optional, String, ForceNew) Public key of certificate.
* `private_key` - (Optional, String, ForceNew) Private key of certificate.

The `domain_certificate` object supports the following:

* `certificate` - (Required, List, ForceNew) Certificate info.
* `domain` - (Required, String, ForceNew) The name of domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



