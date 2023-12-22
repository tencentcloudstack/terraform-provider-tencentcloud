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

func ResourceTencentCloudWafAntiInfoLeak() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafAntiInfoLeakCreate,
		Read:   resourceTencentCloudWafAntiInfoLeakRead,
		Update: resourceTencentCloudWafAntiInfoLeakUpdate,
		Delete: resourceTencentCloudWafAntiInfoLeakDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"action_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ANTI_INFO_LEAK_ACTION_TYPE),
				Description:  "Rule Action. 0: alarm; 1: replacement; 2: only displaying the first four digits; 3: only displaying the last four digits; 4: blocking.",
			},
			"strategies": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Strategies detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(STRATEGIES_FIELD),
							Description:  "Matching Fields. support: returncode, keywords, information.",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Matching Content. If field is returncode support: 400, 403, 404, 4xx, 500, 501, 502, 504, 5xx; If field is information support: idcard, phone, bankcard; If field is keywords users input matching content themselves.",
						},
					},
				},
			},
			"uri": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Uri.",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      ANTI_INFO_LEAK_RULE_STATUS_1,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ANTI_INFO_LEAK_RULE_STATUS),
				Description:  "status.",
			},
		},
	}
}

func resourceTencentCloudWafAntiInfoLeakCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_info_leak.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = waf.NewAddAntiInfoLeakRulesRequest()
		response = waf.NewAddAntiInfoLeakRulesResponse()
		ruleId   string
		domain   string
	)

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("action_type"); ok {
		request.ActionType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			strategyForAntiInfoLeak := waf.StrategyForAntiInfoLeak{}
			if v, ok := dMap["field"]; ok {
				strategyForAntiInfoLeak.Field = helper.String(v.(string))
			}

			if v, ok := dMap["content"]; ok {
				strategyForAntiInfoLeak.Content = helper.String(v.(string))
			}

			strategyForAntiInfoLeak.CompareFunc = helper.String("contains")
			request.Strategies = append(request.Strategies, &strategyForAntiInfoLeak)
		}
	}

	if v, ok := d.GetOk("uri"); ok {
		request.Uri = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().AddAntiInfoLeakRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf antiInfoLeak failed, reason:%+v", logId, err)
		return err
	}

	ruleIdInt := *response.Response.RuleId
	ruleId = strconv.FormatInt(ruleIdInt, 10)
	d.SetId(strings.Join([]string{ruleId, domain}, tccommon.FILED_SP))

	// set status
	if v, ok := d.GetOkExists("status"); ok {
		status := v.(int)
		if status != ANTI_INFO_LEAK_RULE_STATUS_1 {
			modifyAntiInfoLeakRuleStatusRequest := waf.NewModifyAntiInfoLeakRuleStatusRequest()
			idUInt, _ := strconv.ParseUint(ruleId, 10, 64)
			modifyAntiInfoLeakRuleStatusRequest.Domain = &domain
			modifyAntiInfoLeakRuleStatusRequest.RuleId = &idUInt
			modifyAntiInfoLeakRuleStatusRequest.Status = helper.IntUint64(v.(int))
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAntiInfoLeakRuleStatus(modifyAntiInfoLeakRuleStatusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyAntiInfoLeakRuleStatusRequest.GetAction(), modifyAntiInfoLeakRuleStatusRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf antiInfoLeak status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafAntiInfoLeakRead(d, meta)
}

func resourceTencentCloudWafAntiInfoLeakRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_info_leak.read")()
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
	ruleId := idSplit[0]
	domain := idSplit[1]

	antiInfoLeak, err := service.DescribeWafAntiInfoLeakById(ctx, ruleId, domain)
	if err != nil {
		return err
	}

	if antiInfoLeak == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafAntiInfoLeak` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if antiInfoLeak.Name != nil {
		_ = d.Set("name", antiInfoLeak.Name)
	}

	if antiInfoLeak.Action != nil {
		actionInt, _ := strconv.Atoi(*antiInfoLeak.Action)
		_ = d.Set("action_type", actionInt)
	}

	if antiInfoLeak.Strategies != nil {
		strategiesList := []interface{}{}
		for _, strategies := range antiInfoLeak.Strategies {
			strategiesMap := map[string]interface{}{}

			if strategies.Field != nil {
				strategiesMap["field"] = strategies.Field
			}

			if strategies.Content != nil {
				strategiesMap["content"] = strategies.Content
			}

			strategiesList = append(strategiesList, strategiesMap)
		}

		_ = d.Set("strategies", strategiesList)
	}

	if antiInfoLeak.Uri != nil {
		_ = d.Set("uri", antiInfoLeak.Uri)
	}

	if antiInfoLeak.Status != nil {
		_ = d.Set("status", antiInfoLeak.Status)
	}

	return nil
}

func resourceTencentCloudWafAntiInfoLeakUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_info_leak.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                               = tccommon.GetLogId(tccommon.ContextNil)
		modifyAntiInfoLeakRulesRequest      = waf.NewModifyAntiInfoLeakRulesRequest()
		modifyAntiInfoLeakRuleStatusRequest = waf.NewModifyAntiInfoLeakRuleStatusRequest()
	)

	immutableArgs := []string{"domain", "uri"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	ruleId := idSplit[0]
	domain := idSplit[1]

	ruleIdUInt, _ := strconv.ParseUint(ruleId, 10, 64)
	modifyAntiInfoLeakRulesRequest.RuleId = &ruleIdUInt
	modifyAntiInfoLeakRulesRequest.Domain = &domain

	if v, ok := d.GetOk("name"); ok {
		modifyAntiInfoLeakRulesRequest.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("action_type"); ok {
		modifyAntiInfoLeakRulesRequest.ActionType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("strategies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			strategyForAntiInfoLeak := waf.StrategyForAntiInfoLeak{}
			if v, ok := dMap["field"]; ok {
				strategyForAntiInfoLeak.Field = helper.String(v.(string))
			}

			if v, ok := dMap["content"]; ok {
				strategyForAntiInfoLeak.Content = helper.String(v.(string))
			}

			strategyForAntiInfoLeak.CompareFunc = helper.String("contains")
			modifyAntiInfoLeakRulesRequest.Strategies = append(modifyAntiInfoLeakRulesRequest.Strategies, &strategyForAntiInfoLeak)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAntiInfoLeakRules(modifyAntiInfoLeakRulesRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyAntiInfoLeakRulesRequest.GetAction(), modifyAntiInfoLeakRulesRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf antiInfoLeak failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			modifyAntiInfoLeakRuleStatusRequest.Status = helper.IntUint64(v.(int))
		}

		modifyAntiInfoLeakRuleStatusRequest.Domain = &domain
		modifyAntiInfoLeakRuleStatusRequest.RuleId = &ruleIdUInt
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyAntiInfoLeakRuleStatus(modifyAntiInfoLeakRuleStatusRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyAntiInfoLeakRuleStatusRequest.GetAction(), modifyAntiInfoLeakRuleStatusRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update waf antiFake status failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudWafAntiInfoLeakRead(d, meta)
}

func resourceTencentCloudWafAntiInfoLeakDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_anti_info_leak.delete")()
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
	ruleId := idSplit[0]
	domain := idSplit[1]

	if err := service.DeleteWafAntiInfoLeakById(ctx, ruleId, domain); err != nil {
		return err
	}

	return nil
}
