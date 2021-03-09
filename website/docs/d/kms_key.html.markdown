---
subcategory: "KMS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_key"
sidebar_current: "docs-tencentcloud-datasource-kms_key"
description: |-
  Use this data source to query detailed information of KMS key
---

# tencentcloud_kms_key

Use this data source to query detailed information of KMS key

## Example Usage

```hcl
data "tencentcloud_kms_key" "foo" {
  search_key_alias = "test"
  key_state        = "All"
  origin           = "TENCENT_KMS"
  key_usage        = "ALL"
}
```

## Argument Reference

The following arguments are supported:

* `key_state` - (Optional) State of CMK.Available values include `All`, `Enabled`, `Disabled`, `PendingDelete`, `PendingImport`, `Archived`.
* `key_usage` - (Optional) Usage of CMK.Available values include `ALL`, `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`.Default value is `ENCRYPT_DECRYPT`.
* `order_type` - (Optional) Order to sort the CMK create time.`0` - desc, `1` - asc.Default value is `0`.
* `origin` - (Optional) Origin of CMK.`TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user, `ALL` - All CMK.Default value is `ALL`.
* `result_output_file` - (Optional) Used to save results.
* `role` - (Optional) Role of the CMK creator.`0` - created by user, `1` - created by cloud product.Default value is `0`.
* `search_key_alias` - (Optional) Words used to match the results,and the words can be: key_id and alias.
* `tags` - (Optional) Tags to filter CMK.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_list` - A list of KMS keys.
  * `alias` - Name of CMK.
  * `create_time` - Create time of CMK.
  * `creator_uin` - Uin of CMK Creator.
  * `deletion_date` - Delete time of CMK.
  * `description` - Description of CMK.
  * `key_id` - ID of CMK.
  * `key_rotation_enabled` - Specify whether to enable key rotation.
  * `key_state` - State of CMK.Available values include `Enabled`, `Disabled`, `PendingDelete`, `PendingImport`, `Archived`.
  * `key_usage` - Usage of CMK.Available values include `ENCRYPT_DECRYPT`, `ASYMMETRIC_DECRYPT_RSA_2048`, `ASYMMETRIC_DECRYPT_SM2`, `ASYMMETRIC_SIGN_VERIFY_SM2`, `ASYMMETRIC_SIGN_VERIFY_RSA_2048`, `ASYMMETRIC_SIGN_VERIFY_ECC`.
  * `next_rotate_time` - Next rotate time of CMK when key_rotation_enabled is true.
  * `origin` - Origin of CMK.`TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user.
  * `owner` - Creator of CMK.
  * `valid_to` - Valid when Origin is EXTERNAL, it means the effective date of the key material.


