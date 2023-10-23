---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_basic_alarms"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_basic_alarms"
description: |-
  Use this data source to query detailed information of monitor basic_alarms
---

# tencentcloud_monitor_alarm_basic_alarms

Use this data source to query detailed information of monitor basic_alarms

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_basic_alarms" "alarms" {
  module             = "monitor"
  start_time         = 1696990903
  end_time           = 1697098903
  occur_time_order   = "DESC"
  project_ids        = [0]
  view_names         = ["cvm_device"]
  alarm_status       = [1]
  instance_group_ids = [5497073]
  metric_names       = ["cpu_usage"]
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Interface module name, current value monitor.
* `alarm_status` - (Optional, Set: [`Int`]) Filter based on alarm status.
* `end_time` - (Optional, Int) End time, default to current timestamp.
* `instance_group_ids` - (Optional, Set: [`Int`]) Filter based on instance group ID.
* `metric_names` - (Optional, Set: [`String`]) Filter by indicator name.
* `obj_like` - (Optional, String) Filter based on alarm objects.
* `occur_time_order` - (Optional, String) Sort by occurrence time, taking ASC or DESC values.
* `project_ids` - (Optional, Set: [`Int`]) Filter based on project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, Int) Start time, default to one day is timestamp.
* `view_names` - (Optional, Set: [`String`]) Filter based on policy type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alarms` - Alarm List.
  * `alarm_status` - Alarm status, ALARM indicates not recovered; OK indicates that it has been restored; NO_ DATA indicates insufficient data; NO_ CONF indicates that it has expired.
  * `alarm_type` - Alarm type, 0 represents indicator alarm, 2 represents product event alarm, and 3 represents platform event alarm.
  * `content` - Alarm content.
  * `dimensions` - Alarm object dimension information.
  * `duration` - Duration in seconds.
  * `first_occur_time` - Time of occurrence.
  * `group_id` - Policy Group ID.
  * `group_name` - Policy Group Name.
  * `id` - The ID of this alarm.
  * `instance_group` - Instance Group Information.
    * `instance_group_id` - Instance Group ID.
    * `instance_group_name` - Instance Group Name.
  * `last_occur_time` - End time.
  * `metric_id` - Indicator ID.
  * `metric_name` - Indicator Name.
  * `notify_way` - Notification method.
  * `obj_id` - Alarm object ID.
  * `obj_name` - Alarm Object.
  * `project_id` - Project ID.
  * `project_name` - Entry name.
  * `region` - Region.
  * `status` - Alarm status ID, 0 indicates not recovered; 1 indicates that it has been restored; 2,3,5 indicates insufficient data; 4 indicates it has expired.
  * `view_name` - Policy Type.
  * `vpc` - VPC, only CVM has.
* `warning` - Remarks.


