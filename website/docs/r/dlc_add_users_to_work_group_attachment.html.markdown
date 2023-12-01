---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_add_users_to_work_group_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_add_users_to_work_group_attachment"
description: |-
  Provides a resource to create a dlc add_users_to_work_group_attachment
---

# tencentcloud_dlc_add_users_to_work_group_attachment

Provides a resource to create a dlc add_users_to_work_group_attachment

## Example Usage

```hcl
resource "tencentcloud_dlc_add_users_to_work_group_attachment" "add_users_to_work_group_attachment" {
  add_info {
    work_group_id = 23184
    user_ids      = [100032676511]
  }
}
}
```

## Argument Reference

The following arguments are supported:

* `add_info` - (Required, List, ForceNew) Work group and user information to operate on.

The `add_info` object supports the following:

* `user_ids` - (Required, Set) User id set, matched with CAM side uin.
* `work_group_id` - (Required, Int) Work group id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc add_users_to_work_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment add_users_to_work_group_attachment_id
```

