package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoSecurityJsInjectionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityJsInjectionRuleCreate,
		Read:   resourceTencentCloudTeoSecurityJsInjectionRuleRead,
		Update: resourceTencentCloudTeoSecurityJsInjectionRuleUpdate,
		Delete: resourceTencentCloudTeoSecurityJsInjectionRuleDelete,
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
				Description: "JavaScript injection rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule name.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Rule priority. The smaller the value, the earlier the execution. Range 0-100, default 0.",
						},
						"condition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Match condition content. Must conform to expression syntax.",
						},
						"inject_js": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "JavaScript injection option. Valid values: `no-injection`, `inject-sdk-only`.",
						},
					},
				},
			},

			"js_injection_rule_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Rule ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityJsInjectionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewCreateSecurityJSInjectionRuleRequest()
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("js_injection_rules"); ok {
		for _, item := range v.([]interface{}) {
			ruleMap := item.(map[string]interface{})
			rule := teov20220901.JSInjectionRule{}
			if v, ok := ruleMap["name"].(string); ok && v != "" {
				rule.Name = helper.String(v)
			}
			if v, ok := ruleMap["priority"].(int); ok && v != 0 {
				rule.Priority = helper.IntInt64(v)
			}
			if v, ok := ruleMap["condition"].(string); ok && v != "" {
				rule.Condition = helper.String(v)
			}
			if v, ok := ruleMap["inject_js"].(string); ok && v != "" {
				rule.InjectJS = helper.String(v)
			}
			request.JSInjectionRules = append(request.JSInjectionRules, &rule)
		}
	}

	response := teov20220901.NewCreateSecurityJSInjectionRuleResponse()
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityJSInjectionRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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

	d.SetId(zoneId)

	if response.Response.JSInjectionRuleIds != nil {
		ruleIds := make([]string, 0, len(response.Response.JSInjectionRuleIds))
		for _, id := range response.Response.JSInjectionRuleIds {
			if id != nil {
				ruleIds = append(ruleIds, *id)
			}
		}
		_ = d.Set("js_injection_rule_ids", ruleIds)
	}

	return resourceTencentCloudTeoSecurityJsInjectionRuleRead(d, meta)
}

func resourceTencentCloudTeoSecurityJsInjectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	zoneId := d.Id()

	_ = d.Set("zone_id", zoneId)

	var allRules []*teov20220901.JSInjectionRule
	offset := int64(0)
	limit := int64(100)

	for {
		request := teov20220901.NewDescribeSecurityJSInjectionRuleRequest()
		request.ZoneId = helper.String(zoneId)
		request.Limit = helper.Int64(limit)
		request.Offset = helper.Int64(offset)

		var response *teov20220901.DescribeSecurityJSInjectionRuleResponse
		reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeSecurityJSInjectionRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s read teo security js injection rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if response == nil || response.Response == nil {
			return fmt.Errorf("Describe teo security js injection rule failed, Response is nil.")
		}

		if response.Response.JSInjectionRules != nil {
			allRules = append(allRules, response.Response.JSInjectionRules...)
		}

		if response.Response.TotalCount == nil {
			break
		}

		totalCount := *response.Response.TotalCount
		offset += limit
		if offset >= totalCount {
			break
		}
	}

	if len(allRules) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_teo_security_js_injection_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	rulesList := make([]map[string]interface{}, 0, len(allRules))
	ruleIds := make([]string, 0, len(allRules))

	for _, rule := range allRules {
		ruleMap := map[string]interface{}{}
		if rule.RuleId != nil {
			ruleMap["rule_id"] = rule.RuleId
			ruleIds = append(ruleIds, *rule.RuleId)
		}
		if rule.Name != nil {
			ruleMap["name"] = rule.Name
		}
		if rule.Priority != nil {
			ruleMap["priority"] = int(*rule.Priority)
		}
		if rule.Condition != nil {
			ruleMap["condition"] = rule.Condition
		}
		if rule.InjectJS != nil {
			ruleMap["inject_js"] = rule.InjectJS
		}
		rulesList = append(rulesList, ruleMap)
	}

	_ = d.Set("js_injection_rules", rulesList)
	_ = d.Set("js_injection_rule_ids", ruleIds)

	return nil
}

func resourceTencentCloudTeoSecurityJsInjectionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	zoneId := d.Id()

	needChange := false
	mutableArgs := []string{"js_injection_rules"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifySecurityJSInjectionRuleRequest()
		request.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("js_injection_rules"); ok {
			for _, item := range v.([]interface{}) {
				ruleMap := item.(map[string]interface{})
				rule := teov20220901.JSInjectionRule{}
				if v, ok := ruleMap["rule_id"].(string); ok && v != "" {
					rule.RuleId = helper.String(v)
				}
				if v, ok := ruleMap["name"].(string); ok && v != "" {
					rule.Name = helper.String(v)
				}
				if v, ok := ruleMap["priority"].(int); ok {
					rule.Priority = helper.IntInt64(v)
				}
				if v, ok := ruleMap["condition"].(string); ok && v != "" {
					rule.Condition = helper.String(v)
				}
				if v, ok := ruleMap["inject_js"].(string); ok && v != "" {
					rule.InjectJS = helper.String(v)
				}
				request.JSInjectionRules = append(request.JSInjectionRules, &rule)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityJSInjectionRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify teo security js injection rule failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo security js injection rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoSecurityJsInjectionRuleRead(d, meta)
}

func resourceTencentCloudTeoSecurityJsInjectionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_js_injection_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	zoneId := d.Id()

	ruleIds := make([]*string, 0)
	if v, ok := d.GetOk("js_injection_rule_ids"); ok {
		for _, item := range v.([]interface{}) {
			if s, ok := item.(string); ok && s != "" {
				ruleIds = append(ruleIds, helper.String(s))
			}
		}
	}

	if len(ruleIds) == 0 {
		log.Printf("[WARN]%s no js_injection_rule_ids found for deletion, resource may already be deleted.\n", logId)
		return nil
	}

	request := teov20220901.NewDeleteSecurityJSInjectionRuleRequest()
	request.ZoneId = helper.String(zoneId)
	request.JSInjectionRuleIds = ruleIds

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityJSInjectionRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo security js injection rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo security js injection rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
