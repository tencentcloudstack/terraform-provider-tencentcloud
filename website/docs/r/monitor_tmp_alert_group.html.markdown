---
subcategory: "Managed Service for Prometheus(TMP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_alert_group"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_alert_group"
description: |-
  Provides a resource to create a monitor tmp_alert_group
---

# tencentcloud_monitor_tmp_alert_group

Provides a resource to create a monitor tmp_alert_group

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_alert_group" "tmp_alert_group" {
  amp_receivers = [
    "notice-om017kc2",
  ]
  group_name      = "tf-test"
  instance_id     = "prom-ip429jis"
  repeat_interval = "5m"

  custom_receiver {
    type = "amp"
  }

  rules {
    duration  = "1m"
    expr      = "up{job=\"prometheus-agent\"} != 1"
    rule_name = "Agent health check"
    state     = 2

    annotations = {
      "summary"     = "Agent health check"
      "description" = "Agent {{$labels.instance}} is deactivated, please pay attention!"
    }

    labels = {
      "severity" = "critical"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `amp_receivers` - (Optional, Set: [`String`]) Tencent cloud notification template id list.
* `custom_receiver` - (Optional, List) User custom notification template, such as webhook, alertmanager.
* `group_name` - (Optional, String) Unique alert group name.
* `instance_id` - (Optional, String) Instance id.
* `repeat_interval` - (Optional, String) Alert message send interval, default 1 hour.
* `rules` - (Optional, List) A list of alert rules.

The `allowed_time_ranges` object of `custom_receiver` supports the following:

* `end` - (Optional, String) Time range end, seconds since 0 o'clock.
* `start` - (Optional, String) Time range start, seconds since 0 o'clock.

The `custom_receiver` object supports the following:

* `allowed_time_ranges` - (Optional, List) Time ranges which allow alert message send.
* `cluster_id` - (Optional, String) Only effect when alertmanager in user cluster, this cluster id.
* `cluster_type` - (Optional, String) Only effect when alertmanager in user cluster, this cluster type (tke|eks|tdcc).
* `type` - (Optional, String) Custom receiver type, webhook|alertmanager.
* `url` - (Optional, String) Custom receiver address, can be accessed by process in prometheus instance subnet.

The `rules` object supports the following:

* `annotations` - (Optional, Map) Annotation of alert rule. `summary`, `description` is special annotation in prometheus, mapping `Alarm Object`, `Alarm Information` in alarm message.
* `duration` - (Optional, String) Rule alarm duration.
* `expr` - (Optional, String) Prometheus alert expression.
* `labels` - (Optional, Map) Labels of alert rule.
* `rule_name` - (Optional, String) Alert rule name.
* `state` - (Optional, Int) Rule state. `2`-enable, `3`-disable, default `2`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `group_id` - Alarm group id.


## Import

monitor tmp_alert_group can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_alert_group.tmp_alert_group instance_id#group_id
```

