package config

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigStartConfigRuleEvaluationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigStartConfigRuleEvaluationOperationCreate,
		Read:   resourceTencentCloudConfigStartConfigRuleEvaluationOperationRead,
		Delete: resourceTencentCloudConfigStartConfigRuleEvaluationOperationDelete,
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Config rule ID to trigger evaluation for.",
			},

			"compliance_pack_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Compliance pack ID to trigger evaluation for.",
			},
		},
	}
}

func resourceTencentCloudConfigStartConfigRuleEvaluationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_start_config_rule_evaluation_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = configv20220802.NewStartConfigRuleEvaluationRequest()
	)

	if v, ok := d.GetOk("rule_id"); ok {
		request.RuleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compliance_pack_id"); ok {
		request.CompliancePackId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().StartConfigRuleEvaluationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s trigger config rule evaluation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(helper.BuildToken())
	return resourceTencentCloudConfigStartConfigRuleEvaluationOperationRead(d, meta)
}

func resourceTencentCloudConfigStartConfigRuleEvaluationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_start_config_rule_evaluation_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudConfigStartConfigRuleEvaluationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_start_config_rule_evaluation_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
