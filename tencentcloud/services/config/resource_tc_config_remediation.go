package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigRemediation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigRemediationCreate,
		Read:   resourceTencentCloudConfigRemediationRead,
		Update: resourceTencentCloudConfigRemediationUpdate,
		Delete: resourceTencentCloudConfigRemediationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Config rule ID to bind the remediation setting to.",
			},

			"remediation_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remediation type. Valid value: SCF (cloud function, custom remediation).",
			},

			"remediation_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remediation template ID (e.g. SCF function resource path: qcs::scf:ap-guangzhou:uin/functions/xxx).",
			},

			"invoke_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remediation execution mode. Valid values: MANUAL_EXECUTION (manual), AUTO_EXECUTION (automatic), NON_EXECUTION (disabled), NOT_CONFIG (not configured).",
			},

			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template source. Valid value: CUSTOM (custom template).",
			},

			// Computed
			"remediation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Remediation setting ID.",
			},

			"owner_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Owner account UIN.",
			},

			"remediation_source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Remediation source type returned from API.",
			},
		},
	}
}

func resourceTencentCloudConfigRemediationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_remediation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = configv20220802.NewCreateRemediationRequest()
		response = configv20220802.NewCreateRemediationResponse()
	)

	if v, ok := d.GetOk("rule_id"); ok {
		request.RuleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remediation_type"); ok {
		request.RemediationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remediation_template_id"); ok {
		request.RemediationTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invoke_type"); ok {
		request.InvokeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		request.SourceType = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().CreateRemediationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create config remediation failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create config remediation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.RemediationId == nil {
		return fmt.Errorf("RemediationId is nil")
	}

	d.SetId(*response.Response.RemediationId)
	return resourceTencentCloudConfigRemediationRead(d, meta)
}

func resourceTencentCloudConfigRemediationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_remediation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		remediationId = d.Id()
	)

	// rule_id may be empty on import; use empty string to list all and match by remediationId
	ruleId := ""
	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
	}

	service := ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	respData, err := service.DescribeConfigRemediationById(ctx, ruleId, remediationId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_remediation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	if respData.RemediationType != nil {
		_ = d.Set("remediation_type", respData.RemediationType)
	}

	if respData.RemediationTemplateId != nil {
		_ = d.Set("remediation_template_id", respData.RemediationTemplateId)
	}

	if respData.InvokeType != nil {
		_ = d.Set("invoke_type", respData.InvokeType)
	}

	if respData.RemediationId != nil {
		_ = d.Set("remediation_id", respData.RemediationId)
	}

	if respData.OwnerUin != nil {
		_ = d.Set("owner_uin", respData.OwnerUin)
	}

	if respData.RemediationSourceType != nil {
		_ = d.Set("remediation_source_type", respData.RemediationSourceType)
	}

	return nil
}

func resourceTencentCloudConfigRemediationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_remediation.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		remediationId = d.Id()
	)

	mutableArgs := []string{"remediation_type", "remediation_template_id", "invoke_type", "source_type"}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := configv20220802.NewUpdateRemediationRequest()
		request.RemediationId = &remediationId

		if v, ok := d.GetOk("remediation_type"); ok {
			request.RemediationType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remediation_template_id"); ok {
			request.RemediationTemplateId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("invoke_type"); ok {
			request.InvokeType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("source_type"); ok {
			request.SourceType = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateRemediationWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config remediation failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudConfigRemediationRead(d, meta)
}

func resourceTencentCloudConfigRemediationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_remediation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		remediationId = d.Id()
		request       = configv20220802.NewDeleteRemediationsRequest()
	)

	request.RemediationIds = []*string{&remediationId}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteRemediationsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete config remediation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
