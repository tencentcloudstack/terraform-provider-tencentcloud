---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_stop_instances"
sidebar_current: "docs-tencentcloud-resource-as_stop_instances"
description: |-
  Provides a resource to create a as stop_instances
---

# tencentcloud_as_stop_instances

Provides a resource to create a as stop_instances

## Example Usage

```hcl
resource "tencentcloud_as_stop_instances" "stop_instances" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.scaling_group.id
  instance_ids          = ["ins-xxxx"]
  stopped_mode          = "STOP_CHARGING"
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Launch configuration ID.
* `instance_ids` - (Required, Set: [`String`], ForceNew) List of cvm instances to stop.
* `stopped_mode` - (Optional, String, ForceNew) Billing method of a pay-as-you-go instance after shutdown. Available values: `KEEP_CHARGING`,`STOP_CHARGING`. Default `KEEP_CHARGING`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



