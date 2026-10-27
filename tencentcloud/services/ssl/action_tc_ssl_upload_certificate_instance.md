Provides an action to upload new certificate public/private key content and update the associated cloud resources (such as clb) without changing the certificate ID via the UploadUpdateCertificateInstance API. The action waits for the asynchronous deploy task to complete. This is a one-time operation; no cloud-side state is persisted.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** Currently, this product only supports updating cloud resources of the `clb` type.

Example Usage

```hcl
action "tencentcloud_ssl_upload_certificate_instance" "example" {
  config {
    old_certificate_id       = "f8xnynbd"
    resource_types           = ["clb"]
    certificate_public_key   = "-----BEGIN CERTIFICATE-----\nxxx\n-----END CERTIFICATE-----"
    certificate_private_key = "-----BEGIN RSA PRIVATE KEY-----\nxxx\n-----END RSA PRIVATE KEY-----"

    resource_types_regions {
      resource_type = "clb"
      regions       = ["ap-guangzhou"]
    }

    timeouts {
      invoke = "30m"
    }
  }
}
```
