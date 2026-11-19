---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_plan_for_zone"
sidebar_current: "docs-tencentcloud-action-teo_plan_for_zone"
description: |-
  Provides a resource to purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. This is a one-time operation action; the plan purchase is irreversible and cannot be cancelled on resource destroy.
---

# tencentcloud_teo_plan_for_zone

Provides a resource to purchase an EdgeOne (TEO) plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. This is a one-time operation action; the plan purchase is irreversible and cannot be cancelled on resource destroy.

## Example Usage

### Option 1: Invoke on demand via CLI

```hcl
sh
terraform apply - invoke = action.tencentcloud_teo_plan_for_zone.example
```



```hcl
action "tencentcloud_teo_plan_for_zone" "example" {
  config {
    zone_id   = "zone-27h0vbm5w1e"
    plan_type = "sta_global"
  }
}
```

### Option 2: Triggered by the zone resource lifecycle

```hcl
resource "tencentcloud_teo_zone" "example" {
  # ...
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.tencentcloud_teo_plan_for_zone.example]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `plan_type` - (Required, String) Plan type to purchase. Valid values: `sta`, `sta_with_bot`, `sta_cm`, `sta_cm_with_bot`, `sta_global`, `sta_global_with_bot`, `ent`, `ent_with_bot`, `ent_cm`, `ent_cm_with_bot`, `ent_global`, `ent_global_with_bot`.
* `zone_id` - (Required, String) Zone ID.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke: [CreatePlanForZone](https://www.tencentcloud.com/document/api/1145/50546)

