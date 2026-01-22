package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsIntlv20201028 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPrivateDnsForwardRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivateDnsForwardRuleCreate,
		Read:   resourceTencentCloudPrivateDnsForwardRuleRead,
		Update: resourceTencentCloudPrivateDnsForwardRuleUpdate,
		Delete: resourceTencentCloudPrivateDnsForwardRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Forwarding rule name.",
			},

			"rule_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Forwarding rule type. DOWN: From cloud to off-cloud; UP: From off-cloud to cloud.",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private domain ID, which can be viewed on the private domain list page.",
			},

			"end_point_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint ID.",
			},
		},
	}
}

func resourceTencentCloudPrivateDnsForwardRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_forward_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = privatednsIntlv20201028.NewCreateForwardRuleRequest()
		response = privatednsIntlv20201028.NewCreateForwardRuleResponse()
		ruleId   string
	)

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_type"); ok {
		request.RuleType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_id"); ok {
		request.EndPointId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().CreateForwardRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create private dns forward rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create private dns forward rule failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.RuleId == nil {
		return fmt.Errorf("RuleId is nil.")
	}

	ruleId = *response.Response.RuleId
	d.SetId(ruleId)
	return resourceTencentCloudPrivateDnsForwardRuleRead(d, meta)
}

func resourceTencentCloudPrivateDnsForwardRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_forward_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ruleId  = d.Id()
	)

	respData, err := service.DescribePrivateDnsForwardRuleById(ctx, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_forward_rule` [%s] not found, please check if it has been deleted.\n", logId, ruleId)
		return nil
	}

	if respData.RuleName != nil {
		_ = d.Set("rule_name", respData.RuleName)
	}

	if respData.RuleType != nil {
		_ = d.Set("rule_type", respData.RuleType)
	}

	if respData.EndPointId != nil {
		_ = d.Set("end_point_id", respData.EndPointId)
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	return nil
}

func resourceTencentCloudPrivateDnsForwardRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_forward_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		ruleId = d.Id()
	)

	immutableArgs := []string{"rule_type", "zone_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"rule_name", "end_point_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := privatednsIntlv20201028.NewModifyForwardRuleRequest()
		request.RuleId = helper.String(ruleId)
		if v, ok := d.GetOk("rule_name"); ok {
			request.RuleName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("end_point_id"); ok {
			request.EndPointId = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().ModifyForwardRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update private dns forward rule failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudPrivateDnsForwardRuleRead(d, meta)
}

func resourceTencentCloudPrivateDnsForwardRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_forward_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = privatednsIntlv20201028.NewDeleteForwardRuleRequest()
		ruleId  = d.Id()
	)

	request.RuleIdSet = []*string{helper.String(ruleId)}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsIntlV20201028Client().DeleteForwardRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, PRIVATEDNS_CUSTOM_RETRY_SDK_ERROR...)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete private dns forward rule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
