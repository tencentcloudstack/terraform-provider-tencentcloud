/*
Use this data source to query detailed information of as last_activity

Example Usage

```hcl
data "tencentcloud_as_last_activity" "last_activity" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAsLastActivity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsLastActivityRead,
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID list of an auto scaling group.",
			},

			"activity_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information set of eligible scaling activities. Scaling groups without scaling activities are not returned. For example, if there are 50 auto scaling group IDs but only 45 records are returned, it indicates that 5 of the auto scaling groups do not have scaling activities.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling group ID.",
						},
						"activity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling activity ID.",
						},
						"activity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the scaling activity. Value range: SCALE_OUT, SCALE_IN, ATTACH_INSTANCES, REMOVE_INSTANCES, DETACH_INSTANCES, TERMINATE_INSTANCES_UNEXPECTEDLY, REPLACE_UNHEALTHY_INSTANCE, START_INSTANCES, STOP_INSTANCES, INVOKE_COMMAND.",
						},
						"status_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling activity status. Value range: INIT, RUNNING, SUCCESSFUL, PARTIALLY_SUCCESSFUL, FAILED, CANCELLED.",
						},
						"status_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the scaling activity status.",
						},
						"cause": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cause of the scaling activity.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the scaling activity.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of the scaling activity.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of the scaling activity.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the scaling activity.",
						},
						"activity_related_instance_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information set of the instances related to the scaling activity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the instance in the scaling activity. Value range: INIT, RUNNING, SUCCESSFUL, FAILED.",
									},
								},
							},
						},
						"status_message_simplified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Brief description of the scaling activity status.",
						},
						"lifecycle_action_result_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Result of the lifecycle hook action in the scaling activity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lifecycle_hook_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the lifecycle hook.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the instance.",
									},
									"invocation_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution task ID. You can query the result by using the DescribeInvocations API of TAT.",
									},
									"invoke_command_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Result of command invocation, value range: SUCCESSFUL, FAILED, NONE.",
									},
									"notification_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notification result, which indicates whether it is successful to notify CMQ/TDMQ, value range: SUCCESSFUL, FAILED, NONE.",
									},
									"lifecycle_action_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Result of the lifecycle hook action, value range: CONTINUE, ABANDON.",
									},
									"result_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reason of the result, value range: HEARTBEAT_TIMEOUT: Heartbeat timed out. The setting of DefaultResult is used. NOTIFICATION_FAILURE: Failed to send the notification. The setting of DefaultResult is used. CALL_INTERFACE: Calls the CompleteLifecycleAction to set the result ANOTHER_ACTION_ABANDON: It has been set to ABANDON by another operation. COMMAND_CALL_FAILURE: Failed to call the command. The DefaultResult is applied. COMMAND_EXEC_FINISH: Command completed COMMAND_CALL_FAILURE: Failed to execute the command. The DefaultResult is applied. COMMAND_EXEC_RESULT_CHECK_FAILURE: Failed to check the command result. The DefaultResult is applied.",
									},
								},
							},
						},
						"detailed_status_message_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed description of scaling activity status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error type.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AZ information.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance billing mode.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID.",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error message.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type.",
									},
								},
							},
						},
						"invocation_result_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Result of the command execution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"invocation_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution activity ID. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"invocation_task_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution task ID. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"command_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Command ID. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"task_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution Status. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"error_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution exception information. Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAsLastActivityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_last_activity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("auto_scaling_group_ids"); ok {
		autoScalingGroupIdsSet := v.(*schema.Set).List()
		paramMap["AutoScalingGroupIds"] = helper.InterfacesStringsPoint(autoScalingGroupIdsSet)
	}

	service := AsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var activitySet []*as.Activity

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAsLastActivity(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		activitySet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(activitySet))
	tmpList := make([]map[string]interface{}, 0, len(activitySet))

	if activitySet != nil {
		for _, activity := range activitySet {
			activityMap := map[string]interface{}{}

			if activity.AutoScalingGroupId != nil {
				activityMap["auto_scaling_group_id"] = activity.AutoScalingGroupId
			}

			if activity.ActivityId != nil {
				activityMap["activity_id"] = activity.ActivityId
			}

			if activity.ActivityType != nil {
				activityMap["activity_type"] = activity.ActivityType
			}

			if activity.StatusCode != nil {
				activityMap["status_code"] = activity.StatusCode
			}

			if activity.StatusMessage != nil {
				activityMap["status_message"] = activity.StatusMessage
			}

			if activity.Cause != nil {
				activityMap["cause"] = activity.Cause
			}

			if activity.Description != nil {
				activityMap["description"] = activity.Description
			}

			if activity.StartTime != nil {
				activityMap["start_time"] = activity.StartTime
			}

			if activity.EndTime != nil {
				activityMap["end_time"] = activity.EndTime
			}

			if activity.CreatedTime != nil {
				activityMap["created_time"] = activity.CreatedTime
			}

			if activity.ActivityRelatedInstanceSet != nil {
				activityRelatedInstanceSetList := []interface{}{}
				for _, activityRelatedInstanceSet := range activity.ActivityRelatedInstanceSet {
					activityRelatedInstanceSetMap := map[string]interface{}{}

					if activityRelatedInstanceSet.InstanceId != nil {
						activityRelatedInstanceSetMap["instance_id"] = activityRelatedInstanceSet.InstanceId
					}

					if activityRelatedInstanceSet.InstanceStatus != nil {
						activityRelatedInstanceSetMap["instance_status"] = activityRelatedInstanceSet.InstanceStatus
					}

					activityRelatedInstanceSetList = append(activityRelatedInstanceSetList, activityRelatedInstanceSetMap)
				}

				activityMap["activity_related_instance_set"] = activityRelatedInstanceSetList
			}

			if activity.StatusMessageSimplified != nil {
				activityMap["status_message_simplified"] = activity.StatusMessageSimplified
			}

			if activity.LifecycleActionResultSet != nil {
				lifecycleActionResultSetList := []interface{}{}
				for _, lifecycleActionResultSet := range activity.LifecycleActionResultSet {
					lifecycleActionResultSetMap := map[string]interface{}{}

					if lifecycleActionResultSet.LifecycleHookId != nil {
						lifecycleActionResultSetMap["lifecycle_hook_id"] = lifecycleActionResultSet.LifecycleHookId
					}

					if lifecycleActionResultSet.InstanceId != nil {
						lifecycleActionResultSetMap["instance_id"] = lifecycleActionResultSet.InstanceId
					}

					if lifecycleActionResultSet.InvocationId != nil {
						lifecycleActionResultSetMap["invocation_id"] = lifecycleActionResultSet.InvocationId
					}

					if lifecycleActionResultSet.InvokeCommandResult != nil {
						lifecycleActionResultSetMap["invoke_command_result"] = lifecycleActionResultSet.InvokeCommandResult
					}

					if lifecycleActionResultSet.NotificationResult != nil {
						lifecycleActionResultSetMap["notification_result"] = lifecycleActionResultSet.NotificationResult
					}

					if lifecycleActionResultSet.LifecycleActionResult != nil {
						lifecycleActionResultSetMap["lifecycle_action_result"] = lifecycleActionResultSet.LifecycleActionResult
					}

					if lifecycleActionResultSet.ResultReason != nil {
						lifecycleActionResultSetMap["result_reason"] = lifecycleActionResultSet.ResultReason
					}

					lifecycleActionResultSetList = append(lifecycleActionResultSetList, lifecycleActionResultSetMap)
				}

				activityMap["lifecycle_action_result_set"] = lifecycleActionResultSetList
			}

			if activity.DetailedStatusMessageSet != nil {
				detailedStatusMessageSetList := []interface{}{}
				for _, detailedStatusMessageSet := range activity.DetailedStatusMessageSet {
					detailedStatusMessageSetMap := map[string]interface{}{}

					if detailedStatusMessageSet.Code != nil {
						detailedStatusMessageSetMap["code"] = detailedStatusMessageSet.Code
					}

					if detailedStatusMessageSet.Zone != nil {
						detailedStatusMessageSetMap["zone"] = detailedStatusMessageSet.Zone
					}

					if detailedStatusMessageSet.InstanceId != nil {
						detailedStatusMessageSetMap["instance_id"] = detailedStatusMessageSet.InstanceId
					}

					if detailedStatusMessageSet.InstanceChargeType != nil {
						detailedStatusMessageSetMap["instance_charge_type"] = detailedStatusMessageSet.InstanceChargeType
					}

					if detailedStatusMessageSet.SubnetId != nil {
						detailedStatusMessageSetMap["subnet_id"] = detailedStatusMessageSet.SubnetId
					}

					if detailedStatusMessageSet.Message != nil {
						detailedStatusMessageSetMap["message"] = detailedStatusMessageSet.Message
					}

					if detailedStatusMessageSet.InstanceType != nil {
						detailedStatusMessageSetMap["instance_type"] = detailedStatusMessageSet.InstanceType
					}

					detailedStatusMessageSetList = append(detailedStatusMessageSetList, detailedStatusMessageSetMap)
				}

				activityMap["detailed_status_message_set"] = detailedStatusMessageSetList
			}

			if activity.InvocationResultSet != nil {
				invocationResultSetList := []interface{}{}
				for _, invocationResultSet := range activity.InvocationResultSet {
					invocationResultSetMap := map[string]interface{}{}

					if invocationResultSet.InstanceId != nil {
						invocationResultSetMap["instance_id"] = invocationResultSet.InstanceId
					}

					if invocationResultSet.InvocationId != nil {
						invocationResultSetMap["invocation_id"] = invocationResultSet.InvocationId
					}

					if invocationResultSet.InvocationTaskId != nil {
						invocationResultSetMap["invocation_task_id"] = invocationResultSet.InvocationTaskId
					}

					if invocationResultSet.CommandId != nil {
						invocationResultSetMap["command_id"] = invocationResultSet.CommandId
					}

					if invocationResultSet.TaskStatus != nil {
						invocationResultSetMap["task_status"] = invocationResultSet.TaskStatus
					}

					if invocationResultSet.ErrorMessage != nil {
						invocationResultSetMap["error_message"] = invocationResultSet.ErrorMessage
					}

					invocationResultSetList = append(invocationResultSetList, invocationResultSetMap)
				}

				activityMap["invocation_result_set"] = invocationResultSetList
			}

			ids = append(ids, *activity.AutoScalingGroupId)
			tmpList = append(tmpList, activityMap)
		}

		_ = d.Set("activity_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
