---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_upload_certificate_instance"
sidebar_current: "docs-tencentcloud-action-ssl_upload_certificate_instance"
description: |-
  Provides an action to upload new certificate public/private key content and update the associated cloud resources (such as clb) without changing the certificate ID via the UploadUpdateCertificateInstance API. The action waits for the asynchronous deploy task to complete. This is a one-time operation; no cloud-side state is persisted.
---

# tencentcloud_ssl_upload_certificate_instance

Provides an action to upload new certificate public/private key content and update the associated cloud resources (such as clb) without changing the certificate ID via the UploadUpdateCertificateInstance API. The action waits for the asynchronous deploy task to complete. This is a one-time operation; no cloud-side state is persisted.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** Currently, this product only supports updating cloud resources of the `clb` type.

## Example Usage

```hcl
action "tencentcloud_ssl_upload_certificate_instance" "example" {
  config {
    old_certificate_id      = "f8xnynbd"
    resource_types          = ["clb"]
    certificate_public_key  = "-----BEGIN CERTIFICATE-----\nxxx\n-----END CERTIFICATE-----"
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

## Argument Reference

The following arguments are supported:

* `certificate_private_key` - (Required, String) Certificate private key. If you upload the certificate public key, CertificateId does not need to be passed.
* `certificate_public_key` - (Required, String) Certificate public key. If you upload the certificate public key, CertificateId does not need to be passed.
* `old_certificate_id` - (Required, String) Update the original certificate ID. The cloud resources bound to this certificate ID are queried, and then these cloud resources are updated with the new certificate.
* `resource_types` - (Required, Set) The resource type that needs to be deployed. Currently, only `clb` is supported.
* `certificate_id` - (Optional, String) Certificate ID.
* `resource_types_regions` - (Optional, List of Object) List of regions where cloud resources need to be deployed.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

The `resource_types_regions` object supports the following:

* `regions` - (Optional, Set) Region list.
* `resource_type` - (Optional, String) Cloud resource type.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [DescribeHostUploadUpdateRecord](https://www.tencentcloud.com/document/api/1007/70505)
  - [UploadUpdateCertificateInstance](https://www.tencentcloud.com/document/api/1007/70503)

