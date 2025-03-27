package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoSecurityPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityPolicyConfigCreate,
		Read:   resourceTencentCloudTeoSecurityPolicyConfigRead,
		Update: resourceTencentCloudTeoSecurityPolicyConfigUpdate,
		Delete: resourceTencentCloudTeoSecurityPolicyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},

			"security_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Security policy configuration. it is recommended to use for custom policies and managed rule configurations of Web protection. it supports configuring security policies with expression grammar.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Custom rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of the custom rule.",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.",
												},
												"action": {
													Type:        schema.TypeList,
													Required:    true,
													MaxItems:    1,
													Description: "Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific actions for safe execution. valid values:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.",
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameter when Name is ReturnCustomPage.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response status code.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response custom page ID.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameter when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Redirect URL.",
																		},
																	},
																},
															},
														},
													},
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.",
												},
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The ID of a custom rule. <br> the rule ID supports different rule configuration operations: <br> - add a new rule: ID is empty or the ID parameter is not specified; <br> - modify an existing rule: specify the rule ID that needs to be updated/modified; <br> - delete an existing rule: existing Rules not included in the Rules list of the CustomRules parameter will be deleted.",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type of custom rule. valid values: <li>BasicAccessRule: basic access control;</li> <li>PreciseMatchRule: exact matching rule, default;</li> <li>ManagedAccessRule: expert customized rule, for output only.</li> the default value is PreciseMatchRule.",
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports `rule_type` is `PreciseMatchRule`.",
												},
											},
										},
									},
								},
							},
						},
						"managed_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Managed rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Indicates whether the managed rule is enabled. valid values: <li>on: enabled. all managed rules take effect as configured;</li> <li>off: disabled. all managed rules do not take effect.</li>.",
									},
									"detection_only": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Indicates whether the evaluation mode is Enabled. it is valid only when the Enabled parameter is set to on. valid values: <li>on: Enabled. all managed rules take effect in observation mode.</li> <li>off: disabled. all managed rules take effect according to the actual configuration.</li>.",
									},
									"semantic_analysis": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether the managed rule semantic analysis option is Enabled is valid only when the Enabled parameter is on. valid values: <li>on: enable. perform semantic analysis on requests before processing them;</li> <li>off: disable. process requests directly without semantic analysis.</li> <br/>default off.",
									},
									"auto_update": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Managed rule automatic update option.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_update_to_latest_version": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Indicates whether to enable automatic update to the latest version. valid values: <li>on: enabled</li> <li>off: disabled</li>.",
												},
												"ruleset_version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The currently used version, in the format compliant with ISO 8601 standard, such as 2023-12-21T12:00:32Z. it is empty by default and is only an output parameter.",
												},
											},
										},
									},
									"managed_rule_groups": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "Configuration of the managed rule group. if this structure is passed as an empty array or the GroupId is not included in the list, it will be processed based on the default method.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Group name of the managed rule. if the rule group for the configuration is not specified, it will be processed based on the default configuration. refer to product documentation for the specific value of GroupId.",
												},
												"sensitivity_level": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Protection level of the managed rule group. valid values: <li>loose: lenient, only contains ultra-high risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>normal: normal, contains ultra-high risk and high-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>strict: strict, contains ultra-high risk, high-risk and medium-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>extreme: super strict, contains ultra-high risk, high-risk, medium-risk and low-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>custom: custom, refined strategy. configure the disposal method for each individual rule. at this point, the Action field is invalid. use RuleActions to configure the refined strategy for each individual rule.</li>.",
												},
												"action": {
													Type:        schema.TypeList,
													Required:    true,
													MaxItems:    1,
													Description: "Handling actions for managed rule groups. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process requests and record security events in logs;</li> <li>Disabled: not enabled, do not scan requests and skip this rule.</li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific actions for safe execution. valid values:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.",
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameter when Name is ReturnCustomPage.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response status code.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response custom page ID.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameter when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Redirect URL.",
																		},
																	},
																},
															},
														},
													},
												},
												"rule_actions": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specific configuration of rule items under the managed rule group. the configuration is effective only when SensitivityLevel is custom.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"rule_id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific items under the managed rule group, which are used to rewrite the configuration content of this individual rule item. refer to product documentation for details.",
															},
															"action": {
																Type:        schema.TypeList,
																Required:    true,
																MaxItems:    1,
																Description: "Specify the handling action for the managed rule item in RuleId. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process the request and record the security event in logs;</li> <li>Disabled: Disabled, do not scan the request and skip this rule.</li>.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specific actions for safe execution. valid values:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.",
																		},
																		"block_ip_action_parameters": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Additional parameter when Name is BlockIP.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"duration": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.",
																					},
																				},
																			},
																		},
																		"return_custom_page_action_parameters": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Additional parameter when Name is ReturnCustomPage.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"response_code": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Response status code.",
																					},
																					"error_page_id": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Response custom page ID.",
																					},
																				},
																			},
																		},
																		"redirect_action_parameters": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Additional parameter when Name is Redirect.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Redirect URL.",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"meta_data": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Managed rule group information, for output only.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"group_detail": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Managed rule group description, for output only.",
															},
															"group_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Managed rule group name, for output only.",
															},
															"rule_details": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "All sub-rule information under the current managed rule group, for output only.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"rule_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Managed rule Id.",
																		},
																		"risk_level": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Protection level of managed rules. valid values: <li>low: low risk. this rule has a relatively low risk and is applicable to access scenarios in a very strict control environment. this level of rule may generate considerable false alarms.</li> <li>medium: medium risk. this means the risk of this rule is normal and it is suitable for protection scenarios with stricter requirements.</li> <li>high: high risk. this indicates that the risk of this rule is relatively high and it will not generate false alarms in most scenarios.</li> <li>extreme: ultra-high risk. this represents that the risk of this rule is extremely high and it will not generate false alarms basically.</li>.",
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Rule description.",
																		},
																		"tags": {
																			Type:        schema.TypeSet,
																			Computed:    true,
																			Description: "Rule tag. some types of rules do not have tags.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"rule_version": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Rule ownership version.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"entity": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ZoneDefaultPolicy", "Template", "Host"}),
				Description:  "Security policy type. the following parameter values can be used: <li>ZoneDefaultPolicy: used to specify a site-level policy;</li> <li>Template: used to specify a policy Template. you need to simultaneously specify the TemplateId parameter;</li> <li>Host: used to specify a domain-level policy (note: when using a domain name to specify a dns service policy, only dns services or policy templates that have applied a domain-level policy are supported).</li>.",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the specified domain. when the Entity parameter value is Host, use the domain-level policy specified by this parameter. for example: use www.example.com to configure the domain-level policy of the domain.",
			},

			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specify the policy Template ID. use this parameter to specify the ID of the policy Template when the Entity parameter value is Template.",
			},
		},
	}
}

func resourceTencentCloudTeoSecurityPolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId     string
		entity     string
		host       string
		templateId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
	}

	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}

	if v, ok := d.GetOk("template_id"); ok {
		templateId = v.(string)
	}

	if entity == "ZoneDefaultPolicy" && host == "" && templateId == "" {
		d.SetId(strings.Join([]string{zoneId, entity}, tccommon.FILED_SP))
	} else if entity == "Host" && host != "" && templateId == "" {
		d.SetId(strings.Join([]string{zoneId, entity, host}, tccommon.FILED_SP))
	} else if entity == "Template" && host == "" && templateId != "" {
		d.SetId(strings.Join([]string{zoneId, entity, templateId}, tccommon.FILED_SP))
	} else {
		return fmt.Errorf("If `entity` is `ZoneDefaultPolicy`, Please do not set `host` and `template_id`; If `entity` is `Host`, Only support set `host`; If `entity` is `Template`, Only support set `template_id`.")
	}

	return resourceTencentCloudTeoSecurityPolicyConfigUpdate(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId     string
		entity     string
		host       string
		templateId string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 2 || len(idSplit) == 3) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId = idSplit[0]
	entity = idSplit[1]
	if entity == "ZoneDefaultPolicy" && len(idSplit) == 2 {

	} else if entity == "Host" && len(idSplit) == 3 {
		host = idSplit[2]
	} else if entity == "Template" && len(idSplit) == 3 {
		templateId = idSplit[2]
	} else {
		return fmt.Errorf("`entity` is illegal, %s.", entity)
	}

	respData, err := service.DescribeTeoSecurityPolicyConfigById(ctx, zoneId, entity, host, templateId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_security_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("entity", entity)
	_ = d.Set("host", host)
	_ = d.Set("template_id", templateId)

	securityPolicyList := make([]map[string]interface{}, 0, 1)
	securityPolicyMap := map[string]interface{}{}
	if respData.CustomRules != nil {
		customRulesMap := map[string]interface{}{}
		rulesList := make([]map[string]interface{}, 0, len(respData.CustomRules.Rules))
		if respData.CustomRules.Rules != nil {
			for _, rules := range respData.CustomRules.Rules {
				rulesMap := map[string]interface{}{}
				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Condition != nil {
					rulesMap["condition"] = rules.Condition
				}

				actionMap := map[string]interface{}{}
				if rules.Action != nil {
					if rules.Action.Name != nil {
						actionMap["name"] = rules.Action.Name
					}

					blockIPActionParametersMap := map[string]interface{}{}
					if rules.Action.BlockIPActionParameters != nil {
						if rules.Action.BlockIPActionParameters.Duration != nil {
							blockIPActionParametersMap["duration"] = rules.Action.BlockIPActionParameters.Duration
						}

						actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
					}

					returnCustomPageActionParametersMap := map[string]interface{}{}
					if rules.Action.ReturnCustomPageActionParameters != nil {
						if rules.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = rules.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if rules.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = rules.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					redirectActionParametersMap := map[string]interface{}{}
					if rules.Action.RedirectActionParameters != nil {
						if rules.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = rules.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					rulesMap["action"] = []interface{}{actionMap}
				}

				if rules.Enabled != nil {
					rulesMap["enabled"] = rules.Enabled
				}

				if rules.Id != nil {
					rulesMap["id"] = rules.Id
				}

				if rules.RuleType != nil {
					rulesMap["rule_type"] = rules.RuleType
				}

				if rules.Priority != nil {
					rulesMap["priority"] = rules.Priority
				}

				rulesList = append(rulesList, rulesMap)
			}

			customRulesMap["rules"] = rulesList
		} else {
			customRulesMap["rules"] = rulesList
		}

		securityPolicyMap["custom_rules"] = []interface{}{customRulesMap}
	}

	if respData.ManagedRules != nil {
		managedRulesMap := map[string]interface{}{}
		if respData.ManagedRules.Enabled != nil {
			managedRulesMap["enabled"] = respData.ManagedRules.Enabled
		}

		if respData.ManagedRules.DetectionOnly != nil {
			managedRulesMap["detection_only"] = respData.ManagedRules.DetectionOnly
		}

		if respData.ManagedRules.SemanticAnalysis != nil {
			managedRulesMap["semantic_analysis"] = respData.ManagedRules.SemanticAnalysis
		}

		autoUpdateMap := map[string]interface{}{}
		if respData.ManagedRules.AutoUpdate != nil {
			if respData.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion != nil {
				autoUpdateMap["auto_update_to_latest_version"] = respData.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion
			}

			if respData.ManagedRules.AutoUpdate.RulesetVersion != nil {
				autoUpdateMap["ruleset_version"] = respData.ManagedRules.AutoUpdate.RulesetVersion
			}

			managedRulesMap["auto_update"] = []interface{}{autoUpdateMap}
		}

		managedRuleGroupsList := make([]map[string]interface{}, 0, len(respData.ManagedRules.ManagedRuleGroups))
		if respData.ManagedRules.ManagedRuleGroups != nil {
			for _, managedRuleGroups := range respData.ManagedRules.ManagedRuleGroups {
				managedRuleGroupsMap := map[string]interface{}{}

				if managedRuleGroups.GroupId != nil {
					managedRuleGroupsMap["group_id"] = managedRuleGroups.GroupId
				}

				if managedRuleGroups.SensitivityLevel != nil {
					managedRuleGroupsMap["sensitivity_level"] = managedRuleGroups.SensitivityLevel
				}

				actionMap := map[string]interface{}{}
				if managedRuleGroups.Action != nil {
					if managedRuleGroups.Action.Name != nil {
						actionMap["name"] = managedRuleGroups.Action.Name
					}

					blockIPActionParametersMap := map[string]interface{}{}
					if managedRuleGroups.Action.BlockIPActionParameters != nil {
						if managedRuleGroups.Action.BlockIPActionParameters.Duration != nil {
							blockIPActionParametersMap["duration"] = managedRuleGroups.Action.BlockIPActionParameters.Duration
						}

						actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
					}

					returnCustomPageActionParametersMap := map[string]interface{}{}
					if managedRuleGroups.Action.ReturnCustomPageActionParameters != nil {
						if managedRuleGroups.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = managedRuleGroups.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if managedRuleGroups.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = managedRuleGroups.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					redirectActionParametersMap := map[string]interface{}{}
					if managedRuleGroups.Action.RedirectActionParameters != nil {
						if managedRuleGroups.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = managedRuleGroups.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					managedRuleGroupsMap["action"] = []interface{}{actionMap}
				}

				ruleActionsList := make([]map[string]interface{}, 0, len(managedRuleGroups.RuleActions))
				if managedRuleGroups.RuleActions != nil {
					for _, ruleActions := range managedRuleGroups.RuleActions {
						ruleActionsMap := map[string]interface{}{}
						if ruleActions.RuleId != nil {
							ruleActionsMap["rule_id"] = ruleActions.RuleId
						}

						actionMap := map[string]interface{}{}
						if ruleActions.Action != nil {
							if ruleActions.Action.Name != nil {
								actionMap["name"] = ruleActions.Action.Name
							}

							blockIPActionParametersMap := map[string]interface{}{}
							if ruleActions.Action.BlockIPActionParameters != nil {
								if ruleActions.Action.BlockIPActionParameters.Duration != nil {
									blockIPActionParametersMap["duration"] = ruleActions.Action.BlockIPActionParameters.Duration
								}

								actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
							}

							returnCustomPageActionParametersMap := map[string]interface{}{}
							if ruleActions.Action.ReturnCustomPageActionParameters != nil {
								if ruleActions.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
									returnCustomPageActionParametersMap["response_code"] = ruleActions.Action.ReturnCustomPageActionParameters.ResponseCode
								}

								if ruleActions.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
									returnCustomPageActionParametersMap["error_page_id"] = ruleActions.Action.ReturnCustomPageActionParameters.ErrorPageId
								}

								actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
							}

							redirectActionParametersMap := map[string]interface{}{}
							if ruleActions.Action.RedirectActionParameters != nil {
								if ruleActions.Action.RedirectActionParameters.URL != nil {
									redirectActionParametersMap["url"] = ruleActions.Action.RedirectActionParameters.URL
								}

								actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
							}

							ruleActionsMap["action"] = []interface{}{actionMap}
						}

						ruleActionsList = append(ruleActionsList, ruleActionsMap)
					}

					managedRuleGroupsMap["rule_actions"] = ruleActionsList
				}

				metaDataMap := map[string]interface{}{}
				if managedRuleGroups.MetaData != nil {
					if managedRuleGroups.MetaData.GroupDetail != nil {
						metaDataMap["group_detail"] = managedRuleGroups.MetaData.GroupDetail
					}

					if managedRuleGroups.MetaData.GroupName != nil {
						metaDataMap["group_name"] = managedRuleGroups.MetaData.GroupName
					}

					ruleDetailsList := make([]map[string]interface{}, 0, len(managedRuleGroups.MetaData.RuleDetails))
					if managedRuleGroups.MetaData.RuleDetails != nil {
						for _, ruleDetails := range managedRuleGroups.MetaData.RuleDetails {
							ruleDetailsMap := map[string]interface{}{}
							if ruleDetails.RuleId != nil {
								ruleDetailsMap["rule_id"] = ruleDetails.RuleId
							}

							if ruleDetails.RiskLevel != nil {
								ruleDetailsMap["risk_level"] = ruleDetails.RiskLevel
							}

							if ruleDetails.Description != nil {
								ruleDetailsMap["description"] = ruleDetails.Description
							}

							if ruleDetails.Tags != nil {
								ruleDetailsMap["tags"] = ruleDetails.Tags
							}

							if ruleDetails.RuleVersion != nil {
								ruleDetailsMap["rule_version"] = ruleDetails.RuleVersion
							}

							ruleDetailsList = append(ruleDetailsList, ruleDetailsMap)
						}

						metaDataMap["rule_details"] = ruleDetailsList
					}
					managedRuleGroupsMap["meta_data"] = []interface{}{metaDataMap}
				}

				managedRuleGroupsList = append(managedRuleGroupsList, managedRuleGroupsMap)
			}

			managedRulesMap["managed_rule_groups"] = managedRuleGroupsList
		}

		securityPolicyMap["managed_rules"] = []interface{}{managedRulesMap}
	}

	securityPolicyList = append(securityPolicyList, securityPolicyMap)
	_ = d.Set("security_policy", securityPolicyList)
	return nil
}

func resourceTencentCloudTeoSecurityPolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = teov20220901.NewModifySecurityPolicyRequest()
		zoneId     string
		entity     string
		host       string
		templateId string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 2 || len(idSplit) == 3) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId = idSplit[0]
	entity = idSplit[1]
	if entity == "ZoneDefaultPolicy" && len(idSplit) == 2 {

	} else if entity == "Host" && len(idSplit) == 3 {
		host = idSplit[2]
	} else if entity == "Template" && len(idSplit) == 3 {
		templateId = idSplit[2]
	} else {
		return fmt.Errorf("`entity` is illegal, %s.", entity)
	}

	request.ZoneId = &zoneId
	request.Entity = &entity
	request.TemplateId = &templateId
	request.Host = &host
	request.SecurityConfig = &teov20220901.SecurityConfig{}
	if securityPolicyMap, ok := helper.InterfacesHeadMap(d, "security_policy"); ok {
		securityPolicy := teov20220901.SecurityPolicy{}
		if customRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["custom_rules"]); ok {
			customRules := teov20220901.CustomRules{}
			if v, ok := customRulesMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					customRule := teov20220901.CustomRule{}
					if v, ok := rulesMap["name"].(string); ok && v != "" {
						customRule.Name = helper.String(v)
					}

					if v, ok := rulesMap["condition"].(string); ok && v != "" {
						customRule.Condition = helper.String(v)
					}

					if actionMap, ok := helper.ConvertInterfacesHeadToMap(rulesMap["action"]); ok {
						securityAction := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"].(string); ok && v != "" {
							securityAction.Name = helper.String(v)
						}

						if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
							blockIPActionParameters := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
								blockIPActionParameters.Duration = helper.String(v)
							}

							securityAction.BlockIPActionParameters = &blockIPActionParameters
						}

						if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
							returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
								returnCustomPageActionParameters.ResponseCode = helper.String(v)
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
								returnCustomPageActionParameters.ErrorPageId = helper.String(v)
							}

							securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
						}

						if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
							redirectActionParameters := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
								redirectActionParameters.URL = helper.String(v)
							}

							securityAction.RedirectActionParameters = &redirectActionParameters
						}

						customRule.Action = &securityAction
					}

					if v, ok := rulesMap["enabled"].(string); ok && v != "" {
						customRule.Enabled = helper.String(v)
					}

					if v, ok := rulesMap["id"].(string); ok && v != "" {
						customRule.Id = helper.String(v)
					}

					if v, ok := rulesMap["rule_type"].(string); ok && v != "" {
						customRule.RuleType = helper.String(v)
					}

					if v, ok := rulesMap["priority"].(int); ok {
						customRule.Priority = helper.IntInt64(v)
					}

					customRules.Rules = append(customRules.Rules, &customRule)
				}
			}

			securityPolicy.CustomRules = &customRules
		}

		if managedRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["managed_rules"]); ok {
			managedRules := teov20220901.ManagedRules{}
			if v, ok := managedRulesMap["enabled"].(string); ok && v != "" {
				managedRules.Enabled = helper.String(v)
			}

			if v, ok := managedRulesMap["detection_only"].(string); ok && v != "" {
				managedRules.DetectionOnly = helper.String(v)
			}

			if v, ok := managedRulesMap["semantic_analysis"].(string); ok && v != "" {
				managedRules.SemanticAnalysis = helper.String(v)
			}

			if autoUpdateMap, ok := helper.ConvertInterfacesHeadToMap(managedRulesMap["auto_update"]); ok {
				managedRuleAutoUpdate := teov20220901.ManagedRuleAutoUpdate{}
				if v, ok := autoUpdateMap["auto_update_to_latest_version"].(string); ok && v != "" {
					managedRuleAutoUpdate.AutoUpdateToLatestVersion = helper.String(v)
				}

				managedRules.AutoUpdate = &managedRuleAutoUpdate
			}

			if v, ok := managedRulesMap["managed_rule_groups"]; ok {
				for _, item := range v.(*schema.Set).List() {
					managedRuleGroupsMap := item.(map[string]interface{})
					managedRuleGroup := teov20220901.ManagedRuleGroup{}
					if v, ok := managedRuleGroupsMap["group_id"].(string); ok && v != "" {
						managedRuleGroup.GroupId = helper.String(v)
					}

					if v, ok := managedRuleGroupsMap["sensitivity_level"].(string); ok && v != "" {
						managedRuleGroup.SensitivityLevel = helper.String(v)
					}

					if actionMap, ok := helper.ConvertInterfacesHeadToMap(managedRuleGroupsMap["action"]); ok {
						securityAction2 := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"].(string); ok && v != "" {
							securityAction2.Name = helper.String(v)
						}

						if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
							blockIPActionParameters2 := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
								blockIPActionParameters2.Duration = helper.String(v)
							}

							securityAction2.BlockIPActionParameters = &blockIPActionParameters2
						}

						if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
							returnCustomPageActionParameters2 := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
								returnCustomPageActionParameters2.ResponseCode = helper.String(v)
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
								returnCustomPageActionParameters2.ErrorPageId = helper.String(v)
							}

							securityAction2.ReturnCustomPageActionParameters = &returnCustomPageActionParameters2
						}

						if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
							redirectActionParameters2 := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
								redirectActionParameters2.URL = helper.String(v)
							}

							securityAction2.RedirectActionParameters = &redirectActionParameters2
						}

						managedRuleGroup.Action = &securityAction2
					}

					if v, ok := managedRuleGroupsMap["rule_actions"]; ok {
						for _, item := range v.([]interface{}) {
							ruleActionsMap := item.(map[string]interface{})
							managedRuleAction := teov20220901.ManagedRuleAction{}
							if v, ok := ruleActionsMap["rule_id"].(string); ok && v != "" {
								managedRuleAction.RuleId = helper.String(v)
							}

							if actionMap, ok := helper.ConvertInterfacesHeadToMap(ruleActionsMap["action"]); ok {
								securityAction3 := teov20220901.SecurityAction{}
								if v, ok := actionMap["name"].(string); ok && v != "" {
									securityAction3.Name = helper.String(v)
								}

								if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
									blockIPActionParameters3 := teov20220901.BlockIPActionParameters{}
									if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
										blockIPActionParameters3.Duration = helper.String(v)
									}

									securityAction3.BlockIPActionParameters = &blockIPActionParameters3
								}

								if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
									returnCustomPageActionParameters3 := teov20220901.ReturnCustomPageActionParameters{}
									if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
										returnCustomPageActionParameters3.ResponseCode = helper.String(v)
									}

									if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
										returnCustomPageActionParameters3.ErrorPageId = helper.String(v)
									}

									securityAction3.ReturnCustomPageActionParameters = &returnCustomPageActionParameters3
								}

								if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
									redirectActionParameters3 := teov20220901.RedirectActionParameters{}
									if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
										redirectActionParameters3.URL = helper.String(v)
									}

									securityAction3.RedirectActionParameters = &redirectActionParameters3
								}

								managedRuleAction.Action = &securityAction3
							}

							managedRuleGroup.RuleActions = append(managedRuleGroup.RuleActions, &managedRuleAction)
						}
					}

					managedRules.ManagedRuleGroups = append(managedRules.ManagedRuleGroups, &managedRuleGroup)
				}
			}

			securityPolicy.ManagedRules = &managedRules
		}

		request.SecurityPolicy = &securityPolicy
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo security policy failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify teo security policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoSecurityPolicyConfigRead(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
