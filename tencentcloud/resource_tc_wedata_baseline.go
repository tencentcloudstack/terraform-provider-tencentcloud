package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataBaseline() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataBaselineCreate,
		Read:   resourceTencentCloudWedataBaselineRead,
		Update: resourceTencentCloudWedataBaselineUpdate,
		Delete: resourceTencentCloudWedataBaselineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},
			"baseline_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Baseline Name.",
			},
			"baseline_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "D or H; representing daily baseline and hourly baseline respectively.",
			},
			"create_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Creator ID.",
			},
			"create_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Creator Name.",
			},
			"in_charge_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Baseline Owner ID.",
			},
			"in_charge_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Baseline Owner Name.",
			},
			"promise_tasks": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Promise Tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Project ID.",
						},
						"task_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task Name.",
						},
						"task_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task ID.",
						},
						"task_cycle": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task Scheduling Cycle.",
						},
						"workflow_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Workflow Name.",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Workflow ID.",
						},
						"task_in_charge_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task Owner Name.",
						},
						"task_in_charge_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Task Owner ID.",
						},
					},
				},
			},
			"promise_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service Assurance Time.",
			},
			"warning_margin": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Warning Margin in minutes.",
			},
			"is_new_alarm": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Is it a newly created alarm rule.",
			},
			"alarm_rule_dto": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Existing Alarm Rule Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm Rule ID.",
						},
						"alarm_level_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Important;Urgent;Normal.",
						},
					},
				},
			},
			"baseline_create_alarm_rule_request": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Description of the New Alarm Rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Project NameNote: This field may return null, indicating no valid value.",
						},
						"creator_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Creator NameNote: This field may return null, indicating no valid value.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Creator UINNote: This field may return null, indicating no valid value.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Rule NameNote: This field may return null, indicating no valid value.",
						},
						"monitor_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "Monitoring Type, 1. Task, 2. Workflow, 3. Project, 4. Baseline (default is 1. Task)Note: This field may return null, indicating no valid value.",
						},
						"monitor_object_ids": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Monitoring ObjectsNote: This field may return null, indicating no valid value.",
						},
						"alarm_types": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Alarm Types, 1. Failure Alarm, 2. Timeout Alarm, 3. Success Alarm, 4. Baseline Violation, 5. Baseline Warning, 6. Baseline Task Failure (default is 1. Failure Alarm)Note: This field may return null, indicating no valid value.",
						},
						"alarm_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "Alarm Level, 1. Normal, 2. Important, 3. Urgent (default is 1. Normal)Note: This field may return null, indicating no valid value.",
						},
						"alarm_ways": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Alarm Methods, 1. Email, 2. SMS, 3. WeChat, 4. Voice, 5. Enterprise WeChat, 6. HTTP, 7. Enterprise WeChat Group; Alarm method code list (default is 1. Email)Note: This field may return null, indicating no valid value.",
						},
						"alarm_recipient_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "Alarm Recipient Type: 1. Specified Personnel, 2. Task Owner, 3. Duty Roster (default is 1. Specified Personnel)Note: This field may return null, indicating no valid value.",
						},
						"alarm_recipients": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Alarm RecipientsNote: This field may return null, indicating no valid value.",
						},
						"alarm_recipient_ids": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Alarm Recipient IDsNote: This field may return null, indicating no valid value.",
						},
						"ext_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Extended Information, 1. Estimated Runtime (default), 2. Estimated Completion Time, 3. Estimated Scheduling Time, 4. Incomplete within the Cycle; Value Types: 1. Specified Value, 2. Historical Average (default is 1. Specified Value)Note: This field may return null, indicating no valid value.",
						},
					},
				},
			},
			"baseline_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Baseline ID.",
			},
		},
	}
}

func resourceTencentCloudWedataBaselineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = wedata.NewCreateBaselineRequest()
		response   = wedata.NewCreateBaselineResponse()
		projectId  string
		baselineId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("baseline_name"); ok {
		request.BaselineName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("baseline_type"); ok {
		request.BaselineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_uin"); ok {
		request.CreateUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_name"); ok {
		request.CreateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("in_charge_uin"); ok {
		request.InChargeUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("in_charge_name"); ok {
		request.InChargeName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("promise_tasks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			baselineTaskInfo := wedata.BaselineTaskInfo{}
			if v, ok := dMap["project_id"]; ok {
				baselineTaskInfo.ProjectId = helper.String(v.(string))
			}

			if v, ok := dMap["task_name"]; ok {
				baselineTaskInfo.TaskName = helper.String(v.(string))
			}

			if v, ok := dMap["task_id"]; ok {
				baselineTaskInfo.TaskId = helper.String(v.(string))
			}

			if v, ok := dMap["task_cycle"]; ok {
				baselineTaskInfo.TaskCycle = helper.String(v.(string))
			}

			if v, ok := dMap["workflow_name"]; ok {
				baselineTaskInfo.WorkflowName = helper.String(v.(string))
			}

			if v, ok := dMap["workflow_id"]; ok {
				baselineTaskInfo.WorkflowId = helper.String(v.(string))
			}

			if v, ok := dMap["task_in_charge_name"]; ok {
				baselineTaskInfo.TaskInChargeName = helper.String(v.(string))
			}

			if v, ok := dMap["task_in_charge_uin"]; ok {
				baselineTaskInfo.TaskInChargeUin = helper.String(v.(string))
			}

			request.PromiseTasks = append(request.PromiseTasks, &baselineTaskInfo)
		}
	}

	if v, ok := d.GetOk("promise_time"); ok {
		request.PromiseTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("warning_margin"); ok {
		request.WarningMargin = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("is_new_alarm"); ok {
		request.IsNewAlarm = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "alarm_rule_dto"); ok {
		alarmRuleDto := wedata.AlarmRuleDto{}
		if v, ok := dMap["alarm_rule_id"]; ok {
			alarmRuleDto.AlarmRuleId = helper.String(v.(string))
		}

		if v, ok := dMap["alarm_level_type"]; ok {
			alarmRuleDto.AlarmLevelType = helper.String(v.(string))
		}

		request.AlarmRuleDto = &alarmRuleDto
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "baseline_create_alarm_rule_request"); ok {
		createAlarmRuleRequest := wedata.CreateAlarmRuleRequest{}
		if v, ok := dMap["project_id"]; ok {
			createAlarmRuleRequest.ProjectId = helper.String(v.(string))
		}

		if v, ok := dMap["creator_id"]; ok {
			createAlarmRuleRequest.CreatorId = helper.String(v.(string))
		}

		if v, ok := dMap["creator"]; ok {
			createAlarmRuleRequest.Creator = helper.String(v.(string))
		}

		if v, ok := dMap["rule_name"]; ok {
			createAlarmRuleRequest.RuleName = helper.String(v.(string))
		}

		if v, ok := dMap["monitor_type"]; ok {
			createAlarmRuleRequest.MonitorType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["monitor_object_ids"]; ok {
			monitorObjectIdsSet := v.(*schema.Set).List()
			for i := range monitorObjectIdsSet {
				monitorObjectIds := monitorObjectIdsSet[i].(string)
				createAlarmRuleRequest.MonitorObjectIds = append(createAlarmRuleRequest.MonitorObjectIds, &monitorObjectIds)
			}
		}

		if v, ok := dMap["alarm_types"]; ok {
			alarmTypesSet := v.(*schema.Set).List()
			for i := range alarmTypesSet {
				alarmTypes := alarmTypesSet[i].(string)
				createAlarmRuleRequest.AlarmTypes = append(createAlarmRuleRequest.AlarmTypes, &alarmTypes)
			}
		}

		if v, ok := dMap["alarm_level"]; ok {
			createAlarmRuleRequest.AlarmLevel = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["alarm_ways"]; ok {
			alarmWaysSet := v.(*schema.Set).List()
			for i := range alarmWaysSet {
				alarmWays := alarmWaysSet[i].(string)
				createAlarmRuleRequest.AlarmWays = append(createAlarmRuleRequest.AlarmWays, &alarmWays)
			}
		}

		if v, ok := dMap["alarm_recipient_type"]; ok {
			createAlarmRuleRequest.AlarmRecipientType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["alarm_recipients"]; ok {
			alarmRecipientsSet := v.(*schema.Set).List()
			for i := range alarmRecipientsSet {
				alarmRecipients := alarmRecipientsSet[i].(string)
				createAlarmRuleRequest.AlarmRecipients = append(createAlarmRuleRequest.AlarmRecipients, &alarmRecipients)
			}
		}

		if v, ok := dMap["alarm_recipient_ids"]; ok {
			alarmRecipientIdsSet := v.(*schema.Set).List()
			for i := range alarmRecipientIdsSet {
				alarmRecipientIds := alarmRecipientIdsSet[i].(string)
				createAlarmRuleRequest.AlarmRecipientIds = append(createAlarmRuleRequest.AlarmRecipientIds, &alarmRecipientIds)
			}
		}

		if v, ok := dMap["ext_info"]; ok {
			createAlarmRuleRequest.ExtInfo = helper.String(v.(string))
		}

		request.BaselineCreateAlarmRuleRequest = &createAlarmRuleRequest
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateBaseline(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata baseline failed, reason:%+v", logId, err)
		return err
	}

	baselineInt := *response.Response.Data.BaselineId
	baselineId = helper.Int64ToStr(baselineInt)
	d.SetId(strings.Join([]string{projectId, baselineId}, FILED_SP))

	return resourceTencentCloudWedataBaselineRead(d, meta)
}

func resourceTencentCloudWedataBaselineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	baselineId := idSplit[1]

	baseline, err := service.DescribeWedataBaselineById(ctx, projectId, baselineId)
	if err != nil {
		return err
	}

	if baseline == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataBaseline` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("baseline_id", baselineId)

	if baseline.BaselineDto.BaselineName != nil {
		_ = d.Set("baseline_name", baseline.BaselineDto.BaselineName)
	}

	if baseline.BaselineDto.BaselineType != nil {
		_ = d.Set("baseline_type", baseline.BaselineDto.BaselineType)
	}

	if baseline.BaselineDto.UserUin != nil {
		_ = d.Set("create_uin", baseline.BaselineDto.UserUin)
	}

	if baseline.BaselineDto.UserName != nil {
		_ = d.Set("create_name", baseline.BaselineDto.UserName)
	}

	if baseline.BaselineDto.InChargeUin != nil {
		_ = d.Set("in_charge_uin", baseline.BaselineDto.InChargeUin)
	}

	if baseline.BaselineDto.InChargeName != nil {
		_ = d.Set("in_charge_name", baseline.BaselineDto.InChargeName)
	}

	if baseline.BaselineDto.PromiseTasks != nil {
		promiseTasksList := []interface{}{}
		for _, promiseTasks := range baseline.BaselineDto.PromiseTasks {
			promiseTasksMap := map[string]interface{}{}

			if promiseTasks.ProjectId != nil {
				promiseTasksMap["project_id"] = promiseTasks.ProjectId
			}

			if promiseTasks.TaskName != nil {
				promiseTasksMap["task_name"] = promiseTasks.TaskName
			}

			if promiseTasks.TaskId != nil {
				promiseTasksMap["task_id"] = promiseTasks.TaskId
			}

			if promiseTasks.TaskCycle != nil {
				promiseTasksMap["task_cycle"] = promiseTasks.TaskCycle
			}

			if promiseTasks.WorkflowName != nil {
				promiseTasksMap["workflow_name"] = promiseTasks.WorkflowName
			}

			if promiseTasks.WorkflowId != nil {
				promiseTasksMap["workflow_id"] = promiseTasks.WorkflowId
			}

			if promiseTasks.TaskInChargeName != nil {
				promiseTasksMap["task_in_charge_name"] = promiseTasks.TaskInChargeName
			}

			if promiseTasks.TaskInChargeUin != nil {
				promiseTasksMap["task_in_charge_uin"] = promiseTasks.TaskInChargeUin
			}

			promiseTasksList = append(promiseTasksList, promiseTasksMap)
		}

		_ = d.Set("promise_tasks", promiseTasksList)

	}

	if baseline.BaselineDto.PromiseTime != nil {
		_ = d.Set("promise_time", baseline.BaselineDto.PromiseTime)
	}

	if baseline.BaselineDto.WarningMargin != nil {
		_ = d.Set("warning_margin", baseline.BaselineDto.WarningMargin)
	}

	if baseline.IsNewAlarm != nil {
		_ = d.Set("is_new_alarm", baseline.IsNewAlarm)
	}

	if baseline.BaselineDto.AlarmRule != nil {
		alarmRuleDtoMap := map[string]interface{}{}

		if baseline.BaselineDto.AlarmRule.AlarmRuleId != nil {
			alarmRuleDtoMap["alarm_rule_id"] = baseline.BaselineDto.AlarmRule.AlarmRuleId
		}

		if baseline.BaselineDto.AlarmRule.AlarmLevelType != nil {
			alarmRuleDtoMap["alarm_level_type"] = baseline.BaselineDto.AlarmRule.AlarmLevelType
		}

		_ = d.Set("alarm_rule_dto", []interface{}{alarmRuleDtoMap})
	}

	if baseline.BaselineCreateAlarmRuleRequest != nil {
		baselineCreateAlarmRuleRequestMap := map[string]interface{}{}

		if baseline.BaselineCreateAlarmRuleRequest.ProjectId != nil {
			baselineCreateAlarmRuleRequestMap["project_id"] = baseline.BaselineCreateAlarmRuleRequest.ProjectId
		}

		if baseline.BaselineCreateAlarmRuleRequest.CreatorId != nil {
			baselineCreateAlarmRuleRequestMap["creator_id"] = baseline.BaselineCreateAlarmRuleRequest.CreatorId
		}

		if baseline.BaselineCreateAlarmRuleRequest.Creator != nil {
			baselineCreateAlarmRuleRequestMap["creator"] = baseline.BaselineCreateAlarmRuleRequest.Creator
		}

		if baseline.BaselineCreateAlarmRuleRequest.RuleName != nil {
			baselineCreateAlarmRuleRequestMap["rule_name"] = baseline.BaselineCreateAlarmRuleRequest.RuleName
		}

		if baseline.BaselineCreateAlarmRuleRequest.MonitorType != nil {
			baselineCreateAlarmRuleRequestMap["monitor_type"] = baseline.BaselineCreateAlarmRuleRequest.MonitorType
		}

		if baseline.BaselineCreateAlarmRuleRequest.MonitorObjectIds != nil {
			baselineCreateAlarmRuleRequestMap["monitor_object_ids"] = baseline.BaselineCreateAlarmRuleRequest.MonitorObjectIds
		}

		if baseline.BaselineCreateAlarmRuleRequest.AlarmTypes != nil {
			baselineCreateAlarmRuleRequestMap["alarm_types"] = baseline.BaselineCreateAlarmRuleRequest.AlarmTypes
		}

		if baseline.BaselineCreateAlarmRuleRequest.AlarmLevel != nil {
			baselineCreateAlarmRuleRequestMap["alarm_level"] = baseline.BaselineCreateAlarmRuleRequest.AlarmLevel
		}

		if baseline.BaselineCreateAlarmRuleRequest.AlarmWays != nil {
			baselineCreateAlarmRuleRequestMap["alarm_ways"] = baseline.BaselineCreateAlarmRuleRequest.AlarmWays
		}

		if baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientType != nil {
			baselineCreateAlarmRuleRequestMap["alarm_recipient_type"] = baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientType
		}

		if baseline.BaselineCreateAlarmRuleRequest.AlarmRecipients != nil {
			baselineCreateAlarmRuleRequestMap["alarm_recipients"] = baseline.BaselineCreateAlarmRuleRequest.AlarmRecipients
		}

		if baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientIds != nil {
			baselineCreateAlarmRuleRequestMap["alarm_recipient_ids"] = baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientIds
		}

		if baseline.BaselineCreateAlarmRuleRequest.ExtInfo != nil {
			baselineCreateAlarmRuleRequestMap["ext_info"] = baseline.BaselineCreateAlarmRuleRequest.ExtInfo
		}

		_ = d.Set("baseline_create_alarm_rule_request", []interface{}{baselineCreateAlarmRuleRequestMap})
	}

	return nil
}

func resourceTencentCloudWedataBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = wedata.NewEditBaselineRequest()
	)

	immutableArgs := []string{"project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	baselineId := idSplit[1]

	request.ProjectId = &projectId
	request.BaselineId = &baselineId

	if v, ok := d.GetOk("baseline_name"); ok {
		request.BaselineName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("baseline_type"); ok {
		request.BaselineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("in_charge_uin"); ok {
		request.InChargeUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("in_charge_name"); ok {
		request.InChargeName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("promise_tasks"); ok {
		for _, item := range v.([]interface{}) {
			baselineTaskInfo := wedata.BaselineTaskInfo{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["project_id"]; ok {
				baselineTaskInfo.ProjectId = helper.String(v.(string))
			}

			if v, ok := dMap["task_name"]; ok {
				baselineTaskInfo.TaskName = helper.String(v.(string))
			}

			if v, ok := dMap["task_id"]; ok {
				baselineTaskInfo.TaskId = helper.String(v.(string))
			}

			if v, ok := dMap["task_cycle"]; ok {
				baselineTaskInfo.TaskCycle = helper.String(v.(string))
			}

			if v, ok := dMap["workflow_name"]; ok {
				baselineTaskInfo.WorkflowName = helper.String(v.(string))
			}

			if v, ok := dMap["workflow_id"]; ok {
				baselineTaskInfo.WorkflowId = helper.String(v.(string))
			}

			if v, ok := dMap["task_in_charge_name"]; ok {
				baselineTaskInfo.TaskInChargeName = helper.String(v.(string))
			}

			if v, ok := dMap["task_in_charge_uin"]; ok {
				baselineTaskInfo.TaskInChargeUin = helper.String(v.(string))
			}

			request.PromiseTasks = append(request.PromiseTasks, &baselineTaskInfo)
		}
	}

	if v, ok := d.GetOk("promise_time"); ok {
		request.PromiseTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("warning_margin"); ok {
		request.WarningMargin = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("is_new_alarm"); ok {
		request.IsNewAlarm = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "alarm_rule_dto"); ok {
		alarmRuleDto := wedata.AlarmRuleDto{}
		if v, ok := dMap["alarm_rule_id"]; ok {
			alarmRuleDto.AlarmRuleId = helper.String(v.(string))
		}

		if v, ok := dMap["alarm_level_type"]; ok {
			alarmRuleDto.AlarmLevelType = helper.String(v.(string))
		}

		request.AlarmRuleDto = &alarmRuleDto
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().EditBaseline(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update wedata baseline failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataBaselineRead(d, meta)
}

func resourceTencentCloudWedataBaselineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	baselineId := idSplit[1]

	if err := service.DeleteWedataBaselineById(ctx, projectId, baselineId); err != nil {
		return err
	}

	return nil
}
