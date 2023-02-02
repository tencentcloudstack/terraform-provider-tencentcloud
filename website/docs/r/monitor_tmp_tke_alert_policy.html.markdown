---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_alert_policy"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_alert_policy"
description: |-
  Provides a resource to create a tke tmpAlertPolicy
---

# tencentcloud_monitor_tmp_tke_alert_policy

Provides a resource to create a tke tmpAlertPolicy

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_alert_policy" "basic" {
  instance_id = "prom-xxxxxx"
  alert_rule {
    name = "alert_rule-test"
    rules {
      name     = "rules-test"
      rule     = "(count(kube_node_status_allocatable_cpu_cores) by (cluster) -1)   / count(kube_node_status_allocatable_cpu_cores) by (cluster)"
      template = "The CPU requested by the Pod in the cluster {{ $labels.cluster }} is overloaded, and the current CPU application ratio is {{ $value | humanizePercentage }}"
      for      = "5m"
      labels {
        name  = "severity"
        value = "warning"
      }
    }
    notification {
      type    = "amp"
      enabled = true
      alert_manager {
        url = "xxx"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `alert_rule` - (Required, List) Alarm notification channels.
* `instance_id` - (Required, String) Instance Id.

The `alert_manager` object supports the following:

* `url` - (Required, String) Alertmanager url.
* `cluster_id` - (Optional, String) The ID of the cluster where the alertmanager is deployed. Note: This field may return null, indicating that a valid value could not be retrieved.
* `cluster_type` - (Optional, String) Alertmanager is deployed in the cluster type. Note: This field may return null, indicating that a valid value could not be retrieved.

The `alert_rule` object supports the following:

* `name` - (Required, String) Policy name.
* `rules` - (Required, List) A list of rules.
* `cluster_id` - (Optional, String) If the alarm policy is derived from the CRD resource definition of the user cluster, the ClusterId is the cluster ID to which it belongs.
* `id` - (Optional, String) Alarm policy ID. Note: This field may return null, indicating that a valid value could not be retrieved.
* `notification` - (Optional, List) Alarm channels, which may be returned using null in the template.
* `template_id` - (Optional, String) If the alarm is sent from a template, the TemplateId is the template id.
* `updated_at` - (Optional, String) Last modified time.

The `annotations` object supports the following:

* `name` - (Required, String) Name of map.
* `value` - (Required, String) Value of map.

The `labels` object supports the following:

* `name` - (Required, String) Name of map.
* `value` - (Required, String) Value of map.

The `notification` object supports the following:

* `enabled` - (Required, Bool) Whether it is enabled.
* `type` - (Required, String) The channel type, which defaults to amp, supports the following `amp`, `webhook`, `alertmanager`.
* `alert_manager` - (Optional, List) If Type is alertmanager, the field is required. Note: This field may return null, indicating that a valid value could not be retrieved..
* `notify_way` - (Optional, Set) Alarm notification method. At present, there are SMS, EMAIL, CALL, WECHAT methods.
* `phone_arrive_notice` - (Optional, Bool) Telephone alerts reach notifications.
* `phone_circle_interval` - (Optional, Int) Effective end timeTelephone alarm wheel interval. Units: Seconds.
* `phone_circle_times` - (Optional, Int) PhoneCircleTimes.
* `phone_inner_interval` - (Optional, Int) Telephone alarm wheel intervals. Units: Seconds.
* `phone_notify_order` - (Optional, Set) Telephone alarm sequence.
* `receiver_groups` - (Optional, Set) Alert Receiving Group (User Group).
* `repeat_interval` - (Optional, String) Convergence time.
* `time_range_end` - (Optional, String) Effective end time.
* `time_range_start` - (Optional, String) The time from which it takes effect.
* `web_hook` - (Optional, String) If Type is webhook, the field is required. Note: This field may return null, indicating that a valid value could not be retrieved.

The `rules` object supports the following:

* `for` - (Required, String) Time of duration.
* `labels` - (Required, List) Extra labels.
* `name` - (Required, String) Rule name.
* `rule` - (Required, String) Prometheus statement.
* `template` - (Required, String) Alert sending template.
* `annotations` - (Optional, List) Refer to annotations in prometheus rule.
* `describe` - (Optional, String) A description of the rule.
* `rule_state` - (Optional, Int) Alarm rule status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



