package cdb

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlUserTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlUserTaskRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.",
			},

			"async_request_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Asynchronous task request ID, the AsyncRequestId returned by executing cloud database-related operations.",
			},

			"task_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Task type. If no value is passed, all task types will be queried. Supported values include: `ROLLBACK` - database rollback; `SQL OPERATION` - SQL operation; `IMPORT DATA` - data import; `MODIFY PARAM` - parameter setting; `INITIAL` - initialize the cloud database instance; `REBOOT` - restarts the cloud database instance; `OPEN GTID` - open the cloud database instance GTID; `UPGRADE RO` - read-only instance upgrade; `BATCH ROLLBACK` - database batch rollback; `UPGRADE MASTER` - master upgrade; `DROP TABLES` - delete cloud database tables; `SWITCH DR TO MASTER` - The disaster recovery instance.",
			},

			"task_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Task status. If no value is passed, all task statuses will be queried. Supported values include: `UNDEFINED` - undefined; `INITIAL` - initialization; `RUNNING` - running; `SUCCEED` - the execution was successful; `FAILED` - execution failed; `KILLED` - terminated; `REMOVED` - removed; `PAUSED` - Paused.",
			},

			"start_time_begin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The start time of the first task, used for range query, the time format is as follows: 2017-12-31 10:40:01.",
			},

			"start_time_end": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The start time of the last task, used for range query, the time format is as follows: 2017-12-31 10:40:01.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The returned instance task information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "error code.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "error message.",
						},
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance task ID.",
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance task progress.",
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance task status, possible values include:UNDEFINED - undefined;INITIAL - initialization;RUNNING - running;SUCCEED - the execution was successful;FAILED - execution failed;KILLED - terminated;REMOVED - removed;PAUSED - Paused.WAITING - waiting (cancellable).",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance task type, possible values include:ROLLBACK - database rollback;SQL OPERATION - SQL operation;IMPORT DATA - data import;MODIFY PARAM - parameter setting;INITIAL - initialize the cloud database instance;REBOOT - restarts the cloud database instance;OPEN GTID - open the cloud database instance GTID;UPGRADE RO - read-only instance upgrade;BATCH ROLLBACK - database batch rollback;UPGRADE MASTER - master upgrade;DROP TABLES - delete cloud database tables;SWITCH DR TO MASTER - The disaster recovery instance.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance task start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance task end time.",
						},
						"instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The instance ID associated with the task. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"async_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request ID of the asynchronous task.",
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

func dataSourceTencentCloudMysqlUserTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_user_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("async_request_id"); ok {
		paramMap["AsyncRequestId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_types"); ok {
		taskTypesSet := v.(*schema.Set).List()
		paramMap["TaskTypes"] = helper.InterfacesStringsPoint(taskTypesSet)
	}

	if v, ok := d.GetOk("task_status"); ok {
		taskStatusSet := v.(*schema.Set).List()
		paramMap["TaskStatus"] = helper.InterfacesStringsPoint(taskStatusSet)

	}

	if v, ok := d.GetOk("start_time_begin"); ok {
		paramMap["StartTimeBegin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time_end"); ok {
		paramMap["StartTimeEnd"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.TaskDetail
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlUserTaskByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))
	if items != nil {
		for _, taskDetail := range items {
			taskDetailMap := map[string]interface{}{}

			if taskDetail.Code != nil {
				taskDetailMap["code"] = taskDetail.Code
			}

			if taskDetail.Message != nil {
				taskDetailMap["message"] = taskDetail.Message
			}

			if taskDetail.JobId != nil {
				taskDetailMap["job_id"] = taskDetail.JobId
			}

			if taskDetail.Progress != nil {
				taskDetailMap["progress"] = taskDetail.Progress
			}

			if taskDetail.TaskStatus != nil {
				taskDetailMap["task_status"] = taskDetail.TaskStatus
			}

			if taskDetail.TaskType != nil {
				taskDetailMap["task_type"] = taskDetail.TaskType
			}

			if taskDetail.StartTime != nil {
				taskDetailMap["start_time"] = taskDetail.StartTime
			}

			if taskDetail.EndTime != nil {
				taskDetailMap["end_time"] = taskDetail.EndTime
			}

			if taskDetail.InstanceIds != nil {
				taskDetailMap["instance_ids"] = taskDetail.InstanceIds
			}

			if taskDetail.AsyncRequestId != nil {
				taskDetailMap["async_request_id"] = taskDetail.AsyncRequestId
			}

			ids = append(ids, strconv.FormatInt(*taskDetail.JobId, 10))
			tmpList = append(tmpList, taskDetailMap)
		}

		_ = d.Set("items", tmpList)
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
