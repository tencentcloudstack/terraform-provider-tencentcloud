---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_user_roles"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_user_roles"
description: |-
  Use this data source to query detailed information of dlc describe_user_roles
---

# tencentcloud_dlc_describe_user_roles

Use this data source to query detailed information of dlc describe_user_roles

## Example Usage

```hcl
data "tencentcloud_dlc_describe_user_roles" "describe_user_roles" {
  fuzzy = "1"
}
```

## Argument Reference

The following arguments are supported:

* `fuzzy` - (Optional, String) List according to ARN blur.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) The return results are sorted according to this field.
* `sorting` - (Optional, String) Positive or inverted, such as DESC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_roles` - User role information.
  * `app_id` - User app ID.
  * `arn` - Role Permissions.
  * `cos_permission_list` - COS authorization path listNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cos_path` - COS pathNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `permissions` - Permissions [Read, WRITE]Note: This field may return NULL, indicating that the valid value cannot be obtained.
  * `creator` - Creator UinNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `desc` - Character description information.
  * `modify_time` - Recently modify the time stamp.
  * `permission_json` - CAM strategy jsonNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `role_id` - Character ID.
  * `role_name` - Role NameNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `uin` - User ID.


