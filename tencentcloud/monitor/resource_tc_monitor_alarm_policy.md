Provides a alarm policy resource for monitor.

Example Usage

cvm_device alarm policy

```hcl
resource "tencentcloud_monitor_alarm_notice" "foo" {
  name                  = "tf-alarm_notice"
  notice_type           = "ALL"
  notice_language       = "zh-CN"

  user_notices    {
      receiver_type              = "USER"
      start_time                 = 0
      end_time                   = 1
      notice_way                 = ["SMS","EMAIL"]
      user_ids                   = [10001]
      group_ids                  = []
      phone_order                = [10001]
      phone_circle_times         = 2
      phone_circle_interval      = 50
      phone_inner_interval       = 60
      need_phone_arrive_notice   = 1
      phone_call_type            = "CIRCLE"
      weekday                    =[1,2,3,4,5,6,7]
  }

  url_notices {
      url    = "https://www.mytest.com/validate"
      end_time =  0
      start_time = 1
      weekday = [1,2,3,4,5,6,7]
  }

}

resource "tencentcloud_monitor_alarm_policy" "foo" {
  policy_name = "tf-policy"
  monitor_type = "MT_QCE"
  enable = 1
  project_id = 0
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
  notice_ids = [tencentcloud_monitor_alarm_notice.foo.id]

  trigger_tasks {
    type = "AS"
    task_config = "{\"Region\":\"ap-guangzhou\",\"Group\":\"asg-0z312312x\",\"Policy\":\"asp-ganig28\"}"
  }
}
```

k8s_cluster alarm policy

```hcl
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "k8s_cluster"
  notice_ids   = [
    "notice-l9ziyxw6",
  ]
  policy_name  = "TkeClusterNew"
  project_id   = 1244035

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
              Value    = [
                "ap-guangzhou",
              ]
            },
            {
              Key      = "tke_cluster_instance_id"
              Operator = "in"
              Value    = [
                "cls-czhtobea",
              ]
            },
          ],
        ]
        )
        type       = "DIMENSION"
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
              Value    = [
                "ap-guangzhou",
              ]
            },
            {
              Key      = "tke_cluster_instance_id"
              Operator = "in"
              Value    = [
                "cls-czhtobea",
              ]
            },
          ],
        ]
        )
        type       = "DIMENSION"
      }
    }
  }
}
```

cvm_device alarm policy binding cvm by tag
```
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "cvm_device"
  notice_ids   = [
    "notice-l9ziyxw6",
  ]
  policy_name  = "policy"
  project_id   = 0

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

Import

Alarm policy instance can be imported, e.g.

```
$ terraform import tencentcloud_monitor_alarm_policy.policy policy-id
```