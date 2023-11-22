/*
Provides a resource to create a waf cc

Example Usage

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain      = "www.demo.com"
  name        = "terraform"
  status      = 1
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafCc() *schema.Resource {
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
				Description: "Session match mode, 0 use session, 1 use ip.",
			},
			"limit": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CC detection threshold.",
			},
			"interval": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Interval.",
			},
			"url": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Check URL.",
			},
			"match_func": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Match method, 0 equal, 1 contains, 2 prefix.",
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
			//"options_arr": {
			//	Optional:    true,
			//	Type:        schema.TypeString,
			//	Description: "Option param.",
			//},
			"edition": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(EDITION_TYPE),
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
				Description: "Advance mode use session id.",
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
	defer logElapsed("resource.tencentcloud_waf_cc.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
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

	//if v, ok := d.GetOk("options_arr"); ok {
	//	request.OptionsArr = helper.String(v.(string))
	//}

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
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().UpsertCCRule(request)
		if e != nil {
			return retryError(e)
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
	d.SetId(strings.Join([]string{domain, ruleId, name}, FILED_SP))

	return resourceTencentCloudWafCcRead(d, meta)
}

func resourceTencentCloudWafCcRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_cc.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
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

	//if cc.Options != nil {
	//	_ = d.Set("options_arr", cc.Options)
	//}
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
	defer logElapsed("resource.tencentcloud_waf_cc.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = waf.NewUpsertCCRuleRequest()
	)

	immutableArgs := []string{"domain", "name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
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

	//if v, ok := d.GetOk("options_arr"); ok {
	//	request.OptionsArr = helper.String(v.(string))
	//}

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().UpsertCCRule(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_waf_cc.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
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
