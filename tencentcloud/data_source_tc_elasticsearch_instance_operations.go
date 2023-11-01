/*
Use this data source to query detailed information of elasticsearch instance operations

Example Usage

```hcl
data "tencentcloud_elasticsearch_instance_operations" "instance_operations" {
	instance_id = "es-xxxxxx"
	start_time = "2018-01-01 00:00:00"
	end_time = "2023-10-31 10:12:45"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudElasticsearchInstanceOperations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchInstanceOperationsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time, e.g. 2019-03-07 16:30:39.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time, e.g. 2019-03-30 20:18:03.",
			},

			"operations": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Operation records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Id.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type.",
						},
						"detail": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Operation details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"old_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Instance original configuration information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Key.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value.",
												},
											},
										},
									},
									"new_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Configuration information after instance update.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Key.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value.",
												},
											},
										},
									},
								},
							},
						},
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation result.",
						},
						"tasks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Task information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task name.",
									},
									"progress": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Task progress.",
									},
									"finish_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task completion time.",
									},
									"sub_tasks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Subtask.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subtask name.",
												},
												"result": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Subtask result.",
												},
												"err_msg": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subtask error message.",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subtask type.",
												},
												"status": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Subtask status, 1: success; 0: processing; -1: failure.",
												},
												"failed_indices": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The index name of the failed upgrade check.",
												},
												"finish_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subtask end time.",
												},
												"level": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Subtask level, 1: warning; 2: failed.",
												},
											},
										},
									},
									"elapsed_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Elapsed time.",
									},
									"process_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Progress info.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"completed": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Completed quantity.",
												},
												"remain": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Remaining quantity.",
												},
												"total": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total quantity.",
												},
												"task_type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Task type. 60: restart task 70: fragment migration task 80: node modification task.",
												},
											},
										},
									},
								},
							},
						},
						"progress": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Operation progress.",
						},
						"sub_account_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operator uin.",
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

func dataSourceTencentCloudElasticsearchInstanceOperationsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_elasticsearch_instance_operations.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(instanceId)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var operations []*elasticsearch.Operation

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchInstanceOperationsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		operations = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(operations))

	if operations != nil {
		for _, operation := range operations {
			operationMap := map[string]interface{}{}

			if operation.Id != nil {
				operationMap["id"] = operation.Id
			}

			if operation.StartTime != nil {
				operationMap["start_time"] = operation.StartTime
			}

			if operation.Type != nil {
				operationMap["type"] = operation.Type
			}

			if operation.Detail != nil {
				detailMap := map[string]interface{}{}

				if operation.Detail.OldInfo != nil {
					oldInfoList := []interface{}{}
					for _, oldInfo := range operation.Detail.OldInfo {
						oldInfoMap := map[string]interface{}{}

						if oldInfo.Key != nil {
							oldInfoMap["key"] = oldInfo.Key
						}

						if oldInfo.Value != nil {
							oldInfoMap["value"] = oldInfo.Value
						}

						oldInfoList = append(oldInfoList, oldInfoMap)
					}

					detailMap["old_info"] = oldInfoList
				}

				if operation.Detail.NewInfo != nil {
					newInfoList := []interface{}{}
					for _, newInfo := range operation.Detail.NewInfo {
						newInfoMap := map[string]interface{}{}

						if newInfo.Key != nil {
							newInfoMap["key"] = newInfo.Key
						}

						if newInfo.Value != nil {
							newInfoMap["value"] = newInfo.Value
						}

						newInfoList = append(newInfoList, newInfoMap)
					}

					detailMap["new_info"] = newInfoList
				}

				operationMap["detail"] = []interface{}{detailMap}
			}

			if operation.Result != nil {
				operationMap["result"] = operation.Result
			}

			if operation.Tasks != nil {
				tasksList := []interface{}{}
				for _, tasks := range operation.Tasks {
					tasksMap := map[string]interface{}{}

					if tasks.Name != nil {
						tasksMap["name"] = tasks.Name
					}

					if tasks.Progress != nil {
						tasksMap["progress"] = tasks.Progress
					}

					if tasks.FinishTime != nil {
						tasksMap["finish_time"] = tasks.FinishTime
					}

					if tasks.SubTasks != nil {
						subTasksList := []interface{}{}
						for _, subTasks := range tasks.SubTasks {
							subTasksMap := map[string]interface{}{}

							if subTasks.Name != nil {
								subTasksMap["name"] = subTasks.Name
							}

							if subTasks.Result != nil {
								subTasksMap["result"] = subTasks.Result
							}

							if subTasks.ErrMsg != nil {
								subTasksMap["err_msg"] = subTasks.ErrMsg
							}

							if subTasks.Type != nil {
								subTasksMap["type"] = subTasks.Type
							}

							if subTasks.Status != nil {
								subTasksMap["status"] = subTasks.Status
							}

							if subTasks.FailedIndices != nil {
								subTasksMap["failed_indices"] = subTasks.FailedIndices
							}

							if subTasks.FinishTime != nil {
								subTasksMap["finish_time"] = subTasks.FinishTime
							}

							if subTasks.Level != nil {
								subTasksMap["level"] = subTasks.Level
							}

							subTasksList = append(subTasksList, subTasksMap)
						}

						tasksMap["sub_tasks"] = subTasksList
					}

					if tasks.ElapsedTime != nil {
						tasksMap["elapsed_time"] = tasks.ElapsedTime
					}

					if tasks.ProcessInfo != nil {
						processInfoMap := map[string]interface{}{}

						if tasks.ProcessInfo.Completed != nil {
							processInfoMap["completed"] = tasks.ProcessInfo.Completed
						}

						if tasks.ProcessInfo.Remain != nil {
							processInfoMap["remain"] = tasks.ProcessInfo.Remain
						}

						if tasks.ProcessInfo.Total != nil {
							processInfoMap["total"] = tasks.ProcessInfo.Total
						}

						if tasks.ProcessInfo.TaskType != nil {
							processInfoMap["task_type"] = tasks.ProcessInfo.TaskType
						}

						tasksMap["process_info"] = []interface{}{processInfoMap}
					}

					tasksList = append(tasksList, tasksMap)
				}

				operationMap["tasks"] = tasksList
			}

			if operation.Progress != nil {
				operationMap["progress"] = operation.Progress
			}

			if operation.SubAccountUin != nil {
				operationMap["sub_account_uin"] = operation.SubAccountUin
			}

			tmpList = append(tmpList, operationMap)
		}

		_ = d.Set("operations", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
