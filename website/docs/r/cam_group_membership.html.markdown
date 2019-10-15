---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_group_membership"
sidebar_current: "docs-tencentcloud-resource-cam_group_membership"
description: |-
  Provides a resource to create a CAM group membership.
---

# tencentcloud_cam_group_membership

Provides a resource to create a CAM group membership.

## Example Usage

```hcl
resource "tencentcloud_cam_group_membership" "foo" {
  group_id = "12515263"
  user_ids = ["cam-test", "cam-test2"]
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) Id of cam group.
* `user_ids` - (Required) Id set of the cam group members.


## Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.foo 12515263
```

