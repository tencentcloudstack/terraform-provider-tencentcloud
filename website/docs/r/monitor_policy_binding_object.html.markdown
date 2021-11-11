---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_policy_binding_object"
sidebar_current: "docs-tencentcloud-resource-monitor_policy_binding_object"
description: |-
  Provides a resource for bind objects to a alarm policy resource.
---

# tencentcloud_monitor_policy_binding_object

Provides a resource for bind objects to a alarm policy resource.

## Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}
resource "tencentcloud_monitor_alarm_policy" "policy" {
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

#for cvm
resource "tencentcloud_monitor_policy_binding_object" "binding" {
  policy_id = tencentcloud_monitor_alarm_policy.policy.id

  dimensions {
    dimensions_json = "{\"unInstanceId\":\"${data.tencentcloud_instances.instances.instance_list[0].instance_id}\"}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dimensions` - (Required, ForceNew) A list objects. Each element contains the following attributes:
* `policy_id` - (Required, ForceNew) Alarm policy ID for binding objects.

The `dimensions` object supports the following:

* `dimensions_json` - (Required, ForceNew) Represents a collection of dimensions of an object instance, json format.eg:'{"unInstanceId":"ins-ot3cq4bi"}'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



