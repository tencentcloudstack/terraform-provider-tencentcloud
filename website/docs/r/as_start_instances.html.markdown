---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_start_instances"
sidebar_current: "docs-tencentcloud-resource-as_start_instances"
description: |-
  Provides a resource to create a as start_instances
---

# tencentcloud_as_start_instances

Provides a resource to create a as start_instances

## Example Usage

```hcl
resource "tencentcloud_as_start_instances" "start_instances" {
  auto_scaling_group_id = ""
  instance_ids          = ""
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Launch configuration ID.
* `instance_ids` - (Required, Set: [`String`], ForceNew) List of cvm instances to start.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



