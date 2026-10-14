---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_notice_callbacks"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_notice_callbacks"
description: |-
  Use this data source to query detailed information of monitor alarm_notice_callbacks
---

# tencentcloud_monitor_alarm_notice_callbacks

Use this data source to query detailed information of monitor alarm_notice_callbacks

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_notice_callbacks" "alarm_notice_callbacks" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `url_notices` - Alarm callback notification.
  * `end_time` - The number of seconds from the end of the notification day.
  * `is_valid` - Verified 0=No 1=Yes.
  * `start_time` - The number of seconds starting from the day of notification start time.
  * `url` - Callback URL (limited to 256 characters).
  * `validation_code` - Verification code.
  * `weekday` - Notification period 1-7 represents Monday to Sunday.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeAlarmNoticeCallbacks](https://www.tencentcloud.com/document/product/248/39301) - This API is used to obtain all the callback URLs of an alarm notification template.

