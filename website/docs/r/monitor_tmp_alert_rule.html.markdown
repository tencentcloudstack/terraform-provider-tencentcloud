---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_alert_rule"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_alert_rule"
description: |-
  Provides a resource to create a monitor tmpAlertRule
---

# tencentcloud_monitor_tmp_alert_rule

Provides a resource to create a monitor tmpAlertRule

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_alert_rule" "basic" {
  instance_id = "prom-c89b3b3u"
  rule_name   = "test-rule_name"
  receivers   = ["Consumer-6vkna7pevq"]
  expr        = "increase(mysql_global_status_slow_queries[1m]) > 0"
  duration    = "4m"
  rule_state  = 2
}
```

## Argument Reference

The following arguments are supported:

* `expr` - (Required, String) Rule expression.
* `instance_id` - (Required, String) Instance id.
* `receivers` - (Required, Set: [`String`]) Alarm notification template id list.
* `rule_name` - (Required, String) Rule name.
* `annotations` - (Optional, List) Rule alarm duration.
* `duration` - (Optional, String) Rule alarm duration.
* `labels` - (Optional, List) Rule alarm duration.
* `rule_state` - (Optional, Int) Rule state code.
* `type` - (Optional, String) Alarm Policy Template Classification.

The `annotations` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

The `labels` object supports the following:

* `key` - (Required, String) key.
* `value` - (Required, String) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor tmpAlertRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_alert_rule.tmpAlertRule instanceId#Rule_id
```

