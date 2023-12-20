package crs

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisInstanceTaskList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisInstanceTaskListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Project Id.",
			},

			"task_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Task type.",
			},

			"begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Termination time.",
			},

			"task_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Task status.",
			},

			"result": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Task status.",
			},

			"operate_uin": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Operator Uin.",
			},

			"tasks": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Task details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time.",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task type.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of instance.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of instance.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project ID.",
						},
						"progress": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Task progress.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end time.",
						},
						"result": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task status.",
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

func dataSourceTencentCloudRedisInstanceTaskListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_instance_task_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		paramMap["ProjectIds"] = helper.InterfacesIntInt64Point(projectIdsSet)
	}

	if v, ok := d.GetOk("task_types"); ok {
		taskTypesSet := v.(*schema.Set).List()
		paramMap["TaskTypes"] = helper.InterfacesStringsPoint(taskTypesSet)
	}

	if v, ok := d.GetOk("begin_time"); ok {
		paramMap["BeginTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_status"); ok {
		taskStatusSet := v.(*schema.Set).List()
		paramMap["TaskStatus"] = helper.InterfacesIntInt64Point(taskStatusSet)
	}

	if v, ok := d.GetOk("result"); ok {
		resultSet := v.(*schema.Set).List()
		paramMap["Result"] = helper.InterfacesIntInt64Point(resultSet)
	}

	if v, ok := d.GetOk("operate_uin"); ok {
		operateUinSet := v.(*schema.Set).List()
		paramMap["OperateUin"] = helper.InterfacesStringsPoint(operateUinSet)
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tasks []*redis.TaskInfoDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisInstanceTaskListByFilter(ctx, paramMap)
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
		for _, taskInfoDetail := range tasks {
			taskInfoDetailMap := map[string]interface{}{}

			if taskInfoDetail.TaskId != nil {
				taskInfoDetailMap["task_id"] = taskInfoDetail.TaskId
			}

			if taskInfoDetail.StartTime != nil {
				taskInfoDetailMap["start_time"] = taskInfoDetail.StartTime
			}

			if taskInfoDetail.TaskType != nil {
				taskInfoDetailMap["task_type"] = taskInfoDetail.TaskType
			}

			if taskInfoDetail.InstanceName != nil {
				taskInfoDetailMap["instance_name"] = taskInfoDetail.InstanceName
			}

			if taskInfoDetail.InstanceId != nil {
				taskInfoDetailMap["instance_id"] = taskInfoDetail.InstanceId
			}

			if taskInfoDetail.ProjectId != nil {
				taskInfoDetailMap["project_id"] = taskInfoDetail.ProjectId
			}

			if taskInfoDetail.Progress != nil {
				taskInfoDetailMap["progress"] = taskInfoDetail.Progress
			}

			if taskInfoDetail.EndTime != nil {
				taskInfoDetailMap["end_time"] = taskInfoDetail.EndTime
			}

			if taskInfoDetail.Result != nil {
				taskInfoDetailMap["result"] = taskInfoDetail.Result
			}

			ids = append(ids, strconv.FormatInt(*taskInfoDetail.TaskId, 10))
			tmpList = append(tmpList, taskInfoDetailMap)
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
