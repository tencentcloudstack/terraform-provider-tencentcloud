/*
Provides a resource to create a tke tmpAlertPolicy

Example Usage

```hcl

resource "tencentcloud_monitor_tmp_tke_alert_policy" "tmpAlertPolicy" {
  instance_id = "xxxxx"
  alert_rule {
    name = "xxx"
    rules {
      name = "xx"
      rule = "xx"
      template = "xx"
      for = "xx"
      labels {
        name  = "xx"
        value = "xx"
      }
      annotations {
        name  = "xx"
        value = "xx"
      }
    }
    notification {
      type = "xx"
      enabled = true
      alert_manager {
        url         = "xx"
        cluster_id   = "xx"
        cluster_type = "xx"
      }
    }
  }
}

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
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
		request  = tke.NewCreatePrometheusAlertPolicyRequest()
		response *tke.CreatePrometheusAlertPolicyResponse
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
				RulesMap := item.(map[string]interface{})
				prometheusAlertRule := tke.PrometheusAlertRule{}
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
						label := tke.Label{}
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
			prometheusNotificationItem := tke.PrometheusNotificationItem{}
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
				prometheusAlertManagerConfig := tke.PrometheusAlertManagerConfig{}
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

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().CreatePrometheusAlertPolicy(request)
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

	d.SetId(tmpAlertPolicyId)
	return resourceTencentCloudTkeTmpAlertPolicyRead(d, meta)
}

func resourceTencentCloudTkeTmpAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmpAlertPolicy.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTkeTmpAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_alert_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewModifyPrometheusAlertPolicyRequest()

	request.InstanceId = helper.String(d.Id())

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("alert_rule") {
		return fmt.Errorf("`alert_rule` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().ModifyPrometheusAlertPolicy(request)
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

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpAlertPolicyId := d.Id()

	if err := service.DeleteTkeTmpAlertPolicyById(ctx, tmpAlertPolicyId); err != nil {
		return err
	}

	return nil
}
