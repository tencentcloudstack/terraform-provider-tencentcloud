package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrReplicationInstanceCreateTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrReplicationInstanceCreateTasksRead,
		Schema: map[string]*schema.Schema{
			"replication_registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "synchronization instance Id, see RegistryId in DescribeReplicationInstances.",
			},

			"replication_region_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "synchronization instance region Id, see ReplicationRegionId in DescribeReplicationInstances.",
			},

			"task_detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "task details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task name.",
						},
						"task_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task UUID.",
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task status.",
						},
						"task_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status information. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task start name.",
						},
						"finished_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task end time. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "overall task status.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTcrReplicationInstanceCreateTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_replication_instance_create_tasks.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("replication_registry_id"); ok {
		paramMap["ReplicationRegistryId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("replication_region_id"); v != nil {
		paramMap["ReplicationRegionId"] = helper.IntUint64(v.(int))
	}

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		result *tcr.DescribeReplicationInstanceCreateTasksResponseParams
		e      error
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeTcrReplicationInstanceCreateTasksByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	taskDetails := result.TaskDetail
	status := result.Status
	ids := []string{}

	if taskDetails != nil {
		tmpList := make([]map[string]interface{}, 0, len(taskDetails))
		for _, taskDetail := range taskDetails {
			taskDetailMap := map[string]interface{}{}

			if taskDetail.TaskName != nil {
				taskDetailMap["task_name"] = taskDetail.TaskName
			}

			if taskDetail.TaskUUID != nil {
				taskDetailMap["task_uuid"] = taskDetail.TaskUUID
			}

			if taskDetail.TaskStatus != nil {
				taskDetailMap["task_status"] = taskDetail.TaskStatus
			}

			if taskDetail.TaskMessage != nil {
				taskDetailMap["task_message"] = taskDetail.TaskMessage
			}

			if taskDetail.CreatedTime != nil {
				taskDetailMap["created_time"] = taskDetail.CreatedTime
			}

			if taskDetail.FinishedTime != nil {
				taskDetailMap["finished_time"] = taskDetail.FinishedTime
			}

			ids = append(ids, *taskDetail.TaskUUID)
			tmpList = append(tmpList, taskDetailMap)
		}

		_ = d.Set("task_detail", tmpList)
	}

	if status != nil {
		_ = d.Set("status", status)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
