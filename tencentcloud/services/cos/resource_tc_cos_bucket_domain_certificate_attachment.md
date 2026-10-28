Provides a resource to attach/detach the corresponding certificate for the domain name in specified cos bucket.

~> **NOTE:** The current resource does not support cdc.

Example Usage

Use cert_id

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
      }
    }
  }
}
```

Use cert and key

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