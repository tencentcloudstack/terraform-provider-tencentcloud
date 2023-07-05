---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_policy"
sidebar_current: "docs-tencentcloud-resource-monitor_alarm_policy"
description: |-
  Provides a alarm policy resource for monitor.
---

# tencentcloud_monitor_alarm_policy

Provides a alarm policy resource for monitor.

## Example Usage

### cvm_device alarm policy

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

### k8s_cluster alarm policy

```hcl
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "k8s_cluster"
  notice_ids = [
    "notice-l9ziyxw6",
  ]
  policy_name = "TkeClusterNew"
  project_id  = 1244035

  conditions {
    is_union_rule = 0

    rules {
      continue_period  = 3
      description      = "Allocatable Pods"
      is_power_notice  = 0
      metric_name      = "K8sClusterAllocatablePodsTotal"
      notice_frequency = 3600
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "Count"
      value            = "10"

      filter {
        dimensions = jsonencode(
          [
            [
              {
                Key      = "region"
                Operator = "eq"
                Value = [
                  "ap-guangzhou",
                ]
              },
              {
                Key      = "tke_cluster_instance_id"
                Operator = "in"
                Value = [
                  "cls-czhtobea",
                ]
              },
            ],
          ]
        )
        type = "DIMENSION"
      }
    }
    rules {
      continue_period  = 3
      description      = "Total CPU Cores"
      is_power_notice  = 0
      metric_name      = "K8sClusterCpuCoreTotal"
      notice_frequency = 3600
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "Core"
      value            = "2"

      filter {
        dimensions = jsonencode(
          [
            [
              {
                Key      = "region"
                Operator = "eq"
                Value = [
                  "ap-guangzhou",
                ]
              },
              {
                Key      = "tke_cluster_instance_id"
                Operator = "in"
                Value = [
                  "cls-czhtobea",
                ]
              },
            ],
          ]
        )
        type = "DIMENSION"
      }
    }
  }
}
```

### cvm_device alarm policy binding cvm by tag

```hcl
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "cvm_device"
  notice_ids = [
    "notice-l9ziyxw6",
  ]
  policy_name = "policy"
  project_id  = 0

  conditions {
    is_union_rule = 0

    rules {
      continue_period  = 5
      description      = "CPUUtilization"
      is_power_notice  = 0
      metric_name      = "CpuUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "PublicBandwidthUtilization"
      is_power_notice  = 0
      metric_name      = "Outratio"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "MemoryUtilization"
      is_power_notice  = 0
      metric_name      = "MemUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "DiskUtilization"
      is_power_notice  = 0
      metric_name      = "CvmDiskUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
  }

  event_conditions {
    continue_period  = 0
    description      = "DiskReadonly"
    is_power_notice  = 0
    metric_name      = "disk_readonly"
    notice_frequency = 0
    period           = 0
  }

  policy_tag {
    key   = "test-tag"
    value = "unit-test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `monitor_type` - (Required, String, ForceNew) The type of monitor.
* `namespace` - (Required, String, ForceNew) The type of alarm.
* `policy_name` - (Required, String) The name of policy.
* `conditions` - (Optional, List) A list of metric trigger condition.
* `conditon_template_id` - (Optional, Int, ForceNew) ID of trigger condition template.
* `enable` - (Optional, Int) Whether to enable, default is `1`.
* `event_conditions` - (Optional, List) A list of event trigger condition.
* `notice_ids` - (Optional, List: [`String`]) List of notification rule IDs.
* `policy_tag` - (Optional, List, ForceNew) Policy tag to bind object.
* `project_id` - (Optional, Int, ForceNew) Project ID. For products with different projects, a value other than -1 must be passed in.
* `remark` - (Optional, String) The remark of policy group.
* `trigger_tasks` - (Optional, List) Triggered task list.

The `conditions` object supports the following:

* `is_union_rule` - (Optional, Int) The and or relation of indicator alarm rule.
* `rules` - (Optional, List) A list of metric trigger condition.

The `event_conditions` object supports the following:

* `continue_period` - (Optional, Int) Number of periods.
* `description` - (Optional, String) Metric display name, which is used in the output parameter.
* `filter` - (Optional, List) Filter condition for one single trigger rule. Must set it when create tke-xxx rules.
* `is_power_notice` - (Optional, Int) Whether the alarm frequency increases exponentially.
* `metric_name` - (Optional, String) Metric name or event name.
* `notice_frequency` - (Optional, Int) Alarm interval in seconds.
* `operator` - (Optional, String) Operator.
* `period` - (Optional, Int) Statistical period in seconds.
* `rule_type` - (Optional, String) Trigger condition type.
* `unit` - (Optional, String) Unit, which is used in the output parameter.
* `value` - (Optional, String) Threshold.

The `filter` object supports the following:

* `dimensions` - (Optional, String) JSON string generated by serializing the AlarmPolicyDimension two-dimensional array.
* `type` - (Optional, String) Filter condition type. Valid values: DIMENSION (uses dimensions for filtering).

The `policy_tag` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

The `rules` object supports the following:

* `continue_period` - (Optional, Int) Number of periods.
* `description` - (Optional, String) Metric display name, which is used in the output parameter.
* `filter` - (Optional, List) Filter condition for one single trigger rule. Must set it when create tke-xxx rules.
* `is_power_notice` - (Optional, Int) Whether the alarm frequency increases exponentially.
* `metric_name` - (Optional, String) Metric name or event name.
* `notice_frequency` - (Optional, Int) Alarm interval in seconds.
* `operator` - (Optional, String) Operator.
* `period` - (Optional, Int) Statistical period in seconds.
* `rule_type` - (Optional, String) Trigger condition type.
* `unit` - (Optional, String) Unit, which is used in the output parameter.
* `value` - (Optional, String) Threshold.

The `trigger_tasks` object supports the following:

* `task_config` - (Required, String) Configuration information in JSON format.
* `type` - (Required, String) Triggered task type.

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

