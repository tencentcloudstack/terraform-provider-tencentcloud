---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_remove_instances"
sidebar_current: "docs-tencentcloud-resource-as_remove_instances"
description: |-
  Provides a resource to create a as remove_instances
---

# tencentcloud_as_remove_instances

Provides a resource to create a as remove_instances

## Example Usage

```hcl
resource "tencentcloud_as_remove_instances" "remove_instances" {
  auto_scaling_group_id = ""
  instance_ids          = ""
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Launch configuration ID.
* `instance_ids` - (Required, Set: [`String`], ForceNew) List of cvm instances to remove.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

as remove_instances can be imported using the id, e.g.

```
terraform import tencentcloud_as_remove_instances.remove_instances remove_instances_id
```

