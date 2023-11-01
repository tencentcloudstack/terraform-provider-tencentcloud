---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_bind_work_groups_to_user_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_bind_work_groups_to_user_attachment"
description: |-
  Provides a resource to create a dlc bind_work_groups_to_user
---

# tencentcloud_dlc_bind_work_groups_to_user_attachment

Provides a resource to create a dlc bind_work_groups_to_user

## Example Usage

```hcl
resource "tencentcloud_dlc_bind_work_groups_to_user_attachment" "bind_work_groups_to_user" {
  add_info {
    user_id        = "100032772113"
    work_group_ids = [23184, 23181]
  }
}
```

## Argument Reference

The following arguments are supported:

* `add_info` - (Required, List, ForceNew) Bind user and workgroup information.

The `add_info` object supports the following:

* `user_id` - (Required, String) User id, matched with CAM side uin.
* `work_group_ids` - (Required, Set) Work group id set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc bind_work_groups_to_user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user bind_work_groups_to_user_id
```

