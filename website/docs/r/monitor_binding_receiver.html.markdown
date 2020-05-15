---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_binding_receiver"
sidebar_current: "docs-tencentcloud-resource-monitor_binding_receiver"
description: |-
  Provides a resource for bind receivers to a policy group resource.
---

# tencentcloud_monitor_binding_receiver

Provides a resource for bind receivers to a policy group resource.

## Example Usage

```hcl
data "tencentcloud_cam_groups" "groups" {
  //You should first create a user group with CAM
}

resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
}

resource "tencentcloud_monitor_binding_receiver" "receiver" {
  group_id = tencentcloud_monitor_policy_group.group.id
  receivers {
    start_time          = 0
    end_time            = 86399
    notify_way          = ["SMS"]
    receiver_type       = "group"
    receiver_group_list = [data.tencentcloud_cam_groups.groups.group_list[0].group_id]
    receive_language    = "en-US"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) Policy group id for binding receivers.
* `receivers` - (Optional) A list of receivers(will overwrite the configuration of the server or other resources). Each element contains the following attributes:

The `receivers` object supports the following:

* `notify_way` - (Required) Method of warning notification.Optional `CALL`,`EMAIL`,`SITE`,`SMS`,`WECHAT`.
* `receiver_type` - (Required) Receive type. Optional `group`,`user`.
* `end_time` - (Optional) End of alarm period. Meaning with `start_time`.
* `receive_language` - (Optional) Alert sending language. Optional `en-US`,`zh-CN`.
* `receiver_group_list` - (Optional) Alarm receive group id list.
* `receiver_user_list` - (Optional) Alarm receiver id list.
* `start_time` - (Optional) Alarm period start time.Range [0,86399], which removes the date after it is converted to Beijing time as a Unix timestamp, for example 7200 means '10:0:0'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



