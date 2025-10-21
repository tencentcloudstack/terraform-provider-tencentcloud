---
subcategory: "Cloud Access Management(CAM)"
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
# query by group_id
data "tencentcloud_cam_groups" "foo" {
  group_id = tencentcloud_cam_group.foo.id
}

# query by name
data "tencentcloud_cam_groups" "bar" {
  name = "cam-group-test"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, String) ID of CAM group to be queried.
* `name` - (Optional, String) Name of the CAM group to be queried.
* `remark` - (Optional, String) Description of the cam group to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_list` - A list of CAM groups. Each element contains the following attributes:
  * `create_time` - Create time of the CAM group.
  * `group_id` - ID of the CAM group.
  * `name` - Name of CAM group.
  * `remark` - Description of CAM group.


