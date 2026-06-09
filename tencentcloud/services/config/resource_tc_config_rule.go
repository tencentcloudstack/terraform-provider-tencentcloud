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

func ResourceTencentCloudConfigRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigRuleCreate,
		Read:   resourceTencentCloudConfigRuleRead,
		Update: resourceTencentCloudConfigRuleUpdate,
		Delete: resourceTencentCloudConfigRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule template identifier. For system preset rules use the identifier name; for custom rules use the cloud function ARN (region:functionName).",
			},

			"identifier_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule template type. Valid values: SYSTEM (system preset), CUSTOMIZE (custom).",
			},

			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},

			"resource_type": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Supported resource type list (e.g. QCS::CAM::User).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"trigger_type": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Trigger type list, up to 2 entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger message type. Valid values: ScheduledNotification, ConfigurationItemChangeNotification.",
						},
						"maximum_execution_frequency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum execution frequency (only for ScheduledNotification). e.g. TwentyFour_Hours.",
						},
					},
				},
			},

			"risk_level": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},

			"input_parameter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rule input parameter list.",
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

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description (0~1024 characters).",
			},

			"regions_scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Region scope list; rule only applies to resources in the specified regions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags_scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag scope list; rule only applies to resources with the specified tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"exclude_resource_ids_scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Resource ID list excluded from rule evaluation.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule status. Valid values: ACTIVE (enabled), UN_ACTIVE (disabled).",
			},

			// Computed
			"config_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Config rule ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"compliance_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compliance result.",
			},

			"config_rule_invoked_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last rule evaluation time.",
			},
		},
	}
}

func resourceTencentCloudConfigRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = configv20220802.NewAddConfigRuleRequest()
		response = configv20220802.NewAddConfigRuleResponse()
	)

	if v, ok := d.GetOk("identifier"); ok {
		request.Identifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("identifier_type"); ok {
		request.IdentifierType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.ResourceType = append(request.ResourceType, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("trigger_type"); ok {
		request.TriggerType = buildConfigRuleTriggerTypes(v.([]interface{}))
	}

	if v, ok := d.GetOk("risk_level"); ok {
		request.RiskLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("input_parameter"); ok {
		request.InputParameter = buildInputParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("regions_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.RegionsScope = append(request.RegionsScope, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("tags_scope"); ok {
		request.TagsScope = buildConfigRuleTagsScope(v.([]interface{}))
	}

	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.ExcludeResourceIdsScope = append(request.ExcludeResourceIdsScope, helper.String(item.(string)))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().AddConfigRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create config rule failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create config rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.RuleId == nil {
		return fmt.Errorf("RuleId is nil")
	}

	ruleId := *response.Response.RuleId
	d.SetId(ruleId)

	// Handle initial status if provided
	if v, ok := d.GetOk("status"); ok {
		if err := updateConfigRuleStatus(ctx, meta, ruleId, v.(string), logId); err != nil {
			return err
		}
	}

	return resourceTencentCloudConfigRuleRead(d, meta)
}

func resourceTencentCloudConfigRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ruleId  = d.Id()
	)

	respData, err := service.DescribeConfigRuleById(ctx, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Identifier != nil {
		_ = d.Set("identifier", respData.Identifier)
	}

	if respData.IdentifierType != nil {
		_ = d.Set("identifier_type", respData.IdentifierType)
	}

	if respData.RuleName != nil {
		_ = d.Set("rule_name", respData.RuleName)
	}

	if respData.ResourceType != nil {
		resourceTypes := make([]string, 0, len(respData.ResourceType))
		for _, v := range respData.ResourceType {
			if v != nil {
				resourceTypes = append(resourceTypes, *v)
			}
		}

		_ = d.Set("resource_type", resourceTypes)
	}

	if respData.TriggerType != nil {
		triggerTypes := make([]map[string]interface{}, 0, len(respData.TriggerType))
		for _, tt := range respData.TriggerType {
			ttMap := map[string]interface{}{}
			if tt.MessageType != nil {
				ttMap["message_type"] = tt.MessageType
			}

			if tt.MaximumExecutionFrequency != nil {
				ttMap["maximum_execution_frequency"] = tt.MaximumExecutionFrequency
			}

			triggerTypes = append(triggerTypes, ttMap)
		}

		_ = d.Set("trigger_type", triggerTypes)
	}

	if respData.RiskLevel != nil {
		_ = d.Set("risk_level", int(*respData.RiskLevel))
	}

	if respData.InputParameter != nil {
		inputParams := make([]map[string]interface{}, 0, len(respData.InputParameter))
		for _, param := range respData.InputParameter {
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

		_ = d.Set("input_parameter", inputParams)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.RegionsScope != nil {
		regions := make([]string, 0, len(respData.RegionsScope))
		for _, v := range respData.RegionsScope {
			if v != nil {
				regions = append(regions, *v)
			}
		}

		_ = d.Set("regions_scope", regions)
	}

	if respData.TagsScope != nil {
		tagsScope := make([]map[string]interface{}, 0, len(respData.TagsScope))
		for _, tag := range respData.TagsScope {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}

			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}

			tagsScope = append(tagsScope, tagMap)
		}

		_ = d.Set("tags_scope", tagsScope)
	}

	if respData.ExcludeResourceIdsScope != nil {
		excludeIds := make([]string, 0, len(respData.ExcludeResourceIdsScope))
		for _, v := range respData.ExcludeResourceIdsScope {
			if v != nil {
				excludeIds = append(excludeIds, *v)
			}
		}

		_ = d.Set("exclude_resource_ids_scope", excludeIds)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ConfigRuleId != nil {
		_ = d.Set("config_rule_id", respData.ConfigRuleId)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ComplianceResult != nil {
		_ = d.Set("compliance_result", respData.ComplianceResult)
	}

	if respData.ConfigRuleInvokedTime != nil {
		_ = d.Set("config_rule_invoked_time", respData.ConfigRuleInvokedTime)
	}

	return nil
}

func resourceTencentCloudConfigRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		ruleId = d.Id()
	)

	contentArgs := []string{"rule_name", "trigger_type", "risk_level", "input_parameter", "description", "regions_scope", "tags_scope", "exclude_resource_ids_scope"}
	needContentUpdate := false
	for _, v := range contentArgs {
		if d.HasChange(v) {
			needContentUpdate = true
			break
		}
	}

	if needContentUpdate {
		request := configv20220802.NewUpdateConfigRuleRequest()
		request.RuleId = &ruleId

		if v, ok := d.GetOk("rule_name"); ok {
			request.RuleName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("trigger_type"); ok {
			request.TriggerType = buildConfigRuleTriggerTypes(v.([]interface{}))
		}

		if v, ok := d.GetOk("risk_level"); ok {
			request.RiskLevel = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("input_parameter"); ok {
			request.InputParameter = buildInputParameters(v.([]interface{}))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("regions_scope"); ok {
			rawList := v.([]interface{})
			for _, item := range rawList {
				request.RegionsScope = append(request.RegionsScope, helper.String(item.(string)))
			}
		}

		if v, ok := d.GetOk("tags_scope"); ok {
			request.TagsScope = buildConfigRuleTagsScope(v.([]interface{}))
		}

		if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
			rawList := v.([]interface{})
			for _, item := range rawList {
				request.ExcludeResourceIdsScope = append(request.ExcludeResourceIdsScope, helper.String(item.(string)))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateConfigRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update config rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		if err := updateConfigRuleStatus(ctx, meta, ruleId, status, logId); err != nil {
			return err
		}
	}

	return resourceTencentCloudConfigRuleRead(d, meta)
}

func resourceTencentCloudConfigRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		ruleId  = d.Id()
		request = configv20220802.NewDeleteConfigRuleRequest()
	)

	// Only disable if currently ACTIVE in state; UN_ACTIVE can be deleted directly
	if v, ok := d.GetOk("status"); ok && v.(string) == "ACTIVE" {
		if disableErr := updateConfigRuleStatus(ctx, meta, ruleId, "UN_ACTIVE", logId); disableErr != nil {
			log.Printf("[WARN]%s disable config rule before delete failed, reason:%+v", logId, disableErr)
		}
	}

	request.RuleId = &ruleId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteConfigRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete config rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func updateConfigRuleStatus(ctx context.Context, meta interface{}, ruleId, status, logId string) error {
	if status == "ACTIVE" {
		request := configv20220802.NewOpenConfigRuleRequest()
		request.RuleId = &ruleId

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().OpenConfigRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s open config rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	} else if status == "UN_ACTIVE" {
		request := configv20220802.NewCloseConfigRuleRequest()
		request.RuleId = &ruleId

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().CloseConfigRuleWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s close config rule failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return nil
}

func buildConfigRuleTriggerTypes(rawList []interface{}) []*configv20220802.TriggerType {
	triggerTypes := make([]*configv20220802.TriggerType, 0, len(rawList))
	for _, item := range rawList {
		ttMap := item.(map[string]interface{})
		tt := &configv20220802.TriggerType{}

		if v, ok := ttMap["message_type"].(string); ok && v != "" {
			tt.MessageType = helper.String(v)
		}

		if v, ok := ttMap["maximum_execution_frequency"].(string); ok && v != "" {
			tt.MaximumExecutionFrequency = helper.String(v)
		}

		triggerTypes = append(triggerTypes, tt)
	}

	return triggerTypes
}

func buildConfigRuleTagsScope(rawList []interface{}) []*configv20220802.Tag {
	tags := make([]*configv20220802.Tag, 0, len(rawList))
	for _, item := range rawList {
		tagMap := item.(map[string]interface{})
		tag := &configv20220802.Tag{}

		if v, ok := tagMap["tag_key"].(string); ok && v != "" {
			tag.TagKey = helper.String(v)
		}

		if v, ok := tagMap["tag_value"].(string); ok && v != "" {
			tag.TagValue = helper.String(v)
		}

		tags = append(tags, tag)
	}

	return tags
}
