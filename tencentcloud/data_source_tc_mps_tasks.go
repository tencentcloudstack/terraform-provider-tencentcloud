package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMpsTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsTasksRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Filter condition: task status, optional values: WAITING, PROCESSING, FINISH.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 10, maximum value: 100.",
			},

			"scroll_token": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Page turning flag, used when pulling in batches: when a single request cannot pull all the data, the interface will return a ScrollToken, and the next request will carry this Token, and it will be obtained from the next record.",
			},

			"task_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Task list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task ID.",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task type, including:WorkflowTask, EditMediaTask, LiveProcessTask.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, in ISO date format. Refer to https://cloud.tencent.com/document/product/862/37710#52.",
						},
						"begin_process_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Begin process time, in ISO date format. Refer to https://cloud.tencent.com/document/product/862/37710#52. If the task has not started yet, this field is: 0000-00-00T00:00:00Z.",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task finish time, in ISO date format. Refer to https://cloud.tencent.com/document/product/862/37710#52. If the task has not been completed, this field is: 0000-00-00T00:00:00Z.",
						},
						"sub_task_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Sub task types.",
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

func dataSourceTencentCloudMpsTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_tasks.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("scroll_token"); ok {
		paramMap["ScrollToken"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var taskSet []*mps.TaskSimpleInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsTasksByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		taskSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(taskSet))
	tmpList := make([]map[string]interface{}, 0, len(taskSet))

	if taskSet != nil {
		for _, taskSimpleInfo := range taskSet {
			taskSimpleInfoMap := map[string]interface{}{}

			if taskSimpleInfo.TaskId != nil {
				taskSimpleInfoMap["task_id"] = taskSimpleInfo.TaskId
			}

			if taskSimpleInfo.TaskType != nil {
				taskSimpleInfoMap["task_type"] = taskSimpleInfo.TaskType
			}

			if taskSimpleInfo.CreateTime != nil {
				taskSimpleInfoMap["create_time"] = taskSimpleInfo.CreateTime
			}

			if taskSimpleInfo.BeginProcessTime != nil {
				taskSimpleInfoMap["begin_process_time"] = taskSimpleInfo.BeginProcessTime
			}

			if taskSimpleInfo.FinishTime != nil {
				taskSimpleInfoMap["finish_time"] = taskSimpleInfo.FinishTime
			}

			if taskSimpleInfo.SubTaskTypes != nil {
				taskSimpleInfoMap["sub_task_types"] = taskSimpleInfo.SubTaskTypes
			}

			ids = append(ids, *taskSimpleInfo.TaskId)
			tmpList = append(tmpList, taskSimpleInfoMap)
		}

		_ = d.Set("task_set", tmpList)
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
