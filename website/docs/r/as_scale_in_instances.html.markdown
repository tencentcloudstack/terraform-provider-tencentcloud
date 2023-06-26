---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scale_in_instances"
sidebar_current: "docs-tencentcloud-resource-as_scale_in_instances"
description: |-
  Provides a resource to create a as scale_in_instances
---

# tencentcloud_as_scale_in_instances

Provides a resource to create a as scale_in_instances

## Example Usage

```hcl
resource "tencentcloud_as_scale_in_instances" "scale_in_instances" {
  auto_scaling_group_id = "asg-519acdug"
  scale_in_number       = 1
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Scaling group ID.
* `scale_in_number` - (Required, Int, ForceNew) Number of instances to be reduced.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

as scale_in_instances can be imported using the id, e.g.

```
terraform import tencentcloud_as_scale_in_instances.scale_in_instances scale_in_instances_id
```

