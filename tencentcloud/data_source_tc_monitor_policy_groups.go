/*
Use this data source to query monitor policy groups (There is a lot of data and it is recommended to output to a file)

Example Usage

```hcl
data "tencentcloud_monitor_policy_groups" "groups" {
  policy_view_names = ["REDIS-CLUSTER", "cvm_device"]
}

data "tencentcloud_monitor_policy_groups" "name" {
  name = "test"
}
```

*/
package tencentcloud

import (
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentMonitorPolicyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorPolicyGroupsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy group name for query.",
			},
			"policy_view_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The policy view for query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list policy groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy group id.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy group name.",
						},
						"is_open": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether open or not.",
						},
						"policy_view_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy group view name.",
						},
						"last_edit_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recently edited user uin.",
						},
						"use_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of instances of policy group bindings.",
						},
						"no_shielded_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of unmasked instances of policy group bindings.",
						},
						"is_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If is default policy group or not,0 represents the non-default policy, and 1 represents the default policy.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of threshold rules. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_show_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of this metric.",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Data aggregation cycle (unit second).",
									},
									"metric_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The id of this metric.",
									},
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Threshold rule id.",
									},
									"metric_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unit of this metric.",
									},
									"alarm_notify_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending convergence type. 0 continuous alarm, 1 index alarm.",
									},
									"alarm_notify_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending cycle per second.<0 does not fire, 0 only fires once, and >0 fires every triggerTime second.",
									},
									"calc_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Compare type, 1 means more than, 2  means greater than or equal, 3 means less than, 4 means less than or equal to, 5 means equal, 6 means not equal, 7 means days rose, 8 means days fell, 9 means weeks rose, 10  means weeks fell, 11 means period rise, 12 means period fell.",
									},
									"calc_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Threshold value.",
									},
									"continue_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "How long does the triggering rule last (per second).",
									},
								},
							},
						},
						"event_conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of event rules. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The id of this event metric.",
									},
									"event_show_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of this event metric.",
									},
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Threshold rule id.",
									},
									"alarm_notify_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending convergence type. 0 continuous alarm, 1 index alarm.",
									},
									"alarm_notify_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending cycle per second.<0 does not fire, 0 only fires once, and >0 fires every triggerTime second.",
									},
								},
							},
						},
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
										Description: "Alarm receive group id list.",
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
										Description: "Alarm period start time.Range [0,86399], which removes the date after it is converted to Beijing time as a Unix timestamp, for example 7200 means '10:0:0'.",
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
										Description: `Method of warning notification.Optional ` + helper.SliceFieldSerialize(monitorNotifyWays) + `.`,
									},
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Receive type. Optional 'group' or 'user'.",
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
										Description: "Do need a telephone alarm contact prompt.You don't need 0, you need 1.",
									},
									"send_for": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: `Telephone warning time.Option "OCCUR","RECOVER".`,
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
						"can_set_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it can be set as the default policy.",
						},
						"parent_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parent policy group id.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy group remarks.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project id to which the policy group belongs.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy group update timestamp.",
						},
						"insert_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy group create timestamp.",
						},
					},
				},
			},
		},
	}
}
func dataSourceTencentMonitorPolicyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_policy_groups.read")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribePolicyGroupListRequest()
		response       *monitor.DescribePolicyGroupListResponse
		err            error

		name            = d.Get("name").(string)
		policyViewNames = helper.InterfacesStrings(d.Get("policy_view_names").([]interface{}))

		list            = make([]interface{}, 0, 100)
		offset    int64 = 0
		limit     int64 = 20
		groupList       = make([]*monitor.DescribePolicyGroupListGroup, 0, 10)
		finish    bool
	)

	request.Module = helper.String("monitor")
	request.Offset = &offset
	request.Limit = &limit

	if name != "" {
		request.Like = &name
	}

	if len(policyViewNames) != 0 {
		request.ViewNames = helper.Strings(policyViewNames)
	}

	for {
		if finish {
			break
		}
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if response, err = monitorService.client.UseMonitorClient().DescribePolicyGroupList(request); err != nil {
				return retryError(err, InternalError)
			}
			groupList = append(groupList, response.Response.GroupList...)
			if len(response.Response.GroupList) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			return err
		}
		offset = offset + limit
	}

	for _, group := range groupList {
		var listItem = map[string]interface{}{}
		listItem["group_id"] = group.GroupId
		listItem["group_name"] = group.GroupName
		listItem["is_open"] = group.IsOpen
		listItem["policy_view_name"] = group.ViewName
		listItem["last_edit_uin"] = group.LastEditUin
		listItem["use_sum"] = group.UseSum
		listItem["no_shielded_sum"] = group.NoShieldedSum
		listItem["is_default"] = group.IsDefault
		listItem["can_set_default"] = group.CanSetDefault
		listItem["parent_group_id"] = group.ParentGroupId
		listItem["remark"] = group.Remark
		listItem["project_id"] = group.ProjectId
		listItem["update_time"] = group.UpdateTime
		listItem["insert_time"] = group.InsertTime

		conditions := make([]interface{}, 0, 100)
		for _, item := range group.Conditions {
			conditions = append(conditions, map[string]interface{}{
				"metric_show_name":    item.MetricShowName,
				"period":              item.Period,
				"metric_id":           item.MetricId,
				"rule_id":             item.RuleId,
				"metric_unit":         item.Unit,
				"alarm_notify_type":   item.AlarmNotifyType,
				"alarm_notify_period": item.AlarmNotifyPeriod,
				"calc_type":           item.CalcType,
				"calc_value":          item.CalcValue,
				"continue_time":       item.ContinueTime,
			})
		}
		listItem["conditions"] = conditions

		eventConditions := make([]interface{}, 0, 100)
		for _, item := range group.EventConditions {
			eventConditions = append(eventConditions, map[string]interface{}{
				"event_id":            item.EventId,
				"event_show_name":     item.EventShowName,
				"rule_id":             item.RuleId,
				"alarm_notify_type":   item.AlarmNotifyType,
				"alarm_notify_period": item.AlarmNotifyPeriod,
			})
		}
		listItem["event_conditions"] = eventConditions

		receivers := make([]interface{}, 0, 100)
		for _, item := range group.ReceiverInfos {

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
		listItem["receivers"] = receivers

		list = append(list, listItem)
	}
	if err = d.Set("list", list); err != nil {
		return err
	}

	md := md5.New()
	_, _ = md.Write([]byte(request.ToJsonString()))
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), list)
	}
	return nil
}
