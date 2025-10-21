---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_strategy"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_strategy"
description: |-
  Provides a resource to create a tse cngw_strategy
---

# tencentcloud_tse_cngw_strategy

Provides a resource to create a tse cngw_strategy

~> **NOTE:** Please pay attention to the correctness of the cycle when modifying the `params` of `cron_config`, otherwise the modification will not be successful.

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_strategy" "cngw_strategy" {
  description   = "aaaaa"
  gateway_id    = "gateway-cf8c99c3"
  strategy_name = "test-cron"

  config {
    max_replicas = 2
    behavior {
      scale_down {
        select_policy                = "Max"
        stabilization_window_seconds = 301

        policies {
          period_seconds = 9
          type           = "Pods"
          value          = 1
        }
      }
      scale_up {
        select_policy                = "Max"
        stabilization_window_seconds = 31

        policies {
          period_seconds = 10
          type           = "Pods"
          value          = 1
        }
      }
    }

    metrics {
      resource_name = "cpu"
      target_value  = 1
      type          = "Resource"
    }
  }

  cron_config {
    params {
      crontab         = "0 00 00 * * *"
      period          = "* * *"
      start_at        = "00:00"
      target_replicas = 2
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) gateway ID.
* `strategy_name` - (Required, String) strategy name, up to 20 characters.
* `config` - (Optional, List) configuration of metric scaling.
* `cron_config` - (Optional, List) configuration of timed scaling.
* `description` - (Optional, String) description information, up to 120 characters.

The `behavior` object of `config` supports the following:

* `scale_down` - (Optional, List) configuration of down scale
Note: This field may return null, indicating that a valid value is not available.
* `scale_up` - (Optional, List) configuration of up scale
Note: This field may return null, indicating that a valid value is not available.

The `config` object supports the following:

* `behavior` - (Optional, List) behavior configuration of metric
Note: This field may return null, indicating that a valid value is not available.
* `create_time` - (Optional, String) create time
Note: This field may return null, indicating that a valid value is not available.
* `max_replicas` - (Optional, Int) max number of replica for metric scaling.
* `metrics` - (Optional, List) metric list.
* `modify_time` - (Optional, String) modify time
Note: This field may return null, indicating that a valid value is not available.
* `strategy_id` - (Optional, String) strategy ID
Note: This field may return null, indicating that a valid value is not available.

The `cron_config` object supports the following:

* `params` - (Optional, List) parameter list of timed scaling
Note: This field may return null, indicating that a valid value is not available.
* `strategy_id` - (Optional, String) strategy ID
Note: This field may return null, indicating that a valid value is not available.

The `metrics` object of `config` supports the following:

* `resource_name` - (Optional, String) metric name. Reference value:
- cpu
- memory
Note: This field may return null, indicating that a valid value is not available.
* `target_type` - (Optional, String) target type of metric, currently only supports `Utilization`
Note: This field may return null, indicating that a valid value is not available.
* `target_value` - (Optional, Int) target value of metric
Note: This field may return null, indicating that a valid value is not available.
* `type` - (Optional, String) metric type. Deafault value
- Resource.

The `params` object of `cron_config` supports the following:

* `crontab` - (Optional, String) cron expression of timed scaling, no input required
Note: This field may return null, indicating that a valid value is not available.
* `period` - (Optional, String) period of timed scaling
Note: This field may return null, indicating that a valid value is not available.
* `start_at` - (Optional, String) start time of timed scaling
Note: This field may return null, indicating that a valid value is not available.
* `target_replicas` - (Optional, Int) the number of target nodes for the timed scaling. Do not exceed the max number of replica for metric scaling
Note: This field may return null, indicating that a valid value is not available.

The `policies` object of `scale_down` supports the following:

* `period_seconds` - (Optional, Int) period of scale down
Note: This field may return null, indicating that a valid value is not available.
* `type` - (Optional, String) type, default value: Pods
Note: This field may return null, indicating that a valid value is not available.
* `value` - (Optional, Int) value
Note: This field may return null, indicating that a valid value is not available.

The `policies` object of `scale_up` supports the following:

* `period_seconds` - (Optional, Int) period of scale up
Note: This field may return null, indicating that a valid value is not available.
* `type` - (Optional, String) type, default value: Pods
Note: This field may return null, indicating that a valid value is not available.
* `value` - (Optional, Int) value
Note: This field may return null, indicating that a valid value is not available.

The `scale_down` object of `behavior` supports the following:

* `policies` - (Optional, List) policies of scale down
Note: This field may return null, indicating that a valid value is not available.
* `select_policy` - (Optional, String) type of policy, default value: max
Note: This field may return null, indicating that a valid value is not available.
* `stabilization_window_seconds` - (Optional, Int) stability window time, unit:second, default 300 when scale down
Note: This field may return null, indicating that a valid value is not available.

The `scale_up` object of `behavior` supports the following:

* `policies` - (Optional, List) policies of scale up
Note: This field may return null, indicating that a valid value is not available.
* `select_policy` - (Optional, String) type of policy, default value: max
Note: This field may return null, indicating that a valid value is not available.
* `stabilization_window_seconds` - (Optional, Int) stability window time, unit:second, default 0 when scale up
Note: This field may return null, indicating that a valid value is not available.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `strategy_id` - strategy ID
Note: This field may return null, indicating that a valid value is not available.


## Import

tse cngw_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_strategy.cngw_strategy gateway-cf8c99c3#strategy-a6744ff8
```

