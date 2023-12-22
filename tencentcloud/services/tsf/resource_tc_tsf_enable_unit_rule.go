package tsf

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
)

func ResourceTencentCloudTsfEnableUnitRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfEnableUnitRuleCreate,
		Read:   resourceTencentCloudTsfEnableUnitRuleRead,
		Update: resourceTencentCloudTsfEnableUnitRuleUpdate,
		Delete: resourceTencentCloudTsfEnableUnitRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "api ID.",
			},

			"switch": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "switch, on: `enabled`, off: `disabled`.",
			},
		},
	}
}

func resourceTencentCloudTsfEnableUnitRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_enable_unit_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var id string
	if v, ok := d.GetOk("rule_id"); ok {
		id = v.(string)
	}

	d.SetId(id)

	return resourceTencentCloudTsfEnableUnitRuleUpdate(d, meta)
}

func resourceTencentCloudTsfEnableUnitRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_enable_unit_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	enableUnitRule, err := service.DescribeTsfEnableUnitRuleById(ctx, id)
	if err != nil {
		return err
	}

	if enableUnitRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfEnableUnitRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if enableUnitRule.Id != nil {
		_ = d.Set("rule_id", enableUnitRule.Id)
	}

	if enableUnitRule.Status != nil {
		_ = d.Set("switch", enableUnitRule.Status)
	}

	return nil
}

func resourceTencentCloudTsfEnableUnitRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_enable_unit_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	if v, ok := d.GetOk("switch"); ok {
		if v.(string) == "enabled" {
			request := tsf.NewEnableUnitRuleRequest()
			request.Id = &id
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().EnableUnitRule(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf enableUnitRule failed, reason:%+v", logId, err)
				return err
			}
		}

		if v.(string) == "disabled" {
			request := tsf.NewDisableUnitRuleRequest()
			request.Id = &id
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().DisableUnitRule(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s update tsf disableUnitRule failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudTsfEnableUnitRuleRead(d, meta)
}

func resourceTencentCloudTsfEnableUnitRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_enable_unit_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
