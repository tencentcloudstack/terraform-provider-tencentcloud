/*
Use this data source to query detailed information of tcr tag_retention_execution_task

Example Usage

```hcl
data "tencentcloud_tcr_tag_retention_execution_task" "tag_retention_execution_task" {
  registry_id = "tcr-xxx"
  retention_id = 1
  execution_id = 1
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrTagRetentionExecutionTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrTagRetentionExecutionTaskRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"retention_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Retention id.",
			},

			"execution_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Execution id.",
			},

			"retention_task_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of version retention tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task id.",
						},
						"execution_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The rule execution id.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task end time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The execution status of the task: Failed, Succeed, Stopped, InProgress.",
						},
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of tags.",
						},
						"retained": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of retained tags.",
						},
						"repository": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository name.",
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

func dataSourceTencentCloudTcrTagRetentionExecutionTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_tag_retention_execution_task.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("retention_id"); v != nil {
		paramMap["RetentionId"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("execution_id"); v != nil {
		paramMap["ExecutionId"] = helper.IntInt64(v.(int))
	}

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	var retentionTaskList []*tcr.RetentionTask

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrTagRetentionExecutionTaskByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		retentionTaskList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(retentionTaskList))
	tmpList := make([]map[string]interface{}, 0, len(retentionTaskList))

	if retentionTaskList != nil {
		for _, retentionTask := range retentionTaskList {
			retentionTaskMap := map[string]interface{}{}

			if retentionTask.TaskId != nil {
				retentionTaskMap["task_id"] = retentionTask.TaskId
			}

			if retentionTask.ExecutionId != nil {
				retentionTaskMap["execution_id"] = retentionTask.ExecutionId
			}

			if retentionTask.StartTime != nil {
				retentionTaskMap["start_time"] = retentionTask.StartTime
			}

			if retentionTask.EndTime != nil {
				retentionTaskMap["end_time"] = retentionTask.EndTime
			}

			if retentionTask.Status != nil {
				retentionTaskMap["status"] = retentionTask.Status
			}

			if retentionTask.Total != nil {
				retentionTaskMap["total"] = retentionTask.Total
			}

			if retentionTask.Retained != nil {
				retentionTaskMap["retained"] = retentionTask.Retained
			}

			if retentionTask.Repository != nil {
				retentionTaskMap["repository"] = retentionTask.Repository
			}

			ids = append(ids, *retentionTask.RegistryId)
			tmpList = append(tmpList, retentionTaskMap)
		}

		_ = d.Set("retention_task_list", tmpList)
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
