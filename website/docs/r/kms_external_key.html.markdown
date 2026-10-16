---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_external_key"
sidebar_current: "docs-tencentcloud-resource-kms_external_key"
description: |-
  Provide a resource to create a KMS external key.
---

# tencentcloud_kms_external_key

Provide a resource to create a KMS external key.

## Example Usage

### Create a basic instance.

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias       = "tf-example-kms-externalkey"
  description = "example of kms external key"

  tags = {
    createdBy = "terraform"
  }
}
```

### Create kms instance with HSM

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias          = "tf-example-kms-externalkey"
  description    = "example of kms external key"
  hsm_cluster_id = "cls-hsm-mwpd9cjm"

  tags = {
    createdBy = "terraform"
  }
}
```

### Specify the encryption algorithm and public key.

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias               = "tf-example-kms-externalkey"
  description         = "example of kms external key"
  wrapping_algorithm  = "RSAES_PKCS1_V1_5"
  key_material_base64 = "your_public_key_base64_encoded"
  is_enabled          = true

  tags = {
    createdBy = "terraform"
  }
}
```

### Disable the external kms key.

```hcl
resource "tencentcloud_kms_external_key" "example" {
  alias               = "tf-example-kms-externalkey"
  description         = "example of kms external key"
  wrapping_algorithm  = "RSAES_PKCS1_V1_5"
  key_material_base64 = "your_public_key_base64_encoded"
  is_enabled          = false

  tags = {
    createdBy = "terraform"
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
* `key_material_base64` - (Optional, String) The base64-encoded key material encrypted with the public_key. For regions using the national secret version, the length of the imported key material is required to be 128 bits, and for regions using the FIPS version, the length of the imported key material is required to be 256 bits.
* `pending_delete_window_in_days` - (Optional, Int) Duration in days after which the key is deleted after destruction of the resource, must be between 7 and 30 days. Defaults to 7 days.
* `tags` - (Optional, Map) Tags of CMK.
* `valid_to` - (Optional, Int) This value means the effective timestamp of the key material, 0 means it does not expire. Need to be greater than the current timestamp, the maximum support is 2147443200.
* `wrapping_algorithm` - (Optional, String) The algorithm for encrypting key material. Available values include `RSAES_PKCS1_V1_5`, `RSAES_OAEP_SHA_1` and `RSAES_OAEP_SHA_256`. Default value is `RSAES_PKCS1_V1_5`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `key_state` - State of CMK.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ArchiveKey](https://www.tencentcloud.com/search/?wd=ArchiveKey&type=document&lang=en)
* Create: [CreateKey](https://www.tencentcloud.com/search/?wd=CreateKey&type=document&lang=en)
* Create: [DescribeKey](https://www.tencentcloud.com/search/?wd=DescribeKey&type=document&lang=en)
* Create: [DisableKey](https://www.tencentcloud.com/search/?wd=DisableKey&type=document&lang=en)
* Create: [GetParametersForImport](https://www.tencentcloud.com/search/?wd=GetParametersForImport&type=document&lang=en)
* Create: [ImportKeyMaterial](https://www.tencentcloud.com/search/?wd=ImportKeyMaterial&type=document&lang=en)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read: [DescribeKey](https://www.tencentcloud.com/search/?wd=DescribeKey&type=document&lang=en)
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Update: [DescribeKey](https://www.tencentcloud.com/search/?wd=DescribeKey&type=document&lang=en)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Update: [UpdateAlias](https://www.tencentcloud.com/search/?wd=UpdateAlias&type=document&lang=en)
* Update: [UpdateKeyDescription](https://www.tencentcloud.com/search/?wd=UpdateKeyDescription&type=document&lang=en)
* Delete: [DescribeKey](https://www.tencentcloud.com/search/?wd=DescribeKey&type=document&lang=en)
* Delete: [DisableKey](https://www.tencentcloud.com/search/?wd=DisableKey&type=document&lang=en)
* Delete: [ScheduleKeyDeletion](https://www.tencentcloud.com/search/?wd=ScheduleKeyDeletion&type=document&lang=en)

## Import

KMS external keys can be imported using the id, e.g.

```
$ terraform import tencentcloud_kms_external_key.example 25068921-2101-11f0-bf1f-5254000328e1
```

