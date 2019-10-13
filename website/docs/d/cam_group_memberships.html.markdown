---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_group_memberships"
sidebar_current: "docs-tencentcloud-datasource-cam_group_memberships"
description: |-
  Use this data source to query detailed information of CAM groups
---

# tencentcloud_cam_group_memberships

Use this data source to query detailed information of CAM groups

## Example Usage

```hcl
data "tencentcloud_cam_group_memberships" "foo" {
  group_id = "12515263"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) Id of CAM group to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `membership_list` - A list of CAM group membership. Each element contains the following attributes:
  * `group_id` - Id of CAM group.
  * `user_ids` - Id set of the cam group members.


