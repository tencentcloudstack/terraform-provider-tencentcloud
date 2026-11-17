package cat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatInstantTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatInstantTasksRead,
		Schema: map[string]*schema.Schema{
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History instant tasks list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task ID.",
						},
						"target_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target address.",
						},
						"task_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task type.",
						},
						"probe_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status.",
						},
						"success_rate": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Success rate.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node count.",
						},
						"task_category": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task category.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of instant tasks.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCatInstantTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_instant_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	catService := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tasks []*cat.SingleInstantTask
	var total *uint64
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, totalRet, e := catService.DescribeCatInstantTasksByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if results == nil && totalRet == nil {
			return resource.NonRetryableError(nil)
		}

		tasks = results
		total = totalRet
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cat_instant_tasks failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(tasks))
	tasksList := make([]map[string]interface{}, 0, len(tasks))
	for _, task := range tasks {
		taskMap := map[string]interface{}{}
		if task.TaskId != nil {
			taskMap["task_id"] = task.TaskId
		}
		if task.TargetAddress != nil {
			taskMap["target_address"] = task.TargetAddress
		}
		if task.TaskType != nil {
			taskMap["task_type"] = task.TaskType
		}
		if task.ProbeTime != nil {
			taskMap["probe_time"] = task.ProbeTime
		}
		if task.Status != nil {
			taskMap["status"] = task.Status
		}
		if task.SuccessRate != nil {
			taskMap["success_rate"] = task.SuccessRate
		}
		if task.NodeCount != nil {
			taskMap["node_count"] = task.NodeCount
		}
		if task.TaskCategory != nil {
			taskMap["task_category"] = task.TaskCategory
		}

		if task.TaskId != nil {
			ids = append(ids, *task.TaskId)
		}
		tasksList = append(tasksList, taskMap)
	}

	if len(tasks) == 0 {
		log.Printf("[DATASOURCE] read empty, skip SetId")
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("tasks", tasksList)
	if total != nil {
		_ = d.Set("total", total)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tasksList); e != nil {
			return e
		}
	}

	return nil
}
