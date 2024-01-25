// Code generated by iacg; DO NOT EDIT.
package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsShipperTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsShipperTasksRead,
		Schema: map[string]*schema.Schema{
			"shipper_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "shipper id.",
			},

			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "start time(ms).",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "end time(ms).",
			},

			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task id.",
						},
						"shipper_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "shipper id.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "topic id.",
						},
						"range_start": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "start time of current task (ms).",
						},
						"range_end": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "end time of current task (ms).",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "start time(ms).",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "end time(ms).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status of current shipper task.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "detail info.",
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

func dataSourceTencentCloudClsShipperTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_shipper_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("shipper_id"); ok {
		paramMap["ShipperId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	var respData []*cls.ShipperTaskInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsShipperTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(respData))
	tasksList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, tasks := range respData {
			tasksMap := map[string]interface{}{}

			if tasks.TaskId != nil {
				tasksMap["task_id"] = tasks.TaskId
			}

			if tasks.ShipperId != nil {
				tasksMap["shipper_id"] = tasks.ShipperId
			}
			shipperId := *tasks.ShipperId

			if tasks.TopicId != nil {
				tasksMap["topic_id"] = tasks.TopicId
			}

			if tasks.RangeStart != nil {
				tasksMap["range_start"] = tasks.RangeStart
			}

			if tasks.RangeEnd != nil {
				tasksMap["range_end"] = tasks.RangeEnd
			}

			if tasks.StartTime != nil {
				tasksMap["start_time"] = tasks.StartTime
			}

			if tasks.EndTime != nil {
				tasksMap["end_time"] = tasks.EndTime
			}

			if tasks.Status != nil {
				tasksMap["status"] = tasks.Status
			}

			if tasks.Message != nil {
				tasksMap["message"] = tasks.Message
			}

			ids = append(ids, shipperId)
			tasksList = append(tasksList, tasksMap)
		}

		_ = d.Set("tasks", tasksList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tasksList); e != nil {
			return e
		}
	}

	return nil
}
