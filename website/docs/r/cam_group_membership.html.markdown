---
subcategory: "Cloud Access Management(CAM)"
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
  group_id   = tencentcloud_cam_group.foo.id
  user_names = [tencentcloud_cam_user.foo.name, tencentcloud_cam_user.bar.name]
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) ID of CAM group.
* `user_ids` - (Optional, **Deprecated**) It has been deprecated from version 1.59.5. Use `user_names` instead. ID set of the CAM group members.
* `user_names` - (Optional) User name set as ID of the CAM group members.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.foo 12515263
```

