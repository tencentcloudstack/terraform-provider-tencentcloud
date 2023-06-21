---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_group_status"
sidebar_current: "docs-tencentcloud-resource-as_scaling_group_status"
description: |-
  Provides a resource to set as scaling_group status
---

# tencentcloud_as_scaling_group_status

Provides a resource to set as scaling_group status

## Example Usage

```hcl
resource "tencentcloud_as_scaling_group_status" "scaling_group_status" {
  auto_scaling_group_id = "asg-519acdug"
  enable                = false
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Scaling group ID.
* `enable` - (Required, Bool) If enable auto scaling group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

as scaling_group_status can be imported using the id, e.g.

```
terraform import tencentcloud_as_scaling_group_status.scaling_group_status scaling_group_id
```

