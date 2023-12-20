package cls

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsShipperTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsShipperTasksRead,
		Schema: map[string]*schema.Schema{
			"shipper_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "shipper id.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "start time(ms).",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "end time(ms).",
			},

			"tasks": {
				Computed:    true,
				Type:        schema.TypeList,
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

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("shipper_id"); ok {
		paramMap["ShipperId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tasks []*cls.ShipperTaskInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsShipperTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		tasks = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(tasks))
	tmpList := make([]map[string]interface{}, 0, len(tasks))

	if tasks != nil {
		for _, shipperTaskInfo := range tasks {
			shipperTaskInfoMap := map[string]interface{}{}

			if shipperTaskInfo.TaskId != nil {
				shipperTaskInfoMap["task_id"] = shipperTaskInfo.TaskId
			}

			if shipperTaskInfo.ShipperId != nil {
				shipperTaskInfoMap["shipper_id"] = shipperTaskInfo.ShipperId
			}

			if shipperTaskInfo.TopicId != nil {
				shipperTaskInfoMap["topic_id"] = shipperTaskInfo.TopicId
			}

			if shipperTaskInfo.RangeStart != nil {
				shipperTaskInfoMap["range_start"] = shipperTaskInfo.RangeStart
			}

			if shipperTaskInfo.RangeEnd != nil {
				shipperTaskInfoMap["range_end"] = shipperTaskInfo.RangeEnd
			}

			if shipperTaskInfo.StartTime != nil {
				shipperTaskInfoMap["start_time"] = shipperTaskInfo.StartTime
			}

			if shipperTaskInfo.EndTime != nil {
				shipperTaskInfoMap["end_time"] = shipperTaskInfo.EndTime
			}

			if shipperTaskInfo.Status != nil {
				shipperTaskInfoMap["status"] = shipperTaskInfo.Status
			}

			if shipperTaskInfo.Message != nil {
				shipperTaskInfoMap["message"] = shipperTaskInfo.Message
			}

			ids = append(ids, *shipperTaskInfo.ShipperId)
			tmpList = append(tmpList, shipperTaskInfoMap)
		}

		_ = d.Set("tasks", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
