package waf

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafCustomRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafCustomRuleCreate,
		Read:   resourceTencentCloudWafCustomRuleRead,
		Update: resourceTencentCloudWafCustomRuleUpdate,
		Delete: resourceTencentCloudWafCustomRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule Name.",
			},
			"sort_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Priority, value range 0-100.",
			},
			"redirect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "If the action is a redirect, it represents the redirect address; Other situations can be left blank.",
			},
			"expire_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Expiration time, measured in seconds, such as 1677254399, which means the expiration time is 2023-02-24 23:59:59 0 means never expires.",
			},
			"strategies": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Strategies detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Matching Fields.",
						},
						"compare_func": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logical symbol.",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Matching Content.",
						},
						"arg": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Matching parameters.",
						},
					},
				},
			},
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain name that needs to add policy.",
			},
			"action_type": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CUSTOM_RULE_ACTION_TYPE),
				Description:  "Action type, 1 represents blocking, 2 represents captcha, 3 represents observation, and 4 represents redirection.",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CUSTOM_RULE_STATUS),
				Default:      CUSTOM_RULE_STATUS_1,
				Description:  "The status of the switch, 1 is on, 0 is off, default 1.",
			},
			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "rule ID.",
			},
		},
	}
}

func resourceTencentCloudWafCustomRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = waf.NewAddCustomRuleRequest()
		response      = waf.NewAddCustomRuleResponse()
		statusRequest = waf.NewModifyCustomRuleStatusRequest()
		domain        string
		ruleIdStr     string
		status        string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_id"); ok {
		request.SortId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("redirect"); ok {
		request.Redirect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expire_time"); ok {
		request.ExpireTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			strategy := waf.Strategy{}
			if v, ok := dMap["field"]; ok {
				strategy.Field = helper.String(v.(string))
			}

			if v, ok := dMap["compare_func"]; ok {
				strategy.CompareFunc = helper.String(v.(string))
			}

			if v, ok := dMap["content"]; ok {
				strategy.Content = helper.String(v.(string))
			}

			if v, ok := dMap["arg"]; ok {
				strategy.Arg = helper.String(v.(string))
			}

			request.Strategies = append(request.Strategies, &strategy)
		}
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("action_type"); ok {
		request.ActionType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().AddCustomRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf CustomRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId := *response.Response.RuleId
	ruleIdStr = strconv.FormatInt(ruleId, 10)

	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}

	if status == CUSTOM_RULE_STATUS_0 {
		statusRequest.Domain = &domain
		tmpRuleId := uint64(ruleId)
		statusRequest.RuleId = &tmpRuleId
		statusRequest.Status = helper.IntUint64(CUSTOM_RULE_STATUS_0_INT)
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyCustomRuleStatus(statusRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify waf CustomRule status failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(strings.Join([]string{domain, ruleIdStr}, tccommon.FILED_SP))
	return resourceTencentCloudWafCustomRuleRead(d, meta)
}

func resourceTencentCloudWafCustomRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	customRule, err := service.DescribeWafCustomRuleById(ctx, domain, ruleId)
	if err != nil {
		return err
	}

	if customRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafCustomRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if customRule.Name != nil {
		_ = d.Set("name", customRule.Name)
	}

	if customRule.SortId != nil {
		_ = d.Set("sort_id", customRule.SortId)
	}

	if customRule.Redirect != nil {
		_ = d.Set("redirect", customRule.Redirect)
	}

	if customRule.ExpireTime != nil {
		_ = d.Set("expire_time", customRule.ExpireTime)
	}

	if customRule.Strategies != nil {
		strategiesList := []interface{}{}
		for _, strategies := range customRule.Strategies {
			strategiesMap := map[string]interface{}{}

			if strategies.Field != nil {
				strategiesMap["field"] = strategies.Field
			}

			if strategies.CompareFunc != nil {
				strategiesMap["compare_func"] = strategies.CompareFunc
			}

			if strategies.Content != nil {
				strategiesMap["content"] = strategies.Content
			}

			if strategies.Arg != nil {
				strategiesMap["arg"] = strategies.Arg
			}

			strategiesList = append(strategiesList, strategiesMap)
		}

		_ = d.Set("strategies", strategiesList)

	}

	_ = d.Set("domain", domain)

	if customRule.ActionType != nil {
		_ = d.Set("action_type", customRule.ActionType)
	}

	if customRule.Status != nil {
		_ = d.Set("status", customRule.Status)
	}

	if customRule.RuleId != nil {
		_ = d.Set("rule_id", customRule.RuleId)
	}

	return nil
}

func resourceTencentCloudWafCustomRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = waf.NewModifyCustomRuleRequest()
		statusRequest = waf.NewModifyCustomRuleStatusRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	immutableArgs := []string{"domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.Domain = &domain
	ruleIdInt, _ := strconv.ParseInt(ruleId, 10, 64)
	ruleIdUInt := uint64(ruleIdInt)
	request.RuleId = &ruleIdUInt

	if v, ok := d.GetOk("name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("action_type"); ok {
		request.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("redirect"); ok {
		request.Redirect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_id"); ok {
		sortIdInt, _ := strconv.ParseInt(v.(string), 10, 64)
		sortIdUInt := uint64(sortIdInt)
		request.SortId = &sortIdUInt
	}

	if v, ok := d.GetOk("expire_time"); ok {
		expireTimeInt, _ := strconv.ParseInt(v.(string), 10, 64)
		expireTimeUInt := uint64(expireTimeInt)
		request.ExpireTime = &expireTimeUInt
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			strategy := waf.Strategy{}
			if v, ok := dMap["field"]; ok {
				strategy.Field = helper.String(v.(string))
			}
			if v, ok := dMap["compare_func"]; ok {
				strategy.CompareFunc = helper.String(v.(string))
			}
			if v, ok := dMap["content"]; ok {
				strategy.Content = helper.String(v.(string))
			}
			if v, ok := dMap["arg"]; ok {
				strategy.Arg = helper.String(v.(string))
			}
			request.Strategies = append(request.Strategies, &strategy)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyCustomRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf CustomRule failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			statusRequest.Domain = &domain
			statusRequest.RuleId = &ruleIdUInt
			if status == CUSTOM_RULE_STATUS_0 {
				statusRequest.Status = helper.IntUint64(CUSTOM_RULE_STATUS_0_INT)
			} else {
				statusRequest.Status = helper.IntUint64(CUSTOM_RULE_STATUS_1_INT)
			}

			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyCustomRuleStatus(statusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify waf CustomRule status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafCustomRuleRead(d, meta)
}

func resourceTencentCloudWafCustomRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_custom_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	domain := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteWafCustomRuleById(ctx, domain, ruleId); err != nil {
		return err
	}

	return nil
}
