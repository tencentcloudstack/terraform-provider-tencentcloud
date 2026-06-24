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

func ResourceTencentCloudWafApiSecSensitiveCustomApiExcludeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleCreate,
		Read:   resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleRead,
		Update: resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleUpdate,
		Delete: resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleDelete,
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
			"match_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Match type, regex, prefix, suffix, contain match mode.",
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Match content.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update timestamp.",
			},
		},
	}
}

func resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		domain   = d.Get("domain").(string)
		ruleName = d.Get("rule_name").(string)
	)

	if err := modifyWafApiSecSensitiveCustomApiExcludeRule(ctx, d, meta); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{domain, ruleName}, tccommon.FILED_SP))
	return resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.read")()
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
	request.IsQueryApiExcludeRule = helper.Bool(true)

	respData, err := service.DescribeWafApiSecSensitiveRuleListByFilter(ctx, request)
	if err != nil {
		return err
	}

	var ruleData *waf.ApiSecExcludeRule
	if respData != nil {
		for _, item := range respData.ApiExcludeRule {
			if item != nil && item.RuleName != nil && *item.RuleName == ruleName {
				ruleData = item
				break
			}
		}
	}

	if ruleData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("rule_name", ruleName)

	if ruleData.Status != nil {
		_ = d.Set("status", ruleData.Status)
	}

	if ruleData.MatchType != nil {
		_ = d.Set("match_type", ruleData.MatchType)
	}

	if ruleData.Content != nil {
		_ = d.Set("content", ruleData.Content)
	}

	if ruleData.UpdateTime != nil {
		_ = d.Set("update_time", ruleData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.update")()
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
	mutableArgs := []string{"status", "match_type", "content"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if err := modifyWafApiSecSensitiveCustomApiExcludeRule(ctx, d, meta); err != nil {
			return err
		}
	}

	return resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveCustomApiExcludeRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.delete")()
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
	request.ApiExcludeRuleName = helper.Strings([]string{ruleName})

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete waf api sec sensitive custom api exclude rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf api sec sensitive custom api exclude rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func modifyWafApiSecSensitiveCustomApiExcludeRule(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	ruleName := d.Get("rule_name").(string)
	status := d.Get("status").(int)

	request := waf.NewModifyApiSecSensitiveRuleRequest()
	request.Domain = helper.String(d.Get("domain").(string))
	request.RuleName = helper.String(ruleName)
	request.Status = helper.IntUint64(status)

	excludeRule := waf.ApiSecExcludeRule{
		RuleName: helper.String(ruleName),
		Status:   helper.IntInt64(status),
	}

	if v, ok := d.GetOk("match_type"); ok {
		excludeRule.MatchType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		excludeRule.Content = helper.String(v.(string))
	}

	request.CustomApiExcludeRule = &excludeRule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify waf api sec sensitive custom api exclude rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify waf api sec sensitive custom api exclude rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
