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

func ResourceTencentCloudWafCc() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafCcCreate,
		Read:   resourceTencentCloudWafCcRead,
		Update: resourceTencentCloudWafCcUpdate,
		Delete: resourceTencentCloudWafCcDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule Name.",
			},
			"status": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Rule Status, 0 rule close, 1 rule open.",
			},
			"advance": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Advanced mode (whether to use session detection). 0(disabled) 1(enabled).",
			},
			"limit": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CC detection threshold.",
			},
			"interval": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CC detection cycle.",
			},
			"url": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Detection URL.",
			},
			"match_func": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Match method, 0(equal), 1(prefix), 2(contains), 3(not equal), 6(suffix), 7(not contains).",
			},
			"action_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule Action, 20 log, 21 captcha, 22 deny, 23 accurate deny.",
			},
			"priority": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Rule Priority.",
			},
			"valid_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Action ValidTime, minute unit. Min: 60, Max: 604800.",
			},
			"options_arr": {
				Optional: true,
				Type:     schema.TypeString,
				Description: `JSON serialized string of CC matching conditions, example:[{\"key\":\"Method\",\"args\":[\"=R0VU\"],\"match\":\"0\",\"encodeflag\":true}]
        Key optional values are Method, Post, Referer, Cookie, User-Agent, CustomHeader
        Match optional values are, when Key is Method, optional values are 0 (equal), 3 (not equal).
        When the key is Post, the optional values are 0 (equal to), 3 (not equal to), when the key is Cookie, the optional values are 0 (equal to), 2 (included), 3 (not equal to), 7 (not included),
        When the key is Referer, the optional values are 0 (equal to), 3 (not equal to), 1 (prefix match), 6 (suffix match), 2 (included), 7 (not included), 12 (exists), 5 (not included), 4 (content is empty),
        When the key is Cookie, the optional values are 0 (equal to), 3 (not equal to), 2 (included), 7 (not included),
        When the key is User-Agent, the optional values are 0 (equal to), 3 (not equal to), 1 (prefix match), 6 (suffix match), 2 (included), 7 (not included), 12 (exists), 5 (not included), 4 (content is empty),
        When the key is CustomHeader, the optional values are 0 (equal to), 3 (not equal to), 2 (included), 7 (not included), 12 (exists), 5 (not included), 4 (content is empty).
        args is used to indicate the matching content. You need to set encodeflag to true. When the Key is Post, Cookie, or CustomHeader, use the equal sign = to concatenate the Key and Value, and encode them using Base64, similar to YWJj=YWJj. When the Key is Referer or User-Agent, use the equal sign = to concatenate the Value, similar to =YWJj.`,
			},
			"edition": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(EDITION_TYPE),
				Description:  "WAF edition. clb-waf means clb-waf, sparta-waf means saas-waf.",
			},
			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Operate Type.",
			},
			"event_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Event ID.",
			},
			"session_applied": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Session ID that needs to be enabled for the rule.",
			},
			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Rule ID.",
			},
		},
	}
}

func resourceTencentCloudWafCcCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_cc.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = waf.NewUpsertCCRuleRequest()
		response = waf.NewUpsertCCRuleResponse()
		domain   string
		ruleId   string
		name     string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
		name = v.(string)
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("advance"); ok {
		request.Advance = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("match_func"); ok {
		request.MatchFunc = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("action_type"); ok {
		request.ActionType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("valid_time"); ok {
		request.ValidTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("options_arr"); ok {
		request.OptionsArr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("edition"); ok {
		request.Edition = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("event_id"); ok {
		request.EventId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_applied"); ok {
		sessionAppliedSet := v.(*schema.Set).List()
		for i := range sessionAppliedSet {
			sessionApplied := sessionAppliedSet[i].(int)
			request.SessionApplied = append(request.SessionApplied, helper.IntInt64(sessionApplied))
		}
	}

	request.RuleId = helper.IntInt64(0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().UpsertCCRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.RuleId == nil {
			e = fmt.Errorf("waf cc failed not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf cc failed, reason:%+v", logId, err)
		return err
	}

	ruleIdInt := *response.Response.RuleId
	ruleId = strconv.FormatInt(ruleIdInt, 10)
	d.SetId(strings.Join([]string{domain, ruleId, name}, tccommon.FILED_SP))

	return resourceTencentCloudWafCcRead(d, meta)
}

func resourceTencentCloudWafCcRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_cc.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	domain := idSplit[0]
	ruleId := idSplit[1]

	cc, err := service.DescribeWafCcById(ctx, domain, ruleId)
	if err != nil {
		return err
	}

	if cc == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafCc` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if cc.Name != nil {
		_ = d.Set("name", cc.Name)
	}

	if cc.Status != nil {
		_ = d.Set("status", cc.Status)
	}

	if cc.Advance != nil {
		advanceStr := strconv.FormatUint(*cc.Advance, 10)
		_ = d.Set("advance", advanceStr)
	}

	if cc.Limit != nil {
		limitStr := strconv.FormatUint(*cc.Limit, 10)
		_ = d.Set("limit", limitStr)
	}

	if cc.Interval != nil {
		intervalStr := strconv.FormatUint(*cc.Interval, 10)
		_ = d.Set("interval", intervalStr)
	}

	if cc.Url != nil {
		_ = d.Set("url", cc.Url)
	}

	if cc.MatchFunc != nil {
		_ = d.Set("match_func", cc.MatchFunc)
	}

	if cc.ActionType != nil {
		actionTypeStr := strconv.FormatUint(*cc.ActionType, 10)
		_ = d.Set("action_type", actionTypeStr)
	}

	if cc.Priority != nil {
		_ = d.Set("priority", cc.Priority)
	}

	if cc.ValidTime != nil {
		_ = d.Set("valid_time", cc.ValidTime)
	}

	if cc.Options != nil {
		_ = d.Set("options_arr", cc.Options)
	}
	//
	//if cc.Edition != nil {
	//	_ = d.Set("edition", cc.Edition)
	//}
	//
	//if cc.Type != nil {
	//	_ = d.Set("type", cc.Type)
	//}

	if cc.EventId != nil {
		_ = d.Set("event_id", cc.EventId)
	}

	if cc.SessionApplied != nil {
		_ = d.Set("session_applied", cc.SessionApplied)
	}

	if cc.RuleId != nil {
		ruleIdStr := strconv.FormatUint(*cc.RuleId, 10)
		_ = d.Set("rule_id", ruleIdStr)
	}

	return nil
}

func resourceTencentCloudWafCcUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_cc.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = waf.NewUpsertCCRuleRequest()
	)

	immutableArgs := []string{"domain", "name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	domain := idSplit[0]
	ruleId := idSplit[1]
	name := idSplit[2]

	request.Domain = &domain
	ruleIdInt, _ := strconv.ParseInt(ruleId, 10, 64)
	request.RuleId = &ruleIdInt
	request.Name = &name

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("advance"); ok {
		request.Advance = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("match_func"); ok {
		request.MatchFunc = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("action_type"); ok {
		request.ActionType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("valid_time"); ok {
		request.ValidTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("options_arr"); ok {
		request.OptionsArr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("edition"); ok {
		request.Edition = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("event_id"); ok {
		request.EventId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_applied"); ok {
		sessionAppliedSet := v.(*schema.Set).List()
		for i := range sessionAppliedSet {
			sessionApplied := sessionAppliedSet[i].(int)
			request.SessionApplied = append(request.SessionApplied, helper.IntInt64(sessionApplied))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().UpsertCCRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf cc failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafCcRead(d, meta)
}

func resourceTencentCloudWafCcDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_cc.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	domain := idSplit[0]
	ruleId := idSplit[1]
	name := idSplit[2]

	if err := service.DeleteWafCcById(ctx, domain, ruleId, name); err != nil {
		return err
	}

	return nil
}
