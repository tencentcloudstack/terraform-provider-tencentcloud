---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_key"
sidebar_current: "docs-tencentcloud-resource-kms_key"
description: |-
  Provide a resource to create a KMS key.
---

# tencentcloud_kms_key

Provide a resource to create a KMS key.

## Example Usage

### Create and enable a instance.

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                         = "tf-example-kms-key"
  description                   = "example of kms key"
  key_rotation_enabled          = false
  is_enabled                    = true
  pending_delete_window_in_days = 7

  tags = {
    createdBy = "Terraform"
  }
}
```

### Create kms instance with HSM

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                         = "tf-example-kms-key"
  description                   = "example of kms key"
  key_rotation_enabled          = false
  is_enabled                    = true
  pending_delete_window_in_days = 7
  hsm_cluster_id                = "cls-hsm-mwpd9cjm"

  tags = {
    createdBy = "Terraform"
  }
}
```

### Specify the Key Usage as an asymmetry method.

```hcl
resource "tencentcloud_kms_key" "example" {
  alias       = "tf-example-kms-key"
  description = "example of kms key"
  key_usage   = "ASYMMETRIC_DECRYPT_RSA_2048"
  is_enabled  = false
}
```

### Disable the kms key instance.

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = false

  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Required, String) Name of CMK. The name can only contain English letters, numbers, underscore and hyphen '-'. The first character must be a letter or number.
* `description` - (Optional, String) Description of CMK. The maximum is 1024 bytes.
* `hsm_cluster_id` - (Optional, String) The HSM cluster ID corresponding to KMS Advanced Edition (only valid for KMS Exclusive/Managed Edition service instances).
* `is_archived` - (Optional, Bool) Specify whether to archive key. Default value is `false`. This field is conflict with `is_enabled`, valid when key_state is `Enabled`, `Disabled`, `Archived`.
* `is_enabled` - (Optional, Bool) Specify whether to enable key. Default value is `false`. This field is conflict with `is_archived`, valid when key_state is `Enabled`, `Disabled`, `Archived`.
* `key_rotation_enabled` - (Optional, Bool) Specify whether to enable key rotation, valid when key_usage is `ENCRYPT_DECRYPT`. Default value is `false`.
* `key_usage` - (Optional, String, ForceNew) Usage of CMK. Available values include `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`. Default value is `ENCRYPT_DECRYPT`.
* `pending_delete_window_in_days` - (Optional, Int) Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.
* `tags` - (Optional, Map) Tags of CMK.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `key_state` - State of CMK.


## Import

KMS keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_key.example 287e8f40-7cbb-11eb-9a3a-5254004f7f94
```

