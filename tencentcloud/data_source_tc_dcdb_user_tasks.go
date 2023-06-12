/*
Use this data source to query detailed information of dcdb user_tasks

Example Usage

```hcl
data "tencentcloud_dcdb_user_tasks" "user_tasks" {
  statuses =
  instance_ids =
  flow_types =
  start_time = ""
  end_time = ""
  u_task_ids =
  }
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbUserTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbUserTasksRead,
		Schema: map[string]*schema.Schema{
			"statuses": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "task status. 0-starting; 1-running; 2-success; 3-failed.",
			},

			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "list of instance id.",
			},

			"flow_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "task type, 0-rollback task; 1-create instance task; 2-expansion task; 3-migration task; 4-delete instance task; 5-restart task.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "task creation time.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "task end time.",
			},

			"u_task_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "list of task id.",
			},

			"flow_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "task id.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "user app id.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "task status. 0-starting; 1-running; 2-success; 3-failed.",
						},
						"user_task_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "task type, 0-rollback task; 1-create instance task; 2-expansion task; 3-migration task; 4-delete instance task; 5-restart task.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task creation time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task end time.",
						},
						"err_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task error message.",
						},
						"input_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task input data.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance name.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "region id.",
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

func dataSourceTencentCloudDcdbUserTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_user_tasks.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		instanceIds []string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("statuses"); ok {
		statusesSet := v.(*schema.Set).List()
		tmpList := make([]interface{}, len(statusesSet))
		for i := range statusesSet {
			statuses := statusesSet[i].(int)
			tmpList = append(tmpList, helper.IntInt64(statuses))
		}
		paramMap["Statuses"] = tmpList
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		instanceIds = helper.InterfacesStrings(instanceIdsSet)
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("flow_types"); ok {
		flowTypesSet := v.(*schema.Set).List()
		tmpList := make([]interface{}, len(flowTypesSet))
		for i := range flowTypesSet {
			flowTypes := flowTypesSet[i].(int)
			tmpList = append(tmpList, helper.IntInt64(flowTypes))
		}
		paramMap["FlowTypes"] = tmpList
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("u_task_ids"); ok {
		uTaskIdsSet := v.(*schema.Set).List()
		tmpList := make([]interface{}, len(uTaskIdsSet))
		for i := range uTaskIdsSet {
			uTaskIds := uTaskIdsSet[i].(int)
			tmpList = append(tmpList, helper.IntInt64(uTaskIds))
		}
		paramMap["UTaskIds"] = tmpList
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var flowSet []*dcdb.UserTaskInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbUserTasksByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		flowSet = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(flowSet))

	if flowSet != nil {
		for _, userTaskInfo := range flowSet {
			userTaskInfoMap := map[string]interface{}{}

			if userTaskInfo.Id != nil {
				userTaskInfoMap["id"] = userTaskInfo.Id
			}

			if userTaskInfo.AppId != nil {
				userTaskInfoMap["app_id"] = userTaskInfo.AppId
			}

			if userTaskInfo.Status != nil {
				userTaskInfoMap["status"] = userTaskInfo.Status
			}

			if userTaskInfo.UserTaskType != nil {
				userTaskInfoMap["user_task_type"] = userTaskInfo.UserTaskType
			}

			if userTaskInfo.CreateTime != nil {
				userTaskInfoMap["create_time"] = userTaskInfo.CreateTime
			}

			if userTaskInfo.EndTime != nil {
				userTaskInfoMap["end_time"] = userTaskInfo.EndTime
			}

			if userTaskInfo.ErrMsg != nil {
				userTaskInfoMap["err_msg"] = userTaskInfo.ErrMsg
			}

			if userTaskInfo.InputData != nil {
				userTaskInfoMap["input_data"] = userTaskInfo.InputData
			}

			if userTaskInfo.InstanceId != nil {
				userTaskInfoMap["instance_id"] = userTaskInfo.InstanceId
			}

			if userTaskInfo.InstanceName != nil {
				userTaskInfoMap["instance_name"] = userTaskInfo.InstanceName
			}

			if userTaskInfo.RegionId != nil {
				userTaskInfoMap["region_id"] = userTaskInfo.RegionId
			}

			tmpList = append(tmpList, userTaskInfoMap)
		}

		_ = d.Set("flow_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(instanceIds))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
