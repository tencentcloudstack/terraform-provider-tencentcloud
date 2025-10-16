---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task"
sidebar_current: "docs-tencentcloud-resource-wedata_task"
description: |-
  Provides a resource to create a wedata wedata_task
---

# tencentcloud_wedata_task

Provides a resource to create a wedata wedata_task

## Example Usage

```hcl
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
  folder_name        = "tftest"
}

resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id         = 2905622749543821312
  workflow_name      = "test_workflow1"
  parent_folder_path = "${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.parent_folder_path}${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.folder_name}"
  workflow_type      = "cycle"
}

resource "tencentcloud_wedata_task" "wedata_task" {
  project_id = 2905622749543821312
  task_base_attribute {
    task_name    = "tfTask"
    task_type_id = 30
    workflow_id  = tencentcloud_wedata_workflow.wedata_workflow.workflow_id
  }
  task_configuration {
    code_content = base64encode("Hello World")
    task_ext_configuration_list {
      param_key   = "bucket"
      param_value = "wedata-fusion-bjjr-1257305158"
    }
    task_ext_configuration_list {
      param_key   = "ftp.file.name"
      param_value = "/datastudio/project/2905622749543821312/tftest/test_workflow1/tfTask.py"
    }
    task_ext_configuration_list {
      param_key   = "tenantId"
      param_value = "1257305158"
    }
    task_ext_configuration_list {
      param_key   = "region"
      param_value = "ap-beijing-fsi"
    }
    task_ext_configuration_list {
      param_key   = "extraInfo"
      param_value = "{\"fromMapping\":false}"
    }
    task_ext_configuration_list {
      param_key   = "ssmDynamicSkSwitch"
      param_value = "ON"
    }
    task_ext_configuration_list {
      param_key   = "calendar_open"
      param_value = "0"
    }
    task_ext_configuration_list {
      param_key   = "specLabelConfItems"
      param_value = "eyJzcGVjTGFiZWxDb25mSXRlbXMiOltdfQ=="
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTL"
      param_value = "-1"
    }
  }
  task_scheduler_configuration {
    cycle_type = "DAY_CYCLE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `task_base_attribute` - (Required, List) Basic task attributes.
* `task_configuration` - (Required, List) Task configuration.
* `task_scheduler_configuration` - (Required, List) Task scheduling configuration.

The `dependency_strategy` object of `upstream_dependency_config_list` supports the following:

* `polling_null_strategy` - (Optional, String) Strategy for waiting for upstream task instances: EXECUTING; WAITING.
* `task_dependency_executing_strategies` - (Optional, Set) This field is required only when PollingNullStrategy is EXECUTING, List type: NOT_EXIST (default, when minute depends on minute/hour depends on hour, parent instance is not within the downstream instance scheduling time range); PARENT_EXPIRED (parent instance failed); PARENT_TIMEOUT (parent instance timed out). If any of the above scenarios is met, the parent task instance dependency judgment can be passed; otherwise, it is necessary to wait for the parent instance.
* `task_dependency_executing_timeout_value` - (Optional, Int) This field is required only when TaskDependencyExecutingStrategies contains PARENT_TIMEOUT, the timeout time for downstream tasks depending on parent instance execution, unit: minutes.

The `event_listener_list` object of `task_scheduler_configuration` supports the following:

* `event_broadcast_type` - (Required, String) Event broadcast type: SINGLE, BROADCAST.
* `event_name` - (Required, String) Event name.
* `event_sub_type` - (Required, String) Event cycle: SECOND, MIN, HOUR, DAY.
* `properties_list` - (Optional, List) Extended information.

The `param_task_in_list` object of `task_scheduler_configuration` supports the following:

* `from_param_key` - (Required, String) Parent task parameter key.
* `from_task_id` - (Required, String) Parent task ID.
* `param_desc` - (Required, String) Parameter description: format is project_identifier.task_name.parameter_name; example: project_wedata_1.sh_250820_104107.pp_out.
* `param_key` - (Required, String) Parameter name.

The `param_task_out_list` object of `task_scheduler_configuration` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter definition.

The `properties_list` object of `event_listener_list` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

The `task_base_attribute` object supports the following:

* `task_name` - (Required, String) Task name.
* `task_type_id` - (Required, String) Task type ID:

* 21:JDBC SQL
* 23:TDSQL-PostgreSQL
* 26:OfflineSynchronization
* 30:Python
* 31:PySpark
* 32:DLC SQL
* 33:Impala
* 34:Hive SQL
* 35:Shell
* 36:Spark SQL
* 38:Shell Form Mode
* 39:Spark
* 40:TCHouse-P
* 41:Kettle
* 42:Tchouse-X
* 43:TCHouse-X SQL
* 46:DLC Spark
* 47:TiOne
* 48:Trino
* 50:DLC PySpark
* 92:MapReduce
* 130:Branch Node
* 131:Merged Node
* 132:Notebook
* 133:SSH
* 134:StarRocks
* 137:For-each
* 138:Setats SQL.
* `workflow_id` - (Required, String) Workflow ID.
* `owner_uin` - (Optional, String) Task owner ID, defaults to current user.
* `task_description` - (Optional, String) Task description.

The `task_configuration` object supports the following:

* `broker_ip` - (Optional, String) Specified running node.
* `bundle_id` - (Optional, String) ID used by Bundle.
* `bundle_info` - (Optional, String) Bundle information.
* `code_content` - (Optional, String) Base64 encoded code content.
* `data_cluster` - (Optional, String) Cluster ID.
* `resource_group` - (Optional, String) Resource group ID: Need to obtain ExecutorGroupId via DescribeNormalSchedulerExecutorGroups.
* `source_service_id` - (Optional, String) Source data source ID, separated by `;`, need to obtain via DescribeDataSourceWithoutInfo.
* `target_service_id` - (Optional, String) Target data source ID, separated by `;`, need to obtain via DescribeDataSourceWithoutInfo.
* `task_ext_configuration_list` - (Optional, Set) Task extended attribute configuration list.
* `task_scheduling_parameter_list` - (Optional, Set) Scheduling parameters.
* `yarn_queue` - (Optional, String) Resource pool queue name, need to obtain via DescribeProjectClusterQueues.

The `task_ext_configuration_list` object of `task_configuration` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

The `task_output_registry_list` object of `task_scheduler_configuration` supports the following:

* `data_flow_type` - (Required, String) Input/output table type
      Input stream
 UPSTREAM,
      Output stream
  DOWNSTREAM.
* `database_name` - (Required, String) Database name.
* `datasource_id` - (Required, String) Data source ID.
* `partition_name` - (Required, String) Partition name.
* `table_name` - (Required, String) Table name.
* `table_physical_id` - (Required, String) Table physical unique ID.
* `db_guid` - (Optional, String) Database unique identifier.
* `table_guid` - (Optional, String) Table unique identifier.

The `task_scheduler_configuration` object supports the following:

* `allow_redo_type` - (Optional, String) Rerun & backfill configuration, defaults to ALL; ALL: can rerun or backfill after success or failure; FAILURE: cannot rerun or backfill after success, can rerun or backfill after failure; NONE: cannot rerun or backfill after success or failure.
* `calendar_id` - (Optional, String) Calendar scheduling calendar ID.
* `calendar_open` - (Optional, String) Calendar scheduling: Values are 0 and 1, 1 for enabled, 0 for disabled, defaults to 0.
* `crontab_expression` - (Optional, String) Cron expression, defaults to 0 0 0 * * `?` *.
* `cycle_type` - (Optional, String) Cycle type: Defaults to DAY_CYCLE.

Supported types are

* ONEOFF_CYCLE: One-time
* YEAR_CYCLE: Yearly
* MONTH_CYCLE: Monthly
* WEEK_CYCLE: Weekly
* DAY_CYCLE: Daily
* HOUR_CYCLE: Hourly
* MINUTE_CYCLE: Minutely
* CRONTAB_CYCLE: Crontab expression type.
* `end_time` - (Optional, String) End date, defaults to 2099-12-31 23:59:59.
* `event_listener_list` - (Optional, List) Event array.
* `execution_end_time` - (Optional, String) Execution time right-closed interval, default 23:59.
* `execution_start_time` - (Optional, String) Execution time left-closed interval, default 00:00.
* `execution_ttl` - (Optional, String) Timeout handling strategy - running time timeout (unit: minutes) defaults to -1.
* `init_strategy` - (Optional, String) **Instance generation strategy**
* T_PLUS_0: T+0 generation, default strategy
* T_PLUS_1: T+1 generation.
* `max_retry_attempts` - (Optional, String) Retry strategy - maximum number of attempts, default: 4.
* `param_task_in_list` - (Optional, List) Input parameter array.
* `param_task_out_list` - (Optional, List) Output parameter array.
* `retry_wait` - (Optional, String) Retry strategy - retry waiting time, unit: minutes: default: 5.
* `run_priority` - (Optional, String) Task scheduling priority: 4 for high, 5 for medium, 6 for low, default: 6.
* `schedule_run_type` - (Optional, String) Scheduling type: 0 Normal scheduling 1 Empty run scheduling, defaults to 0.
* `schedule_time_zone` - (Optional, String) Time zone, defaults to UTC+8.
* `self_depend` - (Optional, String) Self-dependency, default value serial, values: parallel, serial, orderly.
* `start_time` - (Optional, String) Effective date, defaults to 00:00:00 of current date.
* `task_output_registry_list` - (Optional, List) Output registration.
* `upstream_dependency_config_list` - (Optional, List) Upstream dependency array.
* `wait_execution_total_ttl` - (Optional, String) Timeout handling strategy - total waiting time timeout (unit: minutes) defaults to -1.

The `task_scheduling_parameter_list` object of `task_configuration` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

The `upstream_dependency_config_list` object of `task_scheduler_configuration` supports the following:

* `main_cyclic_config` - (Required, String) Main dependency configuration, values:

* CRONTAB
* DAY
* HOUR
* LIST_DAY
* LIST_HOUR
* LIST_MINUTE
* MINUTE
* MONTH
* RANGE_DAY
* RANGE_HOUR
* RANGE_MINUTE
* WEEK
* YEAR.
* `task_id` - (Required, String) Task ID.
* `dependency_strategy` - (Optional, List) Dependency execution strategy.
* `offset` - (Optional, String) Offset in interval and list modes.
* `subordinate_cyclic_config` - (Optional, String) Secondary dependency configuration, values:
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
* RECENT_DATE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

wedata wedata_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_task.wedata_task wedata_task_id
```

