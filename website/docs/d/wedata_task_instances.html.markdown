---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_instances"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_instances"
description: |-
  Use this data source to query detailed information of wedata task instances
---

# tencentcloud_wedata_task_instances

Use this data source to query detailed information of wedata task instances

## Example Usage

```hcl
data "tencentcloud_wedata_task_instances" "wedata_task_instances" {
  project_id = "1859317240494305280"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `cycle_type` - (Optional, String) Task cycle type * ONEOFF_CYCLE: One-time * YEAR_CYCLE: Year * MONTH_CYCLE: Month * WEEK_CYCLE: Week * DAY_CYCLE: Day * HOUR_CYCLE: Hour * MINUTE_CYCLE: Minute * CRONTAB_CYCLE: Crontab expression type.
* `executor_group_id` - (Optional, String) Executor resource group ID.
* `folder_id` - (Optional, String) Task folder ID.
* `instance_state` - (Optional, String) Instance status - WAIT_EVENT: Waiting for event - WAIT_UPSTREAM: Waiting for upstream - WAIT_RUN: Waiting to run - RUNNING: Running - SKIP_RUNNING: Skipped running - FAILED_RETRY: Failed retry - EXPIRED: Failed - COMPLETED: Success.
* `instance_type` - (Optional, Int) Instance type - 0: Backfill type - 1: Periodic instance - 2: Non-periodic instance.
* `keyword` - (Optional, String) Task name or Task ID. Supports fuzzy search filtering. Multiple values separated by commas.
* `last_update_time_from` - (Optional, String) Instance last update time filter condition.Start time, format yyyy-MM-dd HH:mm:ss.
* `last_update_time_to` - (Optional, String) Instance last update time filter condition.End time, format yyyy-MM-dd HH:mm:ss.
* `owner_uin` - (Optional, String) Task owner ID.
* `result_output_file` - (Optional, String) Used to save results.
* `schedule_time_from` - (Optional, String) Instance scheduled time filter condition Start time, format yyyy-MM-dd HH:mm:ss.
* `schedule_time_to` - (Optional, String) Instance scheduled time filter condition End time, format yyyy-MM-dd HH:mm:ss.
* `sort_column` - (Optional, String) Result sorting field- SCHEDULE_DATE: Sort by scheduled time- START_TIME: Sort by execution start time- END_TIME: Sort by execution end time- COST_TIME: Sort by execution duration.
* `sort_type` - (Optional, String) Sorting order: - ASC; - DESC.
* `start_time_from` - (Optional, String) Instance execution start time filter condition Start time, format yyyy-MM-dd HH:mm:ss.
* `start_time_to` - (Optional, String) Instance execution start time filter condition.End time, format yyyy-MM-dd HH:mm:ss.
* `task_type_id` - (Optional, Int) Task type ID.
* `time_zone` - (Optional, String) Time zone. The time zone of the input time string, default UTC+8.
* `workflow_id` - (Optional, String) Task workflow ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Instance result set.


