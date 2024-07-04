package tmp

import (
	"context"
	"fmt"
	"log"
	"strings"

	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpAlertGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpAlertGroupCreate,
		Read:   resourceTencentCloudMonitorTmpAlertGroupRead,
		Update: resourceTencentCloudMonitorTmpAlertGroupUpdate,
		Delete: resourceTencentCloudMonitorTmpAlertGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Alarm group id.",
			},

			"group_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Unique alert group name.",
			},

			"amp_receivers": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tencent cloud notification template id list.",
			},

			"custom_receiver": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "User custom notification template, such as webhook, alertmanager.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom receiver type, webhook|alertmanager.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom receiver address, can be accessed by process in prometheus instance subnet.",
						},
						"allowed_time_ranges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Time ranges which allow alert message send.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Time range start, seconds since 0 o'clock.",
									},
									"end": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Time range end, seconds since 0 o'clock.",
									},
								},
							},
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Only effect when alertmanager in user cluster, this cluster id.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Only effect when alertmanager in user cluster, this cluster type (tke|eks|tdcc).",
						},
					},
				},
			},

			"repeat_interval": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Alert message send interval, default 1 hour.",
			},

			"rules": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A list of alert rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alert rule name.",
						},
						"labels": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Labels of alert rule.",
						},
						"annotations": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Annotation of alert rule. `summary`, `description` is special annotation in prometheus, mapping `Alarm Object`, `Alarm Information` in alarm message.",
						},
						"duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule alarm duration.",
						},
						"expr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prometheus alert expression.",
						},
						"state": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule state. `2`-enable, `3`-disable, default `2`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorTmpAlertGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = monitor.NewCreatePrometheusAlertGroupRequest()
		response   = monitor.NewCreatePrometheusAlertGroupResponse()
		instanceId string
		groupId    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("amp_receivers"); ok {
		aMPReceiversSet := v.(*schema.Set).List()
		for i := range aMPReceiversSet {
			if aMPReceiversSet[i] != nil {
				aMPReceivers := aMPReceiversSet[i].(string)
				request.AMPReceivers = append(request.AMPReceivers, &aMPReceivers)
			}
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "custom_receiver"); ok {
		prometheusAlertCustomReceiver := monitor.PrometheusAlertCustomReceiver{}
		if v, ok := dMap["type"]; ok {
			prometheusAlertCustomReceiver.Type = helper.String(v.(string))
		}

		if v, ok := dMap["url"]; ok {
			prometheusAlertCustomReceiver.Url = helper.String(v.(string))
		}

		if v, ok := dMap["allowed_time_ranges"]; ok {
			for _, item := range v.([]interface{}) {
				allowedTimeRangesMap := item.(map[string]interface{})
				prometheusAlertAllowTimeRange := monitor.PrometheusAlertAllowTimeRange{}
				if v, ok := allowedTimeRangesMap["start"]; ok {
					prometheusAlertAllowTimeRange.Start = helper.String(v.(string))
				}

				if v, ok := allowedTimeRangesMap["end"]; ok {
					prometheusAlertAllowTimeRange.End = helper.String(v.(string))
				}

				prometheusAlertCustomReceiver.AllowedTimeRanges = append(prometheusAlertCustomReceiver.AllowedTimeRanges, &prometheusAlertAllowTimeRange)
			}
		}

		if v, ok := dMap["cluster_id"]; ok {
			prometheusAlertCustomReceiver.ClusterId = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_type"]; ok {
			prometheusAlertCustomReceiver.ClusterType = helper.String(v.(string))
		}

		request.CustomReceiver = &prometheusAlertCustomReceiver
	}

	if v, ok := d.GetOk("repeat_interval"); ok {
		request.RepeatInterval = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			prometheusAlertGroupRuleSet := monitor.PrometheusAlertGroupRuleSet{}
			if v, ok := dMap["rule_name"]; ok {
				prometheusAlertGroupRuleSet.RuleName = helper.String(v.(string))
			}

			if v, ok := dMap["labels"]; ok {
				labelsMap := v.(map[string]interface{})
				for k, v := range labelsMap {
					prometheusRuleKV := monitor.PrometheusRuleKV{}
					prometheusRuleKV.Key = helper.String(k)
					prometheusRuleKV.Value = helper.String(v.(string))
					prometheusAlertGroupRuleSet.Labels = append(prometheusAlertGroupRuleSet.Labels, &prometheusRuleKV)
				}
			}

			if v, ok := dMap["annotations"]; ok {
				annotationsMap := v.(map[string]interface{})
				for k, v := range annotationsMap {
					prometheusRuleKV := monitor.PrometheusRuleKV{}
					prometheusRuleKV.Key = helper.String(k)
					prometheusRuleKV.Value = helper.String(v.(string))
					prometheusAlertGroupRuleSet.Annotations = append(prometheusAlertGroupRuleSet.Annotations, &prometheusRuleKV)
				}
			}
			if v, ok := dMap["duration"]; ok {
				prometheusAlertGroupRuleSet.Duration = helper.String(v.(string))
			}

			if v, ok := dMap["expr"]; ok {
				prometheusAlertGroupRuleSet.Expr = helper.String(v.(string))
			}

			if v, ok := dMap["state"]; ok {
				prometheusAlertGroupRuleSet.State = helper.IntInt64(v.(int))
			}

			request.Rules = append(request.Rules, &prometheusAlertGroupRuleSet)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreatePrometheusAlertGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpAlertGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.GroupId
	d.SetId(instanceId + tccommon.FILED_SP + groupId)

	return resourceTencentCloudMonitorTmpAlertGroupRead(d, meta)
}

func resourceTencentCloudMonitorTmpAlertGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		tmpAlertGroup *monitor.PrometheusAlertGroupSet
	)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := ids[0]
	groupId := ids[1]

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeMonitorTmpAlertGroupById(ctx, instanceId, groupId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		tmpAlertGroup = result
		return nil
	})

	if err != nil {
		return err
	}

	if tmpAlertGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpAlertGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("group_id", groupId)

	if tmpAlertGroup.GroupName != nil {
		_ = d.Set("group_name", tmpAlertGroup.GroupName)
	}

	if tmpAlertGroup.AMPReceivers != nil {
		_ = d.Set("amp_receivers", tmpAlertGroup.AMPReceivers)
	}

	if tmpAlertGroup.CustomReceiver != nil {
		customReceiverMap := map[string]interface{}{}

		if tmpAlertGroup.CustomReceiver.Type != nil {
			customReceiverMap["type"] = tmpAlertGroup.CustomReceiver.Type
		}

		if tmpAlertGroup.CustomReceiver.Url != nil {
			customReceiverMap["url"] = tmpAlertGroup.CustomReceiver.Url
		}

		if tmpAlertGroup.CustomReceiver.AllowedTimeRanges != nil {
			allowedTimeRangesList := []interface{}{}
			for _, allowedTimeRanges := range tmpAlertGroup.CustomReceiver.AllowedTimeRanges {
				allowedTimeRangesMap := map[string]interface{}{}

				if allowedTimeRanges.Start != nil {
					allowedTimeRangesMap["start"] = allowedTimeRanges.Start
				}

				if allowedTimeRanges.End != nil {
					allowedTimeRangesMap["end"] = allowedTimeRanges.End
				}

				allowedTimeRangesList = append(allowedTimeRangesList, allowedTimeRangesMap)
			}

			customReceiverMap["allowed_time_ranges"] = allowedTimeRangesList
		}

		if tmpAlertGroup.CustomReceiver.ClusterId != nil {
			customReceiverMap["cluster_id"] = tmpAlertGroup.CustomReceiver.ClusterId
		}

		if tmpAlertGroup.CustomReceiver.ClusterType != nil {
			customReceiverMap["cluster_type"] = tmpAlertGroup.CustomReceiver.ClusterType
		}

		_ = d.Set("custom_receiver", []interface{}{customReceiverMap})
	}

	if tmpAlertGroup.RepeatInterval != nil {
		_ = d.Set("repeat_interval", tmpAlertGroup.RepeatInterval)
	}

	if tmpAlertGroup.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range tmpAlertGroup.Rules {
			rulesMap := map[string]interface{}{}

			if rules.RuleName != nil {
				rulesMap["rule_name"] = rules.RuleName
			}

			// The api will have inconsistent order, so map is used here.
			if rules.Labels != nil {
				labelsMap := map[string]interface{}{}
				for _, labels := range rules.Labels {
					if labels.Key != nil {
						labelsMap[*labels.Key] = labels.Value
					}
				}

				rulesMap["labels"] = labelsMap
			}

			if rules.Annotations != nil {
				annotationsMap := map[string]interface{}{}
				for _, annotations := range rules.Annotations {
					if annotations.Key != nil {
						annotationsMap[*annotations.Key] = annotations.Value
					}
				}

				rulesMap["annotations"] = annotationsMap
			}

			if rules.Duration != nil {
				rulesMap["duration"] = rules.Duration
			}

			if rules.Expr != nil {
				rulesMap["expr"] = rules.Expr
			}

			if rules.State != nil {
				rulesMap["state"] = rules.State
			}

			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)
	}

	return nil
}

func resourceTencentCloudMonitorTmpAlertGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = monitor.NewUpdatePrometheusAlertGroupRequest()
	)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := ids[0]
	groupId := ids[1]

	request.InstanceId = &instanceId
	request.GroupId = &groupId
	if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("amp_receivers"); ok {
		aMPReceiversSet := v.(*schema.Set).List()
		for i := range aMPReceiversSet {
			if aMPReceiversSet[i] != nil {
				aMPReceivers := aMPReceiversSet[i].(string)
				request.AMPReceivers = append(request.AMPReceivers, &aMPReceivers)
			}
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "custom_receiver"); ok {
		prometheusAlertCustomReceiver := monitor.PrometheusAlertCustomReceiver{}
		if v, ok := dMap["type"]; ok {
			prometheusAlertCustomReceiver.Type = helper.String(v.(string))
		}

		if v, ok := dMap["url"]; ok {
			prometheusAlertCustomReceiver.Url = helper.String(v.(string))
		}

		if v, ok := dMap["allowed_time_ranges"]; ok {
			for _, item := range v.([]interface{}) {
				allowedTimeRangesMap := item.(map[string]interface{})
				prometheusAlertAllowTimeRange := monitor.PrometheusAlertAllowTimeRange{}
				if v, ok := allowedTimeRangesMap["start"]; ok {
					prometheusAlertAllowTimeRange.Start = helper.String(v.(string))
				}

				if v, ok := allowedTimeRangesMap["end"]; ok {
					prometheusAlertAllowTimeRange.End = helper.String(v.(string))
				}

				prometheusAlertCustomReceiver.AllowedTimeRanges = append(prometheusAlertCustomReceiver.AllowedTimeRanges, &prometheusAlertAllowTimeRange)
			}
		}

		if v, ok := dMap["cluster_id"]; ok {
			prometheusAlertCustomReceiver.ClusterId = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_type"]; ok {
			prometheusAlertCustomReceiver.ClusterType = helper.String(v.(string))
		}

		request.CustomReceiver = &prometheusAlertCustomReceiver
	}

	if v, ok := d.GetOk("repeat_interval"); ok {
		request.RepeatInterval = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			prometheusAlertGroupRuleSet := monitor.PrometheusAlertGroupRuleSet{}
			if v, ok := dMap["rule_name"]; ok {
				prometheusAlertGroupRuleSet.RuleName = helper.String(v.(string))
			}

			if v, ok := dMap["labels"]; ok {
				labelsMap := v.(map[string]interface{})
				for k, v := range labelsMap {
					prometheusRuleKV := monitor.PrometheusRuleKV{}
					prometheusRuleKV.Key = helper.String(k)
					prometheusRuleKV.Value = helper.String(v.(string))
					prometheusAlertGroupRuleSet.Labels = append(prometheusAlertGroupRuleSet.Labels, &prometheusRuleKV)
				}
			}

			if v, ok := dMap["annotations"]; ok {
				annotationsMap := v.(map[string]interface{})
				for k, v := range annotationsMap {
					prometheusRuleKV := monitor.PrometheusRuleKV{}
					prometheusRuleKV.Key = helper.String(k)
					prometheusRuleKV.Value = helper.String(v.(string))
					prometheusAlertGroupRuleSet.Annotations = append(prometheusAlertGroupRuleSet.Annotations, &prometheusRuleKV)
				}
			}

			if v, ok := dMap["duration"]; ok {
				prometheusAlertGroupRuleSet.Duration = helper.String(v.(string))
			}

			if v, ok := dMap["expr"]; ok {
				prometheusAlertGroupRuleSet.Expr = helper.String(v.(string))
			}

			if v, ok := dMap["state"]; ok {
				prometheusAlertGroupRuleSet.State = helper.IntInt64(v.(int))
			}

			request.Rules = append(request.Rules, &prometheusAlertGroupRuleSet)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdatePrometheusAlertGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpAlertGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpAlertGroupRead(d, meta)
}

func resourceTencentCloudMonitorTmpAlertGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := ids[0]
	groupId := ids[1]

	if err := service.DeleteMonitorTmpAlertGroupById(ctx, instanceId, groupId); err != nil {
		return err
	}

	return nil
}
