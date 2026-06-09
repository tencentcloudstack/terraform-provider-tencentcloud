package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigCompliancePack() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigCompliancePackCreate,
		Read:   resourceTencentCloudConfigCompliancePackRead,
		Update: resourceTencentCloudConfigCompliancePackUpdate,
		Delete: resourceTencentCloudConfigCompliancePackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"compliance_pack_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Compliance pack name.",
			},

			"risk_level": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},

			"config_rules": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of compliance pack rules.",
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					identifier, _ := m["identifier"].(string)
					return schema.HashString(identifier)
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule identifier (managed rule name or custom rule cloud function ARN).",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Rule name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Rule description.",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Rule risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
						},
						"managed_rule_identifier": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Managed rule identifier (preset rule identity).",
						},
						"config_rule_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Config rule ID.",
						},
						"compliance_pack_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Compliance pack ID that this rule belongs to.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule status.",
						},
						"compliance_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance result. Valid values: COMPLIANT, NON_COMPLIANT.",
						},
						"input_parameter": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Rule input parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter key.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Parameter type: Require or Optional.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
					},
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the compliance pack.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Compliance pack status. Valid values: ACTIVE, UN_ACTIVE.",
			},

			// Computed
			"compliance_pack_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compliance pack ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the compliance pack.",
			},
		},
	}
}

func resourceTencentCloudConfigCompliancePackCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_compliance_pack.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = configv20220802.NewAddCompliancePackRequest()
		response = configv20220802.NewAddCompliancePackResponse()
	)

	if v, ok := d.GetOk("compliance_pack_name"); ok {
		request.CompliancePackName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("risk_level"); ok {
		request.RiskLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_rules"); ok {
		request.ConfigRules = buildCompliancePackRules(v.(*schema.Set).List())
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().AddCompliancePackWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create config compliance pack failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create config compliance pack failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.CompliancePackId == nil {
		return fmt.Errorf("CompliancePackId is nil")
	}

	compliancePackId := *response.Response.CompliancePackId
	d.SetId(compliancePackId)

	// Handle status if provided and not ACTIVE (default)
	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "UN_ACTIVE" {
			if err := updateCompliancePackStatus(ctx, meta, compliancePackId, status, logId); err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudConfigCompliancePackRead(d, meta)
}

func resourceTencentCloudConfigCompliancePackRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_compliance_pack.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	compliancePackId := d.Id()
	respData, err := service.DescribeConfigCompliancePackById(ctx, compliancePackId)
	if err != nil {
		return err
	}

	if respData == nil || respData.CompliancePackId == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_compliance_pack` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.CompliancePackName != nil {
		_ = d.Set("compliance_pack_name", respData.CompliancePackName)
	}

	if respData.RiskLevel != nil {
		_ = d.Set("risk_level", respData.RiskLevel)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.CompliancePackId != nil {
		_ = d.Set("compliance_pack_id", respData.CompliancePackId)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ConfigRules != nil && len(respData.ConfigRules) > 0 {
		configRulesList := flattenComplianceConfigRules(respData.ConfigRules)
		_ = d.Set("config_rules", configRulesList)
	}

	return nil
}

func resourceTencentCloudConfigCompliancePackUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_compliance_pack.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		compliancePackId = d.Id()
	)

	contentArgs := []string{"compliance_pack_name", "risk_level", "config_rules", "description"}
	needContentUpdate := false
	for _, v := range contentArgs {
		if d.HasChange(v) {
			needContentUpdate = true
			break
		}
	}

	if needContentUpdate {
		request := configv20220802.NewUpdateCompliancePackRequest()
		request.CompliancePackId = &compliancePackId

		if v, ok := d.GetOk("compliance_pack_name"); ok {
			request.CompliancePackName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("risk_level"); ok {
			request.RiskLevel = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("config_rules"); ok {
			request.ConfigRules = buildCompliancePackRules(v.(*schema.Set).List())
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateCompliancePackWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config compliance pack failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		if err := updateCompliancePackStatus(ctx, meta, compliancePackId, status, logId); err != nil {
			return err
		}
	}

	return resourceTencentCloudConfigCompliancePackRead(d, meta)
}

func resourceTencentCloudConfigCompliancePackDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_compliance_pack.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return fmt.Errorf("tencentcloud config compliance pack not supported delete, please contact the work order for processing")

	// var (
	// 	logId            = tccommon.GetLogId(tccommon.ContextNil)
	// 	ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	// 	compliancePackId = d.Id()
	// )

	// // Must disable before delete
	// disableErr := updateCompliancePackStatus(ctx, meta, compliancePackId, "UN_ACTIVE", logId)
	// if disableErr != nil {
	// 	log.Printf("[WARN]%s disable config compliance pack before delete failed, reason:%+v", logId, disableErr)
	// }

	// request := configv20220802.NewDeleteCompliancePackRequest()
	// request.CompliancePackId = &compliancePackId

	// reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
	// 	result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteCompliancePackWithContext(ctx, request)
	// 	if e != nil {
	// 		return tccommon.RetryError(e)
	// 	}

	// 	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	// 	return nil
	// })

	// if reqErr != nil {
	// 	log.Printf("[CRITAL]%s delete config compliance pack failed, reason:%+v", logId, reqErr)
	// 	return reqErr
	// }

	// return nil
}

func updateCompliancePackStatus(ctx context.Context, meta interface{}, compliancePackId, status, logId string) error {
	request := configv20220802.NewUpdateCompliancePackStatusRequest()
	request.CompliancePackId = &compliancePackId
	request.Status = &status

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateCompliancePackStatusWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update config compliance pack status failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func buildCompliancePackRules(rawList []interface{}) []*configv20220802.CompliancePackRule {
	rules := make([]*configv20220802.CompliancePackRule, 0, len(rawList))
	for _, item := range rawList {
		ruleMap := item.(map[string]interface{})
		rule := &configv20220802.CompliancePackRule{}

		if v, ok := ruleMap["identifier"].(string); ok && v != "" {
			rule.Identifier = helper.String(v)
		}

		if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
			rule.RuleName = helper.String(v)
		}

		if v, ok := ruleMap["description"].(string); ok && v != "" {
			rule.Description = helper.String(v)
		}

		if v, ok := ruleMap["risk_level"].(int); ok && v != 0 {
			rule.RiskLevel = helper.IntUint64(v)
		}

		if v, ok := ruleMap["managed_rule_identifier"].(string); ok && v != "" {
			rule.ManagedRuleIdentifier = helper.String(v)
		}

		if v, ok := ruleMap["config_rule_id"].(string); ok && v != "" {
			rule.ConfigRuleId = helper.String(v)
		}

		if v, ok := ruleMap["compliance_pack_id"].(string); ok && v != "" {
			rule.CompliancePackId = helper.String(v)
		}

		if v, ok := ruleMap["input_parameter"]; ok {
			rule.InputParameter = buildInputParameters(v.([]interface{}))
		}

		rules = append(rules, rule)
	}

	return rules
}

func buildInputParameters(rawList []interface{}) []*configv20220802.InputParameter {
	params := make([]*configv20220802.InputParameter, 0, len(rawList))
	for _, item := range rawList {
		paramMap := item.(map[string]interface{})
		param := &configv20220802.InputParameter{}

		if v, ok := paramMap["parameter_key"].(string); ok && v != "" {
			param.ParameterKey = helper.String(v)
		}

		if v, ok := paramMap["type"].(string); ok && v != "" {
			param.Type = helper.String(v)
		}

		if v, ok := paramMap["value"].(string); ok && v != "" {
			param.Value = helper.String(v)
		}

		params = append(params, param)
	}

	return params
}

func flattenComplianceConfigRules(rules []*configv20220802.ComplianceConfigRule) []map[string]interface{} {
	rulesList := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleMap := map[string]interface{}{}

		if rule.Identifier != nil {
			ruleMap["identifier"] = rule.Identifier
		}

		if rule.RuleName != nil {
			ruleMap["rule_name"] = rule.RuleName
		}

		if rule.Description != nil {
			ruleMap["description"] = rule.Description
		}

		if rule.RiskLevel != nil {
			ruleMap["risk_level"] = int(*rule.RiskLevel)
		}

		if rule.ConfigRuleId != nil {
			ruleMap["config_rule_id"] = rule.ConfigRuleId
		}

		if rule.Status != nil {
			ruleMap["status"] = rule.Status
		}

		if rule.ComplianceResult != nil {
			ruleMap["compliance_result"] = rule.ComplianceResult
		}

		if rule.InputParameter != nil {
			inputParams := make([]map[string]interface{}, 0, len(rule.InputParameter))
			for _, param := range rule.InputParameter {
				paramMap := map[string]interface{}{}
				if param.ParameterKey != nil {
					paramMap["parameter_key"] = param.ParameterKey
				}

				if param.Type != nil {
					paramMap["type"] = param.Type
				}

				if param.Value != nil {
					paramMap["value"] = param.Value
				}

				inputParams = append(inputParams, paramMap)
			}

			ruleMap["input_parameter"] = inputParams
		}

		rulesList = append(rulesList, ruleMap)
	}

	return rulesList
}
