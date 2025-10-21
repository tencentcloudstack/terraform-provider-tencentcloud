---
subcategory: "Key Management Service(KMS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kms_list_keys"
sidebar_current: "docs-tencentcloud-datasource-kms_list_keys"
description: |-
  Use this data source to query detailed information of kms list_keys
---

# tencentcloud_kms_list_keys

Use this data source to query detailed information of kms list_keys

## Example Usage

```hcl
data "tencentcloud_kms_list_keys" "example" {
  role = 1
}
```

## Argument Reference

The following arguments are supported:

* `hsm_cluster_id` - (Optional, String) HSM cluster ID (only valid for KMS exclusive/managed service instances).
* `result_output_file` - (Optional, String) Used to save results.
* `role` - (Optional, Int) Filter based on the creator role. The default value is 0, which indicates the cmk created by the user himself, and 1, which indicates the cmk automatically created by authorizing other cloud products.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - A list of KMS keys.
  * `key_id` - ID of CMK.


