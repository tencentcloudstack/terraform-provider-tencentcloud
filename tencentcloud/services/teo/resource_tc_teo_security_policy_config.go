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
	request.SecurityConfig = &teov20220901.SecurityConfig{
		RateLimitConfig: &teov20220901.RateLimitConfig{
			RateLimitUserRules: []*teov20220901.RateLimitUserRule{},
			Switch:             helper.String("on"),
		},
	}

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
