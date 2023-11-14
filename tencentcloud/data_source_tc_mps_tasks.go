/*
Use this data source to query detailed information of mps tasks

Example Usage

```hcl
data "tencentcloud_mps_tasks" "tasks" {
  status = ""
    }
```
*/
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
				Description: "Filter: Task status. Valid values: WAITING (waiting), PROCESSING (processing), FINISH (completed).",
			},

			"scroll_token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scrolling identifier which is used for pulling in batches. If a single request cannot pull all the data entries, the API will return `ScrollToken`, and if the next request carries it, the next pull will start from the next entry.",
			},

			"task_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Task overview list.",
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
							Description: "Task type. Valid values:&amp;lt;li&amp;gt; WorkflowTask: Workflow processing task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt; LiveProcessTask: Live stream processing task.&amp;lt;/li&amp;gt;.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of a task in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F).",
						},
						"begin_process_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of task execution in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F). If the task has not been started yet, this field will be `0000-00-00T00:00:00Z`.",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of a task in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F). If the task has not been completed yet, this field will be `0000-00-00T00:00:00Z`.",
						},
						"sub_task_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The subtask type.",
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

	if scrollToken != nil {
		_ = d.Set("scroll_token", scrollToken)
	}

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
