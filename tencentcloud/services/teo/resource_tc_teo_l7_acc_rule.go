package teo

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoL7AccRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoL7AccRuleCreate,
		Read:   resourceTencentCloudTeoL7AccRuleRead,
		Update: resourceTencentCloudTeoL7AccRuleUpdate,
		Delete: resourceTencentCloudTeoL7AccRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone id.",
			},

			"rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Rules content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID. Unique identifier of the rule.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule status. The possible values are: `enable`: enabled; `disable`: disabled.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule name. The name length limit is 255 characters.",
						},
						"description": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Rule annotation. multiple annotations can be added.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"rule_priority": {
							Type: schema.TypeInt,
							// Optional:    true,
							Computed:    true,
							Description: "Rule priority. only used as an output parameter.",
						},
						"branches": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.",
							Elem: &schema.Resource{
								Schema: TencentTeoL7RuleBranchBasicInfo(1),
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoL7AccRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId string
	)
	var (
		request  = teov20220901.NewCreateL7AccRulesRequest()
		response = teov20220901.NewCreateL7AccRulesResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			rulesMap := item.(map[string]interface{})
			ruleEngineItem := teov20220901.RuleEngineItem{}
			if v, ok := rulesMap["status"].(string); ok && v != "" {
				ruleEngineItem.Status = helper.String(v)
			}
			if v, ok := rulesMap["rule_name"].(string); ok && v != "" {
				ruleEngineItem.RuleName = helper.String(v)
			}
			if v, ok := rulesMap["description"]; ok {
				descriptionSet := v.([]interface{})
				for i := range descriptionSet {
					description := descriptionSet[i].(string)
					ruleEngineItem.Description = append(ruleEngineItem.Description, helper.String(description))
				}
			}
			// if v, ok := rulesMap["rule_priority"].(int); ok {
			// 	ruleEngineItem.RulePriority = helper.IntInt64(v)
			// }
			if _, ok := rulesMap["branches"]; ok {
				ruleEngineItem.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(rulesMap)
			}

			request.Rules = append(request.Rules, &ruleEngineItem)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo l7 acc rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(zoneId)

	_ = response
	return resourceTencentCloudTeoL7AccRuleRead(d, meta)
}

func resourceTencentCloudTeoL7AccRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId := d.Id()

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoL7AccRuleById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_l7_acc_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	rulesList := make([]map[string]interface{}, 0, len(respData.Rules))
	if respData.Rules != nil {
		for _, rules := range respData.Rules {
			rulesMap := map[string]interface{}{}

			if rules.RuleId != nil {
				rulesMap["rule_id"] = rules.RuleId
			}

			if rules.Status != nil {
				rulesMap["status"] = rules.Status
			}

			if rules.RuleName != nil {
				rulesMap["rule_name"] = rules.RuleName
			}

			if rules.Description != nil {
				rulesMap["description"] = rules.Description
			}

			if rules.RulePriority != nil {
				rulesMap["rule_priority"] = rules.RulePriority
			}

			if rules.Branches != nil {
				rulesMap["branches"] = resourceTencentCloudTeoL7AccRuleSetBranchs(rules.Branches)
			}
			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)
	}

	_ = zoneId
	return nil
}

func resourceTencentCloudTeoL7AccRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Id()

	needChange := false
	mutableArgs := []string{"rules"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		oldRules, newRules := d.GetChange("rules")
		oldRulesList := oldRules.([]interface{})
		newRulesList := newRules.([]interface{})
		lenOldRules := len(oldRulesList)
		lenNewRules := len(newRulesList)
		if lenOldRules >= lenNewRules {
			needDelIds := make([]*string, 0, lenOldRules)
			for i := 0; i < lenOldRules; i++ {
				oldRulesMap := oldRulesList[i].(map[string]interface{})
				ruleEngineItem := teov20220901.RuleEngineItem{}
				// update
				if i < lenNewRules {
					rulesMap := newRulesList[i].(map[string]interface{})
					request := teov20220901.NewModifyL7AccRuleRequest()
					request.ZoneId = helper.String(zoneId)
					if v, ok := rulesMap["status"].(string); ok && v != "" {
						ruleEngineItem.Status = helper.String(v)
					}

					if v, ok := rulesMap["rule_name"].(string); ok && v != "" {
						ruleEngineItem.RuleName = helper.String(v)
					}

					if v, ok := rulesMap["description"]; ok {
						descriptionSet := v.([]interface{})
						for i := range descriptionSet {
							description := descriptionSet[i].(string)
							ruleEngineItem.Description = append(ruleEngineItem.Description, helper.String(description))
						}
					}

					if _, ok := rulesMap["branches"]; ok {
						ruleEngineItem.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(rulesMap)
					}

					if v, ok := oldRulesMap["rule_id"].(string); ok && v != "" {
						ruleEngineItem.RuleId = helper.String(v)
					}

					request.Rule = &ruleEngineItem
					reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyL7AccRuleWithContext(ctx, request)
						if e != nil {
							return tccommon.RetryError(e)
						} else {
							log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
						}
						return nil
					})

					if reqErr != nil {
						log.Printf("[CRITAL]%s update teo l7 acc rule failed, reason:%+v", logId, reqErr)
						return reqErr
					}
				} else {
					if v, ok := oldRulesMap["rule_id"].(string); ok && v != "" {
						needDelIds = append(needDelIds, helper.String(v))
					}
				}
			}

			// delete
			if len(needDelIds) == 0 {
				return resourceTencentCloudTeoL7AccRuleRead(d, meta)
			}

			request := teov20220901.NewDeleteL7AccRulesRequest()
			request.ZoneId = helper.String(zoneId)
			request.RuleIds = needDelIds
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteL7AccRulesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete teo l7 acc rule failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else {
			needAddRules := make([]*teov20220901.RuleEngineItem, 0, len(newRulesList))
			for i := 0; i < lenNewRules; i++ {
				rulesMap := newRulesList[i].(map[string]interface{})
				ruleEngineItem := teov20220901.RuleEngineItem{}
				// update
				if v, ok := rulesMap["status"].(string); ok && v != "" {
					ruleEngineItem.Status = helper.String(v)
				}

				if v, ok := rulesMap["rule_name"].(string); ok && v != "" {
					ruleEngineItem.RuleName = helper.String(v)
				}

				if v, ok := rulesMap["description"]; ok {
					descriptionSet := v.([]interface{})
					for i := range descriptionSet {
						description := descriptionSet[i].(string)
						ruleEngineItem.Description = append(ruleEngineItem.Description, helper.String(description))
					}
				}

				if _, ok := rulesMap["branches"]; ok {
					ruleEngineItem.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(rulesMap)
				}

				if i < lenOldRules {
					oldRulesMap := oldRulesList[i].(map[string]interface{})
					if v, ok := oldRulesMap["rule_id"].(string); ok && v != "" {
						ruleEngineItem.RuleId = helper.String(v)
					}

					request := teov20220901.NewModifyL7AccRuleRequest()
					request.ZoneId = helper.String(zoneId)
					request.Rule = &ruleEngineItem
					reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyL7AccRuleWithContext(ctx, request)
						if e != nil {
							return tccommon.RetryError(e)
						} else {
							log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
						}
						return nil
					})

					if reqErr != nil {
						log.Printf("[CRITAL]%s update teo l7 acc rule failed, reason:%+v", logId, reqErr)
						return reqErr
					}
				} else {
					needAddRules = append(needAddRules, &ruleEngineItem)
				}
			}

			// add
			if len(needAddRules) == 0 {
				return resourceTencentCloudTeoL7AccRuleRead(d, meta)
			}

			request := teov20220901.NewCreateL7AccRulesRequest()
			request.ZoneId = helper.String(zoneId)
			request.Rules = needAddRules
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRulesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s create teo l7 acc rule failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	_ = zoneId
	return resourceTencentCloudTeoL7AccRuleRead(d, meta)
}

func resourceTencentCloudTeoL7AccRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Id()

	var (
		request  = teov20220901.NewDeleteL7AccRulesRequest()
		response = teov20220901.NewDeleteL7AccRulesResponse()
	)
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("rule_ids"); ok {
		ruleIdsSet := v.([]interface{})
		for i := range ruleIdsSet {
			ruleIds := ruleIdsSet[i].(string)
			request.RuleIds = append(request.RuleIds, helper.String(ruleIds))
		}
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			rulesMap := item.(map[string]interface{})
			if v, ok := rulesMap["rule_id"].(string); ok && v != "" {
				request.RuleIds = append(request.RuleIds, helper.String(v))
			}
		}
	}

	if len(request.RuleIds) == 0 {
		return errors.New("[CRITAL]%s delete teo l7 acc rule failed, rule_ids is empty")
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteL7AccRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo l7 acc rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	_ = response
	_ = zoneId
	return nil
}
