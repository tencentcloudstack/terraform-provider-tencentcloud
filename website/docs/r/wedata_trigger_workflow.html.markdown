---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_workflow"
sidebar_current: "docs-tencentcloud-resource-wedata_trigger_workflow"
description: |-
  Provides a resource to create a wedata trigger workflow
---

# tencentcloud_wedata_trigger_workflow

Provides a resource to create a wedata trigger workflow

## Example Usage

```hcl
resource "tencentcloud_wedata_trigger_workflow" "workflow" {
  bundle_id          = null
  bundle_info        = null
  owner_uin          = 100044349576
  parent_folder_path = "/默认文件夹"
  project_id         = 3108707295180644352
  workflow_desc      = null
  workflow_name      = "tf-test1"
  general_task_params {
    type  = "SPARK_SQL"
    value = "a=b\nb=c\nc=d\nd=e"
  }
  trigger_workflow_scheduler_configurations {
    config_mode                     = "COMMON"
    crontab_expression              = "0 0 * * * ? *"
    cycle_type                      = "DAY_CYCLE"
    end_time                        = "2099-12-31 23:59:59"
    extra_info                      = null
    file_arrival_path               = null
    schedule_time_zone              = "UTC+8"
    scheduler_status                = "ACTIVE"
    start_time                      = "2026-01-09 00:00:00"
    trigger_minimum_interval_second = 0
    trigger_mode                    = "TIME_TRIGGER"
    trigger_wait_time_second        = 0
  }
  workflow_params {
    param_key   = "aaa"
    param_value = "bbb"
  }
  workflow_params {
    param_key   = "bbb"
    param_value = "ccc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `parent_folder_path` - (Required, String) Parent folder path.
* `project_id` - (Required, String) Project ID.
* `workflow_name` - (Required, String) Workflow name.
* `bundle_id` - (Optional, String) Bundle ID.
* `bundle_info` - (Optional, String) Bundle information.
* `general_task_params` - (Optional, List) General task parameter configuration.
* `owner_uin` - (Optional, String) Workflow owner ID.
* `trigger_workflow_scheduler_configurations` - (Optional, List) Unified scheduling configuration.
* `workflow_desc` - (Optional, String) Workflow description.
* `workflow_params` - (Optional, List) Workflow parameters.

The `general_task_params` object supports the following:

* `type` - (Optional, String) General task parameter type, currently only SPARK_SQL is supported.
* `value` - (Optional, String) General task parameter content; multiple parameters are separated by semicolons (;).

The `trigger_workflow_scheduler_configurations` object supports the following:

* `trigger_mode` - (Required, String) Trigger mode: Scheduled trigger:  `TIME_TRIGGER`; Continuous run: `CONTINUE_RUN`; File arrival: `FILE_ARRIVAL`. Notes: For `TIME_TRIGGER` and `CONTINUE_RUN` modes, SchedulerStatus, SchedulerTimeZone, StartTime, EndTime, ConfigMode, CycleType, and CrontabExpression are required; For `FILE_ARRIVAL` mode, FileArrivalPath, TriggerMinimumIntervalSecond, and TriggerWaitTimeSecond are required.
* `config_mode` - (Optional, String) Configuration mode, COMMON or CRON_EXPRESSION.
* `crontab_expression` - (Optional, String) Cron expression.
* `cycle_type` - (Optional, String) Cycle type. Supported values: `ONEOFF_CYCLE`: One-time; `YEAR_CYCLE`: Yearly; `MONTH_CYCLE`: Monthly; `WEEK_CYCLE`: Weekly; `DAY_CYCLE`: Daily; `HOUR_CYCLE`: Hourly; `MINUTE_CYCLE`: Minutely; `CRONTAB_CYCLE`: Crontab expression.
* `end_time` - (Optional, String) Schedule end time.
* `extra_info` - (Optional, String) WorkflowTriggerConfig converted to JSON format, used for reconciliation.
* `file_arrival_path` - (Optional, String) Listening path in the storage system for file arrival mode.
* `schedule_time_zone` - (Optional, String) Scheduler time zone.
* `scheduler_status` - (Optional, String) Trigger status, ACTIVE or PAUSED.
* `start_time` - (Optional, String) Schedule effective start time.
* `trigger_minimum_interval_second` - (Optional, Int) Minimum trigger interval in file arrival mode (seconds).
* `trigger_wait_time_second` - (Optional, Int) Trigger wait time in file arrival mode (seconds).

The `workflow_params` object supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

wedata trigger_workflow can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_trigger_workflow.trigger_workflow project_id#workflow_id
```

