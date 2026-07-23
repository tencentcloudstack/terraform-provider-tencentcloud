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

func ResourceTencentCloudWafApiSecSensitiveWhiteRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafApiSecSensitiveWhiteRuleCreate,
		Read:   resourceTencentCloudWafApiSecSensitiveWhiteRuleRead,
		Update: resourceTencentCloudWafApiSecSensitiveWhiteRuleUpdate,
		Delete: resourceTencentCloudWafApiSecSensitiveWhiteRuleDelete,
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
				Description: "White rule name.",
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Rule switch, 0: off, 1: on.",
			},
			"white_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "White mode. Enum values: 1: whitelist the whole API, 2: whitelist specified fields.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description.",
			},
			"api_name_op": apiSecApiNameOpSchema(),
			"white_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "White field config list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Field name.",
						},
						"field_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Field position.",
						},
						"sensitive_types": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Sensitive data type list.",
						},
					},
				},
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Update timestamp.",
			},
		},
	}
}

func resourceTencentCloudWafApiSecSensitiveWhiteRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_white_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		domain   = d.Get("domain").(string)
		ruleName = d.Get("rule_name").(string)
	)

	if err := modifyWafApiSecSensitiveWhiteRule(ctx, d, meta); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{domain, ruleName}, tccommon.FILED_SP))
	return resourceTencentCloudWafApiSecSensitiveWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveWhiteRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_white_rule.read")()
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
	request.IsQueryApiSensitiveWhiteRule = helper.Bool(true)

	respData, err := service.DescribeWafApiSecSensitiveRuleListByFilter(ctx, request)
	if err != nil {
		return err
	}

	var ruleData *waf.ApiSecSensitiveWhiteRule
	if respData != nil {
		for _, item := range respData.ApiSecSensitiveWhiteRule {
			if item != nil && item.RuleName != nil && *item.RuleName == ruleName {
				ruleData = item
				break
			}
		}
	}

	if ruleData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_api_sec_sensitive_white_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("rule_name", ruleName)

	if ruleData.Status != nil {
		_ = d.Set("status", ruleData.Status)
	}

	if ruleData.WhiteMode != nil {
		_ = d.Set("white_mode", ruleData.WhiteMode)
	}

	if ruleData.Description != nil {
		_ = d.Set("description", ruleData.Description)
	}

	if ruleData.ApiNameOp != nil {
		_ = d.Set("api_name_op", flattenApiSecApiNameOpList(ruleData.ApiNameOp))
	}

	if ruleData.WhiteFields != nil {
		_ = d.Set("white_fields", flattenApiSecSensitiveWhiteFieldList(ruleData.WhiteFields))
	}

	if ruleData.UpdateTime != nil {
		_ = d.Set("update_time", ruleData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudWafApiSecSensitiveWhiteRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_white_rule.update")()
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
	mutableArgs := []string{"status", "white_mode", "description", "api_name_op", "white_fields"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if err := modifyWafApiSecSensitiveWhiteRule(ctx, d, meta); err != nil {
			return err
		}
	}

	return resourceTencentCloudWafApiSecSensitiveWhiteRuleRead(d, meta)
}

func resourceTencentCloudWafApiSecSensitiveWhiteRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_api_sec_sensitive_white_rule.delete")()
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
	request.ApiSecSensitiveWhiteRuleNameList = helper.Strings([]string{ruleName})
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete waf api sec sensitive white rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete waf api sec sensitive white rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func modifyWafApiSecSensitiveWhiteRule(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	ruleName := d.Get("rule_name").(string)
	status := d.Get("status").(int)

	request := waf.NewModifyApiSecSensitiveRuleRequest()
	request.Domain = helper.String(d.Get("domain").(string))
	request.RuleName = helper.String(ruleName)
	request.Status = helper.IntUint64(status)

	whiteRule := waf.ApiSecSensitiveWhiteRule{
		RuleName: helper.String(ruleName),
		Status:   helper.IntInt64(status),
	}

	if v, ok := d.GetOkExists("white_mode"); ok {
		whiteRule.WhiteMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		whiteRule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_name_op"); ok {
		whiteRule.ApiNameOp = buildApiSecApiNameOpList(v.([]interface{}))
	}

	if v, ok := d.GetOk("white_fields"); ok {
		whiteRule.WhiteFields = buildApiSecSensitiveWhiteFieldList(v.([]interface{}))
	}

	request.ApiSecSensitiveWhiteRuleRule = &whiteRule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyApiSecSensitiveRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify waf api sec sensitive white rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify waf api sec sensitive white rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func buildApiSecSensitiveWhiteFieldList(list []interface{}) []*waf.ApiSecSensitiveWhiteField {
	result := make([]*waf.ApiSecSensitiveWhiteField, 0, len(list))
	for _, item := range list {
		dMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		field := waf.ApiSecSensitiveWhiteField{}
		if v, ok := dMap["field_name"]; ok && v.(string) != "" {
			field.FieldName = helper.String(v.(string))
		}

		if v, ok := dMap["field_type"]; ok && v.(string) != "" {
			field.FieldType = helper.String(v.(string))
		}

		if v, ok := dMap["sensitive_types"]; ok {
			field.SensitiveTypes = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		}

		result = append(result, &field)
	}

	return result
}

func flattenApiSecSensitiveWhiteFieldList(list []*waf.ApiSecSensitiveWhiteField) []interface{} {
	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}

		dMap := map[string]interface{}{}
		if item.FieldName != nil {
			dMap["field_name"] = item.FieldName
		}

		if item.FieldType != nil {
			dMap["field_type"] = item.FieldType
		}

		if item.SensitiveTypes != nil {
			dMap["sensitive_types"] = item.SensitiveTypes
		}

		result = append(result, dMap)
	}

	return result
}
