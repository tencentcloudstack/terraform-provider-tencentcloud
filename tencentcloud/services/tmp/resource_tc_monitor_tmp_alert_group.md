Provides a resource to create a monitor tmp_alert_group

Example Usage

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

Import

monitor tmp_alert_group can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_alert_group.tmp_alert_group instance_id#group_id
```