---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_groups"
sidebar_current: "docs-tencentcloud-datasource-cam_groups"
description: |-
  Use this data source to query detailed information of CAM groups
---

# tencentcloud_cam_groups

Use this data source to query detailed information of CAM groups

## Example Usage

```hcl
data "tencentcloud_cam_groups" "foo" {
  group_id = "12515263"
  name     = "cam-role-test"
  remark   = "test"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) Id of CAM group to be queried.
* `name` - (Optional) Name of the CAM group to be queried.
* `remark` - (Optional) Description of the cam group.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_list` - A list of CAM groups. Each element contains the following attributes:
  * `create_time` - Create time of the CAM group.
  * `name` - Name of CAM group.
  * `remark` - Description of CAM group.


