---
subcategory: "KMS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_key"
sidebar_current: "docs-tencentcloud-resource-kms_key"
description: |-
  Provide a resource to create a KMS key.
---

# tencentcloud_kms_key

Provide a resource to create a KMS key.

## Example Usage

```hcl
resource "tencentcloud_kms_key" "foo" {
  alias                = "test"
  description          = "describe key test message."
  key_rotation_enabled = true
  tags = {
    "test-tag" : "key-test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Required) Name of CMK.The name can only contain English letters, numbers, underscore and hyphen '-'.The first character must be a letter or number.
* `description` - (Optional) Description of CMK.The maximum is 1024 bytes.
* `key_id` - (Optional) ID of CMK.
* `key_rotation_enabled` - (Optional) Specify whether to enable key rotation.
* `key_state` - (Optional) State of CMK.Available values include `Enabled`, `Disabled`, `PendingDelete`, `Archived`.
* `key_usage` - (Optional) Usage of CMK.Available values include `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`.Default value is `ENCRYPT_DECRYPT`.
* `pending_delete_window_in_days` - (Optional) Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.
* `tags` - (Optional) Tags of CMK.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

KMS keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_key.foo 287e8f40-7cbb-11eb-9a3a-5254004f7f94
```

