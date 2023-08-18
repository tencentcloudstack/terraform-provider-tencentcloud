---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_work_group"
sidebar_current: "docs-tencentcloud-resource-dlc_work_group"
description: |-
  Provides a resource to create a dlc work_group
---

# tencentcloud_dlc_work_group

Provides a resource to create a dlc work_group

## Example Usage

```hcl
resource "tencentcloud_dlc_work_group" "work_group" {
  work_group_name        = "tf-demo"
  work_group_description = "dlc workgroup test"
}
```

## Argument Reference

The following arguments are supported:

* `work_group_name` - (Required, String, ForceNew) Name of Work Group.
* `work_group_description` - (Optional, String) Description of Work Group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `user_ids` - A collection of user IDs that has been bound to the workgroup.


## Import

dlc work_group can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_work_group.work_group work_group_id
```

