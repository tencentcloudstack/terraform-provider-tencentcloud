---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_users"
sidebar_current: "docs-tencentcloud-datasource-identity_center_users"
description: |-
  Use this data source to query detailed information of identity center users
---

# tencentcloud_identity_center_users

Use this data source to query detailed information of identity center users

## Example Usage

```hcl
data "tencentcloud_identity_center_users" "identity_center_users" {
  zone_id = "z-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Space ID.
* `filter_groups` - (Optional, Set: [`String`]) Filtered user group. IsSelected=1 will be returned for the sub-user associated with this user group.
* `filter` - (Optional, String) Filter criterion, which currently only supports username, email address, userId, and description.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_field` - (Optional, String) Sorting field, which currently only supports CreateTime. The default is the CreateTime field.
* `sort_type` - (Optional, String) Sorting type. Desc: descending order; Asc: ascending order. It should be set along with SortField.
* `user_status` - (Optional, String) User status: Enabled, Disabled.
* `user_type` - (Optional, String) User type. Manual: manually created; Synchronized: externally imported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `users` - User list.


