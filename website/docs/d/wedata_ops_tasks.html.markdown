---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_tasks"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_tasks"
description: |-
  Use this data source to query detailed information of wedata ops tasks
---

# tencentcloud_wedata_ops_tasks

Use this data source to query detailed information of wedata ops tasks

## Example Usage

```hcl
data "tencentcloud_wedata_ops_tasks" "wedata_ops_tasks" {
  project_id        = "1859317240494305280"
  task_type_id      = 34
  workflow_id       = "d7184172-4879-11ee-ba36-b8cef6a5af5c"
  workflow_name     = "test1"
  folder_id         = "cee5780a-4879-11ee-ba36-b8cef6a5af5c"
  executor_group_id = "20230830105723839685"
  cycle_type        = "MINUTE_CYCLE"
  status            = "F"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `cycle_type` - (Optional, String) Task Cycle Type: ONEOFF_CYCLE: One-time, YEAR_CYCLE: Yearly, MONTH_CYCLE: Monthly, WEEK_CYCLE: Weekly, DAY_CYCLE: Daily, HOUR_CYCLE: Hourly, MINUTE_CYCLE: Minute-level, CRONTAB_CYCLE: Crontab expression-based.
* `executor_group_id` - (Optional, String) Executor Group ID.
* `folder_id` - (Optional, String) Folder ID.
* `owner_uin` - (Optional, String) Owner id.
* `result_output_file` - (Optional, String) Used to save results.
* `source_service_id` - (Optional, String) Data source ID.
* `status` - (Optional, String) Task Status: -Y: Running, -F: Stopped, -O: Frozen, -T: Stopping, -INVALID: Invalid.
* `target_service_id` - (Optional, String) Target data source id.
* `task_type_id` - (Optional, String) Task type Id. -20: common data sync, - 25:ETLTaskType, - 26:ETLTaskType, - 30:python, - 31:pyspark, - 34:HiveSQLTaskType, - 35:shell, - 36:SparkSQLTaskType, - 21:JDBCSQLTaskType, - 32:DLCTaskType, - 33:ImpalaTaskType, - 40:CDWTaskType, - 41:kettle, - 46:DLCSparkTaskType, -47: TiOne machine learning, - 48:TrinoTaskType, - 50:DLCPyspark39:spark, - 92:mr, -38: shell script, -70: hivesql script, -1000: common custom business.
* `time_zone` - (Optional, String) Time zone. defaults to UTC+8.
* `workflow_id` - (Optional, String) Workflow ID.
* `workflow_name` - (Optional, String) Workflow name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Task list.


