package apm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apmv20210622 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApmPrometheusRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApmPrometheusRuleCreate,
		Read:   resourceTencentCloudApmPrometheusRuleRead,
		Update: resourceTencentCloudApmPrometheusRuleUpdate,
		Delete: resourceTencentCloudApmPrometheusRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Metric match rule name.",
			},

			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Applications where the rule takes effect. input an empty string for all applications.",
			},

			"metric_match_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Match type: 0 - precision match, 1 - prefix match, 2 - suffix match.",
			},

			"metric_name_rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the rule for customer-defined metric names with cache hit.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business system ID.",
			},
		},
	}
}

func resourceTencentCloudApmPrometheusRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = apmv20210622.NewCreateApmPrometheusRuleRequest()
		response   = apmv20210622.NewCreateApmPrometheusRuleResponse()
		instanceId string
		ruleId     string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_name"); ok {
		request.ServiceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("metric_match_type"); ok {
		request.MetricMatchType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("metric_name_rule"); ok {
		request.MetricNameRule = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().CreateApmPrometheusRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create apm prometheus rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create apm prometheus rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.InanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId

	d.SetId(strings.Join([]string{instanceId, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudApmPrometheusRuleRead(d, meta)
}

func resourceTencentCloudApmPrometheusRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	ruleId := idSplit[1]

	respData, err := service.DescribeApmPrometheusRuleById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `apm_prometheus_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	apmPrometheusRulesList := make([]map[string]interface{}, 0, len(respData.ApmPrometheusRules))
	if respData.ApmPrometheusRules != nil {
		for _, apmPrometheusRules := range respData.ApmPrometheusRules {
			apmPrometheusRulesMap := map[string]interface{}{}

			if apmPrometheusRules.Id != nil {
				apmPrometheusRulesMap["id"] = apmPrometheusRules.Id
			}

			if apmPrometheusRules.Name != nil {
				apmPrometheusRulesMap["name"] = apmPrometheusRules.Name
			}

			if apmPrometheusRules.ServiceName != nil {
				apmPrometheusRulesMap["service_name"] = apmPrometheusRules.ServiceName
			}

			if apmPrometheusRules.Status != nil {
				apmPrometheusRulesMap["status"] = apmPrometheusRules.Status
			}

			if apmPrometheusRules.MetricNameRule != nil {
				apmPrometheusRulesMap["metric_name_rule"] = apmPrometheusRules.MetricNameRule
			}

			if apmPrometheusRules.MetricMatchType != nil {
				apmPrometheusRulesMap["metric_match_type"] = apmPrometheusRules.MetricMatchType
			}

			apmPrometheusRulesList = append(apmPrometheusRulesList, apmPrometheusRulesMap)
		}

		_ = d.Set("apm_prometheus_rules", apmPrometheusRulesList)
	}

	_ = instanceId
	_ = ruleId
	return nil
}

func resourceTencentCloudApmPrometheusRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	ruleId := idSplit[1]

	needChange := false
	mutableArgs := []string{"id", "instance_id", "name", "status", "service_name", "metric_match_type", "metric_name_rule"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := apmv20210622.NewModifyApmPrometheusRuleRequest()

		if v, ok := d.GetOkExists("id"); ok {
			request.Id = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("service_name"); ok {
			request.ServiceName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("metric_match_type"); ok {
			request.MetricMatchType = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("metric_name_rule"); ok {
			request.MetricNameRule = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmV20210622Client().ModifyApmPrometheusRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if reqErr != nil {
			log.Printf("[CRITAL]%s update apm prometheus rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	_ = instanceId
	_ = ruleId
	return resourceTencentCloudApmPrometheusRuleRead(d, meta)
}

func resourceTencentCloudApmPrometheusRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	ruleId := idSplit[1]

	var (
		request  = apmv20210622.NewModifyApmPrometheusRuleRequest()
		response = apmv20210622.NewModifyApmPrometheusRuleResponse()
	)

	if v, ok := d.GetOkExists("id"); ok {
		request.Id = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("service_name"); ok {
		request.ServiceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("metric_match_type"); ok {
		request.MetricMatchType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("metric_name_rule"); ok {
		request.MetricNameRule = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmV20210622Client().ModifyApmPrometheusRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s delete apm prometheus rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	_ = response
	_ = instanceId
	_ = ruleId
	return nil
}
