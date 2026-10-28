package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafApiSecSensitiveCustomRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafApiSecSensitiveCustomRuleCreate,
		Read:   resourceTencentCloudWafApiSecSensitiveCustomRuleRead,
		Update: resourceTencentCloudWafApiSecSensitiveCustomRuleUpdate,
		Delete: resourceTencentCloudWafApiSecSensitiveCustomRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule name.",
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Rule switch, 0: off, 1: on.",
			},
			"position": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Parameter position.",
			},
			"match_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Match condition.",
			},
			"match_value": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Match value.",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Risk level.",
			},
			"match_cond": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Match symbol, pass this value when the match condition is keyword match or character match, multiple values can be passed.",
			},
			"is_pan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether the rule is generalized, default 0 means not generalized.",
			},
		},
	}
}

func resourceTencentCloudWafApiSecSensitiveCustomRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		domain   = d.Get("domain").(string)
		ruleName = d.Get("rule_name").(string)
	)

	if err := modifyWafApiSecSensitiveCustomRule(ctx, d, meta); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{domain, ruleName}, tccommon.FILED_SP))
	return resourceTencentCloudWafApiSecSensitiveCustomRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveCustomRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id: %s", d.Id())
	}

	domain := idSplit[0]
	ruleName := idSplit[1]

	request := waf.NewDescribeApiSecSensitiveRuleListRequest()
	request.Domain = helper.String(domain)
	request.RuleName = helper.String(ruleName)

	respData, err := service.DescribeWafApiSecSensitiveRuleListByFilter(ctx, request)
	if err != nil {
		return err
	}

	var ruleData *waf.ApiSecSensitiveRule
	if respData != nil {
		for _, item := range respData.Data {
			if item != nil && item.RuleName != nil && *item.RuleName == ruleName {
				ruleData = item
				break
			}
		}
	}

	if ruleData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_api_sec_sensitive_custom_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("rule_name", ruleName)

	if ruleData.Status != nil {
		_ = d.Set("status", ruleData.Status)
	}

	if ruleData.CustomRule != nil {
		customRule := ruleData.CustomRule
		if customRule.Position != nil {
			_ = d.Set("position", customRule.Position)
		}

		if customRule.MatchKey != nil {
			_ = d.Set("match_key", customRule.MatchKey)
		}

		if customRule.MatchValue != nil {
			_ = d.Set("match_value", customRule.MatchValue)
		}

		if customRule.Level != nil {
			_ = d.Set("level", customRule.Level)
		}

		if customRule.MatchCond != nil {
			_ = d.Set("match_cond", customRule.MatchCond)
		}

		if customRule.IsPan != nil {
			_ = d.Set("is_pan", customRule.IsPan)
		}
	}

	return nil
}

func resourceTencentCloudWafApiSecSensitiveCustomRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id: %s", d.Id())
	}

	needChange := false
	mutableArgs := []string{"status", "position", "match_key", "match_value", "level", "match_cond", "is_pan"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if err := modifyWafApiSecSensitiveCustomRule(ctx, d, meta); err != nil {
			return err
		}
	}

	return resourceTencentCloudWafApiSecSensitiveCustomRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveCustomRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, id: %s", d.Id())
	}

	domain := idSplit[0]
	ruleName := idSplit[1]

	request := waf.NewModifyApiSecSensitiveRuleRequest()
	request.Domain = helper.String(domain)
	request.RuleName = helper.String(ruleName)
	// Status 3 means delete the rule.
	request.Status = helper.IntUint64(3)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete waf api sec sensitive custom rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf api sec sensitive custom rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func modifyWafApiSecSensitiveCustomRule(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	request := waf.NewModifyApiSecSensitiveRuleRequest()
	request.Domain = helper.String(d.Get("domain").(string))
	request.RuleName = helper.String(d.Get("rule_name").(string))
	request.Status = helper.IntUint64(d.Get("status").(int))

	customRule := waf.ApiSecCustomSensitiveRule{}
	if v, ok := d.GetOk("position"); ok {
		customRule.Position = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("match_key"); ok {
		customRule.MatchKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("match_value"); ok {
		customRule.MatchValue = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("level"); ok {
		customRule.Level = helper.String(v.(string))
	}

	if v, ok := d.GetOk("match_cond"); ok {
		customRule.MatchCond = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	}

	if v, ok := d.GetOkExists("is_pan"); ok {
		customRule.IsPan = helper.IntInt64(v.(int))
	}

	request.CustomRule = &customRule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify waf api sec sensitive custom rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify waf api sec sensitive custom rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
