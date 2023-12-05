Provides a resource for bind objects to a alarm policy resource.

Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}
resource "tencentcloud_monitor_alarm_policy" "policy" {
  policy_name = "hello"
  monitor_type = "MT_QCE"
  enable = 1
  project_id = 1244035
  namespace = "cvm_device"

  conditions {
    is_union_rule = 1
    rules {
      metric_name = "CpuUsage"
      period = 60
      operator = "ge"
      value = "89.9"
      continue_period = 1
      notice_frequency = 3600
      is_power_notice = 0
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
    type = "AS"
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
Import

Monitor Policy Binding Object can be imported, e.g.

```
$ terraform import tencentcloud_monitor_policy_binding_object.binding policyId
```