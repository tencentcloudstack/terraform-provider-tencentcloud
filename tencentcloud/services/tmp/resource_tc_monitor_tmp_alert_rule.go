package tmp

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpAlertRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpAlertRuleRead,
		Create: resourceTencentCloudMonitorTmpAlertRuleCreate,
		Update: resourceTencentCloudMonitorTmpAlertRuleUpdate,
		Delete: resourceTencentCloudMonitorTmpAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},
			"expr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule expression, reference documentation: `https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/`.",
			},
			"receivers": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Alarm notification template id list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rule_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule state code.",
			},
			"duration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule alarm duration.",
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Rule alarm duration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value.",
						},
					},
				},
			},
			"annotations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Rule alarm duration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value.",
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alarm Policy Template Classification.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpAlertRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = monitor.NewCreateAlertRuleRequest()
		response *monitor.CreateAlertRuleResponse
	)

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("expr"); ok {
		request.Expr = helper.String(v.(string))
	}
	if v, ok := d.GetOk("receivers"); ok {
		receivers := v.(*schema.Set).List()
		receiverArr := make([]*string, 0, len(receivers))
		for _, receiver := range receivers {
			receiverArr = append(receiverArr, helper.String(receiver.(string)))
		}
		request.Receivers = receiverArr
	}
	if v, ok := d.GetOk("rule_state"); ok {
		request.RuleState = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("duration"); ok {
		request.Duration = helper.String(v.(string))
	}

	if v, ok := d.GetOk("labels"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			prometheusRuleKV := monitor.PrometheusRuleKV{}
			if v, ok := dMap["key"]; ok {
				prometheusRuleKV.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				prometheusRuleKV.Value = helper.String(v.(string))
			}
			request.Labels = append(request.Labels, &prometheusRuleKV)
		}
	}

	if v, ok := d.GetOk("annotations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			prometheusRuleKV := monitor.PrometheusRuleKV{}
			if v, ok := dMap["key"]; ok {
				prometheusRuleKV.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				prometheusRuleKV.Value = helper.String(v.(string))
			}
			request.Annotations = append(request.Annotations, &prometheusRuleKV)
		}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateAlertRule(request)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains("FailedOperation", ee.Code) {
				return resource.NonRetryableError(ee)
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpAlertRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId := *response.Response.RuleId

	d.SetId(strings.Join([]string{instanceId, ruleId}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorTmpAlertRuleRead(d, meta)
}

func resourceTencentCloudMonitorTmpAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	tmpAlertRule, err := service.DescribeMonitorTmpAlertRuleById(ctx, ids[0], ids[1])

	if err != nil {
		return err
	}

	if tmpAlertRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tmpAlertRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", ids[0])
	if tmpAlertRule.RuleName != nil {
		_ = d.Set("rule_name", tmpAlertRule.RuleName)
	}
	if tmpAlertRule.Expr != nil {
		_ = d.Set("expr", tmpAlertRule.Expr)
	}
	if tmpAlertRule.Receivers != nil {
		list := tmpAlertRule.Receivers
		result := make([]string, 0, len(list))
		for _, v := range list {
			result = append(result, *v)
		}
		_ = d.Set("receivers", result)
	}
	if tmpAlertRule.RuleState != nil {
		_ = d.Set("rule_state", tmpAlertRule.RuleState)
	}
	if tmpAlertRule.Duration != nil {
		_ = d.Set("duration", tmpAlertRule.Duration)
	}
	if tmpAlertRule.Labels != nil {
		labelsList := tmpAlertRule.Labels
		result := make([]map[string]interface{}, 0, len(labelsList))
		for _, v := range labelsList {
			mapping := map[string]interface{}{
				"key":   v.Key,
				"value": v.Value,
			}
			result = append(result, mapping)
		}
		_ = d.Set("labels", result)
	}
	if tmpAlertRule.Annotations != nil {
		annotationsList := tmpAlertRule.Annotations
		result := make([]map[string]interface{}, 0, len(annotationsList))
		for _, v := range annotationsList {
			mapping := map[string]interface{}{
				"key":   v.Key,
				"value": v.Value,
			}
			result = append(result, mapping)
		}
		_ = d.Set("annotations", result)
	}
	if tmpAlertRule.Type != nil {
		_ = d.Set("type", tmpAlertRule.Type)
	}

	return nil
}

func resourceTencentCloudMonitorTmpAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdateAlertRuleRequest()

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = helper.String(ids[0])
	request.RuleId = helper.String(ids[1])

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expr"); ok {
		request.Expr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receivers"); ok {
		receivers := v.(*schema.Set).List()
		receiverArr := make([]*string, 0, len(receivers))
		for _, receiver := range receivers {
			receiverArr = append(receiverArr, helper.String(receiver.(string)))
		}
		request.Receivers = receiverArr
	}

	if v, ok := d.GetOk("rule_state"); ok {
		request.RuleState = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("duration"); ok {
		request.Duration = helper.String(v.(string))
	}

	if v, ok := d.GetOk("labels"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			prometheusRuleKV := monitor.PrometheusRuleKV{}
			if v, ok := dMap["key"]; ok {
				prometheusRuleKV.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				prometheusRuleKV.Value = helper.String(v.(string))
			}
			request.Labels = append(request.Labels, &prometheusRuleKV)
		}
	}

	if v, ok := d.GetOk("annotations"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			prometheusRuleKV := monitor.PrometheusRuleKV{}
			if v, ok := dMap["key"]; ok {
				prometheusRuleKV.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				prometheusRuleKV.Value = helper.String(v.(string))
			}
			request.Annotations = append(request.Annotations, &prometheusRuleKV)
		}
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateAlertRule(request)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains("FailedOperation", ee.Code) {
				return resource.NonRetryableError(ee)
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorTmpAlertRuleRead(d, meta)
}

func resourceTencentCloudMonitorTmpAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_alert_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	if err := service.DeleteMonitorTmpAlertRule(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
