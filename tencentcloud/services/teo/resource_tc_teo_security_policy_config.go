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
										Deprecated:  "It has been deprecated from version 1.81.184. Please use `precise_match_rules` or `basic_access_rules` instead.",
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
									"precise_match_rules": {
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
													Computed:    true,
													Description: "The ID of a custom rule. <br> the rule ID supports different rule configuration operations: <br> - add a new rule: ID is empty or the ID parameter is not specified; <br> - modify an existing rule: specify the rule ID that needs to be updated/modified; <br> - delete an existing rule: existing Rules not included in the Rules list of the CustomRules parameter will be deleted.",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
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
									"basic_access_rules": {
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
													Computed:    true,
													Description: "The ID of a custom rule. <br> the rule ID supports different rule configuration operations: <br> - add a new rule: ID is empty or the ID parameter is not specified; <br> - modify an existing rule: specify the rule ID that needs to be updated/modified; <br> - delete an existing rule: existing Rules not included in the Rules list of the CustomRules parameter will be deleted.",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
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
						"http_ddos_protection": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "HTTP DDOS protection configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"adaptive_frequency_control": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific configuration of adaptive frequency control.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether adaptive frequency control is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
												},
												"sensitivity": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The restriction level of adaptive frequency control. When Enabled is on, this field is required. The values are: <li>Loose: loose; </li><li>Moderate: moderate; </li><li>Strict: strict. </li>.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The handling method of adaptive frequency control. When Enabled is on, this field is required. SecurityAction's Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The specific action of security execution. The values are:\n<li>Deny: intercept, block the request to access site resources;</li>\n<li>Monitor: observe, only record logs;</li>\n<li>Redirect: redirect to URL;</li>\n<li>Disabled: disabled, do not enable the specified rule;</li>\n<li>Allow: allow access, but delay processing requests;</li>\n<li>Challenge: challenge, respond to challenge content;</li>\n<li>BlockIP: to be abandoned, IP ban;</li>\n<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>\n<li>JSChallenge: to be abandoned, JavaScript challenge;</li>\n<li>ManagedChallenge: to be abandoned, managed challenge.</li>.",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Deny.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to extend the blocking of source IP. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nWhen enabled, the client IP that triggers the rule will be blocked continuously. When this option is enabled, the BlockIpDuration parameter must be specified at the same time.\nNote: This option cannot be enabled at the same time as the ReturnCustomPage or Stall options.",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "When BlockIP is on, the IP blocking duration.",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to use custom pages. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nAfter enabling, use custom page content to intercept (respond to) requests. When enabling this option, you must specify the ResponseCode and ErrorPageId parameters at the same time.\nNote: This option cannot be enabled at the same time as the BlockIp or Stall options.",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize the status code of the page.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The PageId of the custom page.",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to ignore the request source suspension. The value is:\n<li>on: Enable;</li>\n<li>off: Disable.</li>\nAfter enabling, it will no longer respond to requests in the current connection session and will not actively disconnect. It is used to fight against crawlers and consume client connection resources.\nNote: This option cannot be enabled at the same time as the BlockIp or ReturnCustomPage options.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The URL to redirect.",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Challenge.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is ReturnCustomPage.",
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
																			Description: "The custom page ID of the response.",
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
									"client_filtering": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific configuration of intelligent client filtering.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether smart client filtering is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The method of intelligent client filtering. When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The specific action of security execution. The values are:\n<li>Deny: intercept, block the request to access site resources;</li>\n<li>Monitor: observe, only record logs;</li>\n<li>Redirect: redirect to URL;</li>\n<li>Disabled: disabled, do not enable the specified rule;</li>\n<li>Allow: allow access, but delay processing requests;</li>\n<li>Challenge: challenge, respond to challenge content;</li>\n<li>BlockIP: to be abandoned, IP ban;</li>\n<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>\n<li>JSChallenge: to be abandoned, JavaScript challenge;</li>\n<li>ManagedChallenge: to be abandoned, managed challenge.</li>.",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Deny.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to extend the blocking of source IP. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nWhen enabled, the client IP that triggers the rule will be blocked continuously. When this option is enabled, the BlockIpDuration parameter must be specified at the same time.\nNote: This option cannot be enabled at the same time as the ReturnCustomPage or Stall options.",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "When BlockIP is on, the IP blocking duration.",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to use custom pages. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nAfter enabling, use custom page content to intercept (respond to) requests. When enabling this option, you must specify the ResponseCode and ErrorPageId parameters at the same time.\nNote: This option cannot be enabled at the same time as the BlockIp or Stall options.",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize the status code of the page.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The PageId of the custom page.",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to ignore the request source suspension. The value is:\n<li>on: Enable;</li>\n<li>off: Disable.</li>\nAfter enabling, it will no longer respond to requests in the current connection session and will not actively disconnect. It is used to fight against crawlers and consume client connection resources.\nNote: This option cannot be enabled at the same time as the BlockIp or ReturnCustomPage options.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The URL to redirect.",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Challenge.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is ReturnCustomPage.",
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
																			Description: "The custom page ID of the response.",
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
									"bandwidth_abuse_defense": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific configuration of traffic fraud prevention.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether the anti-theft feature (only applicable to mainland China) is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The method for preventing traffic fraud (only applicable to mainland China). When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The specific action of security execution. The values are:\n<li>Deny: intercept, block the request to access site resources;</li>\n<li>Monitor: observe, only record logs;</li>\n<li>Redirect: redirect to URL;</li>\n<li>Disabled: disabled, do not enable the specified rule;</li>\n<li>Allow: allow access, but delay processing requests;</li>\n<li>Challenge: challenge, respond to challenge content;</li>\n<li>BlockIP: to be abandoned, IP ban;</li>\n<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>\n<li>JSChallenge: to be abandoned, JavaScript challenge;</li>\n<li>ManagedChallenge: to be abandoned, managed challenge.</li>.",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Deny.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to extend the blocking of source IP. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nWhen enabled, the client IP that triggers the rule will be blocked continuously. When this option is enabled, the BlockIpDuration parameter must be specified at the same time.\nNote: This option cannot be enabled at the same time as the ReturnCustomPage or Stall options.",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "When BlockIP is on, the IP blocking duration.",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to use custom pages. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nAfter enabling, use custom page content to intercept (respond to) requests. When enabling this option, you must specify the ResponseCode and ErrorPageId parameters at the same time.\nNote: This option cannot be enabled at the same time as the BlockIp or Stall options.",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize the status code of the page.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The PageId of the custom page.",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to ignore the request source suspension. The value is:\n<li>on: Enable;</li>\n<li>off: Disable.</li>\nAfter enabling, it will no longer respond to requests in the current connection session and will not actively disconnect. It is used to fight against crawlers and consume client connection resources.\nNote: This option cannot be enabled at the same time as the BlockIp or ReturnCustomPage options.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The URL to redirect.",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Challenge.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is ReturnCustomPage.",
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
																			Description: "The custom page ID of the response.",
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
									"slow_attack_defense": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific configuration of slow attack protection.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether slow attack protection is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The handling method of slow attack protection. When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The specific action of security execution. The values are:\n<li>Deny: intercept, block the request to access site resources;</li>\n<li>Monitor: observe, only record logs;</li>\n<li>Redirect: redirect to URL;</li>\n<li>Disabled: disabled, do not enable the specified rule;</li>\n<li>Allow: allow access, but delay processing requests;</li>\n<li>Challenge: challenge, respond to challenge content;</li>\n<li>BlockIP: to be abandoned, IP ban;</li>\n<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>\n<li>JSChallenge: to be abandoned, JavaScript challenge;</li>\n<li>ManagedChallenge: to be abandoned, managed challenge.</li>.",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Deny.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to extend the blocking of source IP. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nWhen enabled, the client IP that triggers the rule will be blocked continuously. When this option is enabled, the BlockIpDuration parameter must be specified at the same time.\nNote: This option cannot be enabled at the same time as the ReturnCustomPage or Stall options.",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "When BlockIP is on, the IP blocking duration.",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to use custom pages. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nAfter enabling, use custom page content to intercept (respond to) requests. When enabling this option, you must specify the ResponseCode and ErrorPageId parameters at the same time.\nNote: This option cannot be enabled at the same time as the BlockIp or Stall options.",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize the status code of the page.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The PageId of the custom page.",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to ignore the request source suspension. The value is:\n<li>on: Enable;</li>\n<li>off: Disable.</li>\nAfter enabling, it will no longer respond to requests in the current connection session and will not actively disconnect. It is used to fight against crawlers and consume client connection resources.\nNote: This option cannot be enabled at the same time as the BlockIp or ReturnCustomPage options.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The URL to redirect.",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Challenge.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is ReturnCustomPage.",
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
																			Description: "The custom page ID of the response.",
																		},
																	},
																},
															},
														},
													},
												},
												"minimal_request_body_transfer_rate": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Specific configuration of the minimum rate threshold for text transmission. This field is required when Enabled is on.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"minimal_avg_transfer_rate_threshold": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Minimum text transmission rate threshold. The unit only supports bps.",
															},
															"counting_period": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The minimum text transmission rate statistics time range, the possible values are: <li>10s: 10 seconds; </li><li>30s: 30 seconds; </li><li>60s: 60 seconds; </li><li>120s: 120 seconds. </li>.",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Whether the text transmission minimum rate threshold is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
															},
														},
													},
												},
												"request_body_transfer_timeout": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Specific configuration of the text transmission timeout. When Enabled is on, this field is required.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"idle_timeout": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The text transmission timeout period is between 5 and 120, and the unit only supports seconds (s).",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Whether the text transmission timeout is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
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
						"rate_limiting_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Rate limiting rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A list of precise rate limiting definitions. When using ModifySecurityPolicy to modify the Web protection configuration: <br> <li> If the Rules parameter is not specified, or the Rules parameter length is zero: clear all precise rate limiting configurations. </li>. <li> If the RateLimitingRules parameter value is not specified in the SecurityPolicy parameter: keep the existing custom rule configuration and do not modify it. </li>.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The ID of the precise rate limit. <br>The rule ID can support different rule configuration operations: <br> <li> <b>Add</b> a new rule: the ID is empty or the ID parameter is not specified; </li><li> <b>Modify</b> an existing rule: specify the rule ID to be updated/modified; </li><li> <b>Delete</b> an existing rule: in the RateLimitingRules parameter, the existing rules not included in the Rules list will be deleted. </li>.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the precise rate limit.",
												},
												"condition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The specific content of the precise rate limit must conform to the expression syntax. For detailed specifications, see the product documentation.",
												},
												"count_by": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "The matching method of the rate threshold request feature. When Enabled is on, this field is required. <br /><br />When there are multiple conditions, multiple conditions will be combined for statistical calculation. The number of conditions cannot exceed 5. The possible values are: <br/><li><b>http.request.ip</b>: client IP; </li><li><b>http.request.xff_header_ip</b>: client IP (matching XFF header first); </li><li><b>http.request.uri.path</b>: requested access path; </li><li><b>http.request.cookies['session']</b>: cookie named session, where session can be replaced by the parameter you specify; </li><li><b>http.request.headers['user-agent']</b>: HTTP header named user-agent, where user-agent can be replaced by the parameter you specify; </li><li><b>http.request.ja3</b>: requested JA3 fingerprint; </li><li><b>http.request.uri.query['test']</b>: URL query parameter named test, where test can be replaced by the parameter you specify. </li>.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"max_request_threshold": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The cumulative number of interceptions within the time range of the precise rate limit, ranging from 1 to 100000.",
												},
												"counting_period": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The statistical time window, the possible values are: <li>1s: 1 second; </li><li>5s: 5 seconds; </li><li>10s: 10 seconds; </li><li>20s: 20 seconds; </li><li>30s: 30 seconds; </li><li>40s: 40 seconds; </li><li>50s: 50 seconds; </li><li>1m: 1 minute; </li><li>2m: 2 minutes; </li><li>5m: 5 minutes; </li><li>10m: 10 minutes; </li><li>1h: 1 hour. </li>.",
												},
												"action_duration": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Action The duration of the action. The supported units are: <li>s: seconds, with a value of 1 to 120; </li><li>m: minutes, with a value of 1 to 120; </li><li>h: hours, with a value of 1 to 48; </li><li>d: days, with a value of 1 to 30. </li>.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "The precise rate limit handling method. The values are: <li>Monitor: Observe; </li><li>Deny: Intercept, where DenyActionParameters.Name supports Deny and ReturnCustomPage; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name supports JSChallenge and ManagedChallenge; </li><li>Redirect: Redirect to URL; </li>.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The specific action of security execution. The values are:\n<li>Deny: intercept, block the request to access site resources;</li>\n<li>Monitor: observe, only record logs;</li>\n<li>Redirect: redirect to URL;</li>\n<li>Disabled: disabled, do not enable the specified rule;</li>\n<li>Allow: allow access, but delay processing requests;</li>\n<li>Challenge: challenge, respond to challenge content;</li>\n<li>BlockIP: to be abandoned, IP ban;</li>\n<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>\n<li>JSChallenge: to be abandoned, JavaScript challenge;</li>\n<li>ManagedChallenge: to be abandoned, managed challenge.</li>.",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Deny.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to extend the blocking of source IP. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nWhen enabled, the client IP that triggers the rule will be blocked continuously. When this option is enabled, the BlockIpDuration parameter must be specified at the same time.\nNote: This option cannot be enabled at the same time as the ReturnCustomPage or Stall options.",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "When BlockIP is on, the IP blocking duration.",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to use custom pages. The possible values are:\n<li>on: on;</li>\n<li>off: off.</li>\nAfter enabling, use custom page content to intercept (respond to) requests. When enabling this option, you must specify the ResponseCode and ErrorPageId parameters at the same time.\nNote: This option cannot be enabled at the same time as the BlockIp or Stall options.",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize the status code of the page.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The PageId of the custom page.",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to ignore the request source suspension. The value is:\n<li>on: Enable;</li>\n<li>off: Disable.</li>\nAfter enabling, it will no longer respond to requests in the current connection session and will not actively disconnect. It is used to fight against crawlers and consume client connection resources.\nNote: This option cannot be enabled at the same time as the BlockIp or ReturnCustomPage options.",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Redirect.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The URL to redirect.",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Challenge.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is BlockIP.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "To be deprecated, additional parameter when Name is ReturnCustomPage.",
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
																			Description: "The custom page ID of the response.",
																		},
																	},
																},
															},
														},
													},
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The priority of precise rate limiting ranges from 0 to 100, and the default is 0.",
												},
												"enabled": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether the precise rate limit rule is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.",
												},
											},
										},
									},
								},
							},
						},
						"exception_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Exception rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Definition list of exception rules. When using ModifySecurityPolicy to modify the Web protection configuration: <li>If the Rules parameter is not specified, or the length of the Rules parameter is zero: clear all exception rule configurations. </li>.<li>If the ExceptionRules parameter value is not specified in the SecurityPolicy parameter: keep the existing exception rule configurations and do not modify them. </li>.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "The ID of the exception rule. <br>The rule ID can support different rule configuration operations: <br> <li> <b>Add</b> a new rule: the ID is empty or the ID parameter is not specified; </li><li> <b>Modify</b> an existing rule: specify the rule ID to be updated/modified; </li><li> <b>Delete</b> an existing rule: in the ExceptionRules parameter, the existing rules not included in the Rules list will be deleted. </li>.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the exception rule.",
												},
												"condition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The specific content of the exception rule must comply with the expression syntax. For detailed specifications, see the product documentation.",
												},
												"skip_scope": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Exception rule execution options, the values are: <li>WebSecurityModules: Specifies the security protection module for the exception rule. </li>.<li>ManagedRules: Specifies the managed rules. </li>.",
												},
												"skip_option": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The specific type of the skipped request. The possible values are: <li>SkipOnAllRequestFields: skip all requests; </li><li>SkipOnSpecifiedRequestFields: skip specified request fields. </li>. This option is only valid when SkipScope is ManagedRules.",
												},
												"web_security_modules_for_exception": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Specifies the security protection module for the exception rule. It is valid only when SkipScope is WebSecurityModules. The possible values are: <li>websec-mod-managed-rules: managed rules; </li><li>websec-mod-rate-limiting: rate limiting; </li><li>websec-mod-custom-rules: custom rules; </li><li>websec-mod-adaptive-control: adaptive frequency control, intelligent client filtering, slow attack protection, traffic theft protection; </li><li>websec-mod-bot: Bot management. </li>.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"managed_rules_for_exception": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Specifies the specific managed rule for the exception rule. This is only valid when SkipScope is ManagedRules and ManagedRuleGroupsForException cannot be specified.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"managed_rule_groups_for_exception": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "Specifies the managed rule group for the exception rule. This is only valid when SkipScope is ManagedRules and ManagedRulesForException cannot be specified.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"request_fields_for_exception": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the specific configuration of the exception rule to skip the specified request field. This is only valid when SkipScope is ManagedRules and SkipOption is SkipOnSpecifiedRequestFields.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"scope": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific fields to skip. Supported values:<br/>\n<li>body.json: JSON request content; in this case, Condition supports key and value, and TargetField supports key and value, for example, { \"Scope\": \"body.json\", \"Condition\": \"\", \"TargetField\": \"key\" }, which means that all parameters of JSON request content skip WAF scanning;</li>\n<li style=\"margin-top:5px\">cookie: Cookie; in this case, Condition supports key and value, and TargetField supports key and value, for example, { \"Scope\": \"cookie\", \"Condition\": \"${key} in ['account-id'] and ${value} like ['prefix-*']\", \"TargetField\": \"value\" }, which means that the Cookie parameter name is equal to account-id and the parameter value wildcard matches prefix-* to skip WAF scanning;</li>\n<li style=\"margin-top:5px\">header: HTTP header parameter; Condition supports key and value, TargetField supports key and value, for example { \"Scope\": \"header\", \"Condition\": \"${key} like ['x-auth-*']\", \"TargetField\": \"value\" }, which means that the header parameter name wildcard matches x-auth-* and skips WAF scanning; </li>\n<li style=\"margin-top:5px\">uri.query: URL encoded content/query parameter; Condition supports key and value, TargetField supports key and value, for example { \"Scope\": \"uri.query\", \"Condition\": \"${key} in ['action'] and ${value} in ['upload', 'delete']\", \"TargetField\": \"value\" }, which means that the parameter name of the URL encoded content/query parameter is equal to action And the parameter value is equal to upload or delete to skip WAF scanning;</li>\n<li style=\"margin-top:5px\">uri: request path URI; in this case, Condition must be empty, TargetField supports query, path, fullpath, for example, { \"Scope\": \"uri\", \"Condition\": \"\", \"TargetField\": \"query\" }, indicating that the request path URI only query parameters skip WAF scanning;</li>\n<li style=\"margin-top:5px\">body: request body content. In this case, Condition must be empty, TargetField supports fullbody and multipart, for example, { \"Scope\": \"body\", \"Condition\": \"\", \"TargetField\": \"fullbody\" }, indicating that the request body content is the complete request body and skips WAF scanning;</li>.",
															},
															"condition": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The expression of the specific field to be skipped must conform to the expression syntax. <br />\nCondition supports expression configuration syntax: <li> Written according to the matching condition expression syntax of the rule, supporting references to key and value. </li>.<li> Supports in, like operators, and and logical combinations. </li>.\nFor example: <li>${key} in ['x-trace-id']: parameter name is equal to x-trace-id. </li>.<li>${key} in ['x-trace-id'] and ${value} like ['Bearer *']: parameter name is equal to x-trace-id and the parameter value wildcard matches Bearer *. </li>.",
															},
															"target_field": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "When the Scope parameter uses different values, the supported values in the TargetField expression are as follows:\n<li> body.json: supports key and value</li>\n<li> cookie: supports key and value</li>\n<li> header: supports key and value</li>\n<li> uri.query: supports key and value</li>\n<li> uri: supports path, query and fullpath</li>\n<li> body: supports fullbody and multipart</li>.",
															},
														},
													},
												},
												"enabled": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Whether the exception rule is enabled. The values are: <li>on: enabled</li><li>off: disabled</li>.",
												},
											},
										},
									},
								},
							},
						},
						"bot_management": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Bot management configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Whether bot management is enabled. Valid values: `on`, `off`.",
									},
									"custom_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Bot management custom rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Rule ID. If not specified, a new rule will be created. If specified, the existing rule will be updated or deleted.",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule name.",
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether the rule is enabled. Valid values: `on`, `off`.",
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Rule priority (0-100).",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule condition in expression syntax.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Rule action with weight.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Required:    true,
																Description: "Security action configuration.",
																Elem:        securityActionSchema(),
															},
															"weight": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Action weight (10-100, must be multiples of 10). Sum of all weights must equal 100.",
															},
														},
													},
												},
											},
										},
									},
									"basic_bot_settings": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Basic bot settings.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_idc": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Source IDC configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Base action for IDC requests.",
																Elem:        securityActionSchema(),
															},
															"action_overrides": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Action overrides for specific bot types.",
																Elem:        botManagementActionOverrideSchema(),
															},
														},
													},
												},
												"search_engine_bots": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Search engine bots configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Base action for search engine bots.",
																Elem:        securityActionSchema(),
															},
															"action_overrides": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Action overrides for specific bot types.",
																Elem:        botManagementActionOverrideSchema(),
															},
														},
													},
												},
												"known_bot_categories": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Known bot categories configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Base action for known bot categories.",
																Elem:        securityActionSchema(),
															},
															"action_overrides": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Action overrides for specific bot types.",
																Elem:        botManagementActionOverrideSchema(),
															},
														},
													},
												},
												"ip_reputation": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "IP reputation configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Whether IP reputation is enabled. Valid values: `on`, `off`.",
															},
															"ip_reputation_group": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "IP reputation group configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"base_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Base action for IP reputation.",
																			Elem:        securityActionSchema(),
																		},
																		"action_overrides": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Action overrides for specific IP reputation types.",
																			Elem:        botManagementActionOverrideSchema(),
																		},
																	},
																},
															},
														},
													},
												},
												"bot_intelligence": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Bot intelligence configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Whether bot intelligence is enabled. Valid values: `on`, `off`.",
															},
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Rule ID. Output-only.",
															},
															"bot_ratings": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Bot ratings configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"high_risk_bot_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action for high risk bot requests.",
																			Elem:        securityActionSchema(),
																		},
																		"likely_bot_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action for likely bot requests.",
																			Elem:        securityActionSchema(),
																		},
																		"verified_bot_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action for verified bot requests.",
																			Elem:        securityActionSchema(),
																		},
																		"human_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action for human requests.",
																			Elem:        securityActionSchema(),
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
									"client_attestation_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Client attestation rules (beta feature).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Rule ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule name.",
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether the rule is enabled. Valid values: `on`, `off`.",
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Rule priority (0-100).",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule condition in expression syntax.",
												},
												"attester_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Client authentication method ID.",
												},
												"invalid_attestation_action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Action when attestation is invalid.",
													Elem:        securityActionSchema(),
												},
											},
										},
									},
									"browser_impersonation_detection": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Browser impersonation detection rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Rule ID.",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule name.",
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether the rule is enabled. Valid values: `on`, `off`.",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule condition in expression syntax.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Action configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bot_session_validation": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Cookie validation and session tracking.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"issue_new_bot_session_cookie": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to issue new session cookie. Valid values: `on`, `off`.",
																		},
																		"max_new_session_trigger_config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Trigger config for new session.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_new_session_count_interval": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Statistics time window for trigger threshold. Valid values: `5s`, `10s`, `15s`, `30s`, `60s`, `5m`, `10m`, `30m`, `60m`.",
																					},
																					"max_new_session_count_threshold": {
																						Type:        schema.TypeInt,
																						Optional:    true,
																						Description: "Cumulative count for trigger threshold. Range: 1-100000000.",
																					},
																				},
																			},
																		},
																		"session_expired_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action when session is expired.",
																			Elem:        securityActionSchema(),
																		},
																		"session_invalid_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action when session is invalid.",
																			Elem:        securityActionSchema(),
																		},
																		"session_rate_control": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Session rate control.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"enabled": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Whether session rate control is enabled. Valid values: `on`, `off`.",
																					},
																					"high_rate_session_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Action for high-rate session.",
																						Elem:        securityActionSchema(),
																					},
																					"mid_rate_session_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Action for mid-rate session.",
																						Elem:        securityActionSchema(),
																					},
																					"low_rate_session_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Action for low-rate session.",
																						Elem:        securityActionSchema(),
																					},
																				},
																			},
																		},
																	},
																},
															},
															"client_behavior_detection": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Client behavior detection.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"crypto_challenge_intensity": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Proof-of-work challenge intensity. Valid values: `low`, `medium`, `high`.",
																		},
																		"crypto_challenge_delay_before": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Challenge delay before execution. Valid values: `0ms`, `100ms`, `200ms`, `300ms`, `400ms`, `500ms`, `600ms`, `700ms`, `800ms`, `900ms`, `1000ms`.",
																		},
																		"max_challenge_count_interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Statistics time window for trigger threshold. Valid values: `5s`, `10s`, `15s`, `30s`, `60s`, `5m`, `10m`, `30m`, `60m`.",
																		},
																		"max_challenge_count_threshold": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Cumulative count for trigger threshold. Range: 1-100000000.",
																		},
																		"challenge_not_finished_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action when challenge not finished.",
																			Elem:        securityActionSchema(),
																		},
																		"challenge_timeout_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action when challenge timeout.",
																			Elem:        securityActionSchema(),
																		},
																		"bot_client_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Action for bot client.",
																			Elem:        securityActionSchema(),
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
						"bot_management_lite": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Basic Bot management configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"captcha_page_challenge": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "CAPTCHA page challenge configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether CAPTCHA page challenge is enabled. Valid values: `on`, `off`.",
												},
											},
										},
									},
									"ai_crawler_detection": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "AI crawler detection configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether AI crawler detection is enabled. Valid values: `on`, `off`.",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Execution action when Enabled is on. When Enabled is on, this field is required. SecurityAction Name value supports: Deny, Monitor, Allow, Challenge.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The security action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`.",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Deny.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to extend the blocking of source IP. Valid values: `on`, `off`.",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "When BlockIP is on, the IP blocking duration.",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to use custom pages. Valid values: `on`, `off`.",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize the status code of the page.",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The PageId of the custom page.",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Whether to ignore the request source suspension. Valid values: `on`, `off`.",
																		},
																	},
																},
															},
															"allow_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Allow.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"min_delay_time": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Minimum delay response time. Supported unit: seconds, range 0-5.",
																		},
																		"max_delay_time": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Maximum delay response time. Supported unit: seconds, range 5-10.",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional parameters when Name is Challenge.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The specific challenge action to be executed safely. Valid values: `JSChallenge`, `ManagedChallenge`.",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The time interval for repeating the challenge.",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication method ID.",
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
						"default_deny_security_action_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Default deny action configuration. If not specified, the existing configuration is kept.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"managed_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Managed rules default deny action configuration. Supported parameters: `return_custom_page`, `response_code`, `error_page_id`.",
										Elem:        defaultDenyActionParametersSchema(),
									},
									"other_modules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Default deny action configuration for security protection rules other than managed rules (custom rules, rate limiting and Bot management). Supported parameters: `return_custom_page`, `response_code`, `error_page_id`.",
										Elem:        defaultDenyActionParametersSchema(),
									},
								},
							},
						},
					},
				},
			},

			"security_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Security configuration. Classic web protection settings. Note: the DescribeSecurityPolicy API does not return SecurityConfig, so this field is write-only for state consistency. For each sub-configuration, if not specified, the existing API configuration is kept.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"waf_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Managed rules configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"level": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Protection level. Valid values: `loose`, `normal`, `strict`, `stricter`, `custom`.",
									},
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Global WAF mode. Valid values: `block`, `observe`.",
									},
									"waf_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Managed rule detail configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Switch. Valid values: `on`, `off`.",
												},
												"block_rule_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs to block (disable).",
												},
												"observe_rule_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs in observe mode.",
												},
											},
										},
									},
									"ai_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "AI rule engine configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "AI rule mode. Valid values: `smart_status_close`, `smart_status_open`, `smart_status_observe`.",
												},
											},
										},
									},
								},
							},
						},
						"rate_limit_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Rate limit configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"rate_limit_user_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "User-defined rate limit rules.",
										Elem:        rateLimitUserRuleSchema(),
									},
									"rate_limit_template": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Rate limit template.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Template level. Valid values: `sup_loose`, `loose`, `emergency`, `normal`, `strict`, `close`.",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Template action, e.g. `alg`.",
												},
												"rate_limit_template_detail": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Template detail. Output-only.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Template level. Output-only.",
															},
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Template ID. Output-only.",
															},
															"action": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Template action. Output-only.",
															},
															"punish_time": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Penalty time in seconds. Output-only.",
															},
															"threshold": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Rate limit threshold. Output-only.",
															},
															"period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Statistical period in seconds. Output-only.",
															},
														},
													},
												},
											},
										},
									},
									"rate_limit_intelligence": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Intelligent client filtering.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Switch. Valid values: `on`, `off`.",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Action. Valid values: `monitor`, `alg`.",
												},
												"rule_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID. Output-only.",
												},
											},
										},
									},
									"rate_limit_customizes": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Managed customized rate limit rules.",
										Elem:        rateLimitUserRuleSchema(),
									},
								},
							},
						},
						"acl_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Custom rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"acl_user_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "User-defined ACL rules.",
										Elem:        aclUserRuleSchema(),
									},
									"customizes": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Managed customized ACL rules.",
										Elem:        aclUserRuleSchema(),
									},
								},
							},
						},
						"bot_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Bot configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"bot_managed_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Generic bot managed rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Action. Valid values: `drop`, `trans`, `alg`, `monitor`.",
												},
												"rule_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID. Output-only.",
												},
												"trans_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs to allow.",
												},
												"alg_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs with JS challenge.",
												},
												"cap_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs with CAPTCHA.",
												},
												"mon_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs in monitor mode.",
												},
												"drop_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs to drop.",
												},
											},
										},
									},
									"bot_portrait_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "User portrait rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Switch. Valid values: `on`, `off`.",
												},
												"rule_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID. Output-only.",
												},
												"alg_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs with JS challenge.",
												},
												"cap_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs with CAPTCHA.",
												},
												"mon_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs in monitor mode.",
												},
												"drop_managed_ids": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeInt},
													Description: "Rule IDs to drop.",
												},
											},
										},
									},
									"intelligence_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Bot intelligence rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Switch. Valid values: `on`, `off`.",
												},
												"intelligence_rule_items": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Intelligence rule items.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"label": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Intelligence label. Valid values: `evil_bot`, `suspect_bot`, `good_bot`, `normal`.",
															},
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Action. Valid values: `drop`, `trans`, `alg`, `captcha`, `monitor`.",
															},
														},
													},
												},
											},
										},
									},
									"bot_user_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Bot user-defined rules.",
										Elem:        botUserRuleSchema(),
									},
									"alg_detect_rule": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Bot active feature detection rules.",
										Elem:        algDetectRuleSchema(),
									},
									"customizes": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Bot managed customized rules.",
										Elem:        botUserRuleSchema(),
									},
								},
							},
						},
						"switch_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Layer-7 protection master switch.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"web_switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Web master switch. Valid values: `on`, `off`. Does not affect DDoS or Bot switches.",
									},
								},
							},
						},
						"ip_table_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Basic access control.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"ip_table_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "IP table rules.",
										Elem:        ipTableRuleSchema(),
									},
								},
							},
						},
						"except_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Exception rules configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"except_user_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Exception rules detail.",
										Elem:        exceptUserRuleSchema(),
									},
								},
							},
						},
						"drop_page_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Drop page configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"waf_drop_page_detail": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Managed rule drop page.",
										Elem:        dropPageDetailSchema(),
									},
									"acl_drop_page_detail": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Custom rule drop page.",
										Elem:        dropPageDetailSchema(),
									},
								},
							},
						},
						"template_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Template configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Template ID.",
									},
									"template_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Template name.",
									},
								},
							},
						},
						"detect_length_limit_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detect length limit configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"detect_length_limit_rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Detect length limit rules.",
										Elem:        detectLengthLimitRuleSchema(),
									},
								},
							},
						},
						"slow_post_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Slow attack configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Switch. Valid values: `on`, `off`.",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Action. Valid values: `monitor`, `drop`.",
									},
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Rule ID. Output-only.",
									},
									"first_part_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "First packet configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Switch. Valid values: `on`, `off`.",
												},
												"stat_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "First segment statistical duration in seconds (default 5).",
												},
											},
										},
									},
									"slow_rate_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Slow rate configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Switch. Valid values: `on`, `off`.",
												},
												"interval": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Statistical interval in seconds.",
												},
												"threshold": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Rate threshold in bps.",
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
		preciseMatchRulesList := make([]map[string]interface{}, 0, len(respData.CustomRules.Rules))
		basicAccessRulesList := make([]map[string]interface{}, 0, len(respData.CustomRules.Rules))
		if respData.CustomRules.Rules != nil {
			for _, rules := range respData.CustomRules.Rules {
				rulesMap := map[string]interface{}{}
				ruleType := ""
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
					ruleType = *rules.RuleType
				}

				if rules.Priority != nil {
					rulesMap["priority"] = rules.Priority
				}

				if ruleType == "PreciseMatchRule" {
					preciseMatchRulesList = append(preciseMatchRulesList, rulesMap)
				} else if ruleType == "BasicAccessRule" {
					basicAccessRulesList = append(basicAccessRulesList, rulesMap)
				} else {
					continue
				}
			}

			if len(preciseMatchRulesList) > 0 {
				customRulesMap["precise_match_rules"] = preciseMatchRulesList
			}

			if len(basicAccessRulesList) > 0 {
				customRulesMap["basic_access_rules"] = basicAccessRulesList
			}

			if len(preciseMatchRulesList) > 0 || len(basicAccessRulesList) > 0 {
				securityPolicyMap["custom_rules"] = []interface{}{customRulesMap}
			}
		}
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

		if respData.ManagedRules.AutoUpdate != nil {
			autoUpdateMap := map[string]interface{}{}
			if respData.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion != nil {
				autoUpdateMap["auto_update_to_latest_version"] = respData.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion
			}

			if respData.ManagedRules.AutoUpdate.RulesetVersion != nil {
				autoUpdateMap["ruleset_version"] = respData.ManagedRules.AutoUpdate.RulesetVersion
			}

			managedRulesMap["auto_update"] = []interface{}{autoUpdateMap}
		}

		if respData.ManagedRules.ManagedRuleGroups != nil {
			managedRuleGroupsList := make([]map[string]interface{}, 0, len(respData.ManagedRules.ManagedRuleGroups))
			for _, managedRuleGroups := range respData.ManagedRules.ManagedRuleGroups {
				managedRuleGroupsMap := map[string]interface{}{}

				if managedRuleGroups.GroupId != nil {
					managedRuleGroupsMap["group_id"] = managedRuleGroups.GroupId
				}

				if managedRuleGroups.SensitivityLevel != nil {
					managedRuleGroupsMap["sensitivity_level"] = managedRuleGroups.SensitivityLevel
				}

				if managedRuleGroups.Action != nil {
					actionMap := map[string]interface{}{}
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

					if managedRuleGroups.Action.ReturnCustomPageActionParameters != nil {
						returnCustomPageActionParametersMap := map[string]interface{}{}
						if managedRuleGroups.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = managedRuleGroups.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if managedRuleGroups.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = managedRuleGroups.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					if managedRuleGroups.Action.RedirectActionParameters != nil {
						redirectActionParametersMap := map[string]interface{}{}
						if managedRuleGroups.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = managedRuleGroups.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					managedRuleGroupsMap["action"] = []interface{}{actionMap}
				}

				if managedRuleGroups.RuleActions != nil {
					ruleActionsList := make([]map[string]interface{}, 0, len(managedRuleGroups.RuleActions))
					for _, ruleActions := range managedRuleGroups.RuleActions {
						ruleActionsMap := map[string]interface{}{}
						if ruleActions.RuleId != nil {
							ruleActionsMap["rule_id"] = ruleActions.RuleId
						}

						if ruleActions.Action != nil {
							actionMap := map[string]interface{}{}
							if ruleActions.Action.Name != nil {
								actionMap["name"] = ruleActions.Action.Name
							}

							if ruleActions.Action.BlockIPActionParameters != nil {
								blockIPActionParametersMap := map[string]interface{}{}
								if ruleActions.Action.BlockIPActionParameters.Duration != nil {
									blockIPActionParametersMap["duration"] = ruleActions.Action.BlockIPActionParameters.Duration
								}

								actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
							}

							if ruleActions.Action.ReturnCustomPageActionParameters != nil {
								returnCustomPageActionParametersMap := map[string]interface{}{}
								if ruleActions.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
									returnCustomPageActionParametersMap["response_code"] = ruleActions.Action.ReturnCustomPageActionParameters.ResponseCode
								}

								if ruleActions.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
									returnCustomPageActionParametersMap["error_page_id"] = ruleActions.Action.ReturnCustomPageActionParameters.ErrorPageId
								}

								actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
							}

							if ruleActions.Action.RedirectActionParameters != nil {
								redirectActionParametersMap := map[string]interface{}{}
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

				if managedRuleGroups.MetaData != nil {
					metaDataMap := map[string]interface{}{}
					if managedRuleGroups.MetaData.GroupDetail != nil {
						metaDataMap["group_detail"] = managedRuleGroups.MetaData.GroupDetail
					}

					if managedRuleGroups.MetaData.GroupName != nil {
						metaDataMap["group_name"] = managedRuleGroups.MetaData.GroupName
					}

					if managedRuleGroups.MetaData.RuleDetails != nil {
						ruleDetailsList := make([]map[string]interface{}, 0, len(managedRuleGroups.MetaData.RuleDetails))
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

	if respData.HttpDDoSProtection != nil {
		httpDDoSProtectionMap := map[string]interface{}{}

		if respData.HttpDDoSProtection.AdaptiveFrequencyControl != nil {
			adaptiveFrequencyControlMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Enabled != nil {
				adaptiveFrequencyControlMap["enabled"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Enabled
			}

			if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Sensitivity != nil {
				adaptiveFrequencyControlMap["sensitivity"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Sensitivity
			}

			if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.Name
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_i_p_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				adaptiveFrequencyControlMap["action"] = []interface{}{actionMap}
			}

			httpDDoSProtectionMap["adaptive_frequency_control"] = []interface{}{adaptiveFrequencyControlMap}
		}

		if respData.HttpDDoSProtection.ClientFiltering != nil {
			clientFilteringMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.ClientFiltering.Enabled != nil {
				clientFilteringMap["enabled"] = respData.HttpDDoSProtection.ClientFiltering.Enabled
			}

			if respData.HttpDDoSProtection.ClientFiltering.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.ClientFiltering.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.ClientFiltering.Action.Name
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.ClientFiltering.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.ClientFiltering.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				clientFilteringMap["action"] = []interface{}{actionMap}
			}

			httpDDoSProtectionMap["client_filtering"] = []interface{}{clientFilteringMap}
		}

		if respData.HttpDDoSProtection.BandwidthAbuseDefense != nil {
			bandwidthAbuseDefenseMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.BandwidthAbuseDefense.Enabled != nil {
				bandwidthAbuseDefenseMap["enabled"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Enabled
			}

			if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.Name
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				bandwidthAbuseDefenseMap["action"] = []interface{}{actionMap}
			}

			httpDDoSProtectionMap["bandwidth_abuse_defense"] = []interface{}{bandwidthAbuseDefenseMap}
		}

		if respData.HttpDDoSProtection.SlowAttackDefense != nil {
			slowAttackDefenseMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.SlowAttackDefense.Enabled != nil {
				slowAttackDefenseMap["enabled"] = respData.HttpDDoSProtection.SlowAttackDefense.Enabled
			}

			if respData.HttpDDoSProtection.SlowAttackDefense.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.Name
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				slowAttackDefenseMap["action"] = []interface{}{actionMap}
			}

			if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate != nil {
				minimalRequestBodyTransferRateMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.MinimalAvgTransferRateThreshold != nil {
					minimalRequestBodyTransferRateMap["minimal_avg_transfer_rate_threshold"] = respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.MinimalAvgTransferRateThreshold
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.CountingPeriod != nil {
					minimalRequestBodyTransferRateMap["counting_period"] = respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.CountingPeriod
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.Enabled != nil {
					minimalRequestBodyTransferRateMap["enabled"] = respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.Enabled
				}

				slowAttackDefenseMap["minimal_request_body_transfer_rate"] = []interface{}{minimalRequestBodyTransferRateMap}
			}

			if respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout != nil {
				requestBodyTransferTimeoutMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.IdleTimeout != nil {
					requestBodyTransferTimeoutMap["idle_timeout"] = respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.IdleTimeout
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.Enabled != nil {
					requestBodyTransferTimeoutMap["enabled"] = respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.Enabled
				}

				slowAttackDefenseMap["request_body_transfer_timeout"] = []interface{}{requestBodyTransferTimeoutMap}
			}

			httpDDoSProtectionMap["slow_attack_defense"] = []interface{}{slowAttackDefenseMap}
		}

		securityPolicyMap["http_ddos_protection"] = []interface{}{httpDDoSProtectionMap}
	}

	if respData.RateLimitingRules != nil {
		rateLimitingRulesMap := map[string]interface{}{}

		if respData.RateLimitingRules.Rules != nil {
			rulesList := []interface{}{}
			for _, rules := range respData.RateLimitingRules.Rules {
				rulesMap := map[string]interface{}{}

				if rules.Id != nil {
					rulesMap["id"] = rules.Id
				}

				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Condition != nil {
					rulesMap["condition"] = rules.Condition
				}

				if rules.CountBy != nil {
					rulesMap["count_by"] = rules.CountBy
				}

				if rules.MaxRequestThreshold != nil {
					rulesMap["max_request_threshold"] = rules.MaxRequestThreshold
				}

				if rules.CountingPeriod != nil {
					rulesMap["counting_period"] = rules.CountingPeriod
				}

				if rules.ActionDuration != nil {
					rulesMap["action_duration"] = rules.ActionDuration
				}

				if rules.Action != nil {
					actionMap := map[string]interface{}{}

					if rules.Action.Name != nil {
						actionMap["name"] = rules.Action.Name
					}

					if rules.Action.DenyActionParameters != nil {
						denyActionParametersMap := map[string]interface{}{}

						if rules.Action.DenyActionParameters.BlockIp != nil {
							denyActionParametersMap["block_ip"] = rules.Action.DenyActionParameters.BlockIp
						}

						if rules.Action.DenyActionParameters.BlockIpDuration != nil {
							denyActionParametersMap["block_ip_duration"] = rules.Action.DenyActionParameters.BlockIpDuration
						}

						if rules.Action.DenyActionParameters.ReturnCustomPage != nil {
							denyActionParametersMap["return_custom_page"] = rules.Action.DenyActionParameters.ReturnCustomPage
						}

						if rules.Action.DenyActionParameters.ResponseCode != nil {
							denyActionParametersMap["response_code"] = rules.Action.DenyActionParameters.ResponseCode
						}

						if rules.Action.DenyActionParameters.ErrorPageId != nil {
							denyActionParametersMap["error_page_id"] = rules.Action.DenyActionParameters.ErrorPageId
						}

						if rules.Action.DenyActionParameters.Stall != nil {
							denyActionParametersMap["stall"] = rules.Action.DenyActionParameters.Stall
						}

						actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
					}

					if rules.Action.RedirectActionParameters != nil {
						redirectActionParametersMap := map[string]interface{}{}

						if rules.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = rules.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					if rules.Action.ChallengeActionParameters != nil {
						challengeActionParametersMap := map[string]interface{}{}

						if rules.Action.ChallengeActionParameters.ChallengeOption != nil {
							challengeActionParametersMap["challenge_option"] = rules.Action.ChallengeActionParameters.ChallengeOption
						}

						if rules.Action.ChallengeActionParameters.Interval != nil {
							challengeActionParametersMap["interval"] = rules.Action.ChallengeActionParameters.Interval
						}

						if rules.Action.ChallengeActionParameters.AttesterId != nil {
							challengeActionParametersMap["attester_id"] = rules.Action.ChallengeActionParameters.AttesterId
						}

						actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
					}

					if rules.Action.BlockIPActionParameters != nil {
						blockIPActionParametersMap := map[string]interface{}{}

						if rules.Action.BlockIPActionParameters.Duration != nil {
							blockIPActionParametersMap["duration"] = rules.Action.BlockIPActionParameters.Duration
						}

						actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
					}

					if rules.Action.ReturnCustomPageActionParameters != nil {
						returnCustomPageActionParametersMap := map[string]interface{}{}

						if rules.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = rules.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if rules.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = rules.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					rulesMap["action"] = []interface{}{actionMap}
				}

				if rules.Priority != nil {
					rulesMap["priority"] = rules.Priority
				}

				if rules.Enabled != nil {
					rulesMap["enabled"] = rules.Enabled
				}

				rulesList = append(rulesList, rulesMap)
			}

			if len(rulesList) > 0 {
				rateLimitingRulesMap["rules"] = rulesList
				securityPolicyMap["rate_limiting_rules"] = []interface{}{rateLimitingRulesMap}
			}
		}
	}

	if respData.ExceptionRules != nil {
		exceptionRulesMap := map[string]interface{}{}

		if respData.ExceptionRules.Rules != nil {
			rulesList := []interface{}{}
			for _, rules := range respData.ExceptionRules.Rules {
				rulesMap := map[string]interface{}{}

				if rules.Id != nil {
					rulesMap["id"] = rules.Id
				}

				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Condition != nil {
					rulesMap["condition"] = rules.Condition
				}

				if rules.SkipScope != nil {
					rulesMap["skip_scope"] = rules.SkipScope
				}

				if rules.SkipOption != nil {
					rulesMap["skip_option"] = rules.SkipOption
				}

				if rules.WebSecurityModulesForException != nil {
					rulesMap["web_security_modules_for_exception"] = rules.WebSecurityModulesForException
				}

				if rules.ManagedRulesForException != nil {
					rulesMap["managed_rules_for_exception"] = rules.ManagedRulesForException
				}

				if rules.ManagedRuleGroupsForException != nil {
					rulesMap["managed_rule_groups_for_exception"] = rules.ManagedRuleGroupsForException
				}

				if rules.RequestFieldsForException != nil {
					requestFieldsForExceptionList := []interface{}{}
					for _, requestFieldsForException := range rules.RequestFieldsForException {
						requestFieldsForExceptionMap := map[string]interface{}{}

						if requestFieldsForException.Scope != nil {
							requestFieldsForExceptionMap["scope"] = requestFieldsForException.Scope
						}

						if requestFieldsForException.Condition != nil {
							requestFieldsForExceptionMap["condition"] = requestFieldsForException.Condition
						}

						if requestFieldsForException.TargetField != nil {
							requestFieldsForExceptionMap["target_field"] = requestFieldsForException.TargetField
						}

						requestFieldsForExceptionList = append(requestFieldsForExceptionList, requestFieldsForExceptionMap)
					}

					rulesMap["request_fields_for_exception"] = requestFieldsForExceptionList
				}

				if rules.Enabled != nil {
					rulesMap["enabled"] = rules.Enabled
				}

				rulesList = append(rulesList, rulesMap)
			}

			if len(rulesList) > 0 {
				exceptionRulesMap["rules"] = rulesList
				securityPolicyMap["exception_rules"] = []interface{}{exceptionRulesMap}
			}
		}
	}

	// bot_management
	if respData.BotManagement != nil {
		botManagementMap := map[string]interface{}{}

		if respData.BotManagement.Enabled != nil {
			botManagementMap["enabled"] = respData.BotManagement.Enabled
		}

		// custom_rules
		if respData.BotManagement.CustomRules != nil && respData.BotManagement.CustomRules.Rules != nil {
			customRulesList := []interface{}{}
			for _, rule := range respData.BotManagement.CustomRules.Rules {
				ruleMap := map[string]interface{}{}
				if rule.Id != nil {
					ruleMap["id"] = rule.Id
				}
				if rule.Name != nil {
					ruleMap["name"] = rule.Name
				}
				if rule.Enabled != nil {
					ruleMap["enabled"] = rule.Enabled
				}
				if rule.Priority != nil {
					ruleMap["priority"] = rule.Priority
				}
				if rule.Condition != nil {
					ruleMap["condition"] = rule.Condition
				}
				if rule.Action != nil {
					actionList := []interface{}{}
					for _, weightedAction := range rule.Action {
						actionMap := map[string]interface{}{}
						if weightedAction.Weight != nil {
							actionMap["weight"] = weightedAction.Weight
						}
						if weightedAction.SecurityAction != nil {
							actionMap["action"] = []interface{}{flattenSecurityAction(weightedAction.SecurityAction)}
						}
						actionList = append(actionList, actionMap)
					}
					ruleMap["action"] = actionList
				}
				customRulesList = append(customRulesList, ruleMap)
			}
			botManagementMap["custom_rules"] = customRulesList
		}

		// basic_bot_settings
		if respData.BotManagement.BasicBotSettings != nil {
			basicBotSettingsMap := map[string]interface{}{}
			basicBotSettings := respData.BotManagement.BasicBotSettings

			// source_idc
			if basicBotSettings.SourceIDC != nil {
				sourceIDCMap := map[string]interface{}{}
				if basicBotSettings.SourceIDC.BaseAction != nil {
					sourceIDCMap["base_action"] = []interface{}{flattenSecurityAction(basicBotSettings.SourceIDC.BaseAction)}
				}
				if basicBotSettings.SourceIDC.BotManagementActionOverrides != nil {
					overridesList := []interface{}{}
					for _, override := range basicBotSettings.SourceIDC.BotManagementActionOverrides {
						overridesList = append(overridesList, flattenBotManagementActionOverride(override))
					}
					sourceIDCMap["action_overrides"] = overridesList
				}
				basicBotSettingsMap["source_idc"] = []interface{}{sourceIDCMap}
			}

			// search_engine_bots
			if basicBotSettings.SearchEngineBots != nil {
				searchEngineBotsMap := map[string]interface{}{}
				if basicBotSettings.SearchEngineBots.BaseAction != nil {
					searchEngineBotsMap["base_action"] = []interface{}{flattenSecurityAction(basicBotSettings.SearchEngineBots.BaseAction)}
				}
				if basicBotSettings.SearchEngineBots.BotManagementActionOverrides != nil {
					overridesList := []interface{}{}
					for _, override := range basicBotSettings.SearchEngineBots.BotManagementActionOverrides {
						overridesList = append(overridesList, flattenBotManagementActionOverride(override))
					}
					searchEngineBotsMap["action_overrides"] = overridesList
				}
				basicBotSettingsMap["search_engine_bots"] = []interface{}{searchEngineBotsMap}
			}

			// known_bot_categories
			if basicBotSettings.KnownBotCategories != nil {
				knownBotCategoriesMap := map[string]interface{}{}
				if basicBotSettings.KnownBotCategories.BaseAction != nil {
					knownBotCategoriesMap["base_action"] = []interface{}{flattenSecurityAction(basicBotSettings.KnownBotCategories.BaseAction)}
				}
				if basicBotSettings.KnownBotCategories.BotManagementActionOverrides != nil {
					overridesList := []interface{}{}
					for _, override := range basicBotSettings.KnownBotCategories.BotManagementActionOverrides {
						overridesList = append(overridesList, flattenBotManagementActionOverride(override))
					}
					knownBotCategoriesMap["action_overrides"] = overridesList
				}
				basicBotSettingsMap["known_bot_categories"] = []interface{}{knownBotCategoriesMap}
			}

			// ip_reputation
			if basicBotSettings.IPReputation != nil {
				ipReputationMap := map[string]interface{}{}
				if basicBotSettings.IPReputation.Enabled != nil {
					ipReputationMap["enabled"] = basicBotSettings.IPReputation.Enabled
				}
				if basicBotSettings.IPReputation.IPReputationGroup != nil {
					ipReputationGroupMap := map[string]interface{}{}
					if basicBotSettings.IPReputation.IPReputationGroup.BaseAction != nil {
						ipReputationGroupMap["base_action"] = []interface{}{flattenSecurityAction(basicBotSettings.IPReputation.IPReputationGroup.BaseAction)}
					}
					if basicBotSettings.IPReputation.IPReputationGroup.BotManagementActionOverrides != nil {
						overridesList := []interface{}{}
						for _, override := range basicBotSettings.IPReputation.IPReputationGroup.BotManagementActionOverrides {
							overridesList = append(overridesList, flattenBotManagementActionOverride(override))
						}
						ipReputationGroupMap["action_overrides"] = overridesList
					}
					ipReputationMap["ip_reputation_group"] = []interface{}{ipReputationGroupMap}
				}
				basicBotSettingsMap["ip_reputation"] = []interface{}{ipReputationMap}
			}

			// bot_intelligence
			if basicBotSettings.BotIntelligence != nil {
				botIntelligenceMap := map[string]interface{}{}
				if basicBotSettings.BotIntelligence.Enabled != nil {
					botIntelligenceMap["enabled"] = basicBotSettings.BotIntelligence.Enabled
				}
				if basicBotSettings.BotIntelligence.BotRatings != nil {
					botRatingsMap := map[string]interface{}{}
					if basicBotSettings.BotIntelligence.BotRatings.HighRiskBotRequestsAction != nil {
						botRatingsMap["high_risk_bot_requests_action"] = []interface{}{flattenSecurityAction(basicBotSettings.BotIntelligence.BotRatings.HighRiskBotRequestsAction)}
					}
					if basicBotSettings.BotIntelligence.BotRatings.LikelyBotRequestsAction != nil {
						botRatingsMap["likely_bot_requests_action"] = []interface{}{flattenSecurityAction(basicBotSettings.BotIntelligence.BotRatings.LikelyBotRequestsAction)}
					}
					if basicBotSettings.BotIntelligence.BotRatings.VerifiedBotRequestsAction != nil {
						botRatingsMap["verified_bot_requests_action"] = []interface{}{flattenSecurityAction(basicBotSettings.BotIntelligence.BotRatings.VerifiedBotRequestsAction)}
					}
					if basicBotSettings.BotIntelligence.BotRatings.HumanRequestsAction != nil {
						botRatingsMap["human_requests_action"] = []interface{}{flattenSecurityAction(basicBotSettings.BotIntelligence.BotRatings.HumanRequestsAction)}
					}
					botIntelligenceMap["bot_ratings"] = []interface{}{botRatingsMap}
				}
				basicBotSettingsMap["bot_intelligence"] = []interface{}{botIntelligenceMap}
			}

			botManagementMap["basic_bot_settings"] = []interface{}{basicBotSettingsMap}
		}

		// client_attestation_rules
		if respData.BotManagement.ClientAttestationRules != nil && respData.BotManagement.ClientAttestationRules.Rules != nil {
			rulesList := []interface{}{}
			for _, rule := range respData.BotManagement.ClientAttestationRules.Rules {
				ruleMap := map[string]interface{}{}
				if rule.Id != nil {
					ruleMap["id"] = rule.Id
				}
				if rule.Name != nil {
					ruleMap["name"] = rule.Name
				}
				if rule.Enabled != nil {
					ruleMap["enabled"] = rule.Enabled
				}
				if rule.Priority != nil {
					ruleMap["priority"] = rule.Priority
				}
				if rule.Condition != nil {
					ruleMap["condition"] = rule.Condition
				}
				if rule.AttesterId != nil {
					ruleMap["attester_id"] = rule.AttesterId
				}
				if rule.InvalidAttestationAction != nil {
					ruleMap["invalid_attestation_action"] = []interface{}{flattenSecurityAction(rule.InvalidAttestationAction)}
				}
				rulesList = append(rulesList, ruleMap)
			}
			botManagementMap["client_attestation_rules"] = rulesList
		}

		// browser_impersonation_detection
		if respData.BotManagement.BrowserImpersonationDetection != nil && respData.BotManagement.BrowserImpersonationDetection.Rules != nil {
			rulesList := []interface{}{}
			for _, rule := range respData.BotManagement.BrowserImpersonationDetection.Rules {
				ruleMap := map[string]interface{}{}
				if rule.Id != nil {
					ruleMap["id"] = rule.Id
				}
				if rule.Name != nil {
					ruleMap["name"] = rule.Name
				}
				if rule.Enabled != nil {
					ruleMap["enabled"] = rule.Enabled
				}
				if rule.Condition != nil {
					ruleMap["condition"] = rule.Condition
				}
				if rule.Action != nil {
					actionMap := map[string]interface{}{}

					// bot_session_validation
					if rule.Action.BotSessionValidation != nil {
						botSessionMap := map[string]interface{}{}
						if rule.Action.BotSessionValidation.IssueNewBotSessionCookie != nil {
							botSessionMap["issue_new_bot_session_cookie"] = rule.Action.BotSessionValidation.IssueNewBotSessionCookie
						}
						if rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig != nil {
							triggerConfigMap := map[string]interface{}{}
							if rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountInterval != nil {
								triggerConfigMap["max_new_session_count_interval"] = rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountInterval
							}
							if rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountThreshold != nil {
								triggerConfigMap["max_new_session_count_threshold"] = rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountThreshold
							}
							botSessionMap["max_new_session_trigger_config"] = []interface{}{triggerConfigMap}
						}
						if rule.Action.BotSessionValidation.SessionExpiredAction != nil {
							botSessionMap["session_expired_action"] = []interface{}{flattenSecurityAction(rule.Action.BotSessionValidation.SessionExpiredAction)}
						}
						if rule.Action.BotSessionValidation.SessionInvalidAction != nil {
							botSessionMap["session_invalid_action"] = []interface{}{flattenSecurityAction(rule.Action.BotSessionValidation.SessionInvalidAction)}
						}
						if rule.Action.BotSessionValidation.SessionRateControl != nil {
							sessionRateMap := map[string]interface{}{}
							if rule.Action.BotSessionValidation.SessionRateControl.Enabled != nil {
								sessionRateMap["enabled"] = rule.Action.BotSessionValidation.SessionRateControl.Enabled
							}
							if rule.Action.BotSessionValidation.SessionRateControl.HighRateSessionAction != nil {
								sessionRateMap["high_rate_session_action"] = []interface{}{flattenSecurityAction(rule.Action.BotSessionValidation.SessionRateControl.HighRateSessionAction)}
							}
							if rule.Action.BotSessionValidation.SessionRateControl.MidRateSessionAction != nil {
								sessionRateMap["mid_rate_session_action"] = []interface{}{flattenSecurityAction(rule.Action.BotSessionValidation.SessionRateControl.MidRateSessionAction)}
							}
							if rule.Action.BotSessionValidation.SessionRateControl.LowRateSessionAction != nil {
								sessionRateMap["low_rate_session_action"] = []interface{}{flattenSecurityAction(rule.Action.BotSessionValidation.SessionRateControl.LowRateSessionAction)}
							}
							botSessionMap["session_rate_control"] = []interface{}{sessionRateMap}
						}
						actionMap["bot_session_validation"] = []interface{}{botSessionMap}
					}

					// client_behavior_detection
					if rule.Action.ClientBehaviorDetection != nil {
						clientBehaviorMap := map[string]interface{}{}
						if rule.Action.ClientBehaviorDetection.CryptoChallengeIntensity != nil {
							clientBehaviorMap["crypto_challenge_intensity"] = rule.Action.ClientBehaviorDetection.CryptoChallengeIntensity
						}
						if rule.Action.ClientBehaviorDetection.CryptoChallengeDelayBefore != nil {
							clientBehaviorMap["crypto_challenge_delay_before"] = rule.Action.ClientBehaviorDetection.CryptoChallengeDelayBefore
						}
						if rule.Action.ClientBehaviorDetection.MaxChallengeCountInterval != nil {
							clientBehaviorMap["max_challenge_count_interval"] = rule.Action.ClientBehaviorDetection.MaxChallengeCountInterval
						}
						if rule.Action.ClientBehaviorDetection.MaxChallengeCountThreshold != nil {
							clientBehaviorMap["max_challenge_count_threshold"] = rule.Action.ClientBehaviorDetection.MaxChallengeCountThreshold
						}
						if rule.Action.ClientBehaviorDetection.ChallengeNotFinishedAction != nil {
							clientBehaviorMap["challenge_not_finished_action"] = []interface{}{flattenSecurityAction(rule.Action.ClientBehaviorDetection.ChallengeNotFinishedAction)}
						}
						if rule.Action.ClientBehaviorDetection.ChallengeTimeoutAction != nil {
							clientBehaviorMap["challenge_timeout_action"] = []interface{}{flattenSecurityAction(rule.Action.ClientBehaviorDetection.ChallengeTimeoutAction)}
						}
						if rule.Action.ClientBehaviorDetection.BotClientAction != nil {
							clientBehaviorMap["bot_client_action"] = []interface{}{flattenSecurityAction(rule.Action.ClientBehaviorDetection.BotClientAction)}
						}
						actionMap["client_behavior_detection"] = []interface{}{clientBehaviorMap}
					}

					ruleMap["action"] = []interface{}{actionMap}
				}
				rulesList = append(rulesList, ruleMap)
			}
			botManagementMap["browser_impersonation_detection"] = rulesList
		}

		securityPolicyMap["bot_management"] = []interface{}{botManagementMap}
	}

	if respData.BotManagementLite != nil {
		botManagementLiteMap := map[string]interface{}{}

		if respData.BotManagementLite.CAPTCHAPageChallenge != nil {
			captchaPageChallengeMap := map[string]interface{}{}

			if respData.BotManagementLite.CAPTCHAPageChallenge.Enabled != nil {
				captchaPageChallengeMap["enabled"] = respData.BotManagementLite.CAPTCHAPageChallenge.Enabled
			}

			botManagementLiteMap["captcha_page_challenge"] = []interface{}{captchaPageChallengeMap}
		}

		if respData.BotManagementLite.AICrawlerDetection != nil {
			aiCrawlerDetectionMap := map[string]interface{}{}

			if respData.BotManagementLite.AICrawlerDetection.Enabled != nil {
				aiCrawlerDetectionMap["enabled"] = respData.BotManagementLite.AICrawlerDetection.Enabled
			}

			if respData.BotManagementLite.AICrawlerDetection.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.BotManagementLite.AICrawlerDetection.Action.Name != nil {
					actionMap["name"] = respData.BotManagementLite.AICrawlerDetection.Action.Name
				}

				if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIp
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.ResponseCode
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.ErrorPageId
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.BotManagementLite.AICrawlerDetection.Action.AllowActionParameters != nil {
					allowActionParametersMap := map[string]interface{}{}

					if respData.BotManagementLite.AICrawlerDetection.Action.AllowActionParameters.MinDelayTime != nil {
						allowActionParametersMap["min_delay_time"] = respData.BotManagementLite.AICrawlerDetection.Action.AllowActionParameters.MinDelayTime
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.AllowActionParameters.MaxDelayTime != nil {
						allowActionParametersMap["max_delay_time"] = respData.BotManagementLite.AICrawlerDetection.Action.AllowActionParameters.MaxDelayTime
					}

					actionMap["allow_action_parameters"] = []interface{}{allowActionParametersMap}
				}

				if respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters.Interval
					}

					if respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.BotManagementLite.AICrawlerDetection.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				aiCrawlerDetectionMap["action"] = []interface{}{actionMap}
			}

			botManagementLiteMap["ai_crawler_detection"] = []interface{}{aiCrawlerDetectionMap}
		}

		securityPolicyMap["bot_management_lite"] = []interface{}{botManagementLiteMap}
	}

	if respData.DefaultDenySecurityActionParameters != nil {
		defaultDenyMap := map[string]interface{}{}
		if respData.DefaultDenySecurityActionParameters.ManagedRules != nil {
			defaultDenyMap["managed_rules"] = []interface{}{flattenDenyActionParameters(respData.DefaultDenySecurityActionParameters.ManagedRules)}
		}
		if respData.DefaultDenySecurityActionParameters.OtherModules != nil {
			defaultDenyMap["other_modules"] = []interface{}{flattenDenyActionParameters(respData.DefaultDenySecurityActionParameters.OtherModules)}
		}
		securityPolicyMap["default_deny_security_action_parameters"] = []interface{}{defaultDenyMap}
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

	if securityPolicyMap, ok := helper.InterfacesHeadMap(d, "security_policy"); ok {
		securityPolicy := teov20220901.SecurityPolicy{}
		if customRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["custom_rules"]); ok {
			customRules := teov20220901.CustomRules{}
			if v, ok := customRulesMap["rules"]; ok {
				if len(v.([]interface{})) > 0 {
					return fmt.Errorf("`rules` has been deprecated from version 1.81.184. Please use `precise_match_rules` or `basic_access_rules` instead.")
				}
			}

			if v, ok := customRulesMap["precise_match_rules"]; ok {
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

					customRule.RuleType = helper.String("PreciseMatchRule")

					if v, ok := rulesMap["priority"].(int); ok {
						customRule.Priority = helper.IntInt64(v)
					}

					customRules.Rules = append(customRules.Rules, &customRule)
				}
			}

			if v, ok := customRulesMap["basic_access_rules"]; ok {
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

					customRule.RuleType = helper.String("BasicAccessRule")

					if v, ok := rulesMap["priority"].(int); ok {
						customRule.Priority = helper.IntInt64(v)
					}

					customRules.Rules = append(customRules.Rules, &customRule)
				}
			}

			securityPolicy.CustomRules = &customRules
		} else {
			securityPolicy.CustomRules = &teov20220901.CustomRules{
				Rules: []*teov20220901.CustomRule{},
			}
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

		if httpDDoSProtectionMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["http_ddos_protection"]); ok {
			httpDDoSProtection := teov20220901.HttpDDoSProtection{}
			if adaptiveFrequencyControlMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "adaptive_frequency_control"); ok {
				adaptiveFrequencyControl := teov20220901.AdaptiveFrequencyControl{}
				if v, ok := adaptiveFrequencyControlMap["enabled"]; ok {
					adaptiveFrequencyControl.Enabled = helper.String(v.(string))
				}

				if v, ok := adaptiveFrequencyControlMap["sensitivity"]; ok {
					adaptiveFrequencyControl.Sensitivity = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(adaptiveFrequencyControlMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}

					adaptiveFrequencyControl.Action = &securityAction
				}

				httpDDoSProtection.AdaptiveFrequencyControl = &adaptiveFrequencyControl
			}

			if clientFilteringMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "client_filtering"); ok {
				clientFiltering := teov20220901.ClientFiltering{}
				if v, ok := clientFilteringMap["enabled"]; ok {
					clientFiltering.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(clientFilteringMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}

					clientFiltering.Action = &securityAction
				}

				httpDDoSProtection.ClientFiltering = &clientFiltering
			}

			if bandwidthAbuseDefenseMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "bandwidth_abuse_defense"); ok {
				bandwidthAbuseDefense := teov20220901.BandwidthAbuseDefense{}
				if v, ok := bandwidthAbuseDefenseMap["enabled"]; ok {
					bandwidthAbuseDefense.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(bandwidthAbuseDefenseMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}
					bandwidthAbuseDefense.Action = &securityAction
				}

				httpDDoSProtection.BandwidthAbuseDefense = &bandwidthAbuseDefense
			}

			if slowAttackDefenseMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "slow_attack_defense"); ok {
				slowAttackDefense := teov20220901.SlowAttackDefense{}
				if v, ok := slowAttackDefenseMap["enabled"]; ok {
					slowAttackDefense.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(slowAttackDefenseMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}

					slowAttackDefense.Action = &securityAction
				}

				if minimalRequestBodyTransferRateMap, ok := helper.InterfaceToMap(slowAttackDefenseMap, "minimal_request_body_transfer_rate"); ok {
					minimalRequestBodyTransferRate := teov20220901.MinimalRequestBodyTransferRate{}
					if v, ok := minimalRequestBodyTransferRateMap["minimal_avg_transfer_rate_threshold"]; ok {
						minimalRequestBodyTransferRate.MinimalAvgTransferRateThreshold = helper.String(v.(string))
					}

					if v, ok := minimalRequestBodyTransferRateMap["counting_period"]; ok {
						minimalRequestBodyTransferRate.CountingPeriod = helper.String(v.(string))
					}

					if v, ok := minimalRequestBodyTransferRateMap["enabled"]; ok {
						minimalRequestBodyTransferRate.Enabled = helper.String(v.(string))
					}

					slowAttackDefense.MinimalRequestBodyTransferRate = &minimalRequestBodyTransferRate
				}

				if requestBodyTransferTimeoutMap, ok := helper.InterfaceToMap(slowAttackDefenseMap, "request_body_transfer_timeout"); ok {
					requestBodyTransferTimeout := teov20220901.RequestBodyTransferTimeout{}
					if v, ok := requestBodyTransferTimeoutMap["idle_timeout"]; ok {
						requestBodyTransferTimeout.IdleTimeout = helper.String(v.(string))
					}

					if v, ok := requestBodyTransferTimeoutMap["enabled"]; ok {
						requestBodyTransferTimeout.Enabled = helper.String(v.(string))
					}

					slowAttackDefense.RequestBodyTransferTimeout = &requestBodyTransferTimeout
				}

				httpDDoSProtection.SlowAttackDefense = &slowAttackDefense
			}

			securityPolicy.HttpDDoSProtection = &httpDDoSProtection
		}

		if rateLimitingRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["rate_limiting_rules"]); ok {
			rateLimitingRules := teov20220901.RateLimitingRules{}
			if v, ok := rateLimitingRulesMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					rateLimitingRule := teov20220901.RateLimitingRule{}
					if v, ok := rulesMap["id"]; ok {
						rateLimitingRule.Id = helper.String(v.(string))
					}

					if v, ok := rulesMap["name"]; ok {
						rateLimitingRule.Name = helper.String(v.(string))
					}

					if v, ok := rulesMap["condition"]; ok {
						rateLimitingRule.Condition = helper.String(v.(string))
					}

					if v, ok := rulesMap["count_by"]; ok {
						countBySet := v.(*schema.Set).List()
						for i := range countBySet {
							if countBySet[i] != nil {
								countBy := countBySet[i].(string)
								rateLimitingRule.CountBy = append(rateLimitingRule.CountBy, &countBy)
							}
						}
					}

					if v, ok := rulesMap["max_request_threshold"]; ok {
						rateLimitingRule.MaxRequestThreshold = helper.IntInt64(v.(int))
					}

					if v, ok := rulesMap["counting_period"]; ok {
						rateLimitingRule.CountingPeriod = helper.String(v.(string))
					}

					if v, ok := rulesMap["action_duration"]; ok {
						rateLimitingRule.ActionDuration = helper.String(v.(string))
					}

					if actionMap, ok := helper.InterfaceToMap(rulesMap, "action"); ok {
						securityAction := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"]; ok {
							securityAction.Name = helper.String(v.(string))
						}

						if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
							denyActionParameters := teov20220901.DenyActionParameters{}
							if v, ok := denyActionParametersMap["block_ip"]; ok {
								denyActionParameters.BlockIp = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
								denyActionParameters.BlockIpDuration = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["return_custom_page"]; ok {
								denyActionParameters.ReturnCustomPage = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["response_code"]; ok {
								denyActionParameters.ResponseCode = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["error_page_id"]; ok {
								denyActionParameters.ErrorPageId = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["stall"]; ok {
								denyActionParameters.Stall = helper.String(v.(string))
							}

							securityAction.DenyActionParameters = &denyActionParameters
						}

						if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
							redirectActionParameters := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"]; ok {
								redirectActionParameters.URL = helper.String(v.(string))
							}

							securityAction.RedirectActionParameters = &redirectActionParameters
						}

						if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
							challengeActionParameters := teov20220901.ChallengeActionParameters{}
							if v, ok := challengeActionParametersMap["challenge_option"]; ok {
								challengeActionParameters.ChallengeOption = helper.String(v.(string))
							}

							if v, ok := challengeActionParametersMap["interval"]; ok {
								challengeActionParameters.Interval = helper.String(v.(string))
							}

							if v, ok := challengeActionParametersMap["attester_id"]; ok {
								challengeActionParameters.AttesterId = helper.String(v.(string))
							}

							securityAction.ChallengeActionParameters = &challengeActionParameters
						}

						if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
							blockIPActionParameters := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"]; ok {
								blockIPActionParameters.Duration = helper.String(v.(string))
							}

							securityAction.BlockIPActionParameters = &blockIPActionParameters
						}

						if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
							returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
								returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
								returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
							}

							securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
						}

						rateLimitingRule.Action = &securityAction
					}

					if v, ok := rulesMap["priority"]; ok {
						rateLimitingRule.Priority = helper.IntInt64(v.(int))
					}

					if v, ok := rulesMap["enabled"]; ok {
						rateLimitingRule.Enabled = helper.String(v.(string))
					}

					rateLimitingRules.Rules = append(rateLimitingRules.Rules, &rateLimitingRule)
				}

				securityPolicy.RateLimitingRules = &rateLimitingRules
			}
		}

		if exceptionRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["exception_rules"]); ok {
			exceptionRules := teov20220901.ExceptionRules{}
			if v, ok := exceptionRulesMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					exceptionRule := teov20220901.ExceptionRule{}
					if v, ok := rulesMap["id"]; ok {
						exceptionRule.Id = helper.String(v.(string))
					}

					if v, ok := rulesMap["name"]; ok {
						exceptionRule.Name = helper.String(v.(string))
					}

					if v, ok := rulesMap["condition"]; ok {
						exceptionRule.Condition = helper.String(v.(string))
					}

					if v, ok := rulesMap["skip_scope"]; ok {
						exceptionRule.SkipScope = helper.String(v.(string))
					}

					if v, ok := rulesMap["skip_option"]; ok {
						exceptionRule.SkipOption = helper.String(v.(string))
					}

					if v, ok := rulesMap["web_security_modules_for_exception"]; ok {
						webSecurityModulesForExceptionSet := v.(*schema.Set).List()
						for i := range webSecurityModulesForExceptionSet {
							if webSecurityModulesForExceptionSet[i] != nil {
								webSecurityModulesForException := webSecurityModulesForExceptionSet[i].(string)
								exceptionRule.WebSecurityModulesForException = append(exceptionRule.WebSecurityModulesForException, &webSecurityModulesForException)
							}
						}
					}

					if v, ok := rulesMap["managed_rules_for_exception"]; ok {
						managedRulesForExceptionSet := v.(*schema.Set).List()
						for i := range managedRulesForExceptionSet {
							if managedRulesForExceptionSet[i] != nil {
								managedRulesForException := managedRulesForExceptionSet[i].(string)
								exceptionRule.ManagedRulesForException = append(exceptionRule.ManagedRulesForException, &managedRulesForException)
							}
						}
					}

					if v, ok := rulesMap["managed_rule_groups_for_exception"]; ok {
						managedRuleGroupsForExceptionSet := v.(*schema.Set).List()
						for i := range managedRuleGroupsForExceptionSet {
							if managedRuleGroupsForExceptionSet[i] != nil {
								managedRuleGroupsForException := managedRuleGroupsForExceptionSet[i].(string)
								exceptionRule.ManagedRuleGroupsForException = append(exceptionRule.ManagedRuleGroupsForException, &managedRuleGroupsForException)
							}
						}
					}

					if v, ok := rulesMap["request_fields_for_exception"]; ok {
						for _, item := range v.([]interface{}) {
							requestFieldsForExceptionMap := item.(map[string]interface{})
							requestFieldsForException := teov20220901.RequestFieldsForException{}
							if v, ok := requestFieldsForExceptionMap["scope"]; ok {
								requestFieldsForException.Scope = helper.String(v.(string))
							}

							if v, ok := requestFieldsForExceptionMap["condition"]; ok {
								requestFieldsForException.Condition = helper.String(v.(string))
							}

							if v, ok := requestFieldsForExceptionMap["target_field"]; ok {
								requestFieldsForException.TargetField = helper.String(v.(string))
							}

							exceptionRule.RequestFieldsForException = append(exceptionRule.RequestFieldsForException, &requestFieldsForException)
						}
					}

					if v, ok := rulesMap["enabled"]; ok {
						exceptionRule.Enabled = helper.String(v.(string))
					}

					exceptionRules.Rules = append(exceptionRules.Rules, &exceptionRule)
				}

				securityPolicy.ExceptionRules = &exceptionRules
			}
		} else {
			securityPolicy.ExceptionRules = &teov20220901.ExceptionRules{
				Rules: []*teov20220901.ExceptionRule{},
			}
		}

		if botManagementMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["bot_management"]); ok {
			botManagement := teov20220901.BotManagement{}
			if v, ok := botManagementMap["enabled"].(string); ok && v != "" {
				botManagement.Enabled = helper.String(v)
			}

			// custom_rules
			if v, ok := botManagementMap["custom_rules"]; ok {
				customRules := teov20220901.BotManagementCustomRules{}
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					customRule := teov20220901.BotManagementCustomRule{}
					if v, ok := ruleMap["id"].(string); ok {
						customRule.Id = helper.String(v)
					}
					if v, ok := ruleMap["name"].(string); ok && v != "" {
						customRule.Name = helper.String(v)
					}
					if v, ok := ruleMap["enabled"].(string); ok && v != "" {
						customRule.Enabled = helper.String(v)
					}
					if v, ok := ruleMap["priority"].(int); ok {
						customRule.Priority = helper.IntInt64(v)
					}
					if v, ok := ruleMap["condition"].(string); ok && v != "" {
						customRule.Condition = helper.String(v)
					}
					if v, ok := ruleMap["action"]; ok {
						for _, actionItem := range v.([]interface{}) {
							actionMap := actionItem.(map[string]interface{})
							weightedAction := teov20220901.SecurityWeightedAction{}
							if actionMap["action"] != nil {
								if securityActionMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["action"]); ok {
									securityAction := teov20220901.SecurityAction{}
									if v, ok := securityActionMap["name"].(string); ok && v != "" {
										securityAction.Name = helper.String(v)
									}
									if denyMap, ok := helper.ConvertInterfacesHeadToMap(securityActionMap["deny_action_parameters"]); ok {
										denyParams := teov20220901.DenyActionParameters{}
										if v, ok := denyMap["block_ip"].(string); ok && v != "" {
											denyParams.BlockIp = helper.String(v)
										}
										if v, ok := denyMap["block_ip_duration"].(string); ok && v != "" {
											denyParams.BlockIpDuration = helper.String(v)
										}
										if v, ok := denyMap["return_custom_page"].(string); ok && v != "" {
											denyParams.ReturnCustomPage = helper.String(v)
										}
										if v, ok := denyMap["response_code"].(string); ok && v != "" {
											denyParams.ResponseCode = helper.String(v)
										}
										if v, ok := denyMap["error_page_id"].(string); ok && v != "" {
											denyParams.ErrorPageId = helper.String(v)
										}
										if v, ok := denyMap["stall"].(string); ok && v != "" {
											denyParams.Stall = helper.String(v)
										}
										securityAction.DenyActionParameters = &denyParams
									}
									if redirectMap, ok := helper.ConvertInterfacesHeadToMap(securityActionMap["redirect_action_parameters"]); ok {
										redirectParams := teov20220901.RedirectActionParameters{}
										if v, ok := redirectMap["url"].(string); ok && v != "" {
											redirectParams.URL = helper.String(v)
										}
										securityAction.RedirectActionParameters = &redirectParams
									}
									if allowMap, ok := helper.ConvertInterfacesHeadToMap(securityActionMap["allow_action_parameters"]); ok {
										allowParams := teov20220901.AllowActionParameters{}
										if v, ok := allowMap["min_delay_time"].(string); ok && v != "" {
											allowParams.MinDelayTime = helper.String(v)
										}
										if v, ok := allowMap["max_delay_time"].(string); ok && v != "" {
											allowParams.MaxDelayTime = helper.String(v)
										}
										securityAction.AllowActionParameters = &allowParams
									}
									if challengeMap, ok := helper.ConvertInterfacesHeadToMap(securityActionMap["challenge_action_parameters"]); ok {
										challengeParams := teov20220901.ChallengeActionParameters{}
										if v, ok := challengeMap["challenge_option"].(string); ok && v != "" {
											challengeParams.ChallengeOption = helper.String(v)
										}
										if v, ok := challengeMap["interval"].(string); ok && v != "" {
											challengeParams.Interval = helper.String(v)
										}
										if v, ok := challengeMap["attester_id"].(string); ok && v != "" {
											challengeParams.AttesterId = helper.String(v)
										}
										securityAction.ChallengeActionParameters = &challengeParams
									}
									weightedAction.SecurityAction = &securityAction
								}
							}
							if v, ok := actionMap["weight"].(int); ok {
								weightedAction.Weight = helper.IntInt64(v)
							}
							customRule.Action = append(customRule.Action, &weightedAction)
						}
					}
					customRules.Rules = append(customRules.Rules, &customRule)
				}
				botManagement.CustomRules = &customRules
			}

			// basic_bot_settings
			if basicBotSettingsMap, ok := helper.ConvertInterfacesHeadToMap(botManagementMap["basic_bot_settings"]); ok {
				basicBotSettings := teov20220901.BasicBotSettings{}

				if sourceIDCMaps, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["source_idc"]); ok {
					sourceIDC := teov20220901.SourceIDC{}
					if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(sourceIDCMaps["base_action"]); ok {
						sourceIDC.BaseAction = buildSecurityActionFromMap(baseActionMap)
					}
					if v, ok := sourceIDCMaps["action_overrides"]; ok {
						for _, item := range v.([]interface{}) {
							actionOverrideMap := item.(map[string]interface{})
							actionOverride := buildBotManagementActionOverrideFromMap(actionOverrideMap)
							sourceIDC.BotManagementActionOverrides = append(sourceIDC.BotManagementActionOverrides, actionOverride)
						}
					}
					basicBotSettings.SourceIDC = &sourceIDC
				}

				if searchEngineBotsMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["search_engine_bots"]); ok {
					searchEngineBots := teov20220901.SearchEngineBots{}
					if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(searchEngineBotsMap["base_action"]); ok {
						searchEngineBots.BaseAction = buildSecurityActionFromMap(baseActionMap)
					}
					if v, ok := searchEngineBotsMap["action_overrides"]; ok {
						for _, item := range v.([]interface{}) {
							actionOverrideMap := item.(map[string]interface{})
							actionOverride := buildBotManagementActionOverrideFromMap(actionOverrideMap)
							searchEngineBots.BotManagementActionOverrides = append(searchEngineBots.BotManagementActionOverrides, actionOverride)
						}
					}
					basicBotSettings.SearchEngineBots = &searchEngineBots
				}

				if knownBotCategoriesMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["known_bot_categories"]); ok {
					knownBotCategories := teov20220901.KnownBotCategories{}
					if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(knownBotCategoriesMap["base_action"]); ok {
						knownBotCategories.BaseAction = buildSecurityActionFromMap(baseActionMap)
					}
					if v, ok := knownBotCategoriesMap["action_overrides"]; ok {
						for _, item := range v.([]interface{}) {
							actionOverrideMap := item.(map[string]interface{})
							actionOverride := buildBotManagementActionOverrideFromMap(actionOverrideMap)
							knownBotCategories.BotManagementActionOverrides = append(knownBotCategories.BotManagementActionOverrides, actionOverride)
						}
					}
					basicBotSettings.KnownBotCategories = &knownBotCategories
				}

				if ipReputationMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["ip_reputation"]); ok {
					ipReputation := teov20220901.IPReputation{}
					if v, ok := ipReputationMap["enabled"].(string); ok && v != "" {
						ipReputation.Enabled = helper.String(v)
					}
					if ipReputationGroupMap, ok := helper.ConvertInterfacesHeadToMap(ipReputationMap["ip_reputation_group"]); ok {
						ipReputationGroup := teov20220901.IPReputationGroup{}
						if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(ipReputationGroupMap["base_action"]); ok {
							ipReputationGroup.BaseAction = buildSecurityActionFromMap(baseActionMap)
						}
						if v, ok := ipReputationGroupMap["action_overrides"]; ok {
							for _, item := range v.([]interface{}) {
								actionOverrideMap := item.(map[string]interface{})
								actionOverride := buildBotManagementActionOverrideFromMap(actionOverrideMap)
								ipReputationGroup.BotManagementActionOverrides = append(ipReputationGroup.BotManagementActionOverrides, actionOverride)
							}
						}
						ipReputation.IPReputationGroup = &ipReputationGroup
					}
					basicBotSettings.IPReputation = &ipReputation
				}

				if botIntelligenceMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["bot_intelligence"]); ok {
					botIntelligence := teov20220901.BotIntelligence{}
					if v, ok := botIntelligenceMap["enabled"].(string); ok && v != "" {
						botIntelligence.Enabled = helper.String(v)
					}
					if botRatingsMap, ok := helper.ConvertInterfacesHeadToMap(botIntelligenceMap["bot_ratings"]); ok {
						botRatings := teov20220901.BotRatings{}
						if highRiskMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["high_risk_bot_requests_action"]); ok {
							botRatings.HighRiskBotRequestsAction = buildSecurityActionFromMap(highRiskMap)
						}
						if likelyMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["likely_bot_requests_action"]); ok {
							botRatings.LikelyBotRequestsAction = buildSecurityActionFromMap(likelyMap)
						}
						if verifiedMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["verified_bot_requests_action"]); ok {
							botRatings.VerifiedBotRequestsAction = buildSecurityActionFromMap(verifiedMap)
						}
						if humanMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["human_requests_action"]); ok {
							botRatings.HumanRequestsAction = buildSecurityActionFromMap(humanMap)
						}
						botIntelligence.BotRatings = &botRatings
					}
					basicBotSettings.BotIntelligence = &botIntelligence
				}

				botManagement.BasicBotSettings = &basicBotSettings
			}

			// client_attestation_rules (beta feature)
			if v, ok := botManagementMap["client_attestation_rules"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					clientAttestationRule := teov20220901.ClientAttestationRule{}
					if v, ok := ruleMap["id"].(string); ok {
						clientAttestationRule.Id = helper.String(v)
					}
					if v, ok := ruleMap["name"].(string); ok && v != "" {
						clientAttestationRule.Name = helper.String(v)
					}
					if v, ok := ruleMap["enabled"].(string); ok && v != "" {
						clientAttestationRule.Enabled = helper.String(v)
					}
					if v, ok := ruleMap["priority"].(int); ok {
						clientAttestationRule.Priority = helper.IntUint64(v)
					}
					if v, ok := ruleMap["condition"].(string); ok && v != "" {
						clientAttestationRule.Condition = helper.String(v)
					}
					if v, ok := ruleMap["attester_id"].(string); ok && v != "" {
						clientAttestationRule.AttesterId = helper.String(v)
					}
					if invalidActionMap, ok := helper.ConvertInterfacesHeadToMap(ruleMap["invalid_attestation_action"]); ok {
						clientAttestationRule.InvalidAttestationAction = buildSecurityActionFromMap(invalidActionMap)
					}
					botManagement.ClientAttestationRules = &teov20220901.ClientAttestationRules{
						Rules: []*teov20220901.ClientAttestationRule{&clientAttestationRule},
					}
				}
			}

			// browser_impersonation_detection
			if v, ok := botManagementMap["browser_impersonation_detection"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					browserRule := teov20220901.BrowserImpersonationDetectionRule{}
					if v, ok := ruleMap["id"].(string); ok {
						browserRule.Id = helper.String(v)
					}
					if v, ok := ruleMap["name"].(string); ok && v != "" {
						browserRule.Name = helper.String(v)
					}
					if v, ok := ruleMap["enabled"].(string); ok && v != "" {
						browserRule.Enabled = helper.String(v)
					}
					if v, ok := ruleMap["condition"].(string); ok && v != "" {
						browserRule.Condition = helper.String(v)
					}
					if actionMap, ok := helper.ConvertInterfacesHeadToMap(ruleMap["action"]); ok {
						browserAction := teov20220901.BrowserImpersonationDetectionAction{}
						if botSessionMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["bot_session_validation"]); ok {
							botSessionValidation := teov20220901.BotSessionValidation{}
							if v, ok := botSessionMap["issue_new_bot_session_cookie"].(string); ok && v != "" {
								botSessionValidation.IssueNewBotSessionCookie = helper.String(v)
							}
							if triggerMap, ok := helper.ConvertInterfacesHeadToMap(botSessionMap["max_new_session_trigger_config"]); ok {
								triggerConfig := teov20220901.MaxNewSessionTriggerConfig{}
								if v, ok := triggerMap["max_new_session_count_interval"].(string); ok && v != "" {
									triggerConfig.MaxNewSessionCountInterval = helper.String(v)
								}
								if v, ok := triggerMap["max_new_session_count_threshold"].(int); ok {
									triggerConfig.MaxNewSessionCountThreshold = helper.IntInt64(v)
								}
								botSessionValidation.MaxNewSessionTriggerConfig = &triggerConfig
							}
							if sessionExpiredMap, ok := helper.ConvertInterfacesHeadToMap(botSessionMap["session_expired_action"]); ok {
								botSessionValidation.SessionExpiredAction = buildSecurityActionFromMap(sessionExpiredMap)
							}
							if sessionInvalidMap, ok := helper.ConvertInterfacesHeadToMap(botSessionMap["session_invalid_action"]); ok {
								botSessionValidation.SessionInvalidAction = buildSecurityActionFromMap(sessionInvalidMap)
							}
							if sessionRateMap, ok := helper.ConvertInterfacesHeadToMap(botSessionMap["session_rate_control"]); ok {
								sessionRateControl := teov20220901.SessionRateControl{}
								if v, ok := sessionRateMap["enabled"].(string); ok && v != "" {
									sessionRateControl.Enabled = helper.String(v)
								}
								if highRateMap, ok := helper.ConvertInterfacesHeadToMap(sessionRateMap["high_rate_session_action"]); ok {
									sessionRateControl.HighRateSessionAction = buildSecurityActionFromMap(highRateMap)
								}
								if midRateMap, ok := helper.ConvertInterfacesHeadToMap(sessionRateMap["mid_rate_session_action"]); ok {
									sessionRateControl.MidRateSessionAction = buildSecurityActionFromMap(midRateMap)
								}
								if lowRateMap, ok := helper.ConvertInterfacesHeadToMap(sessionRateMap["low_rate_session_action"]); ok {
									sessionRateControl.LowRateSessionAction = buildSecurityActionFromMap(lowRateMap)
								}
								botSessionValidation.SessionRateControl = &sessionRateControl
							}
							browserAction.BotSessionValidation = &botSessionValidation
						}
						if clientBehaviorMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["client_behavior_detection"]); ok {
							clientBehaviorDetection := teov20220901.ClientBehaviorDetection{}
							if v, ok := clientBehaviorMap["crypto_challenge_intensity"].(string); ok && v != "" {
								clientBehaviorDetection.CryptoChallengeIntensity = helper.String(v)
							}
							if v, ok := clientBehaviorMap["crypto_challenge_delay_before"].(string); ok && v != "" {
								clientBehaviorDetection.CryptoChallengeDelayBefore = helper.String(v)
							}
							if v, ok := clientBehaviorMap["max_challenge_count_interval"].(string); ok && v != "" {
								clientBehaviorDetection.MaxChallengeCountInterval = helper.String(v)
							}
							if v, ok := clientBehaviorMap["max_challenge_count_threshold"].(int); ok {
								clientBehaviorDetection.MaxChallengeCountThreshold = helper.IntInt64(v)
							}
							if notFinishedMap, ok := helper.ConvertInterfacesHeadToMap(clientBehaviorMap["challenge_not_finished_action"]); ok {
								clientBehaviorDetection.ChallengeNotFinishedAction = buildSecurityActionFromMap(notFinishedMap)
							}
							if timeoutMap, ok := helper.ConvertInterfacesHeadToMap(clientBehaviorMap["challenge_timeout_action"]); ok {
								clientBehaviorDetection.ChallengeTimeoutAction = buildSecurityActionFromMap(timeoutMap)
							}
							if botClientMap, ok := helper.ConvertInterfacesHeadToMap(clientBehaviorMap["bot_client_action"]); ok {
								clientBehaviorDetection.BotClientAction = buildSecurityActionFromMap(botClientMap)
							}
							browserAction.ClientBehaviorDetection = &clientBehaviorDetection
						}
						browserRule.Action = &browserAction
					}
					botManagement.BrowserImpersonationDetection = &teov20220901.BrowserImpersonationDetection{
						Rules: []*teov20220901.BrowserImpersonationDetectionRule{&browserRule},
					}
				}
			}

			securityPolicy.BotManagement = &botManagement
		}

		if botManagementLiteMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["bot_management_lite"]); ok {
			botManagementLite := teov20220901.BotManagementLite{}
			if captchaPageChallengeMap, ok := helper.InterfaceToMap(botManagementLiteMap, "captcha_page_challenge"); ok {
				captchaPageChallenge := teov20220901.CAPTCHAPageChallenge{}
				if v, ok := captchaPageChallengeMap["enabled"]; ok {
					captchaPageChallenge.Enabled = helper.String(v.(string))
				}

				botManagementLite.CAPTCHAPageChallenge = &captchaPageChallenge
			}

			if aiCrawlerDetectionMap, ok := helper.InterfaceToMap(botManagementLiteMap, "ai_crawler_detection"); ok {
				aiCrawlerDetection := teov20220901.AICrawlerDetection{}
				if v, ok := aiCrawlerDetectionMap["enabled"]; ok {
					aiCrawlerDetection.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(aiCrawlerDetectionMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if allowActionParametersMap, ok := helper.InterfaceToMap(actionMap, "allow_action_parameters"); ok {
						allowActionParameters := teov20220901.AllowActionParameters{}
						if v, ok := allowActionParametersMap["min_delay_time"]; ok {
							allowActionParameters.MinDelayTime = helper.String(v.(string))
						}

						if v, ok := allowActionParametersMap["max_delay_time"]; ok {
							allowActionParameters.MaxDelayTime = helper.String(v.(string))
						}

						securityAction.AllowActionParameters = &allowActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					aiCrawlerDetection.Action = &securityAction
				}

				botManagementLite.AICrawlerDetection = &aiCrawlerDetection
			}

			securityPolicy.BotManagementLite = &botManagementLite
		}

		if ddMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["default_deny_security_action_parameters"]); ok {
			defaultDeny := teov20220901.DefaultDenySecurityActionParameters{}
			if managedRulesMap, ok := helper.ConvertInterfacesHeadToMap(ddMap["managed_rules"]); ok {
				defaultDeny.ManagedRules = buildDenyActionParametersFromMap(managedRulesMap)
			}
			if otherModulesMap, ok := helper.ConvertInterfacesHeadToMap(ddMap["other_modules"]); ok {
				defaultDeny.OtherModules = buildDenyActionParametersFromMap(otherModulesMap)
			}
			securityPolicy.DefaultDenySecurityActionParameters = &defaultDeny
		}

		request.SecurityPolicy = &securityPolicy
	}

	if securityConfigMap, ok := helper.InterfacesHeadMap(d, "security_config"); ok {
		securityConfig := teov20220901.SecurityConfig{}

		if wafConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["waf_config"]); ok {
			wafConfig := teov20220901.WafConfig{}
			if v, ok := wafConfigMap["switch"].(string); ok && v != "" {
				wafConfig.Switch = helper.String(v)
			}
			if v, ok := wafConfigMap["level"].(string); ok && v != "" {
				wafConfig.Level = helper.String(v)
			}
			if v, ok := wafConfigMap["mode"].(string); ok && v != "" {
				wafConfig.Mode = helper.String(v)
			}
			if wafRuleMap, ok := helper.ConvertInterfacesHeadToMap(wafConfigMap["waf_rule"]); ok {
				wafRule := teov20220901.WafRule{}
				if v, ok := wafRuleMap["switch"].(string); ok && v != "" {
					wafRule.Switch = helper.String(v)
				}
				var blockRuleIDs []*int64
				if v, ok := wafRuleMap["block_rule_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							blockRuleIDs = append(blockRuleIDs, helper.IntInt64(item.(int)))
						}
					}
					wafRule.BlockRuleIDs = blockRuleIDs
				}
				var observeRuleIDs []*int64
				if v, ok := wafRuleMap["observe_rule_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							observeRuleIDs = append(observeRuleIDs, helper.IntInt64(item.(int)))
						}
					}
					wafRule.ObserveRuleIDs = observeRuleIDs
				}
				wafConfig.WafRule = &wafRule
			}
			if aiRuleMap, ok := helper.ConvertInterfacesHeadToMap(wafConfigMap["ai_rule"]); ok {
				aiRule := teov20220901.AiRule{}
				if v, ok := aiRuleMap["mode"].(string); ok && v != "" {
					aiRule.Mode = helper.String(v)
				}
				wafConfig.AiRule = &aiRule
			}
			securityConfig.WafConfig = &wafConfig
		}

		if rateLimitConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["rate_limit_config"]); ok {
			rateLimitConfig := teov20220901.RateLimitConfig{}
			if v, ok := rateLimitConfigMap["switch"].(string); ok && v != "" {
				rateLimitConfig.Switch = helper.String(v)
			}
			if v, ok := rateLimitConfigMap["rate_limit_user_rules"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					rateLimitUserRule := teov20220901.RateLimitUserRule{}
					if v, ok := ruleMap["threshold"].(int); ok {
						rateLimitUserRule.Threshold = helper.IntInt64(v)
					}
					if v, ok := ruleMap["period"].(int); ok {
						rateLimitUserRule.Period = helper.IntInt64(v)
					}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						rateLimitUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						rateLimitUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["punish_time"].(int); ok {
						rateLimitUserRule.PunishTime = helper.IntInt64(v)
					}
					if v, ok := ruleMap["punish_time_unit"].(string); ok && v != "" {
						rateLimitUserRule.PunishTimeUnit = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						rateLimitUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["acl_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							rateLimitUserRule.AclConditions = append(rateLimitUserRule.AclConditions, &aclCondition)
						}
					}
					rateLimitConfig.RateLimitUserRules = append(rateLimitConfig.RateLimitUserRules, &rateLimitUserRule)
				}
			}
			if rateLimitTemplateMap, ok := helper.ConvertInterfacesHeadToMap(rateLimitConfigMap["rate_limit_template"]); ok {
				rateLimitTemplate := teov20220901.RateLimitTemplate{}
				if v, ok := rateLimitTemplateMap["mode"].(string); ok && v != "" {
					rateLimitTemplate.Mode = helper.String(v)
				}
				if v, ok := rateLimitTemplateMap["action"].(string); ok && v != "" {
					rateLimitTemplate.Action = helper.String(v)
				}
				rateLimitConfig.RateLimitTemplate = &rateLimitTemplate
			}
			if rateLimitIntelligenceMap, ok := helper.ConvertInterfacesHeadToMap(rateLimitConfigMap["rate_limit_intelligence"]); ok {
				rateLimitIntelligence := teov20220901.RateLimitIntelligence{}
				if v, ok := rateLimitIntelligenceMap["switch"].(string); ok && v != "" {
					rateLimitIntelligence.Switch = helper.String(v)
				}
				if v, ok := rateLimitIntelligenceMap["action"].(string); ok && v != "" {
					rateLimitIntelligence.Action = helper.String(v)
				}
				rateLimitConfig.RateLimitIntelligence = &rateLimitIntelligence
			}
			if v, ok := rateLimitConfigMap["rate_limit_customizes"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					rateLimitUserRule := teov20220901.RateLimitUserRule{}
					if v, ok := ruleMap["threshold"].(int); ok {
						rateLimitUserRule.Threshold = helper.IntInt64(v)
					}
					if v, ok := ruleMap["period"].(int); ok {
						rateLimitUserRule.Period = helper.IntInt64(v)
					}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						rateLimitUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						rateLimitUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["punish_time"].(int); ok {
						rateLimitUserRule.PunishTime = helper.IntInt64(v)
					}
					if v, ok := ruleMap["punish_time_unit"].(string); ok && v != "" {
						rateLimitUserRule.PunishTimeUnit = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						rateLimitUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["acl_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							rateLimitUserRule.AclConditions = append(rateLimitUserRule.AclConditions, &aclCondition)
						}
					}
					rateLimitConfig.RateLimitCustomizes = append(rateLimitConfig.RateLimitCustomizes, &rateLimitUserRule)
				}
			}
			securityConfig.RateLimitConfig = &rateLimitConfig
		}

		if aclConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["acl_config"]); ok {
			aclConfig := teov20220901.AclConfig{}
			if v, ok := aclConfigMap["switch"].(string); ok && v != "" {
				aclConfig.Switch = helper.String(v)
			}
			if v, ok := aclConfigMap["acl_user_rules"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					aclUserRule := teov20220901.AclUserRule{}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						aclUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						aclUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						aclUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["acl_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							aclUserRule.AclConditions = append(aclUserRule.AclConditions, &aclCondition)
						}
					}
					aclConfig.AclUserRules = append(aclConfig.AclUserRules, &aclUserRule)
				}
			}
			if v, ok := aclConfigMap["customizes"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					aclUserRule := teov20220901.AclUserRule{}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						aclUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						aclUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						aclUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["acl_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							aclUserRule.AclConditions = append(aclUserRule.AclConditions, &aclCondition)
						}
					}
					aclConfig.Customizes = append(aclConfig.Customizes, &aclUserRule)
				}
			}
			securityConfig.AclConfig = &aclConfig
		}

		if botConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["bot_config"]); ok {
			botConfig := teov20220901.BotConfig{}
			if v, ok := botConfigMap["switch"].(string); ok && v != "" {
				botConfig.Switch = helper.String(v)
			}
			if botManagedRuleMap, ok := helper.ConvertInterfacesHeadToMap(botConfigMap["bot_managed_rule"]); ok {
				botManagedRule := teov20220901.BotManagedRule{}
				if v, ok := botManagedRuleMap["action"].(string); ok && v != "" {
					botManagedRule.Action = helper.String(v)
				}
				if v, ok := botManagedRuleMap["trans_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botManagedRule.TransManagedIds = append(botManagedRule.TransManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botManagedRuleMap["alg_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botManagedRule.AlgManagedIds = append(botManagedRule.AlgManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botManagedRuleMap["cap_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botManagedRule.CapManagedIds = append(botManagedRule.CapManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botManagedRuleMap["mon_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botManagedRule.MonManagedIds = append(botManagedRule.MonManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botManagedRuleMap["drop_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botManagedRule.DropManagedIds = append(botManagedRule.DropManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				botConfig.BotManagedRule = &botManagedRule
			}
			if botPortraitRuleMap, ok := helper.ConvertInterfacesHeadToMap(botConfigMap["bot_portrait_rule"]); ok {
				botPortraitRule := teov20220901.BotPortraitRule{}
				if v, ok := botPortraitRuleMap["switch"].(string); ok && v != "" {
					botPortraitRule.Switch = helper.String(v)
				}
				if v, ok := botPortraitRuleMap["alg_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botPortraitRule.AlgManagedIds = append(botPortraitRule.AlgManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botPortraitRuleMap["cap_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botPortraitRule.CapManagedIds = append(botPortraitRule.CapManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botPortraitRuleMap["mon_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botPortraitRule.MonManagedIds = append(botPortraitRule.MonManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				if v, ok := botPortraitRuleMap["drop_managed_ids"]; ok {
					for _, item := range v.([]interface{}) {
						if item != nil {
							botPortraitRule.DropManagedIds = append(botPortraitRule.DropManagedIds, helper.IntInt64(item.(int)))
						}
					}
				}
				botConfig.BotPortraitRule = &botPortraitRule
			}
			if intelligenceRuleMap, ok := helper.ConvertInterfacesHeadToMap(botConfigMap["intelligence_rule"]); ok {
				intelligenceRule := teov20220901.IntelligenceRule{}
				if v, ok := intelligenceRuleMap["switch"].(string); ok && v != "" {
					intelligenceRule.Switch = helper.String(v)
				}
				if v, ok := intelligenceRuleMap["intelligence_rule_items"]; ok {
					for _, item := range v.([]interface{}) {
						itemMap := item.(map[string]interface{})
						intelligenceRuleItem := teov20220901.IntelligenceRuleItem{}
						if v, ok := itemMap["label"].(string); ok && v != "" {
							intelligenceRuleItem.Label = helper.String(v)
						}
						if v, ok := itemMap["action"].(string); ok && v != "" {
							intelligenceRuleItem.Action = helper.String(v)
						}
						intelligenceRule.IntelligenceRuleItems = append(intelligenceRule.IntelligenceRuleItems, &intelligenceRuleItem)
					}
				}
				botConfig.IntelligenceRule = &intelligenceRule
			}
			if v, ok := botConfigMap["bot_user_rules"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					botUserRule := teov20220901.BotUserRule{}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						botUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						botUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						botUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["acl_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							botUserRule.AclConditions = append(botUserRule.AclConditions, &aclCondition)
						}
					}
					botConfig.BotUserRules = append(botConfig.BotUserRules, &botUserRule)
				}
			}
			if v, ok := botConfigMap["alg_detect_rule"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					algDetectRule := teov20220901.AlgDetectRule{}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						algDetectRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["switch"].(string); ok && v != "" {
						algDetectRule.Switch = helper.String(v)
					}
					if v, ok := ruleMap["alg_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							algDetectRule.AlgConditions = append(algDetectRule.AlgConditions, &aclCondition)
						}
					}
					if v, ok := ruleMap["alg_detect_session"]; ok {
						if algDetectSessionMap, ok := helper.ConvertInterfacesHeadToMap(v); ok {
							algDetectSession := teov20220901.AlgDetectSession{}
							if v, ok := algDetectSessionMap["name"].(string); ok && v != "" {
								algDetectSession.Name = helper.String(v)
							}
							if v, ok := algDetectSessionMap["detect_mode"].(string); ok && v != "" {
								algDetectSession.DetectMode = helper.String(v)
							}
							if v, ok := algDetectSessionMap["session_analyze_switch"].(string); ok && v != "" {
								algDetectSession.SessionAnalyzeSwitch = helper.String(v)
							}
							if v, ok := algDetectSessionMap["invalid_stat_time"].(int); ok {
								algDetectSession.InvalidStatTime = helper.IntInt64(v)
							}
							if v, ok := algDetectSessionMap["invalid_threshold"].(int); ok {
								algDetectSession.InvalidThreshold = helper.IntInt64(v)
							}
							algDetectRule.AlgDetectSession = &algDetectSession
						}
					}
					if v, ok := ruleMap["alg_detect_js"]; ok {
						for _, jsItem := range v.([]interface{}) {
							jsMap := jsItem.(map[string]interface{})
							algDetectJS := teov20220901.AlgDetectJS{}
							if v, ok := jsMap["name"].(string); ok && v != "" {
								algDetectJS.Name = helper.String(v)
							}
							if v, ok := jsMap["work_level"].(string); ok && v != "" {
								algDetectJS.WorkLevel = helper.String(v)
							}
							if v, ok := jsMap["execute_mode"].(int); ok {
								algDetectJS.ExecuteMode = helper.IntInt64(v)
							}
							if v, ok := jsMap["invalid_stat_time"].(int); ok {
								algDetectJS.InvalidStatTime = helper.IntInt64(v)
							}
							if v, ok := jsMap["invalid_threshold"].(int); ok {
								algDetectJS.InvalidThreshold = helper.IntInt64(v)
							}
							algDetectRule.AlgDetectJS = append(algDetectRule.AlgDetectJS, &algDetectJS)
						}
					}
					botConfig.AlgDetectRule = append(botConfig.AlgDetectRule, &algDetectRule)
				}
			}
			if v, ok := botConfigMap["customizes"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					botUserRule := teov20220901.BotUserRule{}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						botUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						botUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						botUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["acl_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							aclCondition := teov20220901.AclCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								aclCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								aclCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								aclCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								aclCondition.MatchContent = helper.String(v)
							}
							botUserRule.AclConditions = append(botUserRule.AclConditions, &aclCondition)
						}
					}
					botConfig.Customizes = append(botConfig.Customizes, &botUserRule)
				}
			}
			securityConfig.BotConfig = &botConfig
		}

		if switchConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["switch_config"]); ok {
			switchConfig := teov20220901.SwitchConfig{}
			if v, ok := switchConfigMap["web_switch"].(string); ok && v != "" {
				switchConfig.WebSwitch = helper.String(v)
			}
			securityConfig.SwitchConfig = &switchConfig
		}

		if ipTableConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["ip_table_config"]); ok {
			ipTableConfig := teov20220901.IpTableConfig{}
			if v, ok := ipTableConfigMap["switch"].(string); ok && v != "" {
				ipTableConfig.Switch = helper.String(v)
			}
			if v, ok := ipTableConfigMap["ip_table_rules"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					ipTableRule := teov20220901.IpTableRule{}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						ipTableRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["match_from"].(string); ok && v != "" {
						ipTableRule.MatchFrom = helper.String(v)
					}
					if v, ok := ruleMap["operator"].(string); ok && v != "" {
						ipTableRule.Operator = helper.String(v)
					}
					if v, ok := ruleMap["match_content"].(string); ok && v != "" {
						ipTableRule.MatchContent = helper.String(v)
					}
					if v, ok := ruleMap["rule_id"].(int); ok {
						ipTableRule.RuleID = helper.IntInt64(v)
					}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						ipTableRule.RuleName = helper.String(v)
					}
					ipTableConfig.IpTableRules = append(ipTableConfig.IpTableRules, &ipTableRule)
				}
			}
			securityConfig.IpTableConfig = &ipTableConfig
		}

		if exceptConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["except_config"]); ok {
			exceptConfig := teov20220901.ExceptConfig{}
			if v, ok := exceptConfigMap["switch"].(string); ok && v != "" {
				exceptConfig.Switch = helper.String(v)
			}
			if v, ok := exceptConfigMap["except_user_rules"]; ok {
				for _, item := range v.([]interface{}) {
					ruleMap := item.(map[string]interface{})
					exceptUserRule := teov20220901.ExceptUserRule{}
					if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
						exceptUserRule.RuleName = helper.String(v)
					}
					if v, ok := ruleMap["action"].(string); ok && v != "" {
						exceptUserRule.Action = helper.String(v)
					}
					if v, ok := ruleMap["rule_status"].(string); ok && v != "" {
						exceptUserRule.RuleStatus = helper.String(v)
					}
					if v, ok := ruleMap["except_user_rule_conditions"]; ok {
						for _, condItem := range v.([]interface{}) {
							condMap := condItem.(map[string]interface{})
							exceptUserRuleCondition := teov20220901.ExceptUserRuleCondition{}
							if v, ok := condMap["match_from"].(string); ok && v != "" {
								exceptUserRuleCondition.MatchFrom = helper.String(v)
							}
							if v, ok := condMap["match_param"].(string); ok && v != "" {
								exceptUserRuleCondition.MatchParam = helper.String(v)
							}
							if v, ok := condMap["operator"].(string); ok && v != "" {
								exceptUserRuleCondition.Operator = helper.String(v)
							}
							if v, ok := condMap["match_content"].(string); ok && v != "" {
								exceptUserRuleCondition.MatchContent = helper.String(v)
							}
							exceptUserRule.ExceptUserRuleConditions = append(exceptUserRule.ExceptUserRuleConditions, &exceptUserRuleCondition)
						}
					}
					exceptConfig.ExceptUserRules = append(exceptConfig.ExceptUserRules, &exceptUserRule)
				}
			}
			securityConfig.ExceptConfig = &exceptConfig
		}

		if dropPageConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["drop_page_config"]); ok {
			dropPageConfig := teov20220901.DropPageConfig{}
			if v, ok := dropPageConfigMap["switch"].(string); ok && v != "" {
				dropPageConfig.Switch = helper.String(v)
			}
			if wafDropPageDetailMap, ok := helper.ConvertInterfacesHeadToMap(dropPageConfigMap["waf_drop_page_detail"]); ok {
				wafDropPageDetail := teov20220901.DropPageDetail{}
				if v, ok := wafDropPageDetailMap["page_id"].(int); ok {
					wafDropPageDetail.PageId = helper.IntInt64(v)
				}
				if v, ok := wafDropPageDetailMap["status_code"].(int); ok {
					wafDropPageDetail.StatusCode = helper.IntInt64(v)
				}
				if v, ok := wafDropPageDetailMap["name"].(string); ok && v != "" {
					wafDropPageDetail.Name = helper.String(v)
				}
				if v, ok := wafDropPageDetailMap["type"].(string); ok && v != "" {
					wafDropPageDetail.Type = helper.String(v)
				}
				if v, ok := wafDropPageDetailMap["custom_response_id"].(string); ok && v != "" {
					wafDropPageDetail.CustomResponseId = helper.String(v)
				}
				dropPageConfig.WafDropPageDetail = &wafDropPageDetail
			}
			if aclDropPageDetailMap, ok := helper.ConvertInterfacesHeadToMap(dropPageConfigMap["acl_drop_page_detail"]); ok {
				aclDropPageDetail := teov20220901.DropPageDetail{}
				if v, ok := aclDropPageDetailMap["page_id"].(int); ok {
					aclDropPageDetail.PageId = helper.IntInt64(v)
				}
				if v, ok := aclDropPageDetailMap["status_code"].(int); ok {
					aclDropPageDetail.StatusCode = helper.IntInt64(v)
				}
				if v, ok := aclDropPageDetailMap["name"].(string); ok && v != "" {
					aclDropPageDetail.Name = helper.String(v)
				}
				if v, ok := aclDropPageDetailMap["type"].(string); ok && v != "" {
					aclDropPageDetail.Type = helper.String(v)
				}
				if v, ok := aclDropPageDetailMap["custom_response_id"].(string); ok && v != "" {
					aclDropPageDetail.CustomResponseId = helper.String(v)
				}
				dropPageConfig.AclDropPageDetail = &aclDropPageDetail
			}
			securityConfig.DropPageConfig = &dropPageConfig
		}

		if slowPostConfigMap, ok := helper.ConvertInterfacesHeadToMap(securityConfigMap["slow_post_config"]); ok {
			slowPostConfig := teov20220901.SlowPostConfig{}
			if v, ok := slowPostConfigMap["switch"].(string); ok && v != "" {
				slowPostConfig.Switch = helper.String(v)
			}
			if v, ok := slowPostConfigMap["action"].(string); ok && v != "" {
				slowPostConfig.Action = helper.String(v)
			}
			if firstPartConfigMap, ok := helper.ConvertInterfacesHeadToMap(slowPostConfigMap["first_part_config"]); ok {
				firstPartConfig := teov20220901.FirstPartConfig{}
				if v, ok := firstPartConfigMap["switch"].(string); ok && v != "" {
					firstPartConfig.Switch = helper.String(v)
				}
				if v, ok := firstPartConfigMap["stat_time"].(int); ok {
					firstPartConfig.StatTime = helper.IntUint64(v)
				}
				slowPostConfig.FirstPartConfig = &firstPartConfig
			}
			if slowRateConfigMap, ok := helper.ConvertInterfacesHeadToMap(slowPostConfigMap["slow_rate_config"]); ok {
				slowRateConfig := teov20220901.SlowRateConfig{}
				if v, ok := slowRateConfigMap["switch"].(string); ok && v != "" {
					slowRateConfig.Switch = helper.String(v)
				}
				if v, ok := slowRateConfigMap["interval"].(int); ok {
					slowRateConfig.Interval = helper.IntUint64(v)
				}
				if v, ok := slowRateConfigMap["threshold"].(int); ok {
					slowRateConfig.Threshold = helper.IntUint64(v)
				}
				slowPostConfig.SlowRateConfig = &slowRateConfig
			}
			securityConfig.SlowPostConfig = &slowPostConfig
		}

		request.SecurityConfig = &securityConfig
	} else {
		request.SecurityConfig = &teov20220901.SecurityConfig{}
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

// defaultDenyActionParametersSchema returns the schema resource used by
// `default_deny_security_action_parameters.managed_rules` and `.other_modules`.
// Mirrors the SDK `DenyActionParameters` struct.
func defaultDenyActionParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to extend the source IP block. Valid values: `on`, `off`.",
			},
			"block_ip_duration": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP block duration when `block_ip` is `on`.",
			},
			"return_custom_page": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to use a custom page. Valid values: `on`, `off`.",
			},
			"response_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Status code of the custom page.",
			},
			"error_page_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "PageId of the custom page.",
			},
			"stall": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to suspend the request source. Valid values: `on`, `off`.",
			},
		},
	}
}

// buildDenyActionParametersFromMap converts a schema map to *teov20220901.DenyActionParameters.
func buildDenyActionParametersFromMap(m map[string]interface{}) *teov20220901.DenyActionParameters {
	p := &teov20220901.DenyActionParameters{}
	if v, ok := m["block_ip"].(string); ok && v != "" {
		p.BlockIp = helper.String(v)
	}
	if v, ok := m["block_ip_duration"].(string); ok && v != "" {
		p.BlockIpDuration = helper.String(v)
	}
	if v, ok := m["return_custom_page"].(string); ok && v != "" {
		p.ReturnCustomPage = helper.String(v)
	}
	if v, ok := m["response_code"].(string); ok && v != "" {
		p.ResponseCode = helper.String(v)
	}
	if v, ok := m["error_page_id"].(string); ok && v != "" {
		p.ErrorPageId = helper.String(v)
	}
	if v, ok := m["stall"].(string); ok && v != "" {
		p.Stall = helper.String(v)
	}
	return p
}

// flattenDenyActionParameters converts *teov20220901.DenyActionParameters to a schema map.
func flattenDenyActionParameters(p *teov20220901.DenyActionParameters) map[string]interface{} {
	m := map[string]interface{}{}
	if p == nil {
		return m
	}
	if p.BlockIp != nil {
		m["block_ip"] = p.BlockIp
	}
	if p.BlockIpDuration != nil {
		m["block_ip_duration"] = p.BlockIpDuration
	}
	if p.ReturnCustomPage != nil {
		m["return_custom_page"] = p.ReturnCustomPage
	}
	if p.ResponseCode != nil {
		m["response_code"] = p.ResponseCode
	}
	if p.ErrorPageId != nil {
		m["error_page_id"] = p.ErrorPageId
	}
	if p.Stall != nil {
		m["stall"] = p.Stall
	}
	return m
}

// aclConditionSchema returns the schema for AclCondition used in rate_limit_user_rules/
// acl_user_rules/bot_user_rules/alg_detect_rule.acl_conditions.
func aclConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_from":    {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`)."},
			"match_param":   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match parameter. For `header` MatchFrom, the header key."},
			"operator":      {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`)."},
			"match_content": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match content."},
		},
	}
}

// rateLimitUserRuleSchema returns the schema for RateLimitUserRule.
func rateLimitUserRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"threshold":          {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Rate limit threshold in count. Range 0-4294967294."},
			"period":             {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Rate limit statistical period in seconds (10/20/30/40/50/60)."},
			"rule_name":          {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule name."},
			"action":             {Type: schema.TypeString, Optional: true, Computed: true, Description: "Action. Valid values: `monitor`, `drop`, `redirect`, `page`, `alg`."},
			"punish_time":        {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Penalty duration (0-2 days)."},
			"punish_time_unit":   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Penalty duration unit. Valid values: `second`, `minutes`, `hour`."},
			"rule_status":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule status. Valid values: `on`, `off`."},
			"rule_priority":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Rule priority (0-100)."},
			"rule_id":            {Type: schema.TypeInt, Computed: true, Description: "Rule ID. Output-only."},
			"freq_fields":        {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Filter fields. Valid values: `sip`."},
			"update_time":        {Type: schema.TypeString, Computed: true, Description: "Update time. Output-only."},
			"freq_scope":         {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Statistical scope. Valid values: `source_to_eo`, `client_to_eo`."},
			"name":               {Type: schema.TypeString, Optional: true, Computed: true, Description: "Custom response page name. Required when Action is `page`."},
			"custom_response_id": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Custom response ID."},
			"response_code":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Custom response code (100-600, excl. 3xx)."},
			"redirect_url":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "Redirect URL. Required when Action is `redirect`."},
			"acl_conditions": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Rule ACL conditions.",
				Elem:        aclConditionSchema(),
			},
		},
	}
}

// aclUserRuleSchema returns the schema for AclUserRule.
func aclUserRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_name":          {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule name."},
			"action":             {Type: schema.TypeString, Optional: true, Computed: true, Description: "Action. Valid values: `trans`, `drop`, `monitor`, `ban`, `redirect`, `page`, `alg`."},
			"rule_status":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule status. Valid values: `on`, `off`."},
			"rule_priority":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Rule priority (0-100)."},
			"rule_id":            {Type: schema.TypeInt, Computed: true, Description: "Rule ID. Output-only."},
			"update_time":        {Type: schema.TypeString, Computed: true, Description: "Update time. Output-only."},
			"punish_time":        {Type: schema.TypeInt, Optional: true, Computed: true, Description: "IP ban penalty time."},
			"punish_time_unit":   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Penalty time unit. Valid values: `second`, `minutes`, `hour`."},
			"name":               {Type: schema.TypeString, Optional: true, Computed: true, Description: "Custom response page name."},
			"page_id":            {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Custom page instance ID. Deprecated."},
			"custom_response_id": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Custom response ID."},
			"response_code":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Custom response code."},
			"redirect_url":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "Redirect URL."},
			"acl_conditions": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Rule ACL conditions.",
				Elem:        aclConditionSchema(),
			},
		},
	}
}

// botUserRuleSchema returns the schema for BotUserRule.
func botUserRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_name":          {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule name."},
			"action":             {Type: schema.TypeString, Optional: true, Computed: true, Description: "Action. Valid values: `drop`, `monitor`, `trans`, `redirect`, `page`, `alg`, `captcha`, `random`, `silence`, `shortdelay`, `longdelay`."},
			"rule_status":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule status. Valid values: `on`, `off`."},
			"rule_priority":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Rule priority (0-100)."},
			"rule_id":            {Type: schema.TypeInt, Computed: true, Description: "Rule ID. Output-only."},
			"update_time":        {Type: schema.TypeString, Computed: true, Description: "Update time. Output-only."},
			"freq_fields":        {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Filter fields."},
			"freq_scope":         {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Statistical scope."},
			"name":               {Type: schema.TypeString, Optional: true, Computed: true, Description: "Custom response page name."},
			"custom_response_id": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Custom response ID."},
			"response_code":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Custom response code."},
			"redirect_url":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "Redirect URL."},
			"acl_conditions": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Rule ACL conditions.",
				Elem:        aclConditionSchema(),
			},
			"extend_actions": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Random action weighted distribution.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"action":  {Type: schema.TypeString, Optional: true, Computed: true, Description: "Action. Valid values: `monitor`, `alg`, `captcha`, `random`, `silence`, `shortdelay`, `longdelay`."},
					"percent": {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Action probability (0-100)."},
				}},
			},
		},
	}
}

// ipTableRuleSchema returns the schema for IpTableRule.
func ipTableRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Action. Valid values: `drop`, `trans`, `monitor`."},
			"match_from":    {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match field. Valid values: `ip`, `area`, `asn`, `referer`, `ua`, `url`."},
			"operator":      {Type: schema.TypeString, Optional: true, Computed: true, Description: "Operator. Valid values include `match`, `not_match`, `include_area`, `not_include_area`, `asn_match`, `asn_not_match`, `equal`, `not_equal`, `include`, `not_include`, `is_emty`, `not_exists`."},
			"rule_id":       {Type: schema.TypeInt, Computed: true, Description: "Rule ID. Output-only."},
			"update_time":   {Type: schema.TypeString, Computed: true, Description: "Update time. Output-only."},
			"status":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule status. Valid values: `on`, `off`. Default: `on`."},
			"rule_name":     {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule name."},
			"match_content": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match content. Comma-separated for multi-values."},
		},
	}
}

// exceptUserRuleSchema returns the schema for ExceptUserRule.
func exceptUserRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_name":     {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule name (no Chinese characters)."},
			"action":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule action. Only `skip` is supported."},
			"rule_status":   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule status. Valid values: `on`, `off`."},
			"rule_id":       {Type: schema.TypeInt, Computed: true, Description: "Rule ID. Output-only."},
			"update_time":   {Type: schema.TypeString, Computed: true, Description: "Update time."},
			"rule_priority": {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Priority (0-100). Default 0."},
			"except_user_rule_conditions": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Match conditions.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"match_from":    {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match field."},
					"match_param":   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match parameter (e.g. header key when MatchFrom=header)."},
					"operator":      {Type: schema.TypeString, Optional: true, Computed: true, Description: "Operator."},
					"match_content": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Match value."},
				}},
			},
			"except_user_rule_scope": {
				Type: schema.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Description: "Rule effective scope.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"type":    {Type: schema.TypeString, Optional: true, Computed: true, Description: "Scope type. Valid values: `complete`, `partial`."},
					"modules": {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Effective modules. Valid values: `waf`, `rate`, `acl`, `cc`, `bot`."},
					"partial_modules": {
						Type: schema.TypeList, Optional: true, Computed: true,
						Description: "Partial rule ID exceptions.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"module":  {Type: schema.TypeString, Optional: true, Computed: true, Description: "Module name. Valid values: `managed-rule`, `managed-group`, `waf` (deprecated)."},
							"include": {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeInt}, Description: "Rule IDs to include."},
						}},
					},
					"skip_conditions": {
						Type: schema.TypeList, Optional: true, Computed: true,
						Description: "Conditions to skip.",
						Elem: &schema.Resource{Schema: map[string]*schema.Schema{
							"type":               {Type: schema.TypeString, Optional: true, Computed: true, Description: "Skip type. Valid values: `header_fields`, `cookie`, `query_string`, `uri`, `body_raw`, `body_json`."},
							"selector":           {Type: schema.TypeString, Optional: true, Computed: true, Description: "Selector. Valid values: `args`, `path`, `full`, `upload_filename`, `keys`, `values`, `key_value`."},
							"match_from_type":    {Type: schema.TypeString, Optional: true, Computed: true, Description: "Key match type. Valid values: `equal`, `wildcard`."},
							"match_from":         {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Key values."},
							"match_content_type": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Content match type. Valid values: `equal`, `wildcard`."},
							"match_content":      {Type: schema.TypeList, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Content values."},
						}},
					},
				}},
			},
		},
	}
}

// algDetectRuleSchema returns the schema for AlgDetectRule.
func algDetectRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id":     {Type: schema.TypeInt, Computed: true, Description: "Rule ID. Output-only."},
			"rule_name":   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule name."},
			"switch":      {Type: schema.TypeString, Optional: true, Computed: true, Description: "Rule switch."},
			"update_time": {Type: schema.TypeString, Computed: true, Description: "Update time."},
			"alg_conditions": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Custom conditions.",
				Elem:        aclConditionSchema(),
			},
			"alg_detect_session": {
				Type: schema.TypeList, Optional: true, Computed: true, MaxItems: 1,
				Description: "Cookie validation and session behavior analysis.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name":                   {Type: schema.TypeString, Optional: true, Computed: true, Description: "Operation name."},
					"detect_mode":            {Type: schema.TypeString, Optional: true, Computed: true, Description: "Detection mode. Valid values: `detect`, `update_detect`."},
					"session_analyze_switch": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Session behavior analysis switch. Valid values: `off`, `on`."},
					"invalid_stat_time":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Statistical period for missing/expired cookie (5-3600s, default 10)."},
					"invalid_threshold":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Trigger threshold for missing/expired cookie (1-100000000, default 300)."},
				}},
			},
			"alg_detect_js": {
				Type: schema.TypeList, Optional: true, Computed: true,
				Description: "Client behavior detection.",
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name":              {Type: schema.TypeString, Optional: true, Computed: true, Description: "Operation name."},
					"work_level":        {Type: schema.TypeString, Optional: true, Computed: true, Description: "Proof-of-work strength. Valid values: `low`, `middle`, `high` (default `low`)."},
					"execute_mode":      {Type: schema.TypeInt, Optional: true, Computed: true, Description: "JS execution delay in ms (0-1000, default 500)."},
					"invalid_stat_time": {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Statistical period for invalid JS (5-3600s, default 10)."},
					"invalid_threshold": {Type: schema.TypeInt, Optional: true, Computed: true, Description: "Threshold for invalid JS (1-100000000, default 300)."},
				}},
			},
		},
	}
}

// detectLengthLimitRuleSchema returns the schema for DetectLengthLimitRule used in detect_length_limit_config.
func detectLengthLimitRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule ID. Output-only.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule name. Output-only.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule description. Output-only.",
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Rule conditions. Output-only.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Condition name. Valid values: `body_depth`. Output-only.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Condition values. Valid values: `10KB`, `64KB`, `128KB`. Output-only.",
						},
					},
				},
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action. Valid values: `skip`, `scan`. Output-only.",
			},
		},
	}
}

// dropPageDetailSchema returns the schema for DropPageDetail used in drop_page_config.
func dropPageDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"page_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique ID of the block page. The system includes a built-in block page with ID 0.",
			},
			"status_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTP status code for the block page. Range: 100-600, excluding 3xx.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Block page file name or URL.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Page type. Valid values: `page`.",
			},
			"custom_response_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Custom response ID.",
			},
		},
	}
}

// securityActionSchema returns the schema for SecurityAction used in bot_management.
func securityActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.",
			},
			"deny_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional parameters when Name is Deny.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to extend the blocking of source IP. Valid values: `on`, `off`.",
						},
						"block_ip_duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP blocking duration when BlockIP is on.",
						},
						"return_custom_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to use custom pages. Valid values: `on`, `off`.",
						},
						"response_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Customize the status code of the page.",
						},
						"error_page_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The PageId of the custom page.",
						},
						"stall": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to ignore the request source suspension. Valid values: `on`, `off`.",
						},
					},
				},
			},
			"redirect_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional parameters when Name is Redirect.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL to redirect.",
						},
					},
				},
			},
			"allow_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional parameters when Name is Allow.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_delay_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Minimum delay response time. Supported unit: seconds, range 0-5.",
						},
						"max_delay_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum delay response time. Supported unit: seconds, range 5-10.",
						},
					},
				},
			},
			"challenge_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional parameters when Name is Challenge.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"challenge_option": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The time interval for repeating the challenge.",
						},
						"attester_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client authentication method ID.",
						},
					},
				},
			},
		},
	}
}

// botManagementActionOverrideSchema returns the schema for BotManagementActionOverride.
func botManagementActionOverrideSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule ID or category ID for action override.",
			},
			"action": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Action override configuration.",
				Elem:        securityActionSchema(),
			},
		},
	}
}

// buildSecurityActionFromMap converts a schema map to *teov20220901.SecurityAction.
func buildSecurityActionFromMap(m map[string]interface{}) *teov20220901.SecurityAction {
	if len(m) == 0 {
		return nil
	}
	action := &teov20220901.SecurityAction{}
	if v, ok := m["name"].(string); ok && v != "" {
		action.Name = helper.String(v)
	}
	if denyMap, ok := helper.ConvertInterfacesHeadToMap(m["deny_action_parameters"]); ok && len(denyMap) > 0 {
		denyParams := &teov20220901.DenyActionParameters{}
		if v, ok := denyMap["block_ip"].(string); ok && v != "" {
			denyParams.BlockIp = helper.String(v)
		}
		if v, ok := denyMap["block_ip_duration"].(string); ok && v != "" {
			denyParams.BlockIpDuration = helper.String(v)
		}
		if v, ok := denyMap["return_custom_page"].(string); ok && v != "" {
			denyParams.ReturnCustomPage = helper.String(v)
		}
		if v, ok := denyMap["response_code"].(string); ok && v != "" {
			denyParams.ResponseCode = helper.String(v)
		}
		if v, ok := denyMap["error_page_id"].(string); ok && v != "" {
			denyParams.ErrorPageId = helper.String(v)
		}
		if v, ok := denyMap["stall"].(string); ok && v != "" {
			denyParams.Stall = helper.String(v)
		}
		action.DenyActionParameters = denyParams
	}
	if redirectMap, ok := helper.ConvertInterfacesHeadToMap(m["redirect_action_parameters"]); ok && len(redirectMap) > 0 {
		redirectParams := &teov20220901.RedirectActionParameters{}
		if v, ok := redirectMap["url"].(string); ok && v != "" {
			redirectParams.URL = helper.String(v)
		}
		action.RedirectActionParameters = redirectParams
	}
	if allowMap, ok := helper.ConvertInterfacesHeadToMap(m["allow_action_parameters"]); ok && len(allowMap) > 0 {
		allowParams := &teov20220901.AllowActionParameters{}
		if v, ok := allowMap["min_delay_time"].(string); ok && v != "" {
			allowParams.MinDelayTime = helper.String(v)
		}
		if v, ok := allowMap["max_delay_time"].(string); ok && v != "" {
			allowParams.MaxDelayTime = helper.String(v)
		}
		action.AllowActionParameters = allowParams
	}
	if challengeMap, ok := helper.ConvertInterfacesHeadToMap(m["challenge_action_parameters"]); ok && len(challengeMap) > 0 {
		challengeParams := &teov20220901.ChallengeActionParameters{}
		if v, ok := challengeMap["challenge_option"].(string); ok && v != "" {
			challengeParams.ChallengeOption = helper.String(v)
		}
		if v, ok := challengeMap["interval"].(string); ok && v != "" {
			challengeParams.Interval = helper.String(v)
		}
		if v, ok := challengeMap["attester_id"].(string); ok && v != "" {
			challengeParams.AttesterId = helper.String(v)
		}
		action.ChallengeActionParameters = challengeParams
	}
	return action
}

// buildBotManagementActionOverrideFromMap converts a schema map to *teov20220901.BotManagementActionOverrides.
func buildBotManagementActionOverrideFromMap(m map[string]interface{}) *teov20220901.BotManagementActionOverrides {
	if len(m) == 0 {
		return nil
	}
	override := &teov20220901.BotManagementActionOverrides{}
	if v, ok := m["rule_id"].(string); ok && v != "" {
		override.Ids = []*string{helper.String(v)}
	}
	if actionMap, ok := helper.ConvertInterfacesHeadToMap(m["action"]); ok && len(actionMap) > 0 {
		override.Action = buildSecurityActionFromMap(actionMap)
	}
	return override
}

// flattenSecurityAction converts *teov20220901.SecurityAction to a schema map.
func flattenSecurityAction(action *teov20220901.SecurityAction) map[string]interface{} {
	m := map[string]interface{}{}
	if action == nil {
		return m
	}
	if action.Name != nil {
		m["name"] = action.Name
	}
	if action.DenyActionParameters != nil {
		denyMap := map[string]interface{}{}
		if action.DenyActionParameters.BlockIp != nil {
			denyMap["block_ip"] = action.DenyActionParameters.BlockIp
		}
		if action.DenyActionParameters.BlockIpDuration != nil {
			denyMap["block_ip_duration"] = action.DenyActionParameters.BlockIpDuration
		}
		if action.DenyActionParameters.ReturnCustomPage != nil {
			denyMap["return_custom_page"] = action.DenyActionParameters.ReturnCustomPage
		}
		if action.DenyActionParameters.ResponseCode != nil {
			denyMap["response_code"] = action.DenyActionParameters.ResponseCode
		}
		if action.DenyActionParameters.ErrorPageId != nil {
			denyMap["error_page_id"] = action.DenyActionParameters.ErrorPageId
		}
		if action.DenyActionParameters.Stall != nil {
			denyMap["stall"] = action.DenyActionParameters.Stall
		}
		if len(denyMap) > 0 {
			m["deny_action_parameters"] = []interface{}{denyMap}
		}
	}
	if action.RedirectActionParameters != nil {
		redirectMap := map[string]interface{}{}
		if action.RedirectActionParameters.URL != nil {
			redirectMap["url"] = action.RedirectActionParameters.URL
		}
		if len(redirectMap) > 0 {
			m["redirect_action_parameters"] = []interface{}{redirectMap}
		}
	}
	if action.AllowActionParameters != nil {
		allowMap := map[string]interface{}{}
		if action.AllowActionParameters.MinDelayTime != nil {
			allowMap["min_delay_time"] = action.AllowActionParameters.MinDelayTime
		}
		if action.AllowActionParameters.MaxDelayTime != nil {
			allowMap["max_delay_time"] = action.AllowActionParameters.MaxDelayTime
		}
		if len(allowMap) > 0 {
			m["allow_action_parameters"] = []interface{}{allowMap}
		}
	}
	if action.ChallengeActionParameters != nil {
		challengeMap := map[string]interface{}{}
		if action.ChallengeActionParameters.ChallengeOption != nil {
			challengeMap["challenge_option"] = action.ChallengeActionParameters.ChallengeOption
		}
		if action.ChallengeActionParameters.Interval != nil {
			challengeMap["interval"] = action.ChallengeActionParameters.Interval
		}
		if action.ChallengeActionParameters.AttesterId != nil {
			challengeMap["attester_id"] = action.ChallengeActionParameters.AttesterId
		}
		if len(challengeMap) > 0 {
			m["challenge_action_parameters"] = []interface{}{challengeMap}
		}
	}
	return m
}

// flattenBotManagementActionOverride converts *teov20220901.BotManagementActionOverrides to a schema map.
func flattenBotManagementActionOverride(override *teov20220901.BotManagementActionOverrides) map[string]interface{} {
	m := map[string]interface{}{}
	if override == nil {
		return m
	}
	if len(override.Ids) > 0 && override.Ids[0] != nil {
		m["rule_id"] = override.Ids[0]
	}
	if override.Action != nil {
		m["action"] = []interface{}{flattenSecurityAction(override.Action)}
	}
	return m
}
