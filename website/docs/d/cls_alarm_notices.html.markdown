---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_alarm_notices"
sidebar_current: "docs-tencentcloud-datasource-cls_alarm_notices"
description: |-
  Use this data source to query detailed information of cls alarm notices
---

# tencentcloud_cls_alarm_notices

Use this data source to query detailed information of cls alarm notices

## Example Usage

### Query all cls alarm notices

```hcl
data "tencentcloud_cls_alarm_notices" "example" {}
```

### Query by filters

```hcl
data "tencentcloud_cls_alarm_notices" "example" {
  filters {
    key    = "name"
    values = ["tf-example"]
  }

  filters {
    key    = "alarmNoticeId"
    values = ["notice-c2af43ee-1a4b-4c4a-ae3e-f81481280101"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Maximum 10 filters, each with up to 5 values. Multiple values within the same filter use OR logic, multiple filters use AND logic.
* `has_alarm_shield_count` - (Optional, Bool) Whether to query alarm shield count statistics. Default is false.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Required, String) Filter field name. Supported values: `name` (alarm notice group name), `alarmNoticeId` (alarm notice ID), `uid` (receiver user ID), `groupId` (receiver user group ID), `deliverFlag` (delivery status: 1-not enabled, 2-enabled, 3-abnormal).
* `values` - (Required, Set) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alarm_notices` - List of alarm notice configurations.
  * `alarm_notice_id` - Alarm notice ID.
  * `alarm_shield_count` - Alarm shield count statistics.
    * `total_count` - Total count of shielded alarms.
  * `alarm_shield_status` - Alarm shield status (0: not shielded, 1: shielded).
  * `callback_prioritize` - Whether webhook callback takes priority.
  * `create_time` - Creation time.
  * `deliver_flag` - Delivery flag (1: not enabled, 2: enabled, 3: abnormal).
  * `deliver_status` - Delivery status (0: delivered, 1: not delivered).
  * `jump_domain` - Jump domain for alarm callback.
  * `name` - Alarm notice name.
  * `notice_receivers` - List of notice receivers.
    * `end_time` - Allowed notification end time.
    * `index` - Index order.
    * `receiver_channels` - Notification channels (Email, Sms, WeChat, Phone).
    * `receiver_ids` - Receiver IDs.
    * `receiver_type` - Receiver type. Can be Uin or Group.
    * `start_time` - Allowed notification start time.
  * `notice_rules` - List of notice rules.
    * `day_of_week` - Days of week (0-6, 0 is Sunday).
    * `jump_domain` - Jump domain.
    * `notice_receivers` - Notice receivers for this rule.
      * `end_time` - End time.
      * `index` - Index.
      * `receiver_channels` - Notification channels.
      * `receiver_ids` - Receiver IDs.
      * `receiver_type` - Receiver type.
      * `start_time` - Start time.
    * `notify_way` - Notification ways.
    * `receiver_type` - Receiver type.
    * `repeat_interval` - Repeat interval in minutes.
    * `time_range_end` - Effective end time (24-hour format HH:mm:ss).
    * `time_range_start` - Effective start time (24-hour format HH:mm:ss).
    * `web_callbacks` - Webhook callbacks for this rule.
      * `body` - Body.
      * `callback_type` - Callback type.
      * `headers` - Headers.
      * `index` - Index.
      * `method` - HTTP method.
      * `url` - Callback URL.
  * `tags` - Tag list.
    * `key` - Tag key.
    * `value` - Tag value.
  * `update_time` - Last update time.
  * `web_callbacks` - List of webhook callbacks.
    * `body` - Request body.
    * `callback_type` - Callback type. WeCom or Http or DingTalk or Lark or Webhook.
    * `headers` - Request headers.
    * `index` - Index order.
    * `method` - HTTP method. GET or POST.
    * `url` - Callback URL.


