---
subcategory: "KMS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_external_key"
sidebar_current: "docs-tencentcloud-resource-kms_external_key"
description: |-
  Provide a resource to create a KMS external import key.
---

# tencentcloud_kms_external_key

Provide a resource to create a KMS external import key.

## Example Usage

```hcl
resource "tencentcloud_kms_external_key" "foo" {
  alias               = "test"
  description         = "describe key test message."
  wrapping_algorithm  = "RSAES_PKCS1_V1_5"
  key_material_base64 = "MTIzMTIzMTIzMTIzMTIzQQ=="
  valid_to            = 2147443200
}
```

## Argument Reference

The following arguments are supported:

* `alias` - (Required) Name of CMK.The name can only contain English letters, numbers, underscore and hyphen '-'.The first character must be a letter or number.
* `description` - (Optional) Description of CMK.The maximum is 1024 bytes.
* `key_id` - (Optional) ID of CMK.
* `key_material_base64` - (Optional) The base64-encoded key material encrypted with the public_key.For regions using the national secret version, the length of the imported key material is required to be 128 bits, and for regions using the FIPS version, the length of the imported key material is required to be 256 bits.
* `key_state` - (Optional) State of CMK.Available values include `Enabled`, `Disabled`, `PendingDelete`, `PendingImport`, `Archived`.
* `pending_delete_window_in_days` - (Optional) Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.
* `tags` - (Optional) Tags of CMK.
* `valid_to` - (Optional) this value means the effective timestamp of the key material, 0 means it does not expire.Need to be greater than the current time point, the maximum support is 2147443200.
* `wrapping_algorithm` - (Optional) The algorithm for encrypting key material.Available values include `RSAES_PKCS1_V1_5`, `RSAES_OAEP_SHA_1` and `RSAES_OAEP_SHA_256`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

KMS external keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_external_key.foo 287e8f40-7cbb-11eb-9a3a-5254004f7f94
```

