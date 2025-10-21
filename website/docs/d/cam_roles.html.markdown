---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_roles"
sidebar_current: "docs-tencentcloud-datasource-cam_roles"
description: |-
  Use this data source to query detailed information of CAM roles
---

# tencentcloud_cam_roles

Use this data source to query detailed information of CAM roles

## Example Usage

```hcl
# query by role_id
data "tencentcloud_cam_roles" "foo" {
  role_id = tencentcloud_cam_role.foo.id
}

# query by name
data "tencentcloud_cam_roles" "bar" {
  name = "cam-role-test"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, String) The description of the CAM role to be queried.
* `name` - (Optional, String) Name of the CAM policy to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `role_id` - (Optional, String) ID of the CAM role to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_list` - A list of CAM roles. Each element contains the following attributes:
  * `console_login` - Indicate whether the CAM role can be login or not.
  * `create_time` - The create time of the CAM role.
  * `description` - Description of CAM role.
  * `document` - Policy document of CAM role.
  * `name` - Name of CAM role.
  * `role_id` - Id of CAM role.
  * `update_time` - The last update time of the CAM role.


