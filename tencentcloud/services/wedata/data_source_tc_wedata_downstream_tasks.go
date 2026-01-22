package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataDownstreamTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataDownstreamTasksRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task ID.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes the downstream dependency details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task ID.",
						},
						"task_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task name.",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow id.",
						},
						"workflow_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow name.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Project ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task Status:\n\n* N: New\n\n* Y: Scheduling\n\n* F: Offline\n\n* O: Paused\n\n* T: Offlining (in the process of being taken offline)\n\nI* NVALID: Invalid.",
						},
						"task_type_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Task type id.",
						},
						"task_type_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task type description.\n-20: universal data synchronization.\n - 25:  ETLTaskType\n - 26:  ETLTaskType\n - 30:  python\n - 31:  pyspark\n - 34:  hivesql\n - 35:  shell\n - 36:  sparksql\n - 21:  jdbcsql\n - 32:  dlc\n - 33:  ImpalaTaskType\n - 40:  CDWTaskType\n - 41:  kettle\n - 42:  TCHouse-X\n - 43:  TCHouse-X SQL\n - 46:  dlcsparkTaskType\n - 47:  TiOneMachineLearningTaskType\n - 48:  Trino\n - 50:  DLCPyspark\n - 23:  TencentDistributedSQL\n - 39:  spark\n - 92:  MRTaskType\n - 38:  ShellScript\n - 70:  HiveSQLScrip\n-130: branch.\n-131: merge.\n-132: Notebook \n-133: SSH node.\n - 134:  StarRocks\n - 137:  For-each\n-10000: custom business common.",
						},
						"schedule_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies scheduling plan display description information.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task end time.",
						},
						"delay_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Delay time.",
						},
						"cycle_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cycle Type, Default: D\nSupported types:\n* O: One-time\n\n* Y: Yearly\n\n* M: Monthly\n\n* W: Weekly\n\n* D: Daily\n\n* H: Hourly\n\n* I: Minute\n\n* C: Crontab expression type.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Owner ID.",
						},
						"task_action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Elastic cycle configuration.",
						},
						"init_strategy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Initialization strategy for scheduling.",
						},
						"crontab_expression": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "crontab expression.",
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

func dataSourceTencentCloudWedataDownstreamTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_downstream_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.TaskDependDto
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataDownstreamTasksByFilter(ctx, paramMap)
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
	itemsList := make([]map[string]interface{}, 0, len(respData))

	for _, items := range respData {
		itemsMap := map[string]interface{}{}

		if items.TaskId != nil {
			itemsMap["task_id"] = items.TaskId
			ids = append(ids, *items.TaskId)
		}

		if items.TaskName != nil {
			itemsMap["task_name"] = items.TaskName
		}

		if items.WorkflowId != nil {
			itemsMap["workflow_id"] = items.WorkflowId
		}

		if items.WorkflowName != nil {
			itemsMap["workflow_name"] = items.WorkflowName
		}

		if items.ProjectId != nil {
			itemsMap["project_id"] = items.ProjectId
		}

		if items.Status != nil {
			itemsMap["status"] = items.Status
		}

		if items.TaskTypeId != nil {
			itemsMap["task_type_id"] = items.TaskTypeId
		}

		if items.TaskTypeDesc != nil {
			itemsMap["task_type_desc"] = items.TaskTypeDesc
		}

		if items.ScheduleDesc != nil {
			itemsMap["schedule_desc"] = items.ScheduleDesc
		}

		if items.StartTime != nil {
			itemsMap["start_time"] = items.StartTime
		}

		if items.EndTime != nil {
			itemsMap["end_time"] = items.EndTime
		}

		if items.DelayTime != nil {
			itemsMap["delay_time"] = items.DelayTime
		}

		if items.CycleType != nil {
			itemsMap["cycle_type"] = items.CycleType
		}

		if items.OwnerUin != nil {
			itemsMap["owner_uin"] = items.OwnerUin
		}

		if items.TaskAction != nil {
			itemsMap["task_action"] = items.TaskAction
		}

		if items.InitStrategy != nil {
			itemsMap["init_strategy"] = items.InitStrategy
		}

		if items.CrontabExpression != nil {
			itemsMap["crontab_expression"] = items.CrontabExpression
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("data", itemsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
