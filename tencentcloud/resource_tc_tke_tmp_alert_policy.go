/*
Provides a resource to create a tke tmp_alert_policy

Example Usage

```hcl
resource "tencentcloud_tke_tmp_alert_policy" "tmp_alert_policy" {
  instance_id = &lt;nil&gt;
  alert_rule {
		name = &lt;nil&gt;
		rules {
			name = &lt;nil&gt;
			rule = &lt;nil&gt;
			labels {
				name = &lt;nil&gt;
				value = &lt;nil&gt;
			}
			template = &lt;nil&gt;
			for = &lt;nil&gt;
			describe = &lt;nil&gt;
			annotations {
				name = &lt;nil&gt;
				value = &lt;nil&gt;
			}
			rule_state = &lt;nil&gt;
		}
		id = &lt;nil&gt;
		template_id = &lt;nil&gt;
		notification {
			enabled = &lt;nil&gt;
			type = &lt;nil&gt;
			web_hook = &lt;nil&gt;
			alert_manager {
				url = &lt;nil&gt;
				cluster_type = &lt;nil&gt;
				cluster_id = &lt;nil&gt;
			}
			repeat_interval = &lt;nil&gt;
			time_range_start = &lt;nil&gt;
			time_range_end = &lt;nil&gt;
			notify_way = &lt;nil&gt;
			receiver_groups = &lt;nil&gt;
			phone_notify_order = &lt;nil&gt;
			phone_circle_times = &lt;nil&gt;
			phone_inner_interval = &lt;nil&gt;
			phone_circle_interval = &lt;nil&gt;
			phone_arrive_notice = &lt;nil&gt;
		}
		updated_at = &lt;nil&gt;
		cluster_id = &lt;nil&gt;

  }
}
```

Import

tke tmp_alert_policy can be imported using the id, e.g.

```
terraform import tencentcloud_tke_tmp_alert_policy.tmp_alert_policy tmp_alert_policy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTkeTmpAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeTmpAlertPolicyCreate,
		Read:   resourceTencentCloudTkeTmpAlertPolicyRead,
		Update: resourceTencentCloudTkeTmpAlertPolicyUpdate,
		Delete: resourceTencentCloudTkeTmpAlertPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Id.",
			},

			"alert_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
										Type:        schema.TypeList,
										Required:    true,
										Description: "Extra labels.",
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
		request  = tke.NewCreatePrometheusAlertPolicyRequest()
		response = tke.NewCreatePrometheusAlertPolicyResponse()
		id       string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "alert_rule"); ok {
		prometheusAlertPolicyItem := tke.PrometheusAlertPolicyItem{}
		if v, ok := dMap["name"]; ok {
			prometheusAlertPolicyItem.Name = helper.String(v.(string))
		}
		if v, ok := dMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				rulesMap := item.(map[string]interface{})
				prometheusAlertRule := tke.PrometheusAlertRule{}
				if v, ok := rulesMap["name"]; ok {
					prometheusAlertRule.Name = helper.String(v.(string))
				}
				if v, ok := rulesMap["rule"]; ok {
					prometheusAlertRule.Rule = helper.String(v.(string))
				}
				if v, ok := rulesMap["labels"]; ok {
					for _, item := range v.([]interface{}) {
						labelsMap := item.(map[string]interface{})
						label := tke.Label{}
						if v, ok := labelsMap["name"]; ok {
							label.Name = helper.String(v.(string))
						}
						if v, ok := labelsMap["value"]; ok {
							label.Value = helper.String(v.(string))
						}
						prometheusAlertRule.Labels = append(prometheusAlertRule.Labels, &label)
					}
				}
				if v, ok := rulesMap["template"]; ok {
					prometheusAlertRule.Template = helper.String(v.(string))
				}
				if v, ok := rulesMap["for"]; ok {
					prometheusAlertRule.For = helper.String(v.(string))
				}
				if v, ok := rulesMap["describe"]; ok {
					prometheusAlertRule.Describe = helper.String(v.(string))
				}
				if v, ok := rulesMap["annotations"]; ok {
					for _, item := range v.([]interface{}) {
						annotationsMap := item.(map[string]interface{})
						label := tke.Label{}
						if v, ok := annotationsMap["name"]; ok {
							label.Name = helper.String(v.(string))
						}
						if v, ok := annotationsMap["value"]; ok {
							label.Value = helper.String(v.(string))
						}
						prometheusAlertRule.Annotations = append(prometheusAlertRule.Annotations, &label)
					}
				}
				if v, ok := rulesMap["rule_state"]; ok {
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
		if notificationMap, ok := helper.InterfaceToMap(dMap, "notification"); ok {
			prometheusNotificationItem := tke.PrometheusNotificationItem{}
			if v, ok := notificationMap["enabled"]; ok {
				prometheusNotificationItem.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := notificationMap["type"]; ok {
				prometheusNotificationItem.Type = helper.String(v.(string))
			}
			if v, ok := notificationMap["web_hook"]; ok {
				prometheusNotificationItem.WebHook = helper.String(v.(string))
			}
			if alertManagerMap, ok := helper.InterfaceToMap(notificationMap, "alert_manager"); ok {
				prometheusAlertManagerConfig := tke.PrometheusAlertManagerConfig{}
				if v, ok := alertManagerMap["url"]; ok {
					prometheusAlertManagerConfig.Url = helper.String(v.(string))
				}
				if v, ok := alertManagerMap["cluster_type"]; ok {
					prometheusAlertManagerConfig.ClusterType = helper.String(v.(string))
				}
				if v, ok := alertManagerMap["cluster_id"]; ok {
					prometheusAlertManagerConfig.ClusterId = helper.String(v.(string))
				}
				prometheusNotificationItem.AlertManager = &prometheusAlertManagerConfig
			}
			if v, ok := notificationMap["repeat_interval"]; ok {
				prometheusNotificationItem.RepeatInterval = helper.String(v.(string))
			}
			if v, ok := notificationMap["time_range_start"]; ok {
				prometheusNotificationItem.TimeRangeStart = helper.String(v.(string))
			}
			if v, ok := notificationMap["time_range_end"]; ok {
				prometheusNotificationItem.TimeRangeEnd = helper.String(v.(string))
			}
			if v, ok := notificationMap["notify_way"]; ok {
				notifyWaySet := v.(*schema.Set).List()
				for i := range notifyWaySet {
					notifyWay := notifyWaySet[i].(string)
					prometheusNotificationItem.NotifyWay = append(prometheusNotificationItem.NotifyWay, &notifyWay)
				}
			}
			if v, ok := notificationMap["receiver_groups"]; ok {
				receiverGroupsSet := v.(*schema.Set).List()
				for i := range receiverGroupsSet {
					receiverGroups := receiverGroupsSet[i].(string)
					prometheusNotificationItem.ReceiverGroups = append(prometheusNotificationItem.ReceiverGroups, &receiverGroups)
				}
			}
			if v, ok := notificationMap["phone_notify_order"]; ok {
				phoneNotifyOrderSet := v.(*schema.Set).List()
				for i := range phoneNotifyOrderSet {
					phoneNotifyOrder := phoneNotifyOrderSet[i].(int)
					prometheusNotificationItem.PhoneNotifyOrder = append(prometheusNotificationItem.PhoneNotifyOrder, helper.IntUint64(phoneNotifyOrder))
				}
			}
			if v, ok := notificationMap["phone_circle_times"]; ok {
				prometheusNotificationItem.PhoneCircleTimes = helper.IntInt64(v.(int))
			}
			if v, ok := notificationMap["phone_inner_interval"]; ok {
				prometheusNotificationItem.PhoneInnerInterval = helper.IntInt64(v.(int))
			}
			if v, ok := notificationMap["phone_circle_interval"]; ok {
				prometheusNotificationItem.PhoneCircleInterval = helper.IntInt64(v.(int))
			}
			if v, ok := notificationMap["phone_arrive_notice"]; ok {
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().CreatePrometheusAlertPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tke tmpAlertPolicy failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(id)

	return resourceTencentCloudTkeTmpAlertPolicyRead(d, meta)
}

func resourceTencentCloudTkeTmpAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpAlertPolicyId := d.Id()

	tmpAlertPolicy, err := service.DescribeTkeTmpAlertPolicyById(ctx, id)
	if err != nil {
		return err
	}

	if tmpAlertPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TkeTmpAlertPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tmpAlertPolicy.InstanceId != nil {
		_ = d.Set("instance_id", tmpAlertPolicy.InstanceId)
	}

	if tmpAlertPolicy.AlertRule != nil {
		alertRuleMap := map[string]interface{}{}

		if tmpAlertPolicy.AlertRule.Name != nil {
			alertRuleMap["name"] = tmpAlertPolicy.AlertRule.Name
		}

		if tmpAlertPolicy.AlertRule.Rules != nil {
			rulesList := []interface{}{}
			for _, rules := range tmpAlertPolicy.AlertRule.Rules {
				rulesMap := map[string]interface{}{}

				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Rule != nil {
					rulesMap["rule"] = rules.Rule
				}

				if rules.Labels != nil {
					labelsList := []interface{}{}
					for _, labels := range rules.Labels {
						labelsMap := map[string]interface{}{}

						if labels.Name != nil {
							labelsMap["name"] = labels.Name
						}

						if labels.Value != nil {
							labelsMap["value"] = labels.Value
						}

						labelsList = append(labelsList, labelsMap)
					}

					rulesMap["labels"] = []interface{}{labelsList}
				}

				if rules.Template != nil {
					rulesMap["template"] = rules.Template
				}

				if rules.For != nil {
					rulesMap["for"] = rules.For
				}

				if rules.Describe != nil {
					rulesMap["describe"] = rules.Describe
				}

				if rules.Annotations != nil {
					annotationsList := []interface{}{}
					for _, annotations := range rules.Annotations {
						annotationsMap := map[string]interface{}{}

						if annotations.Name != nil {
							annotationsMap["name"] = annotations.Name
						}

						if annotations.Value != nil {
							annotationsMap["value"] = annotations.Value
						}

						annotationsList = append(annotationsList, annotationsMap)
					}

					rulesMap["annotations"] = []interface{}{annotationsList}
				}

				if rules.RuleState != nil {
					rulesMap["rule_state"] = rules.RuleState
				}

				rulesList = append(rulesList, rulesMap)
			}

			alertRuleMap["rules"] = []interface{}{rulesList}
		}

		if tmpAlertPolicy.AlertRule.Id != nil {
			alertRuleMap["id"] = tmpAlertPolicy.AlertRule.Id
		}

		if tmpAlertPolicy.AlertRule.TemplateId != nil {
			alertRuleMap["template_id"] = tmpAlertPolicy.AlertRule.TemplateId
		}

		if tmpAlertPolicy.AlertRule.Notification != nil {
			notificationMap := map[string]interface{}{}

			if tmpAlertPolicy.AlertRule.Notification.Enabled != nil {
				notificationMap["enabled"] = tmpAlertPolicy.AlertRule.Notification.Enabled
			}

			if tmpAlertPolicy.AlertRule.Notification.Type != nil {
				notificationMap["type"] = tmpAlertPolicy.AlertRule.Notification.Type
			}

			if tmpAlertPolicy.AlertRule.Notification.WebHook != nil {
				notificationMap["web_hook"] = tmpAlertPolicy.AlertRule.Notification.WebHook
			}

			if tmpAlertPolicy.AlertRule.Notification.AlertManager != nil {
				alertManagerMap := map[string]interface{}{}

				if tmpAlertPolicy.AlertRule.Notification.AlertManager.Url != nil {
					alertManagerMap["url"] = tmpAlertPolicy.AlertRule.Notification.AlertManager.Url
				}

				if tmpAlertPolicy.AlertRule.Notification.AlertManager.ClusterType != nil {
					alertManagerMap["cluster_type"] = tmpAlertPolicy.AlertRule.Notification.AlertManager.ClusterType
				}

				if tmpAlertPolicy.AlertRule.Notification.AlertManager.ClusterId != nil {
					alertManagerMap["cluster_id"] = tmpAlertPolicy.AlertRule.Notification.AlertManager.ClusterId
				}

				notificationMap["alert_manager"] = []interface{}{alertManagerMap}
			}

			if tmpAlertPolicy.AlertRule.Notification.RepeatInterval != nil {
				notificationMap["repeat_interval"] = tmpAlertPolicy.AlertRule.Notification.RepeatInterval
			}

			if tmpAlertPolicy.AlertRule.Notification.TimeRangeStart != nil {
				notificationMap["time_range_start"] = tmpAlertPolicy.AlertRule.Notification.TimeRangeStart
			}

			if tmpAlertPolicy.AlertRule.Notification.TimeRangeEnd != nil {
				notificationMap["time_range_end"] = tmpAlertPolicy.AlertRule.Notification.TimeRangeEnd
			}

			if tmpAlertPolicy.AlertRule.Notification.NotifyWay != nil {
				notificationMap["notify_way"] = tmpAlertPolicy.AlertRule.Notification.NotifyWay
			}

			if tmpAlertPolicy.AlertRule.Notification.ReceiverGroups != nil {
				notificationMap["receiver_groups"] = tmpAlertPolicy.AlertRule.Notification.ReceiverGroups
			}

			if tmpAlertPolicy.AlertRule.Notification.PhoneNotifyOrder != nil {
				notificationMap["phone_notify_order"] = tmpAlertPolicy.AlertRule.Notification.PhoneNotifyOrder
			}

			if tmpAlertPolicy.AlertRule.Notification.PhoneCircleTimes != nil {
				notificationMap["phone_circle_times"] = tmpAlertPolicy.AlertRule.Notification.PhoneCircleTimes
			}

			if tmpAlertPolicy.AlertRule.Notification.PhoneInnerInterval != nil {
				notificationMap["phone_inner_interval"] = tmpAlertPolicy.AlertRule.Notification.PhoneInnerInterval
			}

			if tmpAlertPolicy.AlertRule.Notification.PhoneCircleInterval != nil {
				notificationMap["phone_circle_interval"] = tmpAlertPolicy.AlertRule.Notification.PhoneCircleInterval
			}

			if tmpAlertPolicy.AlertRule.Notification.PhoneArriveNotice != nil {
				notificationMap["phone_arrive_notice"] = tmpAlertPolicy.AlertRule.Notification.PhoneArriveNotice
			}

			alertRuleMap["notification"] = []interface{}{notificationMap}
		}

		if tmpAlertPolicy.AlertRule.UpdatedAt != nil {
			alertRuleMap["updated_at"] = tmpAlertPolicy.AlertRule.UpdatedAt
		}

		if tmpAlertPolicy.AlertRule.ClusterId != nil {
			alertRuleMap["cluster_id"] = tmpAlertPolicy.AlertRule.ClusterId
		}

		_ = d.Set("alert_rule", []interface{}{alertRuleMap})
	}

	return nil
}

func resourceTencentCloudTkeTmpAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewModifyPrometheusAlertPolicyRequest()

	tmpAlertPolicyId := d.Id()

	request.Id = &id

	immutableArgs := []string{"instance_id", "alert_rule"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().ModifyPrometheusAlertPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tke tmpAlertPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTkeTmpAlertPolicyRead(d, meta)
}

func resourceTencentCloudTkeTmpAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpAlertPolicyId := d.Id()

	if err := service.DeleteTkeTmpAlertPolicyById(ctx, id); err != nil {
		return err
	}

	return nil
}
