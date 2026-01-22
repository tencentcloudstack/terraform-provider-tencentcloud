---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_bind_work_groups_to_user_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_bind_work_groups_to_user_attachment"
description: |-
  Provides a resource to create a DLC bind work groups to user
---

# tencentcloud_dlc_bind_work_groups_to_user_attachment

Provides a resource to create a DLC bind work groups to user

## Example Usage

```hcl
resource "tencentcloud_dlc_bind_work_groups_to_user_attachment" "example" {
  add_info {
    user_id        = "100032772113"
    work_group_ids = [23184, 23181]
  }
}
```

## Argument Reference

The following arguments are supported:

* `add_info` - (Required, List, ForceNew) Information about bound working groups and users.

The `add_info` object supports the following:

* `user_id` - (Required, String, ForceNew) User ID, which matches Uin on the CAM side.
* `work_group_ids` - (Required, Set, ForceNew) Collections of IDs of working groups.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC bind work groups to user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_bind_work_groups_to_user_attachment.example 100032772113
```

