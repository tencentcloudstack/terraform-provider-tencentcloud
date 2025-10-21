---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invocation_task"
sidebar_current: "docs-tencentcloud-datasource-tat_invocation_task"
description: |-
  Use this data source to query detailed information of tat invocation_task
---

# tencentcloud_tat_invocation_task

Use this data source to query detailed information of tat invocation_task

## Example Usage

```hcl
data "tencentcloud_tat_invocation_task" "invocation_task" {
  # invocation_task_ids = ["invt-a8bv0ip7"]
  filters {
    name   = "instance-id"
    values = ["ins-p4pq4gaq"]
  }
  hide_output = true
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions.invocation-id - String - Required: No - (Filter condition) Filter by the execution activity ID.invocation-task-id - String - Required: No - (Filter condition) Filter by the execution task ID.instance-id - String - Required: No - (Filter condition) Filter by the instance ID.command-id - String - Required: No - (Filter condition) Filter by the command ID.Up to 10 Filters are allowed for each request. Each filter can have up to five Filter.Values. InvocationTaskIds and Filters cannot be specified at the same time.
* `hide_output` - (Optional, Bool) Whether to hide the output. Valid values:True (default): Hide the outputFalse: Show the output.
* `invocation_task_ids` - (Optional, Set: [`String`]) List of execution task IDs. Up to 100 IDs are allowed for each request. InvocationTaskIds and Filters cannot be specified at the same time.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter values of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `invocation_task_set` - List of execution tasks.
  * `command_document` - Command details of the execution task.
    * `command_type` - Command type.
    * `content` - Base64-encoded command.
    * `output_cos_bucket_url` - URL of the COS bucket to store the output.
    * `output_cos_key_prefix` - Prefix of the output file name.
    * `timeout` - Timeout period.
    * `username` - The user who executes the command.
    * `working_directory` - Execution path.
  * `command_id` - Command ID.
  * `created_time` - Creation time.
  * `end_time` - End time of the execution task.
  * `error_info` - Error message displayed when the execution task fails.
  * `instance_id` - Instance ID.
  * `invocation_id` - Execution activity ID.
  * `invocation_source` - Invocation source.
  * `invocation_task_id` - Execution task ID.
  * `start_time` - Start time of the execution task.
  * `task_result` - Execution result.
    * `dropped` - Dropped bytes of the command output.
    * `exec_end_time` - Time when the execution is ended.
    * `exec_start_time` - Time when the execution is started.
    * `exit_code` - ExitCode of the execution.
    * `output_upload_cos_error_info` - Error message for uploading logs to COS.
    * `output_url` - COS URL of the logs.
    * `output` - Base64-encoded command output. The maximum length is 24 KB.
  * `task_status` - Execution task status. Valid values:PENDING: PendingDELIVERING: DeliveringDELIVER_DELAYED: Delivery delayedDELIVER_FAILED: Delivery failedSTART_FAILED: Failed to start the commandRUNNING: RunningSUCCESS: SuccessFAILED: Failed to execute the command. The exit code is not 0 after execution.TIMEOUT: Command timed outTASK_TIMEOUT: Task timed outCANCELLING: CancelingCANCELLED: Canceled (canceled before execution)TERMINATED: Terminated (canceled during execution).
  * `updated_time` - Update time.


