---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_last_activity"
sidebar_current: "docs-tencentcloud-datasource-as_last_activity"
description: |-
  Use this data source to query detailed information of as last_activity
---

# tencentcloud_as_last_activity

Use this data source to query detailed information of as last_activity

## Example Usage

```hcl
data "tencentcloud_as_last_activity" "last_activity" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_ids` - (Required, Set: [`String`]) ID list of an auto scaling group.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `activity_set` - Information set of eligible scaling activities. Scaling groups without scaling activities are not returned. For example, if there are 50 auto scaling group IDs but only 45 records are returned, it indicates that 5 of the auto scaling groups do not have scaling activities.
  * `activity_id` - Scaling activity ID.
  * `activity_related_instance_set` - Information set of the instances related to the scaling activity.
    * `instance_id` - Instance ID.
    * `instance_status` - Status of the instance in the scaling activity. Value range: INIT, RUNNING, SUCCESSFUL, FAILED.
  * `activity_type` - Type of the scaling activity. Value range: SCALE_OUT, SCALE_IN, ATTACH_INSTANCES, REMOVE_INSTANCES, DETACH_INSTANCES, TERMINATE_INSTANCES_UNEXPECTEDLY, REPLACE_UNHEALTHY_INSTANCE, START_INSTANCES, STOP_INSTANCES, INVOKE_COMMAND.
  * `auto_scaling_group_id` - Auto scaling group ID.
  * `cause` - Cause of the scaling activity.
  * `created_time` - Creation time of the scaling activity.
  * `description` - Description of the scaling activity.
  * `detailed_status_message_set` - Detailed description of scaling activity status.
    * `code` - Error type.
    * `instance_charge_type` - Instance billing mode.
    * `instance_id` - Instance ID.
    * `instance_type` - Instance type.
    * `message` - Error message.
    * `subnet_id` - Subnet ID.
    * `zone` - AZ information.
  * `end_time` - End time of the scaling activity.
  * `invocation_result_set` - Result of the command execution.
    * `command_id` - Command ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `error_message` - Execution exception information. Note: This field may return null, indicating that no valid values can be obtained.
    * `instance_id` - Instance ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `invocation_id` - Execution activity ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `invocation_task_id` - Execution task ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `task_status` - Execution Status. Note: This field may return null, indicating that no valid values can be obtained.
  * `lifecycle_action_result_set` - Result of the lifecycle hook action in the scaling activity.
    * `instance_id` - ID of the instance.
    * `invocation_id` - Execution task ID. You can query the result by using the DescribeInvocations API of TAT.
    * `invoke_command_result` - Result of command invocation, value range: SUCCESSFUL, FAILED, NONE.
    * `lifecycle_action_result` - Result of the lifecycle hook action, value range: CONTINUE, ABANDON.
    * `lifecycle_hook_id` - ID of the lifecycle hook.
    * `notification_result` - Notification result, which indicates whether it is successful to notify CMQ/TDMQ, value range: SUCCESSFUL, FAILED, NONE.
    * `result_reason` - Reason of the result, value range: HEARTBEAT_TIMEOUT: Heartbeat timed out. The setting of DefaultResult is used. NOTIFICATION_FAILURE: Failed to send the notification. The setting of DefaultResult is used. CALL_INTERFACE: Calls the CompleteLifecycleAction to set the result ANOTHER_ACTION_ABANDON: It has been set to ABANDON by another operation. COMMAND_CALL_FAILURE: Failed to call the command. The DefaultResult is applied. COMMAND_EXEC_FINISH: Command completed COMMAND_CALL_FAILURE: Failed to execute the command. The DefaultResult is applied. COMMAND_EXEC_RESULT_CHECK_FAILURE: Failed to check the command result. The DefaultResult is applied.
  * `start_time` - Start time of the scaling activity.
  * `status_code` - Scaling activity status. Value range: INIT, RUNNING, SUCCESSFUL, PARTIALLY_SUCCESSFUL, FAILED, CANCELLED.
  * `status_message_simplified` - Brief description of the scaling activity status.
  * `status_message` - Description of the scaling activity status.


