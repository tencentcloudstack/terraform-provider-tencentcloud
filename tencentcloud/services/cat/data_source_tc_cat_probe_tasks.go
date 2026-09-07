package cat

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatProbeTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatProbeTasksRead,
		Schema: map[string]*schema.Schema{
			"task_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Task ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Task name.",
			},
			"target_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target address.",
			},
			"task_status": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Task status list.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"pay_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pay mode.",
			},
			"order_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Order state.",
			},
			"task_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Task type list.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"task_category": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Task category list.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Order by column.",
			},
			"ascend": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to sort in ascending order.",
			},
			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"task_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Probe task list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task name.",
						},
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task ID.",
						},
						"task_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task type.",
						},
						"nodes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Probe node list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"node_ip_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe node IP type.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe interval in minutes.",
						},
						"parameters": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Probe parameters.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task status.",
						},
						"target_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target address.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pay mode.",
						},
						"order_state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Order state.",
						},
						"task_category": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task category.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created time.",
						},
						"cron": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cron expression for scheduled task.",
						},
						"cron_state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Scheduled task start status.",
						},
						"tag_info_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag info list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"sub_sync_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a sync account.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of probe tasks.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCatProbeTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_probe_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	catService := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_ids"); ok {
		taskIDs := v.([]interface{})
		tmpList := make([]*string, 0, len(taskIDs))
		for _, item := range taskIDs {
			tmpList = append(tmpList, helper.String(item.(string)))
		}
		paramMap["task_ids"] = tmpList
	}
	if v, ok := d.GetOk("task_name"); ok {
		paramMap["task_name"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("target_address"); ok {
		paramMap["target_address"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("task_status"); ok {
		taskStatus := v.([]interface{})
		tmpList := make([]*int64, 0, len(taskStatus))
		for _, item := range taskStatus {
			tmpList = append(tmpList, helper.Int64(int64(item.(int))))
		}
		paramMap["task_status"] = tmpList
	}
	if v, ok := d.GetOk("pay_mode"); ok {
		paramMap["pay_mode"] = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("order_state"); ok {
		paramMap["order_state"] = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("task_type"); ok {
		taskType := v.([]interface{})
		tmpList := make([]*int64, 0, len(taskType))
		for _, item := range taskType {
			tmpList = append(tmpList, helper.Int64(int64(item.(int))))
		}
		paramMap["task_type"] = tmpList
	}
	if v, ok := d.GetOk("task_category"); ok {
		taskCategory := v.([]interface{})
		tmpList := make([]*int64, 0, len(taskCategory))
		for _, item := range taskCategory {
			tmpList = append(tmpList, helper.Int64(int64(item.(int))))
		}
		paramMap["task_category"] = tmpList
	}
	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("ascend"); ok {
		paramMap["ascend"] = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("tag_filters"); ok {
		tagFilters := v.([]interface{})
		tmpList := make([]*cat.KeyValuePair, 0, len(tagFilters))
		for _, item := range tagFilters {
			tagFilterMap := item.(map[string]interface{})
			keyValuePair := cat.KeyValuePair{}
			if v, ok := tagFilterMap["key"].(string); ok && v != "" {
				keyValuePair.Key = helper.String(v)
			}
			if v, ok := tagFilterMap["value"].(string); ok && v != "" {
				keyValuePair.Value = helper.String(v)
			}
			tmpList = append(tmpList, &keyValuePair)
		}
		paramMap["tag_filters"] = tmpList
	}

	var tasks []*cat.ProbeTask
	var total *int64
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, totalRet, e := catService.DescribeCatProbeTasksByFilter(ctx, paramMap)
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
		log.Printf("[CRITAL]%s read cat_probe_tasks failed, reason:%+v", logId, err)
		return err
	}

	taskSetList := make([]map[string]interface{}, 0, len(tasks))
	for _, task := range tasks {
		taskMap := map[string]interface{}{}
		if task.Name != nil {
			taskMap["name"] = task.Name
		}
		if task.TaskId != nil {
			taskMap["task_id"] = task.TaskId
		}
		if task.TaskType != nil {
			taskMap["task_type"] = task.TaskType
		}
		if task.Nodes != nil {
			nodesList := make([]string, 0, len(task.Nodes))
			for _, node := range task.Nodes {
				if node != nil {
					nodesList = append(nodesList, *node)
				}
			}
			taskMap["nodes"] = nodesList
		}
		if task.NodeIpType != nil {
			taskMap["node_ip_type"] = task.NodeIpType
		}
		if task.Interval != nil {
			taskMap["interval"] = task.Interval
		}
		if task.Parameters != nil {
			taskMap["parameters"] = task.Parameters
		}
		if task.Status != nil {
			taskMap["status"] = task.Status
		}
		if task.TargetAddress != nil {
			taskMap["target_address"] = task.TargetAddress
		}
		if task.PayMode != nil {
			taskMap["pay_mode"] = task.PayMode
		}
		if task.OrderState != nil {
			taskMap["order_state"] = task.OrderState
		}
		if task.TaskCategory != nil {
			taskMap["task_category"] = task.TaskCategory
		}
		if task.CreatedAt != nil {
			taskMap["created_at"] = task.CreatedAt
		}
		if task.Cron != nil {
			taskMap["cron"] = task.Cron
		}
		if task.CronState != nil {
			taskMap["cron_state"] = task.CronState
		}
		if task.TagInfoList != nil {
			tagInfoList := make([]map[string]interface{}, 0, len(task.TagInfoList))
			for _, tagInfo := range task.TagInfoList {
				tagInfoMap := map[string]interface{}{}
				if tagInfo.Key != nil {
					tagInfoMap["key"] = tagInfo.Key
				}
				if tagInfo.Value != nil {
					tagInfoMap["value"] = tagInfo.Value
				}
				tagInfoList = append(tagInfoList, tagInfoMap)
			}
			taskMap["tag_info_list"] = tagInfoList
		}
		if task.SubSyncFlag != nil {
			taskMap["sub_sync_flag"] = task.SubSyncFlag
		}

		taskSetList = append(taskSetList, taskMap)
	}

	if len(tasks) == 0 {
		log.Printf("[DATASOURCE] read empty, skip SetId")
	}

	d.SetId(helper.BuildToken())
	_ = d.Set("task_set", taskSetList)
	if total != nil {
		_ = d.Set("total", total)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskSetList); e != nil {
			return e
		}
	}

	return nil
}
