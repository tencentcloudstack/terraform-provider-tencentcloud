---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_user_type"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_user_type"
description: |-
  Use this data source to query detailed information of dlc describe_user_type
---

# tencentcloud_dlc_describe_user_type

Use this data source to query detailed information of dlc describe_user_type

## Example Usage

```hcl
data "tencentcloud_dlc_describe_user_type" "describe_user_type" {
  user_id = "127382378"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `user_id` - (Optional, String) User id (uin), if left blank, it defaults to the caller's sub-uin.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_type` - User type, only support: ADMIN: ddministrator/COMMON: ordinary user.


