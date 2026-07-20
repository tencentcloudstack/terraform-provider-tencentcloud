---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_bind_zone_to_plan"
sidebar_current: "docs-tencentcloud-resource-teo_bind_zone_to_plan"
description: |-
  Provides a resource to bind a TEO zone to an existing plan. After creating an EdgeOne (TEO) zone, you can use this resource to bind the unbound zone to an existing plan so that the zone takes effect.
---

# tencentcloud_teo_bind_zone_to_plan

Provides a resource to bind a TEO zone to an existing plan. After creating an EdgeOne (TEO) zone, you can use this resource to bind the unbound zone to an existing plan so that the zone takes effect.

## Example Usage

```hcl
resource "tencentcloud_teo_bind_zone_to_plan" "example" {
  zone_id = "zone-12345678"
  plan_id = "edgeone-12345678"
}
```

## Argument Reference

The following arguments are supported:

* `plan_id` - (Required, String, ForceNew) ID of the target plan to be bound.
* `zone_id` - (Required, String, ForceNew) ID of the unbound zone that needs to be bound to a plan.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



