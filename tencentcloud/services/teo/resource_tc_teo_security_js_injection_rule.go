package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoSecurityJSInjectionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityJSInjectionRuleCreate,
		Read:   resourceTencentCloudTeoSecurityJSInjectionRuleRead,
		Update: resourceTencentCloudTeoSecurityJSInjectionRuleUpdate,
		Delete: resourceTencentCloudTeoSecurityJSInjectionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"js_injection_rules": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "JavaScript injection rule configuration. Only one rule is allowed per request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule name.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule priority. Range: 0-100, smaller value means higher priority. Default: 0.",
						},
						"condition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Match condition expression, e.g. `${http.request.host} in ['www.example.com'] and ${http.request.uri.path} in ['/path']`.",
						},
						"inject_j_s": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "JavaScript injection option. Valid values: `no-injection` (do not inject JS); `inject-sdk-only` (inject SDK for all supported authentication methods, currently TC-RCE and TC-CAPTCHA).",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID, e.g. `injection-xxxxxxxx`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityJSInjectionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request           = teo.NewCreateSecurityJSInjectionRuleRequest()
		response          = teo.NewCreateSecurityJSInjectionRuleResponse()
		zoneId            string
		jsInjectionRuleId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("js_injection_rules"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		request.JSInjectionRules = []*teo.JSInjectionRule{buildJSInjectionRuleFromMap(m, "")}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityJSInjectionRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo security js injection rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo security js injection rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.JSInjectionRuleIds) == 0 || response.Response.JSInjectionRuleIds[0] == nil {
		return fmt.Errorf("JSInjectionRuleIds is empty.")
	}

	jsInjectionRuleId = *response.Response.JSInjectionRuleIds[0]
	d.SetId(strings.Join([]string{zoneId, jsInjectionRuleId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoSecurityJSInjectionRuleRead(d, meta)
}

func resourceTencentCloudTeoSecurityJSInjectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	jsInjectionRuleId := idSplit[1]

	respData, err := service.DescribeTeoSecurityJSInjectionRuleById(ctx, zoneId, jsInjectionRuleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_security_js_injection_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	ruleMap := map[string]interface{}{}

	if respData.RuleId != nil {
		ruleMap["rule_id"] = *respData.RuleId
	}

	if respData.Name != nil {
		ruleMap["name"] = *respData.Name
	}

	if respData.Priority != nil {
		ruleMap["priority"] = int(*respData.Priority)
	}

	if respData.Condition != nil {
		ruleMap["condition"] = *respData.Condition
	}

	if respData.InjectJS != nil {
		ruleMap["inject_j_s"] = *respData.InjectJS
	}

	_ = d.Set("js_injection_rules", []interface{}{ruleMap})

	return nil
}

func resourceTencentCloudTeoSecurityJSInjectionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewModifySecurityJSInjectionRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	jsInjectionRuleId := idSplit[1]

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("js_injection_rules"); ok {
		m := v.([]interface{})[0].(map[string]interface{})
		request.JSInjectionRules = []*teo.JSInjectionRule{buildJSInjectionRuleFromMap(m, jsInjectionRuleId)}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityJSInjectionRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo security js injection rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoSecurityJSInjectionRuleRead(d, meta)
}

func resourceTencentCloudTeoSecurityJSInjectionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteSecurityJSInjectionRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.JSInjectionRuleIds = []*string{helper.String(idSplit[1])}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityJSInjectionRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo security js injection rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildJSInjectionRuleFromMap converts a schema map block to *teo.JSInjectionRule.
// Pass ruleId="" when creating (API does not accept RuleId on create).
func buildJSInjectionRuleFromMap(m map[string]interface{}, ruleId string) *teo.JSInjectionRule {
	rule := &teo.JSInjectionRule{}

	if ruleId != "" {
		rule.RuleId = helper.String(ruleId)
	}

	if val, ok := m["name"].(string); ok && val != "" {
		rule.Name = helper.String(val)
	}

	if val, ok := m["priority"].(int); ok {
		rule.Priority = helper.IntInt64(val)
	}

	if val, ok := m["condition"].(string); ok && val != "" {
		rule.Condition = helper.String(val)
	}

	if val, ok := m["inject_j_s"].(string); ok && val != "" {
		rule.InjectJS = helper.String(val)
	}

	return rule
}
