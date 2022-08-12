---
subcategory: "TEM"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_scale_rule"
sidebar_current: "docs-tencentcloud-resource-tem_scale_rule"
description: |-
  Provides a resource to create a tem scaleRule
---

# tencentcloud_tem_scale_rule

Provides a resource to create a tem scaleRule

## Example Usage

```hcl
resource "tencentcloud_tem_scale_rule" "scaleRule" {
  environment_id = "en-853mggjm"
  application_id = "app-3j29aa2p"
  autoscaler {
    autoscaler_name = "test3123"
    description     = "test"
    enabled         = true
    min_replicas    = 1
    max_replicas    = 3
    cron_horizontal_autoscaler {
      name     = "test"
      period   = "* * *"
      priority = 1
      enabled  = true
      schedules {
        start_at        = "03:00"
        target_replicas = 1
      }
    }
    cron_horizontal_autoscaler {
      name     = "test123123"
      period   = "* * *"
      priority = 0
      enabled  = true
      schedules {
        start_at        = "04:13"
        target_replicas = 1
      }
    }
    horizontal_autoscaler {
      metrics      = "CPU"
      enabled      = true
      max_replicas = 3
      min_replicas = 1
      threshold    = 60
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String, ForceNew) application ID.
* `autoscaler` - (Required, List) .
* `environment_id` - (Required, String, ForceNew) environment ID.

The `autoscaler` object supports the following:

* `autoscaler_name` - (Required, String) name.
* `enabled` - (Required, Bool) enable AutoScaler.
* `max_replicas` - (Required, Int) maximal replica number.
* `min_replicas` - (Required, Int) minimal replica number.
* `cron_horizontal_autoscaler` - (Optional, List) scaler based on cron configuration.
* `description` - (Optional, String) description.
* `horizontal_autoscaler` - (Optional, List) scaler based on metrics.

The `cron_horizontal_autoscaler` object supports the following:

* `enabled` - (Required, Bool) enable scaler.
* `name` - (Required, String) name.
* `period` - (Required, String) period.
* `priority` - (Required, Int) priority.
* `schedules` - (Required, List) schedule payload.

The `horizontal_autoscaler` object supports the following:

* `enabled` - (Required, Bool) enable scaler.
* `max_replicas` - (Required, Int) maximal replica number.
* `metrics` - (Required, String) metric name.
* `min_replicas` - (Required, Int) minimal replica number.
* `threshold` - (Required, Int) metric threshold.

The `schedules` object supports the following:

* `start_at` - (Required, String) start time.
* `target_replicas` - (Required, Int) target replica number.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tem scaleRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_scale_rule.scaleRule environmentId#applicationId#scaleRuleId
```

