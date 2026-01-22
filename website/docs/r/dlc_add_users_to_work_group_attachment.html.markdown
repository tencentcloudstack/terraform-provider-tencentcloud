---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_add_users_to_work_group_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_add_users_to_work_group_attachment"
description: |-
  Provides a resource to create a DLC add users to work group attachment
---

# tencentcloud_dlc_add_users_to_work_group_attachment

Provides a resource to create a DLC add users to work group attachment

## Example Usage

```hcl
resource "tencentcloud_dlc_add_users_to_work_group_attachment" "example" {
  add_info {
    work_group_id = 70220
    user_ids      = ["100032717595", "100030773831"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `add_info` - (Required, List, ForceNew) Information about working groups and users to be operated.

The `add_info` object supports the following:

* `user_ids` - (Required, Set, ForceNew) User ID which matches the Uin on the CAM side.
* `work_group_id` - (Required, Int, ForceNew) Working group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC add users to work group attachment can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_add_users_to_work_group_attachment.example '70220#100032717595|100030773831'
```

