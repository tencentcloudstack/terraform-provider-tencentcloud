---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_public_key"
sidebar_current: "docs-tencentcloud-datasource-kms_public_key"
description: |-
  Use this data source to query detailed information of kms public_key
---

# tencentcloud_kms_public_key

Use this data source to query detailed information of kms public_key

## Example Usage

```hcl
data "tencentcloud_kms_public_key" "example" {
  key_id = tencentcloud_kms_key.example.id
}

resource "tencentcloud_kms_key" "example" {
  alias                         = "tf-example-kms-key"
  description                   = "example of kms key"
  key_usage                     = "ASYMMETRIC_DECRYPT_RSA_2048"
  is_enabled                    = true
  pending_delete_window_in_days = 7
}
```

## Argument Reference

The following arguments are supported:

* `key_id` - (Required, String) CMK unique identifier.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_key_pem` - Public key content in PEM format.
* `public_key` - Base64-encoded public key content.


