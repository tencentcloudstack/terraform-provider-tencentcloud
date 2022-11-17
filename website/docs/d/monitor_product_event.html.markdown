---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_product_event"
sidebar_current: "docs-tencentcloud-datasource-monitor_product_event"
description: |-
  Use this data source to query monitor events(There is a lot of data and it is recommended to output to a file)
---

# tencentcloud_monitor_product_event

Use this data source to query monitor events(There is a lot of data and it is recommended to output to a file)

## Example Usage

```hcl
data "tencentcloud_monitor_product_event" "cvm_event_data" {
  start_time      = 1588700283
  is_alarm_config = 0
  product_name    = ["cvm"]
}
```

## Argument Reference

The following arguments are supported:

* `dimensions` - (Optional, List) Dimensional composition of instance objects.
* `end_time` - (Optional, Int) End timestamp for this query, eg:`1588232111`. Default start time is `now-3000`.
* `event_name` - (Optional, List: [`String`]) Event name filtering, such as `guest_reboot` indicates that the machine restart.
* `instance_id` - (Optional, List: [`String`]) Affect objects, such as `ins-19708ino`.
* `is_alarm_config` - (Optional, Int) Alarm status configuration filter, 1means configured, 0(default) means not configured.
* `product_name` - (Optional, List: [`String`]) Product type filtering, such as `cvm` for cloud server.
* `project_id` - (Optional, List: [`String`]) Project ID filter.
* `region_list` - (Optional, List: [`String`]) Region filter, such as `gz`.
* `result_output_file` - (Optional, String) Used to store results.
* `start_time` - (Optional, Int) Start timestamp for this query, eg:`1588230000`. Default start time is `now-3600`.
* `status` - (Optional, List: [`String`]) Event status filter, value range `-`,`alarm`,`recover`, indicating recovered, unrecovered and stateless.
* `type` - (Optional, List: [`String`]) Event type filtering, with value range `abnormal`,`status_change`, indicating state change and abnormal events.

The `dimensions` object supports the following:

* `name` - (Optional, String) Instance dimension name, eg: `deviceWanIp` for internet ip.
* `value` - (Optional, String) Instance dimension value, eg: `119.119.119.119` for internet ip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list events. Each element contains the following attributes:
  * `addition_msg` - A list of addition message. Each element contains the following attributes:
    * `key` - The key of this addition message.
    * `name` - The name of this addition message.
    * `value` - The value of this addition message.
  * `dimensions` - A list of event dimensions. Each element contains the following attributes:
    * `key` - The key of this dimension.
    * `name` - The name of this dimension.
    * `value` - The value of this dimension.
  * `event_cname` - Event chinese name.
  * `event_ename` - Event english name.
  * `event_id` - Event ID.
  * `event_name` - Event short name.
  * `group_info` - A list of group info. Each element contains the following attributes:
    * `group_id` - Policy group ID.
    * `group_name` - Policy group name.
  * `instance_id` - The instance ID of this event.
  * `instance_name` - The name of this instance.
  * `is_alarm_config` - Whether to configure alarm.
  * `product_cname` - Product chinese name.
  * `product_ename` - Product english name.
  * `product_name` - Product short name.
  * `project_id` - Project ID of this instance.
  * `region` - The region of this instance.
  * `start_time` - The start timestamp of this event.
  * `status` - The status of this event.
  * `support_alarm` - Whether to support alarm.
  * `type` - The type of this event.
  * `update_time` - The update timestamp of this event.


