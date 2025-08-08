---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_user_type"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_user_type"
description: |-
  Use this data source to query detailed information of DLC describe user type
---

# tencentcloud_dlc_describe_user_type

Use this data source to query detailed information of DLC describe user type

## Example Usage

```hcl
data "tencentcloud_dlc_describe_user_type" "example" {
  user_id = "100021240183"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `user_id` - (Optional, String) User ID (UIN). If it is not specified, it will be the sub-UIN of the caller by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_type` - Types of users. ADMIN: administrators; COMMON: general users.


