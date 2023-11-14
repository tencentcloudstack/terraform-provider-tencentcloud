/*
Provides a resource to create a wedata baseline_baseline

Example Usage

```hcl
resource "tencentcloud_wedata_baseline_baseline" "baseline_baseline" {
  project_id = ""
  baseline_name = ""
  baseline_type = ""
  create_uin = ""
  create_name = ""
  in_charge_uin = ""
  in_charge_name = ""
  promise_tasks {
		project_id = ""
		task_name = ""
		task_id = ""
		task_cycle = ""
		workflow_name = ""
		workflow_id = ""
		task_in_charge_name = ""
		task_in_charge_uin = ""

  }
  promise_time = ""
  warning_margin =
  is_new_alarm =
  alarm_rule_dto {
		alarm_rule_id = ""
		alarm_level_type = ""

  }
  baseline_create_alarm_rule_request {
		project_id = ""
		creator_id = ""
		creator = ""
		rule_name = ""
		monitor_type =
		monitor_object_ids =
		alarm_types =
		alarm_level =
		alarm_ways =
		alarm_recipient_type =
		alarm_recipients =
		alarm_recipient_ids =
		ext_info = ""

  }
}
```

Import

wedata baseline_baseline can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_baseline_baseline.baseline_baseline baseline_baseline_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWedataBaseline_baseline() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataBaseline_baselineCreate,
		Read:   resourceTencentCloudWedataBaseline_baselineRead,
		Update: resourceTencentCloudWedataBaseline_baselineUpdate,
		Delete: resourceTencentCloudWedataBaseline_baselineDelete,
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
				Description: "1.",
			},

			"alarm_rule_dto": {
				Optional:    true,
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
							Optional:    true,
							Description: "Project NameNote: This field may return null, indicating no valid value.",
						},
						"creator_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creator NameNote: This field may return null, indicating no valid value.",
						},
						"creator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creator UINNote: This field may return null, indicating no valid value.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule NameNote: This field may return null, indicating no valid value.",
						},
						"monitor_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Monitoring Type, 1. Task, 2. Workflow, 3. Project, 4. Baseline (default is 1. Task)Note: This field may return null, indicating no valid value.",
						},
						"monitor_object_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Monitoring ObjectsNote: This field may return null, indicating no valid value.",
						},
						"alarm_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Alarm Types, 1. Failure Alarm, 2. Timeout Alarm, 3. Success Alarm, 4. Baseline Violation, 5. Baseline Warning, 6. Baseline Task Failure (default is 1. Failure Alarm)Note: This field may return null, indicating no valid value.",
						},
						"alarm_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm Level, 1. Normal, 2. Important, 3. Urgent (default is 1. Normal)Note: This field may return null, indicating no valid value.",
						},
						"alarm_ways": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Alarm Methods, 1. Email, 2. SMS, 3. WeChat, 4. Voice, 5. Enterprise WeChat, 6. HTTP, 7. Enterprise WeChat Group; Alarm method code list (default is 1. Email)Note: This field may return null, indicating no valid value.",
						},
						"alarm_recipient_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm Recipient Type: 1. Specified Personnel, 2. Task Owner, 3. Duty Roster (default is 1. Specified Personnel)Note: This field may return null, indicating no valid value.",
						},
						"alarm_recipients": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Alarm RecipientsNote: This field may return null, indicating no valid value.",
						},
						"alarm_recipient_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Alarm Recipient IDsNote: This field may return null, indicating no valid value.",
						},
						"ext_info": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Extended Information, 1. Estimated Runtime (default), 2. Estimated Completion Time, 3. Estimated Scheduling Time, 4. Incomplete within the Cycle; Value Types: 1. Specified Value, 2. Historical Average (default is 1. Specified Value)Note: This field may return null, indicating no valid value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataBaseline_baselineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline_baseline.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = wedata.NewCreateBaselineRequest()
		response   = wedata.NewCreateBaselineResponse()
		baselineId int
	)
	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
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
		log.Printf("[CRITAL]%s create wedata baseline_baseline failed, reason:%+v", logId, err)
		return err
	}

	baselineId = *response.Response.BaselineId
	d.SetId(helper.Int64ToStr(baselineId))

	return resourceTencentCloudWedataBaseline_baselineRead(d, meta)
}

func resourceTencentCloudWedataBaseline_baselineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline_baseline.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	baseline_baselineId := d.Id()

	baseline_baseline, err := service.DescribeWedataBaseline_baselineById(ctx, baselineId)
	if err != nil {
		return err
	}

	if baseline_baseline == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataBaseline_baseline` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if baseline_baseline.ProjectId != nil {
		_ = d.Set("project_id", baseline_baseline.ProjectId)
	}

	if baseline_baseline.BaselineName != nil {
		_ = d.Set("baseline_name", baseline_baseline.BaselineName)
	}

	if baseline_baseline.BaselineType != nil {
		_ = d.Set("baseline_type", baseline_baseline.BaselineType)
	}

	if baseline_baseline.CreateUin != nil {
		_ = d.Set("create_uin", baseline_baseline.CreateUin)
	}

	if baseline_baseline.CreateName != nil {
		_ = d.Set("create_name", baseline_baseline.CreateName)
	}

	if baseline_baseline.InChargeUin != nil {
		_ = d.Set("in_charge_uin", baseline_baseline.InChargeUin)
	}

	if baseline_baseline.InChargeName != nil {
		_ = d.Set("in_charge_name", baseline_baseline.InChargeName)
	}

	if baseline_baseline.PromiseTasks != nil {
		promiseTasksList := []interface{}{}
		for _, promiseTasks := range baseline_baseline.PromiseTasks {
			promiseTasksMap := map[string]interface{}{}

			if baseline_baseline.PromiseTasks.ProjectId != nil {
				promiseTasksMap["project_id"] = baseline_baseline.PromiseTasks.ProjectId
			}

			if baseline_baseline.PromiseTasks.TaskName != nil {
				promiseTasksMap["task_name"] = baseline_baseline.PromiseTasks.TaskName
			}

			if baseline_baseline.PromiseTasks.TaskId != nil {
				promiseTasksMap["task_id"] = baseline_baseline.PromiseTasks.TaskId
			}

			if baseline_baseline.PromiseTasks.TaskCycle != nil {
				promiseTasksMap["task_cycle"] = baseline_baseline.PromiseTasks.TaskCycle
			}

			if baseline_baseline.PromiseTasks.WorkflowName != nil {
				promiseTasksMap["workflow_name"] = baseline_baseline.PromiseTasks.WorkflowName
			}

			if baseline_baseline.PromiseTasks.WorkflowId != nil {
				promiseTasksMap["workflow_id"] = baseline_baseline.PromiseTasks.WorkflowId
			}

			if baseline_baseline.PromiseTasks.TaskInChargeName != nil {
				promiseTasksMap["task_in_charge_name"] = baseline_baseline.PromiseTasks.TaskInChargeName
			}

			if baseline_baseline.PromiseTasks.TaskInChargeUin != nil {
				promiseTasksMap["task_in_charge_uin"] = baseline_baseline.PromiseTasks.TaskInChargeUin
			}

			promiseTasksList = append(promiseTasksList, promiseTasksMap)
		}

		_ = d.Set("promise_tasks", promiseTasksList)

	}

	if baseline_baseline.PromiseTime != nil {
		_ = d.Set("promise_time", baseline_baseline.PromiseTime)
	}

	if baseline_baseline.WarningMargin != nil {
		_ = d.Set("warning_margin", baseline_baseline.WarningMargin)
	}

	if baseline_baseline.IsNewAlarm != nil {
		_ = d.Set("is_new_alarm", baseline_baseline.IsNewAlarm)
	}

	if baseline_baseline.AlarmRuleDto != nil {
		alarmRuleDtoMap := map[string]interface{}{}

		if baseline_baseline.AlarmRuleDto.AlarmRuleId != nil {
			alarmRuleDtoMap["alarm_rule_id"] = baseline_baseline.AlarmRuleDto.AlarmRuleId
		}

		if baseline_baseline.AlarmRuleDto.AlarmLevelType != nil {
			alarmRuleDtoMap["alarm_level_type"] = baseline_baseline.AlarmRuleDto.AlarmLevelType
		}

		_ = d.Set("alarm_rule_dto", []interface{}{alarmRuleDtoMap})
	}

	if baseline_baseline.BaselineCreateAlarmRuleRequest != nil {
		baselineCreateAlarmRuleRequestMap := map[string]interface{}{}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.ProjectId != nil {
			baselineCreateAlarmRuleRequestMap["project_id"] = baseline_baseline.BaselineCreateAlarmRuleRequest.ProjectId
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.CreatorId != nil {
			baselineCreateAlarmRuleRequestMap["creator_id"] = baseline_baseline.BaselineCreateAlarmRuleRequest.CreatorId
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.Creator != nil {
			baselineCreateAlarmRuleRequestMap["creator"] = baseline_baseline.BaselineCreateAlarmRuleRequest.Creator
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.RuleName != nil {
			baselineCreateAlarmRuleRequestMap["rule_name"] = baseline_baseline.BaselineCreateAlarmRuleRequest.RuleName
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.MonitorType != nil {
			baselineCreateAlarmRuleRequestMap["monitor_type"] = baseline_baseline.BaselineCreateAlarmRuleRequest.MonitorType
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.MonitorObjectIds != nil {
			baselineCreateAlarmRuleRequestMap["monitor_object_ids"] = baseline_baseline.BaselineCreateAlarmRuleRequest.MonitorObjectIds
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmTypes != nil {
			baselineCreateAlarmRuleRequestMap["alarm_types"] = baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmTypes
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmLevel != nil {
			baselineCreateAlarmRuleRequestMap["alarm_level"] = baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmLevel
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmWays != nil {
			baselineCreateAlarmRuleRequestMap["alarm_ways"] = baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmWays
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientType != nil {
			baselineCreateAlarmRuleRequestMap["alarm_recipient_type"] = baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientType
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmRecipients != nil {
			baselineCreateAlarmRuleRequestMap["alarm_recipients"] = baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmRecipients
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientIds != nil {
			baselineCreateAlarmRuleRequestMap["alarm_recipient_ids"] = baseline_baseline.BaselineCreateAlarmRuleRequest.AlarmRecipientIds
		}

		if baseline_baseline.BaselineCreateAlarmRuleRequest.ExtInfo != nil {
			baselineCreateAlarmRuleRequestMap["ext_info"] = baseline_baseline.BaselineCreateAlarmRuleRequest.ExtInfo
		}

		_ = d.Set("baseline_create_alarm_rule_request", []interface{}{baselineCreateAlarmRuleRequestMap})
	}

	return nil
}

func resourceTencentCloudWedataBaseline_baselineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline_baseline.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewEditBaselineRequest()

	baseline_baselineId := d.Id()

	request.BaselineId = &baselineId

	immutableArgs := []string{"project_id", "baseline_name", "baseline_type", "create_uin", "create_name", "in_charge_uin", "in_charge_name", "promise_tasks", "promise_time", "warning_margin", "is_new_alarm", "alarm_rule_dto", "baseline_create_alarm_rule_request"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("baseline_name") {
		if v, ok := d.GetOk("baseline_name"); ok {
			request.BaselineName = helper.String(v.(string))
		}
	}

	if d.HasChange("baseline_type") {
		if v, ok := d.GetOk("baseline_type"); ok {
			request.BaselineType = helper.String(v.(string))
		}
	}

	if d.HasChange("in_charge_uin") {
		if v, ok := d.GetOk("in_charge_uin"); ok {
			request.InChargeUin = helper.String(v.(string))
		}
	}

	if d.HasChange("in_charge_name") {
		if v, ok := d.GetOk("in_charge_name"); ok {
			request.InChargeName = helper.String(v.(string))
		}
	}

	if d.HasChange("promise_tasks") {
		if v, ok := d.GetOk("promise_tasks"); ok {
			for _, item := range v.([]interface{}) {
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
	}

	if d.HasChange("promise_time") {
		if v, ok := d.GetOk("promise_time"); ok {
			request.PromiseTime = helper.String(v.(string))
		}
	}

	if d.HasChange("warning_margin") {
		if v, ok := d.GetOkExists("warning_margin"); ok {
			request.WarningMargin = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("is_new_alarm") {
		if v, ok := d.GetOkExists("is_new_alarm"); ok {
			request.IsNewAlarm = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("alarm_rule_dto") {
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
		log.Printf("[CRITAL]%s update wedata baseline_baseline failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataBaseline_baselineRead(d, meta)
}

func resourceTencentCloudWedataBaseline_baselineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_baseline_baseline.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	baseline_baselineId := d.Id()

	if err := service.DeleteWedataBaseline_baselineById(ctx, baselineId); err != nil {
		return err
	}

	return nil
}
