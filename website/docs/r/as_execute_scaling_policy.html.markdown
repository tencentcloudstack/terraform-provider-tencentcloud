---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_execute_scaling_policy"
sidebar_current: "docs-tencentcloud-resource-as_execute_scaling_policy"
description: |-
  Provides a resource to create a as execute_scaling_policy
---

# tencentcloud_as_execute_scaling_policy

Provides a resource to create a as execute_scaling_policy

## Example Usage

```hcl
resource "tencentcloud_as_execute_scaling_policy" "execute_scaling_policy" {
  auto_scaling_policy_id = "asp-519acdug"
  honor_cooldown         = false
  trigger_source         = "API"
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_policy_id` - (Required, String, ForceNew) Auto-scaling policy ID. This parameter is not available to a target tracking policy.
* `honor_cooldown` - (Optional, Bool, ForceNew) Whether to check if the auto scaling group is in the cooldown period. Default value: false.
* `trigger_source` - (Optional, String, ForceNew) Source that triggers the scaling policy. Valid values: API and CLOUD_MONITOR. Default value: API. The value CLOUD_MONITOR is specific to the Cloud Monitor service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

as execute_scaling_policy can be imported using the id, e.g.

```
terraform import tencentcloud_as_execute_scaling_policy.execute_scaling_policy execute_scaling_policy_id
```

