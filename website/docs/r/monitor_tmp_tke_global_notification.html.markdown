---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_tke_global_notification"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_tke_global_notification"
description: |-
  Provides a resource to create a tmp tke global notification
---

# tencentcloud_monitor_tmp_tke_global_notification

Provides a resource to create a tmp tke global notification

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_global_notification" "basic" {
  instance_id = "prom-xxxxxx"
  notification {
    enabled = true
    type    = "webhook"
    alert_manager {
      cluster_id   = ""
      cluster_type = ""
      url          = ""
    }
    web_hook              = ""
    repeat_interval       = "5m"
    time_range_start      = "00:00:00"
    time_range_end        = "23:59:59"
    notify_way            = ["SMS", "EMAIL"]
    receiver_groups       = []
    phone_notify_order    = []
    phone_circle_times    = 0
    phone_inner_interval  = 0
    phone_circle_interval = 0
    phone_arrive_notice   = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance Id.
* `notification` - (Required, List) Alarm notification channels.

The `alert_manager` object supports the following:

* `url` - (Required, String) Alert manager url.
* `cluster_id` - (Optional, String) Cluster id.
* `cluster_type` - (Optional, String) Cluster type.

The `notification` object supports the following:

* `enabled` - (Required, Bool) Alarm notification switch.
* `type` - (Required, String) Alarm notification type, Valid values: `amp`, `webhook`, `alertmanager`.
* `alert_manager` - (Optional, List) Alert manager, if Type is `alertmanager`, this field is required.
* `notify_way` - (Optional, Set) Alarm notification method, Valid values: `SMS`, `EMAIL`, `CALL`, `WECHAT`.
* `phone_arrive_notice` - (Optional, Bool) Phone Alarm Reach Notification, NotifyWay is `CALL`, and this parameter is used.
* `phone_circle_interval` - (Optional, Int) Telephone alarm off-wheel interval, NotifyWay is `CALL`, and this parameter is used.
* `phone_circle_times` - (Optional, Int) Number of phone alerts (user group), NotifyWay is `CALL`, and this parameter is used.
* `phone_inner_interval` - (Optional, Int) Interval between telephone alarm rounds, NotifyWay is `CALL`, and this parameter is used.
* `phone_notify_order` - (Optional, Set) Phone alert sequence, NotifyWay is `CALL`, and this parameter is used.
* `receiver_groups` - (Optional, Set) Alarm receiving group(user group).
* `repeat_interval` - (Optional, String) Convergence time.
* `time_range_end` - (Optional, String) Effective end time.
* `time_range_start` - (Optional, String) Effective start time.
* `web_hook` - (Optional, String) Web hook, if Type is `webhook`, this field is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



