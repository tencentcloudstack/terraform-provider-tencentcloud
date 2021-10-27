/*
Provides a policy group resource for monitor.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_alarm_policy.

Example Usage

```hcl
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "nice_group"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 1
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
  conditions {
    metric_id           = 30
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 2
    calc_value          = 30
    calc_period         = 300
    continue_period     = 2
  }
  event_conditions {
    event_id            = 39
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
  event_conditions {
    event_id            = 40
    alarm_notify_type   = 0
    alarm_notify_period = 300
  }
}
```
Import

Policy group instance can be imported, e.g.

```
$ terraform import tencentcloud_monitor_policy_group.group group-id
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentMonitorPolicyGroup() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.59.18. Please use 'tencentcloud_alarm_policy' instead.",
		Create: resourceTencentMonitorPolicyGroupCreate,
		Read:   resourceTencentMonitorPolicyGroupRead,
		Update: resourceTencentMonitorPolicyGroupUpdate,
		Delete: resourceTencentMonitorPolicyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Policy group name, length should between 1 and 20.",
				ValidateFunc: validateStringLengthInRange(1, 20),
			},
			"policy_view_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Policy view name, eg:`cvm_device`,`BANDWIDTHPACKAGE`, refer to `data.tencentcloud_monitor_policy_conditions(policy_view_name)`.",
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(0, 100),
				Description:  "Policy group's remark information.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "The project id to which the policy group belongs, default is `0`.",
			},
			"is_union_rule": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "The and or relation of indicator alarm rule. Valid values: `0`, `1`. `0` represents or rule (if any rule is met, the alarm will be raised), `1` represents and rule (if all rules are met, the alarm will be raised).The default is 0.",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of threshold rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Id of the metric, refer to `data.tencentcloud_monitor_policy_conditions(metric_id)`.",
						},
						"alarm_notify_type": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedIntValue([]int{0, 1}),
							Description:  "Alarm sending convergence type. `0` continuous alarm, `1` index alarm.",
						},
						"alarm_notify_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm sending cycle per second. <0 does not fire, `0` only fires once, and >0 fires every triggerTime second.",
						},
						"calc_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateIntegerInRange(1, 12),
							Description:  "Compare type. Valid value ranges: [1~12]. `1` means more than, `2` means greater than or equal, `3` means less than, `4` means less than or equal to, `5` means equal, `6` means not equal, `7` means days rose, `8` means days fell, `9` means weeks rose, `10` means weeks fell, `11` means period rise, `12` means period fell, refer to `data.tencentcloud_monitor_policy_conditions(calc_type_keys)`.",
						},
						"calc_value": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "Threshold value, refer to `data.tencentcloud_monitor_policy_conditions(calc_value_*)`.",
						},
						"calc_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Data aggregation cycle (unit of second), if the metric has a default value can not be filled, refer to `data.tencentcloud_monitor_policy_conditions(period_keys)`.",
						},
						"continue_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The rule triggers an alert that lasts for several detection cycles, refer to `data.tencentcloud_monitor_policy_conditions(period_num_keys)`.",
						},
					},
				},
			},
			"event_conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of event rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The ID of this event metric, refer to `data.tencentcloud_monitor_policy_conditions(event_id).",
						},
						"alarm_notify_type": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateAllowedIntValue([]int{0, 1}),
							Description:  "Alarm sending convergence type. `0` continuous alarm, `1` index alarm.",
						},
						"alarm_notify_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm sending cycle per second. <0 does not fire, `0` only fires once, and >0 fires every triggerTime second.",
						},
					},
				},
			},
			// computed value
			"receivers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of receivers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"receiver_group_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Alarm receive group ID list.",
						},
						"receiver_user_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Alarm receiver id list.",
						},
						"uid_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "The phone alerts the receiver uid.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm period start time.Range [0,86400], which removes the date after it is converted to Beijing time as a Unix timestamp, for example 7200 means '10:0:0'.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "End of alarm period. Meaning with `start_time`.",
						},
						"notify_way": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Method of warning notification. Valid values: "SMS", "SITE", "EMAIL", "CALL", "WECHAT".`,
						},
						"receiver_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Receive type. Valid values: group, user. 'group' (receiving group) or 'user' (receiver).",
						},
						"round_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Telephone alarm number.",
						},
						"round_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Telephone alarm interval per round (seconds).",
						},
						"person_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Telephone warning to individual interval (seconds).",
						},
						"need_send_notice": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Do need a telephone alarm contact prompt. You don't need `0`, you need `1`.",
						},
						"send_for": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Telephone warning time. Valid values: "OCCUR","RECOVER".`,
						},
						"recover_notify": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Restore notification mode. Optional "SMS".`,
						},
						"receive_language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alert sending language.",
						},
					},
				},
			},
			"binding_objects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list binding objects(list only those in the `provider.region`). Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object unique id.",
						},
						"dimensions_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Represents a collection of dimensions of an object instance, json format.",
						},
						"is_shielded": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the object is shielded or not, 0 means unshielded and 1 means shielded.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the object is located.",
						},
					},
				},
			},
			"dimension_group": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A list of dimensions for this policy group.",
			},
			"support_regions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Support regions this policy group.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy group update time.",
			},
			"last_edit_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Recently edited user uin.",
			},
		},
	}
}
func resourceTencentMonitorPolicyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_policy_group.create")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewCreatePolicyGroupRequest()
	)

	request.GroupName = helper.String(d.Get("group_name").(string))
	request.ViewName = helper.String(d.Get("policy_view_name").(string))
	if iface, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(iface.(string))
	}
	request.IsUnionRule = helper.IntInt64(d.Get("is_union_rule").(int))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	request.Module = helper.String("monitor")

	if iface, ok := d.GetOk("conditions"); ok {
		request.Conditions = make([]*monitor.CreatePolicyGroupCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			createPolicyGroupCondition := monitor.CreatePolicyGroupCondition{}
			createPolicyGroupCondition.MetricId = helper.IntInt64(m["metric_id"].(int))
			createPolicyGroupCondition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			createPolicyGroupCondition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			if m["calc_type"] != nil {
				createPolicyGroupCondition.CalcType = helper.IntInt64(m["calc_type"].(int))
			}
			if m["calc_value"] != nil {
				createPolicyGroupCondition.CalcValue = helper.Float64(m["calc_value"].(float64))
			}
			if m["calc_period"] != nil {
				createPolicyGroupCondition.CalcPeriod = helper.IntInt64(m["calc_period"].(int))
			}
			if m["continue_period"] != nil {
				createPolicyGroupCondition.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			request.Conditions = append(request.Conditions, &createPolicyGroupCondition)
		}
	}

	if iface, ok := d.GetOk("event_conditions"); ok {
		request.EventConditions = make([]*monitor.CreatePolicyGroupEventCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			createPolicyGroupCondition := monitor.CreatePolicyGroupEventCondition{}
			createPolicyGroupCondition.EventId = helper.IntInt64(m["event_id"].(int))
			createPolicyGroupCondition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			createPolicyGroupCondition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			request.EventConditions = append(request.EventConditions, &createPolicyGroupCondition)
		}
	}

	var groupId *int64
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreatePolicyGroup(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		groupId = response.Response.GroupId

		return nil
	}); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", *groupId))
	return resourceTencentMonitorPolicyGroupRead(d, meta)
}

func resourceTencentMonitorPolicyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_policy_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribePolicyGroupInfoRequest()
		response       *monitor.DescribePolicyGroupInfoResponse
	)

	groupId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("id [%s] is broken", d.Id())
	}

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}
	if info == nil {
		d.SetId("")
		return nil
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = monitorService.client.UseMonitorClient().DescribePolicyGroupInfo(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if response == nil {
		d.SetId("")
		return nil
	}
	var errs []error
	errs = append(errs,
		d.Set("group_name", response.Response.GroupName),
		d.Set("policy_view_name", response.Response.ViewName),
		d.Set("remark", response.Response.Remark),
		d.Set("is_union_rule", response.Response.IsUnionRule),
		d.Set("project_id", response.Response.ProjectId),
	)

	var conditions = make([]interface{}, 0, 100)
	for _, condition := range response.Response.ConditionsConfig {
		m := map[string]interface{}{}
		m["metric_id"] = condition.MetricId
		m["alarm_notify_type"] = condition.AlarmNotifyType
		m["alarm_notify_period"] = condition.AlarmNotifyPeriod
		m["calc_type"] = condition.CalcType
		m["calc_value"] = condition.CalcValue
		m["calc_period"] = condition.Period
		m["continue_period"] = condition.ContinueTime
		conditions = append(conditions, m)
	}
	errs = append(errs, d.Set("conditions", conditions))

	var eventConditions = make([]interface{}, 0, 100)
	for _, condition := range response.Response.EventConfig {
		m := map[string]interface{}{}
		m["event_id"] = condition.EventId
		m["alarm_notify_type"] = condition.AlarmNotifyType
		m["alarm_notify_period"] = condition.AlarmNotifyPeriod
		eventConditions = append(eventConditions, m)
	}
	errs = append(errs, d.Set("event_conditions", eventConditions))

	receivers := make([]interface{}, 0, 100)
	for _, item := range response.Response.ReceiverInfos {

		receiver := map[string]interface{}{
			"start_time":       item.StartTime,
			"end_time":         item.EndTime,
			"receiver_type":    item.ReceiverType,
			"round_number":     item.RoundNumber,
			"round_interval":   item.RoundInterval,
			"person_interval":  item.PersonInterval,
			"need_send_notice": item.NeedSendNotice,
			"receive_language": item.ReceiveLanguage,
		}
		{
			slice := make([]int64, 0, len(item.ReceiverGroupList))
			for _, value := range item.ReceiverGroupList {
				slice = append(slice, *value)
			}
			receiver["receiver_group_list"] = slice
		}

		{
			slice := make([]int64, 0, len(item.ReceiverUserList))
			for _, value := range item.ReceiverUserList {
				slice = append(slice, *value)
			}
			receiver["receiver_user_list"] = slice
		}

		{
			slice := make([]int64, 0, len(item.UidList))
			for _, value := range item.UidList {
				slice = append(slice, *value)
			}
			receiver["uid_list"] = slice
		}

		{
			slice := make([]string, 0, len(item.NotifyWay))
			for _, value := range item.NotifyWay {
				slice = append(slice, *value)
			}
			receiver["notify_way"] = slice
		}

		{
			slice := make([]string, 0, len(item.SendFor))
			for _, value := range item.SendFor {
				slice = append(slice, *value)
			}
			receiver["send_for"] = slice
		}

		{
			slice := make([]string, 0, len(item.RecoverNotify))
			for _, value := range item.RecoverNotify {
				slice = append(slice, *value)
			}
			receiver["recover_notify"] = slice
		}
		receivers = append(receivers, receiver)
	}
	errs = append(errs, d.Set("receivers", receivers))

	errs = append(errs,
		d.Set("group_name", response.Response.GroupName),
		d.Set("policy_view_name", response.Response.ViewName),
		d.Set("remark", response.Response.Remark),
		d.Set("is_union_rule", response.Response.IsUnionRule),
		d.Set("support_regions", response.Response.Region),
		d.Set("dimension_group", response.Response.DimensionGroup),
		d.Set("update_time", response.Response.UpdateTime),
		d.Set("last_edit_uin", response.Response.LastEditUin),
	)

	objects, err := monitorService.DescribeBindingPolicyObjectList(ctx, groupId)
	if err != nil {
		return err
	}
	bindingObjects := make([]interface{}, 0, len(objects))

	for _, event := range objects {
		var listItem = map[string]interface{}{}
		listItem["region"] = event.Region
		listItem["unique_id"] = event.UniqueId
		listItem["dimensions_json"] = event.Dimensions
		listItem["is_shielded"] = event.IsShielded
		listItem["region"] = event.Region
		bindingObjects = append(bindingObjects, listItem)
	}
	errs = append(errs, d.Set("binding_objects", bindingObjects))

	if len(errs) > 0 {
		return errs[0]
	} else {
		return nil
	}
}

func resourceTencentMonitorPolicyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_policy_group.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewModifyPolicyGroupRequest()
	)
	groupId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("id [%s] is broken", d.Id())
	}

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}
	if info == nil {
		d.SetId("")
		return nil
	}
	request.GroupId = &groupId
	request.GroupName = helper.String(d.Get("group_name").(string))
	request.ViewName = helper.String(d.Get("policy_view_name").(string))
	request.IsUnionRule = helper.IntInt64(d.Get("is_union_rule").(int))
	request.Module = helper.String("monitor")

	if iface, ok := d.GetOk("conditions"); ok {
		request.Conditions = make([]*monitor.ModifyPolicyGroupCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			condition := monitor.ModifyPolicyGroupCondition{}
			condition.MetricId = helper.IntInt64(m["metric_id"].(int))
			condition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			condition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			if m["calc_type"] != nil {
				condition.CalcType = helper.IntInt64(m["calc_type"].(int))
			}
			if m["calc_value"] != nil {
				condition.CalcValue = helper.String(fmt.Sprintf("%f", m["calc_value"].(float64)))
			}
			if m["calc_period"] != nil {
				condition.CalcPeriod = helper.IntInt64(m["calc_period"].(int))
			}
			if m["continue_period"] != nil {
				condition.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			request.Conditions = append(request.Conditions, &condition)
		}
	}

	if iface, ok := d.GetOk("event_conditions"); ok {
		request.EventConditions = make([]*monitor.ModifyPolicyGroupEventCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			condition := monitor.ModifyPolicyGroupEventCondition{}
			condition.EventId = helper.IntInt64(m["event_id"].(int))
			condition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			condition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			request.EventConditions = append(request.EventConditions, &condition)
		}
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := monitorService.client.UseMonitorClient().ModifyPolicyGroup(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return resourceTencentMonitorPolicyGroupRead(d, meta)
}

func resourceTencentMonitorPolicyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_policy_group.delete")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDeletePolicyGroupRequest()
	)

	groupId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("id [%s] is broken", d.Id())
	}
	request.GroupId = []*int64{&groupId}
	request.Module = helper.String("monitor")

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().DeletePolicyGroup(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
