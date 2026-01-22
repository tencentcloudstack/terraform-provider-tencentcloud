---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_user_roles"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_user_roles"
description: |-
  Use this data source to query detailed information of DLC describe user roles
---

# tencentcloud_dlc_describe_user_roles

Use this data source to query detailed information of DLC describe user roles

## Example Usage

```hcl
data "tencentcloud_dlc_describe_user_roles" "example" {
  fuzzy   = "1"
  sort_by = "modify-time"
  sorting = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `fuzzy` - (Optional, String) Fuzzy enumeration by arn.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) The field for sorting the returned results.
* `sorting` - (Optional, String) The sorting order, descending or ascending, such as `desc`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_roles` - The user roles.
  * `app_id` - The user's app ID.
  * `arn` - The role permission.
  * `cos_permission_list` - COS authorization path list.
    * `cos_path` - COS path.
    * `permissions` - Permissions [read, write].
  * `creator` - Creator Uin.
  * `desc` - The role description.
  * `modify_time` - The last modified timestamp.
  * `permission_json` - CAM strategy json.
  * `role_id` - The role ID.
  * `role_name` - The role name.
  * `uin` - The user ID.


