---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_work_group"
sidebar_current: "docs-tencentcloud-resource-dlc_work_group"
description: |-
  Provides a resource to create a DLC work group
---

# tencentcloud_dlc_work_group

Provides a resource to create a DLC work group

## Example Usage

```hcl
resource "tencentcloud_dlc_work_group" "example" {
  work_group_name        = "tf-example"
  work_group_description = "DLC workgroup demo"
}
```

## Argument Reference

The following arguments are supported:

* `work_group_name` - (Required, String, ForceNew) Working group name.
* `work_group_description` - (Optional, String) Working group description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `user_ids` - Collection of IDs of users to be bound to working groups.
* `work_group_id` - Working group ID.


## Import

DLC work group can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_work_group.example 135
```

