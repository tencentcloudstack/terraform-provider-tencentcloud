---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_workflow"
sidebar_current: "docs-tencentcloud-resource-wedata_workflow"
description: |-
  Provides a resource to create a wedata wedata_workflow
---

# tencentcloud_wedata_workflow

Provides a resource to create a wedata wedata_workflow

## Example Usage

```hcl
resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id         = 2905622749543821312
  workflow_name      = "test"
  parent_folder_path = "/tfmika"
  workflow_type      = "cycle"
}
```

## Argument Reference

The following arguments are supported:

* `parent_folder_path` - (Required, String) Parent folder path.
* `project_id` - (Required, String, ForceNew) Project id.
* `workflow_name` - (Required, String) Workflow name.
* `bundle_id` - (Optional, String) Bundle Id.
* `bundle_info` - (Optional, String) Bundle Information.
* `owner_uin` - (Optional, String) Workflow Owner ID.
* `workflow_desc` - (Optional, String) Workflow description.
* `workflow_params` - (Optional, Set) workflow parameter.
* `workflow_scheduler_configuration` - (Optional, List) Unified dispatch information.
* `workflow_type` - (Optional, String) Workflow type, value example: cycle cycle workflow;manual manual workflow, passed in cycle by default.

The `workflow_params` object supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

The `workflow_scheduler_configuration` object supports the following:

* `crontab_expression` - (Required, String) Crontab expression.
* `cycle_type` - (Required, String) Cycle type: Supported types are
ONEOFF_CYCLE: One-time
YEAR_CYCLE: Year
MONTH_CYCLE: Month
WEEK_CYCLE: Week
DAY_CYCLE: Day
HOUR_CYCLE: Hour
MINUTE_CYCLE: Minute
CRONTAB_CYCLE: crontab expression type.
* `end_time` - (Required, String) End time.
* `schedule_time_zone` - (Required, String) time zone.
* `self_depend` - (Required, String) Self-dependence, default value serial, values are: parallel, serial, orderly.
* `start_time` - (Required, String) Start time.
* `calendar_id` - (Optional, String) calendar id.
* `calendar_open` - (Optional, String) Do you want to turn on calendar scheduling 1 on 0 off.
* `clear_link` - (Optional, Bool) Workflows have cross-workflow dependencies and are scheduled using cron expressions. If you save unified scheduling, unsupported dependencies will be broken.
* `dependency_workflow` - (Optional, String) Workflow dependence, yes or no.
* `execution_end_time` - (Optional, String) Execution time right-closed interval, example: 23:59, only if the cycle type is MINUTE_CYCLE needs to be filled in.
* `execution_start_time` - (Optional, String) Execution time left-closed interval, example: 00:00, only if the cycle type is MINUTE_CYCLE needs to be filled in.
* `main_cyclic_config` - (Optional, String) Effective when ModifyCycleValue is 1, indicating the default modified upstream dependence-time dimension. The value is: 
* CRONTAB
* DAY
* HOUR
* LIST_DAY
* LIST_HOUR
 * LIST_MINUTE
 * MONTH
* RANGE_DAY
 * RANGE_HOUR
 * RANGE_MINUTE
* WEEK
* YEAR

https://capi.woa.com/object/detail? product=wedata&env=api_dev&version=2025-08-06&name=WorkflowSchedulerConfigurationInfo.
* `modify_cycle_value` - (Optional, String) 0: Do not modify 1: Change the upstream dependency configuration of the task to the default value.
* `subordinate_cyclic_config` - (Optional, String) Effective when ModifyCycleValue is 1, which means that the default modified upstream dependency-instance range
 value is: 
* ALL_DAY_OF_YEAR
* ALL_MONTH_OF_YEAR
* CURRENT
* CURRENT_DAY
* CURRENT_HOUR
* CURRENT_MINUTE
* CURRENT_MONTH
* CURRENT_WEEK
* CURRENT_YEAR
* PREVIOUS_BEGIN_OF_MONTH
* PREVIOUS_DAY
* PREVIOUS_DAY_LATER_OFFSET_HOUR
* PREVIOUS_DAY_LATER_OFFSET_MINUTE
* PREVIOUS_END_OF_MONTH
* PREVIOUS_FRIDAY
* PREVIOUS_HOUR
* PREVIOUS_HOUR_CYCLE
* PREVIOUS_HOUR_LATER_OFFSET_MINUTE
* PREVIOUS_MINUTE_CYCLE
* PREVIOUS_MONTH
* PREVIOUS_WEEK
* PREVIOUS_WEEKEND
* RECENT_DATE

https://capi.woa.com/object/detail? product=wedata&env=api_dev&version=2025-08-06&name=WorkflowSchedulerConfigurationInfo.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `workflow_id` - Workflow id.


## Import

wedata wedata_workflow can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_workflow.wedata_workflow wedata_workflow_id
```

