---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_integration_offline_task"
sidebar_current: "docs-tencentcloud-resource-wedata_integration_offline_task"
description: |-
  Provides a resource to create a wedata integration_offline_task
---

# tencentcloud_wedata_integration_offline_task

Provides a resource to create a wedata integration_offline_task

## Example Usage

```hcl
resource "tencentcloud_wedata_integration_offline_task" "example" {
  project_id  = "1455251608631480391"
  cycle_step  = 1
  delay_time  = 0
  end_time    = "2099-12-31 00:00:00"
  notes       = "terraform example demo."
  start_time  = "2023-12-31 00:00:00"
  task_name   = "tf_example"
  type_id     = 27
  task_action = "0,3,4"
  task_mode   = "1"

  task_info {
    executor_id = "20230313175748567418"
    config {
      name  = "Args"
      value = "args"
    }
    config {
      name  = "dirtyDataThreshold"
      value = "0"
    }
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "syncRateLimitUnit"
      value = "0"
    }
    ext_config {
      name  = "TaskAlarmRegularList"
      value = "73"
    }
    incharge = "demo_user"
    offline_task_add_entity {
      cycle_type         = 3
      crontab_expression = "0 0 1 * * ?"
      retry_wait         = 5
      retriable          = 1
      try_limit          = 5
      self_depend        = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cycle_step` - (Required, Int) Interval time of scheduling, the minimum value: 1.
* `delay_time` - (Required, Int) Execution time, unit is minutes, only available for day/week/month/year scheduling. For example, daily scheduling is executed once every day at 02:00, and the delayTime is 120 minutes.
* `end_time` - (Required, String) Effective end time, the format is yyyy-MM-dd HH:mm:ss.
* `notes` - (Required, String) Description information.
* `project_id` - (Required, String) Project ID.
* `start_time` - (Required, String) Effective start time, the format is yyyy-MM-dd HH:mm:ss.
* `task_action` - (Required, String) Scheduling configuration: flexible period configuration, only available for hourly/weekly/monthly/yearly scheduling. If the hourly task is specified to run at 0:00, 3:00 and 4:00 every day, it is 0,3,4.
* `task_info` - (Required, List) Task Information.
* `task_mode` - (Required, String) Task display mode, 0: canvas mode, 1: form mode.
* `task_name` - (Required, String) Task name.
* `type_id` - (Required, Int) Task type ID, for intgration task the value is 27.

The `config` object supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `execute_context` object supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `ext_config` object supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `mappings` object supports the following:

* `ext_config` - (Optional, List) Node extension configuration information.
* `schema_mappings` - (Optional, List) Schema mapping information.
* `sink_id` - (Optional, String) Sink node ID.
* `source_id` - (Optional, String) Source node ID.
* `source_schema` - (Optional, List) Source node schema information.

The `offline_task_add_entity` object supports the following:

* `crontab_expression` - (Optional, String) Crontab expression.
* `cycle_type` - (Optional, Int) Scheduling type, 0: crontab type, 1: minutes, 2: hours, 3: days, 4: weeks, 5: months, 6: one-time, 7: user-driven, 10: elastic period (week), 11: elastic period (month), 12: year, 13: instant trigger.
* `dependency_workflow` - (Optional, String) Whether to support workflow dependencies: yes / no, default value: no.
* `execution_end_time` - (Optional, String) Scheduling execution end time.
* `execution_start_time` - (Optional, String) Scheduling execution start time.
* `instance_init_strategy` - (Optional, String) Instance initialization strategy.
* `product_name` - (Optional, String) Product name: DATA_INTEGRATION.
* `retriable` - (Optional, Int) Whether to retry.
* `retry_wait` - (Optional, Int) Retry waiting time, unit is minutes.
* `run_priority` - (Optional, Int) Task running priority.
* `self_depend` - (Optional, Int) Self-dependent rules, 1: Ordered serial one at a time, queued execution, 2: Unordered serial one at a time, not queued execution, 3: Parallel, multiple at once.
* `task_auto_submit` - (Optional, Bool) Whether to automatically submit.
* `try_limit` - (Optional, Int) Number of retries.
* `workflow_name` - (Optional, String) The name of the workflow to which the task belongs.

The `properties` object supports the following:

* `name` - (Optional, String) Attributes name.
* `value` - (Optional, String) Attributes value.

The `schema_mappings` object supports the following:

* `sink_schema_id` - (Required, String) Schema ID from sink node.
* `source_schema_id` - (Required, String) Schema ID from source node.

The `source_schema` object supports the following:

* `id` - (Required, String) Schema ID.
* `name` - (Required, String) Schema name.
* `type` - (Required, String) Schema type.
* `alias` - (Optional, String) Schema alias.
* `comment` - (Optional, String) Schema comment.
* `properties` - (Optional, List) Schema extended attributes.
* `value` - (Optional, String) Schema value.

The `task_info` object supports the following:

* `app_id` - (Optional, String) User App Id.
* `config` - (Optional, List) Task configuration.
* `create_time` - (Optional, String) Create time.
* `creator_uin` - (Optional, String) Creator User ID.
* `data_proxy_url` - (Optional, Set) Data proxy url.
* `execute_context` - (Optional, List) Execute context.
* `executor_group_name` - (Optional, String) Executor group name.
* `executor_id` - (Optional, String) Executor resource ID.
* `ext_config` - (Optional, List) Node extension configuration information.
* `has_version` - (Optional, Bool) Whether the task been submitted.
* `in_long_manager_url` - (Optional, String) InLong manager url.
* `in_long_manager_version` - (Optional, String) InLong manager version.
* `in_long_stream_id` - (Optional, String) InLong stream id.
* `incharge` - (Optional, String) Incharge user.
* `input_datasource_type` - (Optional, String) Input datasource type.
* `instance_version` - (Optional, Int) Instance version.
* `last_run_time` - (Optional, String) The last time the task was run.
* `locked` - (Optional, Bool) Whether the task been locked.
* `locker` - (Optional, String) User locked task.
* `mappings` - (Optional, List) Node mapping.
* `num_records_in` - (Optional, Int) Number of reads.
* `num_records_out` - (Optional, Int) Number of writes.
* `num_restarts` - (Optional, Int) Times of restarts.
* `offline_task_add_entity` - (Optional, List) Offline task scheduling configuration.
* `operator_uin` - (Optional, String) Operator User ID.
* `output_datasource_type` - (Optional, String) Output datasource type.
* `owner_uin` - (Optional, String) Owner User ID.
* `read_phase` - (Optional, Int) Reading stage, 0: full amount, 1: partial full amount, 2: all incremental.
* `reader_delay` - (Optional, Float64) Read latency.
* `running_cu` - (Optional, Float64) The amount of resources consumed by real-time task.
* `schedule_task_id` - (Optional, String) Task scheduling id (job id such as oceanus or us).
* `status` - (Optional, Int) Task status 1. Not started | Task initialization, 2. Task starting, 3. Running, 4. Paused, 5. Task stopping, 6. Stopped, 7. Execution failed, 8. deleted, 9. Locked, 404. unknown status.
* `stop_time` - (Optional, String) The time the task was stopped.
* `submit` - (Optional, Bool) Whether the task version has been submitted for operation and maintenance.
* `switch_resource` - (Optional, Int) Resource tiering status, 0: in progress, 1: successful, 2: failed.
* `sync_type` - (Optional, Int) Synchronization type: 1. Whole database synchronization, 2. Single table synchronization.
* `task_alarm_regular_list` - (Optional, Set) Task alarm regular.
* `task_group_id` - (Optional, String) Inlong Task Group ID.
* `task_mode` - (Optional, String) Task display mode, 0: canvas mode, 1: form mode.
* `update_time` - (Optional, String) Update time.
* `workflow_id` - (Optional, String) The workflow id to which the task belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - Task ID.


## Import

wedata integration_offline_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_offline_task.example 1612982498218618880#20231102200955095
```

