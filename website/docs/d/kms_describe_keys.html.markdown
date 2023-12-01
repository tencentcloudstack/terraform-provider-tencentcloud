---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_describe_keys"
sidebar_current: "docs-tencentcloud-datasource-kms_describe_keys"
description: |-
  Use this data source to query detailed information of kms key_lists
---

# tencentcloud_kms_describe_keys

Use this data source to query detailed information of kms key_lists

## Example Usage

```hcl
data "tencentcloud_kms_describe_keys" "example" {
  key_ids = [
    "9ffacc8b-6461-11ee-a54e-525400dd8a7d",
    "bffae4ed-6465-11ee-90b2-5254000ef00e"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `key_ids` - (Required, Set: [`String`]) Query the ID list of CMK, batch query supports up to 100 KeyIds at a time.
* `result_output_file` - (Optional, String) Used to save results.

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
  * `key_state` - State of CMK.
  * `key_usage` - Usage of CMK.
  * `next_rotate_time` - Next rotate time of CMK when key_rotation_enabled is true.
  * `origin` - Origin of CMK. `TENCENT_KMS` - CMK created by KMS, `EXTERNAL` - CMK imported by user.
  * `owner` - Creator of CMK.
  * `valid_to` - Valid when origin is `EXTERNAL`, it means the effective date of the key material.


