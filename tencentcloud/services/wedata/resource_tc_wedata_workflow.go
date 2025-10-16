package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataWorkflowCreate,
		Read:   resourceTencentCloudWedataWorkflowRead,
		Update: resourceTencentCloudWedataWorkflowUpdate,
		Delete: resourceTencentCloudWedataWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project id.",
			},

			"workflow_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Workflow name.",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Parent folder path.",
			},

			"workflow_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workflow type, value example: cycle cycle workflow;manual manual workflow, passed in cycle by default.",
			},

			"workflow_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workflow description.",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Workflow Owner ID.",
			},

			"workflow_params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "workflow parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name.",
						},
						"param_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter value.",
						},
					},
				},
			},

			"workflow_scheduler_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Unified dispatch information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_time_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "time zone.",
						},
						"cycle_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cycle type: Supported types are\nONEOFF_CYCLE: One-time\nYEAR_CYCLE: Year\nMONTH_CYCLE: Month\nWEEK_CYCLE: Week\nDAY_CYCLE: Day\nHOUR_CYCLE: Hour\nMINUTE_CYCLE: Minute\nCRONTAB_CYCLE: crontab expression type.",
						},
						"self_depend": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Self-dependence, default value serial, values are: parallel, serial, orderly.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "End time.",
						},
						"crontab_expression": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Crontab expression.",
						},
						"dependency_workflow": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Workflow dependence, yes or no.",
						},
						"modify_cycle_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "0: Do not modify 1: Change the upstream dependency configuration of the task to the default value.",
						},
						"clear_link": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Workflows have cross-workflow dependencies and are scheduled using cron expressions. If you save unified scheduling, unsupported dependencies will be broken.",
						},
						"main_cyclic_config": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effective when ModifyCycleValue is 1, indicating the default modified upstream dependence-time dimension. The value is: \n* CRONTAB\n* DAY\n* HOUR\n* LIST_DAY\n* LIST_HOUR\n * LIST_MINUTE\n * MONTH\n* RANGE_DAY\n * RANGE_HOUR\n * RANGE_MINUTE\n* WEEK\n* YEAR\n\nhttps://capi.woa.com/object/detail? product=wedata&env=api_dev&version=2025-08-06&name=WorkflowSchedulerConfigurationInfo.",
						},
						"subordinate_cyclic_config": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effective when ModifyCycleValue is 1, which means that the default modified upstream dependency-instance range\n value is: \n* ALL_DAY_OF_YEAR\n* ALL_MONTH_OF_YEAR\n* CURRENT\n* CURRENT_DAY\n* CURRENT_HOUR\n* CURRENT_MINUTE\n* CURRENT_MONTH\n* CURRENT_WEEK\n* CURRENT_YEAR\n* PREVIOUS_BEGIN_OF_MONTH\n* PREVIOUS_DAY\n* PREVIOUS_DAY_LATER_OFFSET_HOUR\n* PREVIOUS_DAY_LATER_OFFSET_MINUTE\n* PREVIOUS_END_OF_MONTH\n* PREVIOUS_FRIDAY\n* PREVIOUS_HOUR\n* PREVIOUS_HOUR_CYCLE\n* PREVIOUS_HOUR_LATER_OFFSET_MINUTE\n* PREVIOUS_MINUTE_CYCLE\n* PREVIOUS_MONTH\n* PREVIOUS_WEEK\n* PREVIOUS_WEEKEND\n* RECENT_DATE\n\nhttps://capi.woa.com/object/detail? product=wedata&env=api_dev&version=2025-08-06&name=WorkflowSchedulerConfigurationInfo.",
						},
						"execution_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution time left-closed interval, example: 00:00, only if the cycle type is MINUTE_CYCLE needs to be filled in.",
						},
						"execution_end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution time right-closed interval, example: 23:59, only if the cycle type is MINUTE_CYCLE needs to be filled in.",
						},
						"calendar_open": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Do you want to turn on calendar scheduling 1 on 0 off.",
						},
						"calendar_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "calendar id.",
						},
					},
				},
			},

			"bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bundle Id.",
			},

			"bundle_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bundle Information.",
			},

			"workflow_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Workflow id.",
			},
		},
	}
}

func resourceTencentCloudWedataWorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId  string
		workflowId string
	)
	var (
		request  = wedatav20250806.NewCreateWorkflowRequest()
		response = wedatav20250806.NewCreateWorkflowResponse()
	)

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(projectId)
	}

	if v, ok := d.GetOk("workflow_name"); ok {
		request.WorkflowName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_type"); ok {
		request.WorkflowType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_desc"); ok {
		request.WorkflowDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("owner_uin"); ok {
		request.OwnerUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_params"); ok {
		for _, item := range v.(*schema.Set).List() {
			workflowParamsMap := item.(map[string]interface{})
			paramInfo := wedatav20250806.ParamInfo{}
			if v, ok := workflowParamsMap["param_key"]; ok {
				paramInfo.ParamKey = helper.String(v.(string))
			}
			if v, ok := workflowParamsMap["param_value"]; ok {
				paramInfo.ParamValue = helper.String(v.(string))
			}
			request.WorkflowParams = append(request.WorkflowParams, &paramInfo)
		}
	}

	if workflowSchedulerConfigurationMap, ok := helper.InterfacesHeadMap(d, "workflow_scheduler_configuration"); ok {
		workflowSchedulerConfigurationInfo := wedatav20250806.WorkflowSchedulerConfigurationInfo{}
		if v, ok := workflowSchedulerConfigurationMap["schedule_time_zone"]; ok {
			workflowSchedulerConfigurationInfo.ScheduleTimeZone = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["cycle_type"]; ok {
			workflowSchedulerConfigurationInfo.CycleType = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["self_depend"]; ok {
			workflowSchedulerConfigurationInfo.SelfDepend = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["start_time"]; ok {
			workflowSchedulerConfigurationInfo.StartTime = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["end_time"]; ok {
			workflowSchedulerConfigurationInfo.EndTime = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["crontab_expression"]; ok {
			workflowSchedulerConfigurationInfo.CrontabExpression = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["dependency_workflow"]; ok {
			workflowSchedulerConfigurationInfo.DependencyWorkflow = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["modify_cycle_value"]; ok && v.(string) != "" {
			workflowSchedulerConfigurationInfo.ModifyCycleValue = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["clear_link"]; ok {
			workflowSchedulerConfigurationInfo.ClearLink = helper.Bool(v.(bool))
		}
		if v, ok := workflowSchedulerConfigurationMap["main_cyclic_config"]; ok && v.(string) != "" {
			workflowSchedulerConfigurationInfo.MainCyclicConfig = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["subordinate_cyclic_config"]; ok && v.(string) != "" {
			workflowSchedulerConfigurationInfo.SubordinateCyclicConfig = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["execution_start_time"]; ok {
			workflowSchedulerConfigurationInfo.ExecutionStartTime = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["execution_end_time"]; ok {
			workflowSchedulerConfigurationInfo.ExecutionEndTime = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["calendar_open"]; ok && v.(string) != "" {
			workflowSchedulerConfigurationInfo.CalendarOpen = helper.String(v.(string))
		}
		if v, ok := workflowSchedulerConfigurationMap["calendar_id"]; ok {
			workflowSchedulerConfigurationInfo.CalendarId = helper.String(v.(string))
		}
		request.WorkflowSchedulerConfiguration = &workflowSchedulerConfigurationInfo
	}

	if v, ok := d.GetOk("bundle_id"); ok {
		request.BundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bundle_info"); ok {
		request.BundleInfo = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateWorkflowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata workflow failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Data != nil && response.Response.Data.WorkflowId != nil {
		workflowId = *response.Response.Data.WorkflowId
		d.SetId(strings.Join([]string{projectId, workflowId}, tccommon.FILED_SP))
	}

	return resourceTencentCloudWedataWorkflowRead(d, meta)
}

func resourceTencentCloudWedataWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	workflowId := idSplit[1]

	var (
		respData *wedatav20250806.WorkflowDetail
		innerErr error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		respData, innerErr = service.DescribeWedataWorkflowById(ctx, projectId, workflowId)
		if innerErr != nil {
			return resource.RetryableError(innerErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `wedata_workflow_folder` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)

	if respData.Path != nil && respData.WorkflowName != nil {
		_ = d.Set("parent_folder_path", strings.TrimSuffix(*respData.Path, fmt.Sprintf("/%s", *respData.WorkflowName)))
	}
	if respData.WorkflowName != nil {
		_ = d.Set("workflow_name", respData.WorkflowName)
	}

	if respData.OwnerUin != nil {
		_ = d.Set("owner_uin", respData.OwnerUin)
	}

	if respData.WorkflowType != nil {
		_ = d.Set("workflow_type", respData.WorkflowType)
	}

	workflowParamsList := make([]map[string]interface{}, 0, len(respData.WorkflowParams))
	if respData.WorkflowParams != nil {
		for _, workflowParams := range respData.WorkflowParams {
			workflowParamsMap := map[string]interface{}{}

			if workflowParams.ParamKey != nil {
				workflowParamsMap["param_key"] = workflowParams.ParamKey
			}

			if workflowParams.ParamValue != nil {
				workflowParamsMap["param_value"] = workflowParams.ParamValue
			}

			workflowParamsList = append(workflowParamsList, workflowParamsMap)
		}

		_ = d.Set("workflow_params", workflowParamsList)
	}

	workflowSchedulerConfigurationMap := map[string]interface{}{}

	if respData.WorkflowSchedulerConfiguration != nil {
		if respData.WorkflowSchedulerConfiguration.ScheduleTimeZone != nil {
			workflowSchedulerConfigurationMap["schedule_time_zone"] = respData.WorkflowSchedulerConfiguration.ScheduleTimeZone
		}

		if respData.WorkflowSchedulerConfiguration.CycleType != nil {
			workflowSchedulerConfigurationMap["cycle_type"] = respData.WorkflowSchedulerConfiguration.CycleType
		}

		if respData.WorkflowSchedulerConfiguration.SelfDepend != nil {
			workflowSchedulerConfigurationMap["self_depend"] = respData.WorkflowSchedulerConfiguration.SelfDepend
		}

		if respData.WorkflowSchedulerConfiguration.StartTime != nil {
			workflowSchedulerConfigurationMap["start_time"] = respData.WorkflowSchedulerConfiguration.StartTime
		}

		if respData.WorkflowSchedulerConfiguration.EndTime != nil {
			workflowSchedulerConfigurationMap["end_time"] = respData.WorkflowSchedulerConfiguration.EndTime
		}

		if respData.WorkflowSchedulerConfiguration.DependencyWorkflow != nil {
			workflowSchedulerConfigurationMap["dependency_workflow"] = respData.WorkflowSchedulerConfiguration.DependencyWorkflow
		}

		if respData.WorkflowSchedulerConfiguration.ExecutionStartTime != nil {
			workflowSchedulerConfigurationMap["execution_start_time"] = respData.WorkflowSchedulerConfiguration.ExecutionStartTime
		}

		if respData.WorkflowSchedulerConfiguration.ExecutionEndTime != nil {
			workflowSchedulerConfigurationMap["execution_end_time"] = respData.WorkflowSchedulerConfiguration.ExecutionEndTime
		}

		if respData.WorkflowSchedulerConfiguration.CrontabExpression != nil {
			workflowSchedulerConfigurationMap["crontab_expression"] = respData.WorkflowSchedulerConfiguration.CrontabExpression
		}

		if respData.WorkflowSchedulerConfiguration.CalendarOpen != nil {
			workflowSchedulerConfigurationMap["calendar_open"] = respData.WorkflowSchedulerConfiguration.CalendarOpen
		}

		if respData.WorkflowSchedulerConfiguration.CalendarName != nil {
			workflowSchedulerConfigurationMap["calendar_name"] = respData.WorkflowSchedulerConfiguration.CalendarName
		}

		if respData.WorkflowSchedulerConfiguration.CalendarId != nil {
			workflowSchedulerConfigurationMap["calendar_id"] = respData.WorkflowSchedulerConfiguration.CalendarId
		}

		_ = d.Set("workflow_scheduler_configuration", []interface{}{workflowSchedulerConfigurationMap})
	}

	if respData.WorkflowDesc != nil {
		_ = d.Set("workflow_desc", respData.WorkflowDesc)
	}

	if respData.BundleId != nil {
		_ = d.Set("bundle_id", respData.BundleId)
	}

	if respData.BundleInfo != nil {
		_ = d.Set("bundle_info", respData.BundleInfo)
	}

	_ = d.Set("workflow_id", workflowId)
	_ = projectId
	return nil
}

func resourceTencentCloudWedataWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"parent_folder_path", "workflow_type"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	workflowId := idSplit[1]

	needChange := false
	mutableArgs := []string{"project_id", "workflow_id", "workflow_name", "owner_uin", "workflow_desc", "workflow_params", "workflow_scheduler_configuration", "bundle_id", "bundle_info"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateWorkflowRequest()
		request.ProjectId = helper.String(projectId)
		request.WorkflowId = helper.String(workflowId)

		if v, ok := d.GetOk("workflow_name"); ok {
			request.WorkflowName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("owner_uin"); ok {
			request.OwnerUin = helper.String(v.(string))
		}

		if v, ok := d.GetOk("workflow_desc"); ok {
			request.WorkflowDesc = helper.String(v.(string))
		}

		if v, ok := d.GetOk("workflow_params"); ok {
			for _, item := range v.(*schema.Set).List() {
				workflowParamsMap := item.(map[string]interface{})
				paramInfo := wedatav20250806.ParamInfo{}
				if v, ok := workflowParamsMap["param_key"]; ok {
					paramInfo.ParamKey = helper.String(v.(string))
				}
				if v, ok := workflowParamsMap["param_value"]; ok {
					paramInfo.ParamValue = helper.String(v.(string))
				}
				request.WorkflowParams = append(request.WorkflowParams, &paramInfo)
			}
		}

		if workflowSchedulerConfigurationMap, ok := helper.InterfacesHeadMap(d, "workflow_scheduler_configuration"); ok {
			workflowSchedulerConfigurationInfo := wedatav20250806.WorkflowSchedulerConfigurationInfo{}
			if v, ok := workflowSchedulerConfigurationMap["schedule_time_zone"]; ok {
				workflowSchedulerConfigurationInfo.ScheduleTimeZone = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["cycle_type"]; ok {
				workflowSchedulerConfigurationInfo.CycleType = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["self_depend"]; ok {
				workflowSchedulerConfigurationInfo.SelfDepend = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["start_time"]; ok {
				workflowSchedulerConfigurationInfo.StartTime = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["end_time"]; ok {
				workflowSchedulerConfigurationInfo.EndTime = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["crontab_expression"]; ok {
				workflowSchedulerConfigurationInfo.CrontabExpression = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["dependency_workflow"]; ok {
				workflowSchedulerConfigurationInfo.DependencyWorkflow = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["modify_cycle_value"]; ok && v.(string) != "" {
				workflowSchedulerConfigurationInfo.ModifyCycleValue = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["clear_link"]; ok {
				workflowSchedulerConfigurationInfo.ClearLink = helper.Bool(v.(bool))
			}
			if v, ok := workflowSchedulerConfigurationMap["main_cyclic_config"]; ok && v.(string) != "" {
				workflowSchedulerConfigurationInfo.MainCyclicConfig = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["subordinate_cyclic_config"]; ok && v.(string) != "" {
				workflowSchedulerConfigurationInfo.SubordinateCyclicConfig = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["execution_start_time"]; ok {
				workflowSchedulerConfigurationInfo.ExecutionStartTime = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["execution_end_time"]; ok {
				workflowSchedulerConfigurationInfo.ExecutionEndTime = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["calendar_open"]; ok && v.(string) != "" {
				workflowSchedulerConfigurationInfo.CalendarOpen = helper.String(v.(string))
			}
			if v, ok := workflowSchedulerConfigurationMap["calendar_id"]; ok {
				workflowSchedulerConfigurationInfo.CalendarId = helper.String(v.(string))
			}
			request.WorkflowSchedulerConfiguration = &workflowSchedulerConfigurationInfo
		}

		if v, ok := d.GetOk("bundle_id"); ok {
			request.BundleId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("bundle_info"); ok {
			request.BundleInfo = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateWorkflowWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update wedata workflow failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = projectId
	_ = workflowId
	return resourceTencentCloudWedataWorkflowRead(d, meta)
}

func resourceTencentCloudWedataWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_workflow.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectId := idSplit[0]
	workflowId := idSplit[1]

	var (
		request  = wedatav20250806.NewDeleteWorkflowRequest()
		response = wedatav20250806.NewDeleteWorkflowResponse()
	)

	request.ProjectId = helper.String(projectId)
	request.WorkflowId = helper.String(workflowId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteWorkflowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete wedata workflow failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = projectId
	return nil
}
