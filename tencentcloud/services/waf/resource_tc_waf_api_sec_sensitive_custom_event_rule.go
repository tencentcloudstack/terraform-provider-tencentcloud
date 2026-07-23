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

func ResourceTencentCloudWafApiSecSensitiveCustomEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafApiSecSensitiveCustomEventRuleCreate,
		Read:   resourceTencentCloudWafApiSecSensitiveCustomEventRuleRead,
		Update: resourceTencentCloudWafApiSecSensitiveCustomEventRuleUpdate,
		Delete: resourceTencentCloudWafApiSecSensitiveCustomEventRuleDelete,
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Event description.",
			},
			"req_frequency": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Access frequency, the first field represents the count, the second field represents the minute.",
			},
			"risk_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Risk level, the value is 100, 200, 300, respectively representing low, medium, high risk.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule source.",
			},
			"api_name_op":     apiSecApiNameOpSchema(),
			"match_rule_list": apiSecSceneRuleEntrySchema("Match rule list."),
			"stat_rule_list":  apiSecSceneRuleEntrySchema("Statistics rule list."),
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update timestamp.",
			},
		},
	}
}

func resourceTencentCloudWafApiSecSensitiveCustomEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_event_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		domain   = d.Get("domain").(string)
		ruleName = d.Get("rule_name").(string)
	)

	if err := modifyWafApiSecSensitiveCustomEventRule(ctx, d, meta); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{domain, ruleName}, tccommon.FILED_SP))
	return resourceTencentCloudWafApiSecSensitiveCustomEventRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveCustomEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_event_rule.read")()
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
	request.IsQueryApiCustomEventRule = helper.Bool(true)

	respData, err := service.DescribeWafApiSecSensitiveRuleListByFilter(ctx, request)
	if err != nil {
		return err
	}

	var ruleData *waf.ApiSecCustomEventRule
	if respData != nil {
		for _, item := range respData.ApiSecCustomEventRule {
			if item != nil && item.RuleName != nil && *item.RuleName == ruleName {
				ruleData = item
				break
			}
		}
	}

	if ruleData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_api_sec_sensitive_custom_event_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("rule_name", ruleName)

	if ruleData.Status != nil {
		_ = d.Set("status", ruleData.Status)
	}

	if ruleData.Description != nil {
		_ = d.Set("description", ruleData.Description)
	}

	if ruleData.ReqFrequency != nil {
		_ = d.Set("req_frequency", ruleData.ReqFrequency)
	}

	if ruleData.RiskLevel != nil {
		_ = d.Set("risk_level", ruleData.RiskLevel)
	}

	if ruleData.Source != nil {
		_ = d.Set("source", ruleData.Source)
	}

	if ruleData.ApiNameOp != nil {
		_ = d.Set("api_name_op", flattenApiSecApiNameOpList(ruleData.ApiNameOp))
	}

	if ruleData.MatchRuleList != nil {
		_ = d.Set("match_rule_list", flattenApiSecSceneRuleEntryList(ruleData.MatchRuleList))
	}

	if ruleData.StatRuleList != nil {
		_ = d.Set("stat_rule_list", flattenApiSecSceneRuleEntryList(ruleData.StatRuleList))
	}

	if ruleData.UpdateTime != nil {
		_ = d.Set("update_time", ruleData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudWafApiSecSensitiveCustomEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_event_rule.update")()
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
	mutableArgs := []string{"status", "description", "req_frequency", "risk_level", "source", "api_name_op", "match_rule_list", "stat_rule_list"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if err := modifyWafApiSecSensitiveCustomEventRule(ctx, d, meta); err != nil {
			return err
		}
	}

	return resourceTencentCloudWafApiSecSensitiveCustomEventRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveCustomEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_custom_event_rule.delete")()
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
	request.ApiSecCustomEventRuleNameList = helper.Strings([]string{ruleName})

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete waf api sec sensitive custom event rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf api sec sensitive custom event rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func modifyWafApiSecSensitiveCustomEventRule(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	ruleName := d.Get("rule_name").(string)
	status := d.Get("status").(int)

	request := waf.NewModifyApiSecSensitiveRuleRequest()
	request.Domain = helper.String(d.Get("domain").(string))
	request.RuleName = helper.String(ruleName)
	request.Status = helper.IntUint64(status)

	eventRule := waf.ApiSecCustomEventRule{
		RuleName: helper.String(ruleName),
		Status:   helper.IntInt64(status),
	}

	if v, ok := d.GetOk("description"); ok {
		eventRule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("req_frequency"); ok {
		for _, item := range v.([]interface{}) {
			eventRule.ReqFrequency = append(eventRule.ReqFrequency, helper.IntInt64(item.(int)))
		}
	}

	if v, ok := d.GetOk("risk_level"); ok {
		eventRule.RiskLevel = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source"); ok {
		eventRule.Source = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_name_op"); ok {
		eventRule.ApiNameOp = buildApiSecApiNameOpList(v.([]interface{}))
	}

	if v, ok := d.GetOk("match_rule_list"); ok {
		eventRule.MatchRuleList = buildApiSecSceneRuleEntryList(v.([]interface{}))
	}

	if v, ok := d.GetOk("stat_rule_list"); ok {
		eventRule.StatRuleList = buildApiSecSceneRuleEntryList(v.([]interface{}))
	}

	request.ApiSecCustomEventRuleRule = &eventRule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify waf api sec sensitive custom event rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify waf api sec sensitive custom event rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
