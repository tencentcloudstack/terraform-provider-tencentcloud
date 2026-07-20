---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_plan_for_zone"
sidebar_current: "docs-tencentcloud-resource-teo_plan_for_zone"
description: |-
  Provides a resource to purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. This is a one-time operation resource; the plan purchase is irreversible and cannot be cancelled on resource destroy.
---

# tencentcloud_teo_plan_for_zone

Provides a resource to purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. This is a one-time operation resource; the plan purchase is irreversible and cannot be cancelled on resource destroy.

## Example Usage

```hcl
resource "tencentcloud_teo_plan_for_zone" "example" {
  zone_id   = "zone-27h0vbm5w1e"
  plan_type = "sta_global"
}
```

### Purchase a standard plan with bot management

```hcl
resource "tencentcloud_teo_plan_for_zone" "example_bot" {
  zone_id   = "zone-27h0vbm5w1e"
  plan_type = "sta_global_with_bot"
}
```

## Argument Reference

The following arguments are supported:

* `plan_type` - (Required, String, ForceNew) Plan type to purchase. Valid values: `sta`, `sta_with_bot`, `sta_cm`, `sta_cm_with_bot`, `sta_global`, `sta_global_with_bot`, `ent`, `ent_with_bot`, `ent_cm`, `ent_cm_with_bot`, `ent_global`, `ent_global_with_bot`.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `deal_names` - List of purchased order/deal names returned by the API.
* `resource_names` - List of purchased resource names returned by the API.


