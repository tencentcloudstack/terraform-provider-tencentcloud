Provides a policy group resource for monitor.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_monitor_alarm_policy.

Example Usage

```hcl
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 1
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
  conditions {
    metric_id           = 30
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 2
    calc_value          = 30
    calc_period         = 300
    continue_period     = 2
  }
  event_conditions {
    event_id            = 39
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
  event_conditions {
    event_id            = 40
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
}
```
Import

Policy group instance can be imported, e.g.

```
$ terraform import tencentcloud_monitor_policy_group.group group-id
```
