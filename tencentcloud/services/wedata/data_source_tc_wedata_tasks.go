package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataTasksRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Task name.",
			},

			"workflow_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workflow ID.",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner ID.",
			},

			"task_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Task type.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Task Status:\n* N: New\n* Y: Scheduling\n* F: Offline\n* O: Paused\n* T: Offlining\n* INVALID: Invalid.",
			},

			"submit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Submission status.",
			},

			"bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bundle id.",
			},

			"create_user_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Creator ID.",
			},

			"modify_time": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Modification time range (yyyy-MM-dd HH:mm:ss). Two time values must be provided in the array.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"create_time": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Creation time range (yyyy-MM-dd HH:MM:ss). Two time values must be provided in the array.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes the task pagination information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"task_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task ID.",
						},
						"task_type_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the task type ID.\n\n* 21:JDBC SQL\n* 23:TDSQL-PostgreSQL\n* 26:OfflineSynchronization\n* 30:Python\n* 31:PySpark\n* 33:Impala\n* 34:Hive SQL\n* 35:Shell\n* 36:Spark SQL\n* 38:Shell Form Mode\n* 39:Spark\n* 40:TCHouse-P\n* 41:Kettle\n* 42:Tchouse-X\n* 43:TCHouse-X SQL\n* 46:DLC Spark\n* 47:TiOne\n* 48:Trino\n* 50:DLC PySpark\n* 92:MapReduce\n* 130:Branch Node\n* 131:Merged Node\n* 132:Notebook\n* 133:SSH\n* 134:StarRocks\n* 137:For-each\n* 138:Setats SQL.",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Workflow ID.",
						},
						"task_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task name.",
						},
						"task_latest_version_no": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last save version number.",
						},
						"task_latest_submit_version_no": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last submit version number.",
						},
						"workflow_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Workflow name.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task Status:\n\n* N: New\n* Y: Scheduling\n* F: Offline\n* O: Paused\n* T: Offlining (in the process of being taken offline)\n* INVALID: Invalid.",
						},
						"submit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Latest submission status of the task. Specifies whether it has been submitted: true/false.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task creation time. example: 2022-02-12 11:13:41.",
						},
						"last_update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last update time. example: 2025-08-13 16:34:06.",
						},
						"last_update_user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last Updated By (Name).",
						},
						"last_ops_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last operation time.",
						},
						"last_ops_user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last operator name.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task owner ID.",
						},
						"task_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task description.",
						},
						"update_user_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last Updated User ID.",
						},
						"create_user_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Created By User ID.",
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

func dataSourceTencentCloudWedataTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_name"); ok {
		paramMap["TaskName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_id"); ok {
		paramMap["WorkflowId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_uin"); ok {
		paramMap["OwnerUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("task_type_id"); ok {
		paramMap["TaskTypeId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("submit"); ok {
		paramMap["Submit"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("bundle_id"); ok {
		paramMap["BundleId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_user_uin"); ok {
		paramMap["CreateUserUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("modify_time"); ok {
		modifyTimeList := []*string{}
		modifyTimeSet := v.(*schema.Set).List()
		for i := range modifyTimeSet {
			modifyTime := modifyTimeSet[i].(string)
			modifyTimeList = append(modifyTimeList, helper.String(modifyTime))
		}
		paramMap["ModifyTime"] = modifyTimeList
	}

	if v, ok := d.GetOk("create_time"); ok {
		createTimeList := []*string{}
		createTimeSet := v.(*schema.Set).List()
		for i := range createTimeSet {
			createTime := createTimeSet[i].(string)
			createTimeList = append(createTimeList, helper.String(createTime))
		}
		paramMap["CreateTime"] = createTimeList
	}

	var respData []*wedatav20250806.TaskBaseAttribute
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataTasksByFilter(ctx, paramMap)
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

		if items.TaskTypeId != nil {
			itemsMap["task_type_id"] = items.TaskTypeId
		}

		if items.WorkflowId != nil {
			itemsMap["workflow_id"] = items.WorkflowId
		}

		if items.TaskName != nil {
			itemsMap["task_name"] = items.TaskName
		}

		if items.TaskLatestVersionNo != nil {
			itemsMap["task_latest_version_no"] = items.TaskLatestVersionNo
		}

		if items.TaskLatestSubmitVersionNo != nil {
			itemsMap["task_latest_submit_version_no"] = items.TaskLatestSubmitVersionNo
		}

		if items.WorkflowName != nil {
			itemsMap["workflow_name"] = items.WorkflowName
		}

		if items.Status != nil {
			itemsMap["status"] = items.Status
		}

		if items.Submit != nil {
			itemsMap["submit"] = items.Submit
		}

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.LastUpdateTime != nil {
			itemsMap["last_update_time"] = items.LastUpdateTime
		}

		if items.LastUpdateUserName != nil {
			itemsMap["last_update_user_name"] = items.LastUpdateUserName
		}

		if items.LastOpsTime != nil {
			itemsMap["last_ops_time"] = items.LastOpsTime
		}

		if items.LastOpsUserName != nil {
			itemsMap["last_ops_user_name"] = items.LastOpsUserName
		}

		if items.OwnerUin != nil {
			itemsMap["owner_uin"] = items.OwnerUin
		}

		if items.TaskDescription != nil {
			itemsMap["task_description"] = items.TaskDescription
		}

		if items.UpdateUserUin != nil {
			itemsMap["update_user_uin"] = items.UpdateUserUin
		}

		if items.CreateUserUin != nil {
			itemsMap["create_user_uin"] = items.CreateUserUin
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
