package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTkeTmpAlertPolicyRead,
		Create: resourceTencentCloudTkeTmpAlertPolicyCreate,
		Update: resourceTencentCloudTkeTmpAlertPolicyUpdate,
		Delete: resourceTencentCloudTkeTmpAlertPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Id.",
			},

			"alert_rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Alarm notification channels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Policy name.",
						},
						"rules": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "A list of rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule name.",
									},
									"rule": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Prometheus statement.",
									},
									"labels": {
										Required:    true,
										Description: "Extra labels.",
										Type:        schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Name of map.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Value of map.",
												},
											},
										},
									},
									"template": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Alert sending template.",
									},
									"for": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time of duration.",
									},
									"describe": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A description of the rule.",
									},
									"annotations": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Refer to annotations in prometheus rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Name of map.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Value of map.",
												},
											},
										},
									},
									"rule_state": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Alarm rule status.",
									},
								},
							},
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm policy ID. Note: This field may return null, indicating that a valid value could not be retrieved.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If the alarm is sent from a template, the TemplateId is the template id.",
						},
						"notification": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Alarm channels, which may be returned using null in the template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether it is enabled.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The channel type, which defaults to amp, supports the following `amp`, `webhook`, `alertmanager`.",
									},
									"web_hook": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If Type is webhook, the field is required. Note: This field may return null, indicating that a valid value could not be retrieved.",
									},
									"alert_manager": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "If Type is alertmanager, the field is required. Note: This field may return null, indicating that a valid value could not be retrieved..",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Alertmanager url.",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Alertmanager is deployed in the cluster type. Note: This field may return null, indicating that a valid value could not be retrieved.",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The ID of the cluster where the alertmanager is deployed. Note: This field may return null, indicating that a valid value could not be retrieved.",
												},
											},
										},
									},
									"repeat_interval": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Convergence time.",
									},
									"time_range_start": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The time from which it takes effect.",
									},
									"time_range_end": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Effective end time.",
									},
									"notify_way": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Alarm notification method. At present, there are SMS, EMAIL, CALL, WECHAT methods.",
									},
									"receiver_groups": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Alert Receiving Group (User Group).",
									},
									"phone_notify_order": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Optional:    true,
										Description: "Telephone alarm sequence.",
									},
									"phone_circle_times": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "PhoneCircleTimes.",
									},
									"phone_inner_interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Telephone alarm wheel intervals. Units: Seconds.",
									},
									"phone_circle_interval": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Effective end timeTelephone alarm wheel interval. Units: Seconds.",
									},
									"phone_arrive_notice": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Telephone alerts reach notifications.",
									},
								},
							},
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last modified time.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If the alarm policy is derived from the CRD resource definition of the user cluster, the ClusterId is the cluster ID to which it belongs.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTkeTmpAlertPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreatePrometheusAlertPolicyRequest()
		response *monitor.CreatePrometheusAlertPolicyResponse
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "alert_rule"); ok {
		prometheusAlertPolicyItem := monitor.PrometheusAlertPolicyItem{}
		if v, ok := dMap["name"]; ok {
			prometheusAlertPolicyItem.Name = helper.String(v.(string))
		}
		if v, ok := dMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				RulesMap := item.(map[string]interface{})
				prometheusAlertRule := monitor.PrometheusAlertRule{}
				if v, ok := RulesMap["name"]; ok {
					prometheusAlertRule.Name = helper.String(v.(string))
				}
				if v, ok := RulesMap["rule"]; ok {
					prometheusAlertRule.Rule = helper.String(v.(string))
				}
				if v, ok := RulesMap["template"]; ok {
					prometheusAlertRule.Template = helper.String(v.(string))
				}
				if v, ok := RulesMap["for"]; ok {
					prometheusAlertRule.For = helper.String(v.(string))
				}
				if v, ok := RulesMap["describe"]; ok {
					prometheusAlertRule.Describe = helper.String(v.(string))
				}
				if v, ok := RulesMap["labels"]; ok {
					for _, item := range v.([]interface{}) {
						labelsMap := item.(map[string]interface{})
						label := monitor.Label{}
						if v, ok := labelsMap["name"]; ok {
							label.Name = helper.String(v.(string))
						}
						if v, ok := labelsMap["value"]; ok {
							label.Value = helper.String(v.(string))
						}
						prometheusAlertRule.Labels = append(prometheusAlertRule.Labels, &label)
					}
				}
				if v, ok := RulesMap["annotations"]; ok {
					for _, item := range v.([]interface{}) {
						AnnotationsMap := item.(map[string]interface{})
						label := monitor.Label{}
						if v, ok := AnnotationsMap["name"]; ok {
							label.Name = helper.String(v.(string))
						}
						if v, ok := AnnotationsMap["value"]; ok {
							label.Value = helper.String(v.(string))
						}
						prometheusAlertRule.Annotations = append(prometheusAlertRule.Annotations, &label)
					}
				}
				if v, ok := RulesMap["rule_state"]; ok {
					prometheusAlertRule.RuleState = helper.IntInt64(v.(int))
				}
				prometheusAlertPolicyItem.Rules = append(prometheusAlertPolicyItem.Rules, &prometheusAlertRule)
			}
		}
		if v, ok := dMap["id"]; ok {
			prometheusAlertPolicyItem.Id = helper.String(v.(string))
		}
		if v, ok := dMap["template_id"]; ok {
			prometheusAlertPolicyItem.TemplateId = helper.String(v.(string))
		}
		if NotificationMap, ok := helper.InterfaceToMap(dMap, "notification"); ok {
			prometheusNotificationItem := monitor.PrometheusNotificationItem{}
			if v, ok := NotificationMap["enabled"]; ok {
				prometheusNotificationItem.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := NotificationMap["type"]; ok {
				prometheusNotificationItem.Type = helper.String(v.(string))
			}
			if v, ok := NotificationMap["web_hook"]; ok {
				prometheusNotificationItem.WebHook = helper.String(v.(string))
			}
			if AlertManagerMap, ok := helper.InterfaceToMap(NotificationMap, "alert_manager"); ok {
				prometheusAlertManagerConfig := monitor.PrometheusAlertManagerConfig{}
				if v, ok := AlertManagerMap["url"]; ok {
					prometheusAlertManagerConfig.Url = helper.String(v.(string))
				}
				if v, ok := AlertManagerMap["cluster_type"]; ok {
					prometheusAlertManagerConfig.ClusterType = helper.String(v.(string))
				}
				if v, ok := AlertManagerMap["cluster_id"]; ok {
					prometheusAlertManagerConfig.ClusterId = helper.String(v.(string))
				}
				prometheusNotificationItem.AlertManager = &prometheusAlertManagerConfig
			}
			if v, ok := NotificationMap["repeat_interval"]; ok {
				prometheusNotificationItem.RepeatInterval = helper.String(v.(string))
			}
			if v, ok := NotificationMap["time_range_start"]; ok {
				prometheusNotificationItem.TimeRangeStart = helper.String(v.(string))
			}
			if v, ok := NotificationMap["time_range_end"]; ok {
				prometheusNotificationItem.TimeRangeEnd = helper.String(v.(string))
			}
			if v, ok := NotificationMap["notify_way"]; ok {
				notifyWaySet := v.(*schema.Set).List()
				for i := range notifyWaySet {
					notifyWay := notifyWaySet[i].(string)
					prometheusNotificationItem.NotifyWay = append(prometheusNotificationItem.NotifyWay, &notifyWay)
				}
			}
			if v, ok := NotificationMap["receiver_groups"]; ok {
				receiverGroupsSet := v.(*schema.Set).List()
				for i := range receiverGroupsSet {
					receiverGroups := receiverGroupsSet[i].(string)
					prometheusNotificationItem.ReceiverGroups = append(prometheusNotificationItem.ReceiverGroups, &receiverGroups)
				}
			}
			if v, ok := NotificationMap["phone_notify_order"]; ok {
				phoneNotifyOrderSet := v.(*schema.Set).List()
				for i := range phoneNotifyOrderSet {
					phoneNotifyOrder := phoneNotifyOrderSet[i].(int)
					prometheusNotificationItem.PhoneNotifyOrder = append(prometheusNotificationItem.PhoneNotifyOrder, helper.IntUint64(phoneNotifyOrder))
				}
			}
			if v, ok := NotificationMap["phone_circle_times"]; ok {
				prometheusNotificationItem.PhoneCircleTimes = helper.IntInt64(v.(int))
			}
			if v, ok := NotificationMap["phone_inner_interval"]; ok {
				prometheusNotificationItem.PhoneInnerInterval = helper.IntInt64(v.(int))
			}
			if v, ok := NotificationMap["phone_circle_interval"]; ok {
				prometheusNotificationItem.PhoneCircleInterval = helper.IntInt64(v.(int))
			}
			if v, ok := NotificationMap["phone_arrive_notice"]; ok {
				prometheusNotificationItem.PhoneArriveNotice = helper.Bool(v.(bool))
			}
			prometheusAlertPolicyItem.Notification = &prometheusNotificationItem
		}
		if v, ok := dMap["updated_at"]; ok {
			prometheusAlertPolicyItem.UpdatedAt = helper.String(v.(string))
		}
		if v, ok := dMap["cluster_id"]; ok {
			prometheusAlertPolicyItem.ClusterId = helper.String(v.(string))
		}
		request.AlertRule = &prometheusAlertPolicyItem
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusAlertPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tke tmpAlertPolicy failed, reason:%+v", logId, err)
		return err
	}

	tmpAlertPolicyId := *response.Response.Id
	instanceId := *request.InstanceId

	d.SetId(strings.Join([]string{instanceId, tmpAlertPolicyId}, FILED_SP))
	return resourceTencentCloudTkeTmpAlertPolicyRead(d, meta)
}

func resourceTencentCloudTkeTmpAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	tmpAlertPolicyId := ids[1]

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpAlertPolicy, err := service.DescribeTkeTmpAlertPolicy(ctx, instanceId, tmpAlertPolicyId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] tmpAlertPolicy[%v]\n", tmpAlertPolicy)
	if tmpAlertPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `AlertPolicy` %s does not exist", tmpAlertPolicyId)
	}

	rules := make([]map[string]interface{}, 0, len(tmpAlertPolicy.Rules))
	for _, v := range tmpAlertPolicy.Rules {
		labelList := make([]map[string]interface{}, 0, len(v.Labels))
		annotations := make([]map[string]interface{}, 0, len(v.Annotations))
		for _, label := range v.Labels {
			labelList = append(labelList, map[string]interface{}{
				"name":  label.Name,
				"value": label.Value,
			})
		}
		for _, annotation := range v.Annotations {
			annotations = append(annotations, map[string]interface{}{
				"name":  annotation.Name,
				"value": annotation.Value,
			})
		}
		rules = append(rules, map[string]interface{}{
			"name":        v.Name,
			"rule":        v.Rule,
			"labels":      labelList,
			"template":    v.Template,
			"for":         v.For,
			"describe":    v.Describe,
			"annotations": annotations,
			"rule_state":  v.RuleState,
		})
	}

	notify := tmpAlertPolicy.Notification
	alertManager := map[string]interface{}{
		"url":          notify.AlertManager.Url,
		"cluster_type": notify.AlertManager.ClusterType,
		"cluster_id":   notify.AlertManager.ClusterId,
	}
	var alertManagers []map[string]interface{}
	alertManagers = append(alertManagers, alertManager)

	var notifyWay []string
	if len(notify.NotifyWay) > 0 {
		for _, v := range notify.NotifyWay {
			notifyWay = append(notifyWay, *v)
		}
	}

	var receiverGroups []string
	if len(notify.ReceiverGroups) > 0 {
		for _, v := range notify.ReceiverGroups {
			receiverGroups = append(receiverGroups, *v)
		}
	}

	var phoneNotifyOrder []uint64
	if len(notify.PhoneNotifyOrder) > 0 {
		for _, v := range notify.PhoneNotifyOrder {
			phoneNotifyOrder = append(phoneNotifyOrder, *v)
		}
	}

	notification := map[string]interface{}{
		"enabled":       notify.Enabled,
		"type":          notify.Type,
		"web_hook":      notify.WebHook,
		"alert_manager": alertManagers,
		//"repeat_interval":       notify.RepeatInterval,
		//"time_range_start":      notify.TimeRangeStart,
		//"time_range_end":        notify.TimeRangeEnd,
		"notify_way":            notifyWay,
		"receiver_groups":       receiverGroups,
		"phone_notify_order":    phoneNotifyOrder,
		"phone_circle_times":    notify.PhoneCircleTimes,
		"phone_inner_interval":  notify.PhoneInnerInterval,
		"phone_circle_interval": notify.PhoneCircleInterval,
		"phone_arrive_notice":   notify.PhoneArriveNotice,
	}
	var notifications []map[string]interface{}
	notifications = append(notifications, notification)

	var alertRules []map[string]interface{}
	alertRule := make(map[string]interface{})
	alertRule["name"] = tmpAlertPolicy.Name
	alertRule["rules"] = rules
	alertRule["template_id"] = tmpAlertPolicy.TemplateId
	alertRule["notification"] = notifications
	alertRule["updated_at"] = tmpAlertPolicy.UpdatedAt
	alertRule["cluster_id"] = tmpAlertPolicy.ClusterId
	alertRules = append(alertRules, alertRule)

	_ = d.Set("alert_rule", alertRules)
	_ = d.Set("instance_id", instanceId)

	return nil
}

func resourceTencentCloudTkeTmpAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewModifyPrometheusAlertPolicyRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	policyId := ids[1]

	request.InstanceId = &instanceId

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "alert_rule"); ok {
		prometheusAlertPolicyItem := monitor.PrometheusAlertPolicyItem{}
		prometheusAlertPolicyItem.Id = &policyId

		if v, ok := dMap["name"]; ok {
			prometheusAlertPolicyItem.Name = helper.String(v.(string))
		}
		if v, ok := dMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				RulesMap := item.(map[string]interface{})
				prometheusAlertRule := monitor.PrometheusAlertRule{}
				if v, ok := RulesMap["name"]; ok {
					prometheusAlertRule.Name = helper.String(v.(string))
				}
				if v, ok := RulesMap["rule"]; ok {
					prometheusAlertRule.Rule = helper.String(v.(string))
				}
				if v, ok := RulesMap["template"]; ok {
					prometheusAlertRule.Template = helper.String(v.(string))
				}
				if v, ok := RulesMap["for"]; ok {
					prometheusAlertRule.For = helper.String(v.(string))
				}
				if v, ok := RulesMap["describe"]; ok {
					prometheusAlertRule.Describe = helper.String(v.(string))
				}
				if v, ok := RulesMap["annotations"]; ok {
					for _, item := range v.([]interface{}) {
						AnnotationsMap := item.(map[string]interface{})
						label := monitor.Label{}
						if v, ok := AnnotationsMap["name"]; ok {
							label.Name = helper.String(v.(string))
						}
						if v, ok := AnnotationsMap["value"]; ok {
							label.Value = helper.String(v.(string))
						}
						prometheusAlertRule.Annotations = append(prometheusAlertRule.Annotations, &label)
					}
				}
				if v, ok := RulesMap["labels"]; ok {
					for _, item := range v.([]interface{}) {
						labelsMap := item.(map[string]interface{})
						label := monitor.Label{}
						if v, ok := labelsMap["name"]; ok {
							label.Name = helper.String(v.(string))
						}
						if v, ok := labelsMap["value"]; ok {
							label.Value = helper.String(v.(string))
						}
						prometheusAlertRule.Labels = append(prometheusAlertRule.Labels, &label)
					}
				}
				if v, ok := RulesMap["rule_state"]; ok {
					prometheusAlertRule.RuleState = helper.IntInt64(v.(int))
				}
				prometheusAlertPolicyItem.Rules = append(prometheusAlertPolicyItem.Rules, &prometheusAlertRule)
			}
		}
		if v, ok := dMap["template_id"]; ok {
			prometheusAlertPolicyItem.TemplateId = helper.String(v.(string))
		}
		if NotificationMap, ok := helper.InterfaceToMap(dMap, "notification"); ok {
			prometheusNotificationItem := monitor.PrometheusNotificationItem{}
			if v, ok := NotificationMap["enabled"]; ok {
				prometheusNotificationItem.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := NotificationMap["type"]; ok {
				prometheusNotificationItem.Type = helper.String(v.(string))
			}
			if v, ok := NotificationMap["web_hook"]; ok {
				prometheusNotificationItem.WebHook = helper.String(v.(string))
			}
			if AlertManagerMap, ok := helper.InterfaceToMap(NotificationMap, "alert_manager"); ok {
				prometheusAlertManagerConfig := monitor.PrometheusAlertManagerConfig{}
				if v, ok := AlertManagerMap["url"]; ok {
					prometheusAlertManagerConfig.Url = helper.String(v.(string))
				}
				if v, ok := AlertManagerMap["cluster_type"]; ok {
					prometheusAlertManagerConfig.ClusterType = helper.String(v.(string))
				}
				if v, ok := AlertManagerMap["cluster_id"]; ok {
					prometheusAlertManagerConfig.ClusterId = helper.String(v.(string))
				}
				prometheusNotificationItem.AlertManager = &prometheusAlertManagerConfig
			}
			if v, ok := NotificationMap["repeat_interval"]; ok {
				prometheusNotificationItem.RepeatInterval = helper.String(v.(string))
			}
			if v, ok := NotificationMap["time_range_start"]; ok {
				prometheusNotificationItem.TimeRangeStart = helper.String(v.(string))
			}
			if v, ok := NotificationMap["time_range_end"]; ok {
				prometheusNotificationItem.TimeRangeEnd = helper.String(v.(string))
			}
			if v, ok := NotificationMap["notify_way"]; ok {
				notifyWaySet := v.(*schema.Set).List()
				for i := range notifyWaySet {
					notifyWay := notifyWaySet[i].(string)
					prometheusNotificationItem.NotifyWay = append(prometheusNotificationItem.NotifyWay, &notifyWay)
				}
			}
			if v, ok := NotificationMap["receiver_groups"]; ok {
				receiverGroupsSet := v.(*schema.Set).List()
				for i := range receiverGroupsSet {
					receiverGroups := receiverGroupsSet[i].(string)
					prometheusNotificationItem.ReceiverGroups = append(prometheusNotificationItem.ReceiverGroups, &receiverGroups)
				}
			}
			if v, ok := NotificationMap["phone_notify_order"]; ok {
				phoneNotifyOrderSet := v.(*schema.Set).List()
				for i := range phoneNotifyOrderSet {
					phoneNotifyOrder := phoneNotifyOrderSet[i].(int)
					prometheusNotificationItem.PhoneNotifyOrder = append(prometheusNotificationItem.PhoneNotifyOrder, helper.IntUint64(phoneNotifyOrder))
				}
			}
			if v, ok := NotificationMap["phone_circle_times"]; ok {
				prometheusNotificationItem.PhoneCircleTimes = helper.IntInt64(v.(int))
			}
			if v, ok := NotificationMap["phone_inner_interval"]; ok {
				prometheusNotificationItem.PhoneInnerInterval = helper.IntInt64(v.(int))
			}
			if v, ok := NotificationMap["phone_circle_interval"]; ok {
				prometheusNotificationItem.PhoneCircleInterval = helper.IntInt64(v.(int))
			}
			if v, ok := NotificationMap["phone_arrive_notice"]; ok {
				prometheusNotificationItem.PhoneArriveNotice = helper.Bool(v.(bool))
			}
			prometheusAlertPolicyItem.Notification = &prometheusNotificationItem
		}
		if v, ok := dMap["updated_at"]; ok {
			prometheusAlertPolicyItem.UpdatedAt = helper.String(v.(string))
		}
		if v, ok := dMap["cluster_id"]; ok {
			prometheusAlertPolicyItem.ClusterId = helper.String(v.(string))
		}
		request.AlertRule = &prometheusAlertPolicyItem
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyPrometheusAlertPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeTmpAlertPolicyRead(d, meta)
}

func resourceTencentCloudTkeTmpAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	tmpAlertPolicyId := ids[1]

	if err := service.DeleteTkeTmpAlertPolicyById(ctx, instanceId, tmpAlertPolicyId); err != nil {
		return err
	}

	return nil
}
