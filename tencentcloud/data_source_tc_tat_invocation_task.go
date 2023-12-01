/*
Use this data source to query detailed information of tat invocation_task

Example Usage

```hcl
data "tencentcloud_tat_invocation_task" "invocation_task" {
  # invocation_task_ids = ["invt-a8bv0ip7"]
  filters {
    name = "instance-id"
    values = ["ins-p4pq4gaq"]
  }
  hide_output = true
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTatInvocationTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatInvocationTaskRead,
		Schema: map[string]*schema.Schema{
			"invocation_task_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of execution task IDs. Up to 100 IDs are allowed for each request. InvocationTaskIds and Filters cannot be specified at the same time.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions.invocation-id - String - Required: No - (Filter condition) Filter by the execution activity ID.invocation-task-id - String - Required: No - (Filter condition) Filter by the execution task ID.instance-id - String - Required: No - (Filter condition) Filter by the instance ID.command-id - String - Required: No - (Filter condition) Filter by the command ID.Up to 10 Filters are allowed for each request. Each filter can have up to five Filter.Values. InvocationTaskIds and Filters cannot be specified at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values of the field.",
						},
					},
				},
			},

			"hide_output": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to hide the output. Valid values:True (default): Hide the outputFalse: Show the output.",
			},

			"invocation_task_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of execution tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invocation_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution activity ID.",
						},
						"invocation_task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution task ID.",
						},
						"command_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command ID.",
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution task status. Valid values:PENDING: PendingDELIVERING: DeliveringDELIVER_DELAYED: Delivery delayedDELIVER_FAILED: Delivery failedSTART_FAILED: Failed to start the commandRUNNING: RunningSUCCESS: SuccessFAILED: Failed to execute the command. The exit code is not 0 after execution.TIMEOUT: Command timed outTASK_TIMEOUT: Task timed outCANCELLING: CancelingCANCELLED: Canceled (canceled before execution)TERMINATED: Terminated (canceled during execution).",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"task_result": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Execution result.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exit_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ExitCode of the execution.",
									},
									"output": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64-encoded command output. The maximum length is 24 KB.",
									},
									"exec_start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time when the execution is started.",
									},
									"exec_end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time when the execution is ended.",
									},
									"dropped": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dropped bytes of the command output.",
									},
									"output_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS URL of the logs.",
									},
									"output_upload_cos_error_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error message for uploading logs to COS.",
									},
								},
							},
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of the execution task.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of the execution task.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"command_document": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Command details of the execution task.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64-encoded command.",
									},
									"command_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Command type.",
									},
									"timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Timeout period.",
									},
									"working_directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution path.",
									},
									"username": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user who executes the command.",
									},
									"output_cos_bucket_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of the COS bucket to store the output.",
									},
									"output_cos_key_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix of the output file name.",
									},
								},
							},
						},
						"error_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error message displayed when the execution task fails.",
						},
						"invocation_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invocation source.",
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

func dataSourceTencentCloudTatInvocationTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tat_invocation_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("invocation_task_ids"); ok {
		invocationTaskIdsSet := v.(*schema.Set).List()
		paramMap["InvocationTaskIds"] = helper.InterfacesStringsPoint(invocationTaskIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tat.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tat.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, _ := d.GetOk("hide_output"); v != nil {
		paramMap["HideOutput"] = helper.Bool(v.(bool))
	}

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	var invocationTaskSet []*tat.InvocationTask

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTatInvocationTaskByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		invocationTaskSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(invocationTaskSet))
	tmpList := make([]map[string]interface{}, 0, len(invocationTaskSet))

	if invocationTaskSet != nil {
		for _, invocationTask := range invocationTaskSet {
			invocationTaskMap := map[string]interface{}{}

			if invocationTask.InvocationId != nil {
				invocationTaskMap["invocation_id"] = invocationTask.InvocationId
			}

			if invocationTask.InvocationTaskId != nil {
				invocationTaskMap["invocation_task_id"] = invocationTask.InvocationTaskId
			}

			if invocationTask.CommandId != nil {
				invocationTaskMap["command_id"] = invocationTask.CommandId
			}

			if invocationTask.TaskStatus != nil {
				invocationTaskMap["task_status"] = invocationTask.TaskStatus
			}

			if invocationTask.InstanceId != nil {
				invocationTaskMap["instance_id"] = invocationTask.InstanceId
			}

			if invocationTask.TaskResult != nil {
				taskResultMap := map[string]interface{}{}

				if invocationTask.TaskResult.ExitCode != nil {
					taskResultMap["exit_code"] = invocationTask.TaskResult.ExitCode
				}

				if invocationTask.TaskResult.Output != nil {
					taskResultMap["output"] = invocationTask.TaskResult.Output
				}

				if invocationTask.TaskResult.ExecStartTime != nil {
					taskResultMap["exec_start_time"] = invocationTask.TaskResult.ExecStartTime
				}

				if invocationTask.TaskResult.ExecEndTime != nil {
					taskResultMap["exec_end_time"] = invocationTask.TaskResult.ExecEndTime
				}

				if invocationTask.TaskResult.Dropped != nil {
					taskResultMap["dropped"] = invocationTask.TaskResult.Dropped
				}

				if invocationTask.TaskResult.OutputUrl != nil {
					taskResultMap["output_url"] = invocationTask.TaskResult.OutputUrl
				}

				if invocationTask.TaskResult.OutputUploadCOSErrorInfo != nil {
					taskResultMap["output_upload_cos_error_info"] = invocationTask.TaskResult.OutputUploadCOSErrorInfo
				}

				invocationTaskMap["task_result"] = []interface{}{taskResultMap}
			}

			if invocationTask.StartTime != nil {
				invocationTaskMap["start_time"] = invocationTask.StartTime
			}

			if invocationTask.EndTime != nil {
				invocationTaskMap["end_time"] = invocationTask.EndTime
			}

			if invocationTask.CreatedTime != nil {
				invocationTaskMap["created_time"] = invocationTask.CreatedTime
			}

			if invocationTask.UpdatedTime != nil {
				invocationTaskMap["updated_time"] = invocationTask.UpdatedTime
			}

			if invocationTask.CommandDocument != nil {
				commandDocumentMap := map[string]interface{}{}

				if invocationTask.CommandDocument.Content != nil {
					commandDocumentMap["content"] = invocationTask.CommandDocument.Content
				}

				if invocationTask.CommandDocument.CommandType != nil {
					commandDocumentMap["command_type"] = invocationTask.CommandDocument.CommandType
				}

				if invocationTask.CommandDocument.Timeout != nil {
					commandDocumentMap["timeout"] = invocationTask.CommandDocument.Timeout
				}

				if invocationTask.CommandDocument.WorkingDirectory != nil {
					commandDocumentMap["working_directory"] = invocationTask.CommandDocument.WorkingDirectory
				}

				if invocationTask.CommandDocument.Username != nil {
					commandDocumentMap["username"] = invocationTask.CommandDocument.Username
				}

				if invocationTask.CommandDocument.OutputCOSBucketUrl != nil {
					commandDocumentMap["output_cos_bucket_url"] = invocationTask.CommandDocument.OutputCOSBucketUrl
				}

				if invocationTask.CommandDocument.OutputCOSKeyPrefix != nil {
					commandDocumentMap["output_cos_key_prefix"] = invocationTask.CommandDocument.OutputCOSKeyPrefix
				}

				invocationTaskMap["command_document"] = []interface{}{commandDocumentMap}
			}

			if invocationTask.ErrorInfo != nil {
				invocationTaskMap["error_info"] = invocationTask.ErrorInfo
			}

			if invocationTask.InvocationSource != nil {
				invocationTaskMap["invocation_source"] = invocationTask.InvocationSource
			}

			ids = append(ids, *invocationTask.InvocationTaskId)
			tmpList = append(tmpList, invocationTaskMap)
		}

		_ = d.Set("invocation_task_set", tmpList)
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
