---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_policy"
sidebar_current: "docs-tencentcloud-resource-monitor_alarm_policy"
description: |-
  Provides a alarm policy resource for monitor.
---

# tencentcloud_monitor_alarm_policy

Provides a alarm policy resource for monitor.

## Example Usage

```hcl
resource "tencentcloud_monitor_alarm_policy" "group" {
  policy_name  = "hello"
  monitor_type = "MT_QCE"
  enable       = 1
  project_id   = 1244035
  namespace    = "cvm_device"
  conditions {
    is_union_rule = 1
    rules {
      metric_name      = "CpuUsage"
      period           = 60
      operator         = "ge"
      value            = "89.9"
      continue_period  = 1
      notice_frequency = 3600
      is_power_notice  = 0
    }
  }
  event_conditions {
    metric_name = "ping_unreachable"
  }
  event_conditions {
    metric_name = "guest_reboot"
  }
  notice_ids = ["notice-l9ziyxw6"]

  trigger_tasks {
    type        = "AS"
    task_config = "{\"Region\":\"ap-guangzhou\",\"Group\":\"asg-0z312312x\",\"Policy\":\"asp-ganig28\"}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `monitor_type` - (Required) The type of monitor.
* `namespace` - (Required) The type of alarm.
* `policy_name` - (Required) The name of policy.
* `conditions` - (Optional) A list of metric trigger condition.
* `conditon_template_id` - (Optional, ForceNew) ID of trigger condition template.
* `enable` - (Optional) Whether to enable, default is `1`.
* `event_conditions` - (Optional) A list of event trigger condition.
* `notice_ids` - (Optional) List of notification rule IDs.
* `project_id` - (Optional) Project ID. For products with different projects, a value other than -1 must be passed in.
* `remark` - (Optional) The remark of policy group.
* `trigger_tasks` - (Optional) Triggered task list.

The `conditions` object supports the following:

* `is_union_rule` - (Optional) The and or relation of indicator alarm rule.
* `rules` - (Optional) A list of metric trigger condition.

The `event_conditions` object supports the following:

* `continue_period` - (Optional) Number of periods.
* `description` - (Optional, ForceNew) Metric display name, which is used in the output parameter.
* `is_power_notice` - (Optional) Whether the alarm frequency increases exponentially.
* `metric_name` - (Optional) Metric name or event name.
* `notice_frequency` - (Optional) Alarm interval in seconds.
* `operator` - (Optional) Operator.
* `period` - (Optional) Statistical period in seconds.
* `rule_type` - (Optional, ForceNew) Trigger condition type.
* `unit` - (Optional, ForceNew) Unit, which is used in the output parameter.
* `value` - (Optional) Threshold.

The `rules` object supports the following:

* `continue_period` - (Optional) Number of periods.
* `description` - (Optional, ForceNew) Metric display name, which is used in the output parameter.
* `is_power_notice` - (Optional) Whether the alarm frequency increases exponentially.
* `metric_name` - (Optional) Metric name or event name.
* `notice_frequency` - (Optional) Alarm interval in seconds.
* `operator` - (Optional) Operator.
* `period` - (Optional) Statistical period in seconds.
* `rule_type` - (Optional, ForceNew) Trigger condition type.
* `unit` - (Optional, ForceNew) Unit, which is used in the output parameter.
* `value` - (Optional) Threshold.

The `trigger_tasks` object supports the following:

* `task_config` - (Required) Configuration information in JSON format.
* `type` - (Required) Triggered task type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The alarm policy create time.
* `update_time` - The alarm policy update time.


## Import

Alarm policy instance can be imported, e.g.

```
$ terraform import tencentcloud_monitor_alarm_policy.policy policy-id
```

