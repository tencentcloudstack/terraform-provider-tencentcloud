---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_task"
sidebar_current: "docs-tencentcloud-resource-wedata_trigger_task"
description: |-
  Provides a resource to create a wedata trigger task
---

# tencentcloud_wedata_trigger_task

Provides a resource to create a wedata trigger task

## Example Usage

```hcl
resource "tencentcloud_wedata_trigger_task" "trigger_task" {
  project_id = jsonencode(3108707295180644352)
  trigger_task_base_attribute {
    owner_uin        = jsonencode(100044349576)
    task_folder_path = "/"
    task_name        = "tf-test-task"
    task_type_id     = jsonencode(35)
    workflow_id      = tencentcloud_wedata_trigger_workflow.workflow.id
  }
  trigger_task_configuration {
    broker_ip      = "any"
    code_content   = base64encode("echo Hello, World")
    resource_group = jsonencode(20241107171437783498)
    task_ext_configuration_list {
      param_key   = "enableKerberosLogin"
      param_value = true
    }
    task_ext_configuration_list {
      param_key   = "executionTTLStrategy"
      param_value = "fail"
    }
    task_ext_configuration_list {
      param_key   = "python_sub_version"
      param_value = "python3"
    }
    task_ext_configuration_list {
      param_key   = "python_type"
      param_value = "python3"
    }
    task_ext_configuration_list {
      param_key   = "specLabelConfItems"
      param_value = "eyJzcGVjTGxxxxxxxfQ=="
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTL"
      param_value = jsonencode(-1)
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTLStrategy"
      param_value = "fail"
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTLStrategy"
      param_value = "fail"
    }
  }
  trigger_task_scheduler_configuration {
    allow_redo_type                 = "ALL"
    execution_ttl_minute            = -1
    max_retry_number                = 4
    retry_wait_minute               = 5
    run_priority_type               = 6
    wait_execution_total_ttl_minute = -1
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `trigger_task_base_attribute` - (Required, List) Basic task attributes.
* `trigger_task_configuration` - (Required, List) Task configuration.
* `trigger_task_scheduler_configuration` - (Required, List) Task scheduling configuration.

The `param_task_in_list` object of `trigger_task_scheduler_configuration` supports the following:

* `from_param_key` - (Required, String) Parent task parameter key.
* `from_task_id` - (Required, String) Parent task ID.
* `param_desc` - (Required, String) Parameter description. Format: project_identifier.task_name.parameter_name; e.g., project_wedata_1.sh_250820_104107.pp_out.
* `param_key` - (Required, String) Parameter name.

The `param_task_out_list` object of `trigger_task_scheduler_configuration` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter definition.

The `task_ext_configuration_list` object of `trigger_task_configuration` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

The `task_ext_configuration_system_list` object of `trigger_task_configuration` supports the following:


The `task_output_registry_list` object of `trigger_task_scheduler_configuration` supports the following:

* `data_flow_type` - (Required, String) Input/output table type: input stream: `UPSTREAM`, output stream: `DOWNSTREAM`.
* `database_name` - (Required, String) Database name.
* `datasource_id` - (Required, String) Data source ID.
* `partition_name` - (Required, String) Partition name.
* `table_name` - (Required, String) Table name.
* `table_physical_id` - (Required, String) Table physical unique ID.
* `db_guid` - (Optional, String) Database unique identifier.
* `table_guid` - (Optional, String) Table unique identifier.

The `task_scheduling_parameter_list` object of `trigger_task_configuration` supports the following:

* `param_key` - (Required, String) Parameter name.
* `param_value` - (Required, String) Parameter value.

The `trigger_task_base_attribute` object supports the following:

* `task_name` - (Required, String) Task name.
* `task_type_id` - (Required, String) Task type ID: `26`: OfflineSynchronization; `30`: Python; `32`: DLC SQL; `35`: Shell; `38`: Shell Form Mode; `46`: DLC Spark; `50`: DLC PySpark; `130`: Branch Node; `131`: Merged Node; `132`: Notebook; `133`: SSH; `137`: For-each; `139`: DLC Spark Streaming; `140`: Run Workflow.
* `workflow_id` - (Required, String) Workflow ID.
* `owner_uin` - (Optional, String) Task owner ID, defaults to the current user.
* `task_description` - (Optional, String) Task description.
* `task_folder_path` - (Optional, String) Task folder path. Do not include the task node type in the path. For example, in a workflow named wf01 under the "General" category, to create a shell task in the tf_01 folder under this category, set the value to /tf_01. If the tf_01 folder does not exist, it must be created first (using the CreateTaskFolder API) before the operation can succeed.

The `trigger_task_configuration` object supports the following:

* `broker_ip` - (Optional, String) Specified execution node.
* `bundle_id` - (Optional, String) Bundle ID in use.
* `bundle_info` - (Optional, String) Bundle information.
* `code_content` - (Optional, String) Base64-encoded code content.
* `data_cluster` - (Optional, String) Cluster ID.
* `resource_group` - (Optional, String) Resource group ID. Obtain ExecutorGroupId via DescribeNormalSchedulerExecutorGroups.
* `source_service_id` - (Optional, String) Source data source IDs, separated by semicolons (;). Obtain via DescribeDataSourceWithoutInfo.
* `source_service_name` - (Optional, String) The source data source name needs to be obtained through DescribeDataSourceWithoutInfo..
* `source_service_type` - (Optional, String) The source data source type needs to be obtained through DescribeDataSourceWithoutInfo.
* `target_service_id` - (Optional, String) Target data source IDs, separated by semicolons (;). Obtain via DescribeDataSourceWithoutInfo.
* `target_service_name` - (Optional, String) The target data source name, which needs to be obtained through DescribeDataSourceWithoutInfo.
* `target_service_type` - (Optional, String) The target data source type needs to be obtained through DescribeDataSourceWithoutInfo.
* `task_ext_configuration_list` - (Optional, Set) Task extended attribute configuration list. `sql.file.name`, `ftp.file.name`, `bucket`, `region`, `tenantId` cannot be customized; they are generated by the system.
* `task_scheduling_parameter_list` - (Optional, List) Scheduling parameters.
* `yarn_queue` - (Optional, String) Resource pool queue name. Obtain via DescribeProjectClusterQueues.

The `trigger_task_scheduler_configuration` object supports the following:

* `allow_redo_type` - (Optional, String) Rerun & backfill configuration. Default: ALL. ALL: rerun or backfill allowed after success or failure; FAILURE: not allowed after success, allowed after failure; NONE: not allowed after success or failure.
* `execution_ttl_minute` - (Optional, Int) Timeout handling policy. Execution timeout in minutes. Default: -1.
* `max_retry_number` - (Optional, Int) Retry policy. Maximum retry attempts. Default: 4.
* `param_task_in_list` - (Optional, List) Input parameter list.
* `param_task_out_list` - (Optional, List) Output parameter list.
* `retry_wait_minute` - (Optional, Int) Retry policy. Retry wait time in minutes. Default: 5.
* `run_priority_type` - (Optional, Int) Task scheduling priority. Run priority: `4`-High; `5`-Medium; `6`-Low. Default: 6.
* `task_output_registry_list` - (Optional, List) Output registry.
* `upstream_dependency_config_list` - (Optional, List) List of upstream dependent tasks.
* `wait_execution_total_ttl_minute` - (Optional, Int) Timeout handling policy. Total wait timeout in minutes. Default: -1.

The `upstream_dependency_config_list` object of `trigger_task_scheduler_configuration` supports the following:

* `task_id` - (Required, String) Task ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

wedata trigger_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_trigger_task.trigger_task project_id#task_id
```

