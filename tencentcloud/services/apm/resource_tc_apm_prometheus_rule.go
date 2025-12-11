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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				ForceNew:    true,
				Description: "Business system ID.",
			},

			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "Rule status. 1 - enabled, 2 - disabled. Default value: 1.",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the indicator matching rule.",
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

	if response.Response.RuleId == nil {
		return fmt.Errorf("RuleId is nil.")
	}

	ruleId = helper.Int64ToStr(*response.Response.RuleId)
	d.SetId(strings.Join([]string{instanceId, ruleId}, tccommon.FILED_SP))

	// set status
	if v, ok := d.GetOkExists("status"); ok {
		if v.(int) == 2 {
			request := apmv20210622.NewModifyApmPrometheusRuleRequest()
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

			request.InstanceId = &instanceId
			request.Id = helper.StrToInt64Point(ruleId)
			request.Status = helper.IntUint64(2)
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmPrometheusRuleWithContext(ctx, request)
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
	}

	return resourceTencentCloudApmPrometheusRuleRead(d, meta)
}

func resourceTencentCloudApmPrometheusRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ApmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	respData, err := service.DescribeApmPrometheusRuleById(ctx, instanceId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_apm_prometheus_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.Name != nil {
		_ = d.Set("name", *respData.Name)
	}

	if respData.ServiceName != nil {
		_ = d.Set("service_name", *respData.ServiceName)
	}

	if respData.MetricMatchType != nil {
		_ = d.Set("metric_match_type", *respData.MetricMatchType)
	}

	if respData.MetricNameRule != nil {
		_ = d.Set("metric_name_rule", *respData.MetricNameRule)
	}

	if respData.Status != nil {
		_ = d.Set("status", *respData.Status)
	}

	if respData.Id != nil {
		_ = d.Set("rule_id", *respData.Id)
	}

	return nil
}

func resourceTencentCloudApmPrometheusRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "service_name", "metric_match_type", "metric_name_rule", "status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := apmv20210622.NewModifyApmPrometheusRuleRequest()
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

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		request.InstanceId = &instanceId
		request.Id = helper.StrToInt64Point(ruleId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmPrometheusRuleWithContext(ctx, request)
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

	return resourceTencentCloudApmPrometheusRuleRead(d, meta)
}

func resourceTencentCloudApmPrometheusRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_apm_prometheus_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = apmv20210622.NewModifyApmPrometheusRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	ruleId := idSplit[1]

	request.InstanceId = &instanceId
	request.Id = helper.StrToInt64Point(ruleId)
	request.Status = helper.IntUint64(3)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseApmClient().ModifyApmPrometheusRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete apm prometheus rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
