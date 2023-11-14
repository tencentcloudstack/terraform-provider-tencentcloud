/*
Provides a resource to create a teo security_policy

Example Usage

```hcl
resource "tencentcloud_teo_security_policy" "security_policy" {
  zone_id = &lt;nil&gt;
  entity = &lt;nil&gt;
  config {
		waf_config {
			switch = &lt;nil&gt;
			level = &lt;nil&gt;
			mode = &lt;nil&gt;
			waf_rules {
				switch = &lt;nil&gt;
				block_rule_i_ds = &lt;nil&gt;
				observe_rule_i_ds = &lt;nil&gt;
			}
			ai_rule {
				mode = &lt;nil&gt;
			}
		}
		rate_limit_config {
			switch = &lt;nil&gt;
			user_rules {
				rule_name = &lt;nil&gt;
				threshold = &lt;nil&gt;
				period = &lt;nil&gt;
				action = &lt;nil&gt;
				punish_time = &lt;nil&gt;
				punish_time_unit = &lt;nil&gt;
				rule_status = &lt;nil&gt;
				freq_fields = &lt;nil&gt;
				conditions {
					match_from = &lt;nil&gt;
					match_param = &lt;nil&gt;
					operator = &lt;nil&gt;
					match_content = &lt;nil&gt;
				}
				rule_priority = &lt;nil&gt;
			}
			template {
				mode = &lt;nil&gt;
				detail {
					mode = &lt;nil&gt;
					i_d = &lt;nil&gt;
					action = &lt;nil&gt;
					punish_time = &lt;nil&gt;
					threshold = &lt;nil&gt;
					period = &lt;nil&gt;
				}
			}
			intelligence {
				switch = &lt;nil&gt;
				action = &lt;nil&gt;
			}
		}
		acl_config {
			switch = &lt;nil&gt;
			user_rules {
				rule_name = &lt;nil&gt;
				action = &lt;nil&gt;
				rule_status = &lt;nil&gt;
				conditions {
					match_from = &lt;nil&gt;
					match_param = &lt;nil&gt;
					operator = &lt;nil&gt;
					match_content = &lt;nil&gt;
				}
				rule_priority = &lt;nil&gt;
				punish_time = &lt;nil&gt;
				punish_time_unit = &lt;nil&gt;
				name = &lt;nil&gt;
				page_id = &lt;nil&gt;
				redirect_url = &lt;nil&gt;
				response_code = &lt;nil&gt;
			}
		}
		bot_config {
			switch = &lt;nil&gt;
			managed_rule {
				action = &lt;nil&gt;
				punish_time = &lt;nil&gt;
				punish_time_unit = &lt;nil&gt;
				name = &lt;nil&gt;
				page_id = &lt;nil&gt;
				redirect_url = &lt;nil&gt;
				response_code = &lt;nil&gt;
				trans_managed_ids = &lt;nil&gt;
				alg_managed_ids = &lt;nil&gt;
				cap_managed_ids = &lt;nil&gt;
				mon_managed_ids = &lt;nil&gt;
				drop_managed_ids = &lt;nil&gt;
			}
			portrait_rule {
				alg_managed_ids = &lt;nil&gt;
				cap_managed_ids = &lt;nil&gt;
				mon_managed_ids = &lt;nil&gt;
				drop_managed_ids = &lt;nil&gt;
				switch = &lt;nil&gt;
			}
			intelligence_rule {
				switch = &lt;nil&gt;
				items {
					label = &lt;nil&gt;
					action = &lt;nil&gt;
				}
			}
		}
		switch_config {
			web_switch = &lt;nil&gt;
		}
		ip_table_config {
			switch = &lt;nil&gt;
			rules {
				action = &lt;nil&gt;
				match_from = &lt;nil&gt;
				match_content = &lt;nil&gt;
			}
		}
		except_config {
			switch = &lt;nil&gt;
			except_user_rules {
				action = &lt;nil&gt;
				rule_status = &lt;nil&gt;
				rule_priority = &lt;nil&gt;
				except_user_rule_conditions {
					match_from = &lt;nil&gt;
					match_param = &lt;nil&gt;
					operator = &lt;nil&gt;
					match_content = &lt;nil&gt;
				}
				except_user_rule_scope {
					modules = &lt;nil&gt;
				}
			}
		}
		drop_page_config {
			switch = &lt;nil&gt;
			waf_drop_page_detail {
				page_id = &lt;nil&gt;
				status_code = &lt;nil&gt;
				name = &lt;nil&gt;
				type = &lt;nil&gt;
			}
			acl_drop_page_detail {
				page_id = &lt;nil&gt;
				status_code = &lt;nil&gt;
				name = &lt;nil&gt;
				type = &lt;nil&gt;
			}
		}

  }
}
```

Import

teo security_policy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_policy.security_policy security_policy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTeoSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityPolicyCreate,
		Read:   resourceTencentCloudTeoSecurityPolicyRead,
		Update: resourceTencentCloudTeoSecurityPolicyUpdate,
		Delete: resourceTencentCloudTeoSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"entity": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subdomain.",
			},

			"config": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Security policy configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"waf_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "WAF (Web Application Firewall) Configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to enable WAF rules. Valid values:- `on`: Enable.- `off`: Disable.",
									},
									"level": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protection level. Valid values: `loose`, `normal`, `strict`, `stricter`, `custom`.",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protection mode. Valid values:- `block`: use block mode globally, you still can set a group of rules to use observe mode.- `observe`: use observe mode globally.",
									},
									"waf_rules": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "WAF Rules Configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether to host the rules&amp;#39; configuration.- `on`: Enable.- `off`: Disable.",
												},
												"block_rule_i_ds": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Required:    true,
													Description: "Block mode rules list. See details in data source `waf_managed_rules`.",
												},
												"observe_rule_i_ds": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Computed:    true,
													Description: "Observe rules list. See details in data source `waf_managed_rules`.",
												},
											},
										},
									},
									"ai_rule": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "AI based rules configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Valid values:- `smart_status_close`: disabled.- `smart_status_open`: blocked.- `smart_status_observe`: observed.",
												},
											},
										},
									},
								},
							},
						},
						"rate_limit_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "RateLimit Configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"user_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Custom configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_i_d": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID.",
												},
												"rule_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule Name.",
												},
												"threshold": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Threshold of the rate limit. Valid value range: 0-4294967294.",
												},
												"period": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Period of the rate limit. Valid values: 10, 20, 30, 40, 50, 60 (in seconds).",
												},
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Valid values: `monitor`, `drop`.",
												},
												"punish_time": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Punish time, Valid value range: 0-2 days.",
												},
												"punish_time_unit": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Time unit of the punish time. Valid values: `second`, `minutes`, `hour`.",
												},
												"rule_status": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Status of the rule. Valid values: `on`, `off`, `hour`.",
												},
												"freq_fields": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Filter words.",
												},
												"conditions": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Conditions of the rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_from": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Items to match. Valid values:- `host`: Host of the request.- `sip`: Client IP.- `ua`: User-Agent.- `cookie`: Session cookie.- `cgi`: CGI script.- `xff`: XFF extension header.- `url`: URL of the request.- `accept`: Accept encoding of the request.- `method`: HTTP method of the request.- `header`: HTTP header of the request.- `sip_proto`: Network protocol of the request.",
															},
															"match_param": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Parameter for match item. For example, when match from header, match parameter can be set to a header key.",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Valid values:- `equal`: string equal.- `not_equal`: string not equal.- `include`: string include.- `not_include`: string not include.- `match`: ip match.- `not_match`: ip not match.- `include_area`: area include.- `is_empty`: field existed but empty.- `not_exists`: field is not existed.- `regexp`: regex match.- `len_gt`: value greater than.- `len_lt`: value less than.- `len_eq`: value equal.- `match_prefix`: string prefix match.- `match_suffix`: string suffix match.- `wildcard`: wildcard match.",
															},
															"match_content": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Content to match.",
															},
														},
													},
												},
												"rule_priority": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Priority of the rule. Valid value range: 1-100.",
												},
												"update_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Last modification date.",
												},
											},
										},
									},
									"template": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Default Template. Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Template Name. Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"detail": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Detail of the template.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Template Name. Note: This field may return null, indicating that no valid value can be obtained.",
															},
															"i_d": {
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Template ID. Note: This field may return null, indicating that no valid value can be obtained.",
															},
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Action to take.",
															},
															"punish_time": {
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Punish time.",
															},
															"threshold": {
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Threshold.",
															},
															"period": {
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Period.",
															},
														},
													},
												},
											},
										},
									},
									"intelligence": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Intelligent client filter.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "- `on`: Enable.- `off`: Disable.",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Action to take. Valid values: `monitor`, `alg`.",
												},
											},
										},
									},
								},
							},
						},
						"acl_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "ACL configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"user_rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Custom configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_i_d": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID.",
												},
												"rule_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Rule name.",
												},
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action to take. Valid values: `trans`, `drop`, `monitor`, `ban`, `redirect`, `page`, `alg`.",
												},
												"rule_status": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Status of the rule. Valid values: `on`, `off`.",
												},
												"conditions": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Conditions of the rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_from": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Items to match. Valid values:- `host`: Host of the request.- `sip`: Client IP.- `ua`: User-Agent.- `cookie`: Session cookie.- `cgi`: CGI script.- `xff`: XFF extension header.- `url`: URL of the request.- `accept`: Accept encoding of the request.- `method`: HTTP method of the request.- `header`: HTTP header of the request.- `sip_proto`: Network protocol of the request.",
															},
															"match_param": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Parameter for match item. For example, when match from header, match parameter can be set to a header key.",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Valid values:- `equal`: string equal.- `not_equal`: string not equal.- `include`: string include.- `not_include`: string not include.- `match`: ip match.- `not_match`: ip not match.- `include_area`: area include.- `is_empty`: field existed but empty.- `not_exists`: field is not existed.- `regexp`: regex match.- `len_gt`: value greater than.- `len_lt`: value less than.- `len_eq`: value equal.- `match_prefix`: string prefix match.- `match_suffix`: string suffix match.- `wildcard`: wildcard match.",
															},
															"match_content": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Content to match.",
															},
														},
													},
												},
												"rule_priority": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Priority of the rule. Valid value range: 0-100.",
												},
												"update_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Last modification date.",
												},
												"punish_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Punish time, Valid value range: 0-2 days.",
												},
												"punish_time_unit": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Time unit of the punish time. Valid values: `second`, `minutes`, `hour`.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the custom response page.",
												},
												"page_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "ID of the custom response page.",
												},
												"redirect_url": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Redirect target URL, must be an sub-domain from one of the account&amp;#39;s site.",
												},
												"response_code": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Response code to use when redirecting.",
												},
											},
										},
									},
								},
							},
						},
						"bot_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Bot Configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"managed_rule": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Preset rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_i_d": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID.",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Action to take. Valid values: `drop`, `trans`, `monitor`, `alg`.",
												},
												"punish_time": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Punish time.",
												},
												"punish_time_unit": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Time unit of the punish time.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the custom response page.",
												},
												"page_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "ID of the custom response page.",
												},
												"redirect_url": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Redirect target URL, must be an sub-domain from one of the account&amp;#39;s site.",
												},
												"response_code": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Response code to use when redirecting.",
												},
												"trans_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `trans`. See details in data source `bot_managed_rules`.",
												},
												"alg_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `alg`. See details in data source `bot_managed_rules`.",
												},
												"cap_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `captcha`. See details in data source `bot_managed_rules`.",
												},
												"mon_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `monitor`. See details in data source `bot_managed_rules`.",
												},
												"drop_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `drop`. See details in data source `bot_managed_rules`.",
												},
											},
										},
									},
									"portrait_rule": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Portrait rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_i_d": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID.",
												},
												"alg_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `alg`. See details in data source `bot_portrait_rules`.",
												},
												"cap_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `captcha`. See details in data source `bot_portrait_rules`.",
												},
												"mon_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `monitor`. See details in data source `bot_portrait_rules`.",
												},
												"drop_managed_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Optional:    true,
													Description: "Rules to enable when action is `drop`. See details in data source `bot_portrait_rules`.",
												},
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "- `on`: Enable.- `off`: Disable.",
												},
											},
										},
									},
									"intelligence_rule": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Bot intelligent rule configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "- `on`: Enable.- `off`: Disable.",
												},
												"items": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Configuration detail.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"label": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Bot label, valid values: `evil_bot`, `suspect_bot`, `good_bot`, `normal`.",
															},
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Action to take. Valid values: `trans`, `monitor`, `alg`, `captcha`, `drop`.",
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
						"switch_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Main switch of 7-layer security.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"web_switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
								},
							},
						},
						"ip_table_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Basic access control.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Rules list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Actions to take. Valid values: `drop`, `trans`, `monitor`.",
												},
												"match_from": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Matching type. Valid values: `ip`, `area`.",
												},
												"match_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Matching content.",
												},
												"rule_i_d": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID.",
												},
												"update_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Last modification date.",
												},
											},
										},
									},
								},
							},
						},
						"except_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Exception rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"except_user_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Exception rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_i_d": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Rule ID.",
												},
												"rule_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule name.",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Action to take. Valid values: `skip`.",
												},
												"rule_status": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Status of the rule. Valid values:- `on`: Enabled.- `off`: Disabled.",
												},
												"update_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Last modification date.",
												},
												"rule_priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Priority of the rule. Valid value range: 0-100.",
												},
												"except_user_rule_conditions": {
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Conditions of the rule.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_from": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Items to match. Valid values:- `host`: Host of the request.- `sip`: Client IP.- `ua`: User-Agent.- `cookie`: Session cookie.- `cgi`: CGI script.- `xff`: XFF extension header.- `url`: URL of the request.- `accept`: Accept encoding of the request.- `method`: HTTP method of the request.- `header`: HTTP header of the request.- `sip_proto`: Network protocol of the request.",
															},
															"match_param": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Parameter for match item. For example, when match from header, match parameter can be set to a header key.",
															},
															"operator": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Valid values:- `equal`: string equal.- `not_equal`: string not equal.- `include`: string include.- `not_include`: string not include.- `match`: ip match.- `not_match`: ip not match.- `include_area`: area include.- `is_empty`: field existed but empty.- `not_exists`: field is not existed.- `regexp`: regex match.- `len_gt`: value greater than.- `len_lt`: value less than.- `len_eq`: value equal.- `match_prefix`: string prefix match.- `match_suffix`: string suffix match.- `wildcard`: wildcard match.",
															},
															"match_content": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Content to match.",
															},
														},
													},
												},
												"except_user_rule_scope": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Scope of the rule in effect.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"modules": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Computed:    true,
																Description: "Modules in which the rule take effect. Valid values: `waf`.",
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
						"drop_page_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Custom drop page configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "- `on`: Enable.- `off`: Disable.",
									},
									"waf_drop_page_detail": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Custom error page of WAF rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"page_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "ID of the custom error page. when set to 0, use system default error page.",
												},
												"status_code": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "HTTP status code to use. Valid range: 100-600.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "File name or URL.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Type of the custom error page. Valid values: `file`, `url`.",
												},
											},
										},
									},
									"acl_drop_page_detail": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Custom error page of ACL rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"page_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "ID of the custom error page. when set to 0, use system default error page.",
												},
												"status_code": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "HTTP status code to use. Valid range: 100-600.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "File name or URL.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Type of the custom error page. Valid values: `file`, `url`.",
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
	}
}

func resourceTencentCloudTeoSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.create")()
	defer inconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	var entity string
	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
	}

	d.SetId(strings.Join([]string{zoneId, entity}, FILED_SP))

	return resourceTencentCloudTeoSecurityPolicyUpdate(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	entity := idSplit[1]

	securityPolicy, err := service.DescribeTeoSecurityPolicyById(ctx, zoneId, entity)
	if err != nil {
		return err
	}

	if securityPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoSecurityPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityPolicy.ZoneId != nil {
		_ = d.Set("zone_id", securityPolicy.ZoneId)
	}

	if securityPolicy.Entity != nil {
		_ = d.Set("entity", securityPolicy.Entity)
	}

	if securityPolicy.Config != nil {
		configMap := map[string]interface{}{}

		if securityPolicy.Config.WafConfig != nil {
			wafConfigMap := map[string]interface{}{}

			if securityPolicy.Config.WafConfig.Switch != nil {
				wafConfigMap["switch"] = securityPolicy.Config.WafConfig.Switch
			}

			if securityPolicy.Config.WafConfig.Level != nil {
				wafConfigMap["level"] = securityPolicy.Config.WafConfig.Level
			}

			if securityPolicy.Config.WafConfig.Mode != nil {
				wafConfigMap["mode"] = securityPolicy.Config.WafConfig.Mode
			}

			if securityPolicy.Config.WafConfig.WafRules != nil {
				wafRulesMap := map[string]interface{}{}

				if securityPolicy.Config.WafConfig.WafRules.Switch != nil {
					wafRulesMap["switch"] = securityPolicy.Config.WafConfig.WafRules.Switch
				}

				if securityPolicy.Config.WafConfig.WafRules.BlockRuleIDs != nil {
					wafRulesMap["block_rule_i_ds"] = securityPolicy.Config.WafConfig.WafRules.BlockRuleIDs
				}

				if securityPolicy.Config.WafConfig.WafRules.ObserveRuleIDs != nil {
					wafRulesMap["observe_rule_i_ds"] = securityPolicy.Config.WafConfig.WafRules.ObserveRuleIDs
				}

				wafConfigMap["waf_rules"] = []interface{}{wafRulesMap}
			}

			if securityPolicy.Config.WafConfig.AiRule != nil {
				aiRuleMap := map[string]interface{}{}

				if securityPolicy.Config.WafConfig.AiRule.Mode != nil {
					aiRuleMap["mode"] = securityPolicy.Config.WafConfig.AiRule.Mode
				}

				wafConfigMap["ai_rule"] = []interface{}{aiRuleMap}
			}

			configMap["waf_config"] = []interface{}{wafConfigMap}
		}

		if securityPolicy.Config.RateLimitConfig != nil {
			rateLimitConfigMap := map[string]interface{}{}

			if securityPolicy.Config.RateLimitConfig.Switch != nil {
				rateLimitConfigMap["switch"] = securityPolicy.Config.RateLimitConfig.Switch
			}

			if securityPolicy.Config.RateLimitConfig.UserRules != nil {
				userRulesList := []interface{}{}
				for _, userRules := range securityPolicy.Config.RateLimitConfig.UserRules {
					userRulesMap := map[string]interface{}{}

					if userRules.RuleID != nil {
						userRulesMap["rule_i_d"] = userRules.RuleID
					}

					if userRules.RuleName != nil {
						userRulesMap["rule_name"] = userRules.RuleName
					}

					if userRules.Threshold != nil {
						userRulesMap["threshold"] = userRules.Threshold
					}

					if userRules.Period != nil {
						userRulesMap["period"] = userRules.Period
					}

					if userRules.Action != nil {
						userRulesMap["action"] = userRules.Action
					}

					if userRules.PunishTime != nil {
						userRulesMap["punish_time"] = userRules.PunishTime
					}

					if userRules.PunishTimeUnit != nil {
						userRulesMap["punish_time_unit"] = userRules.PunishTimeUnit
					}

					if userRules.RuleStatus != nil {
						userRulesMap["rule_status"] = userRules.RuleStatus
					}

					if userRules.FreqFields != nil {
						userRulesMap["freq_fields"] = userRules.FreqFields
					}

					if userRules.Conditions != nil {
						conditionsList := []interface{}{}
						for _, conditions := range userRules.Conditions {
							conditionsMap := map[string]interface{}{}

							if conditions.MatchFrom != nil {
								conditionsMap["match_from"] = conditions.MatchFrom
							}

							if conditions.MatchParam != nil {
								conditionsMap["match_param"] = conditions.MatchParam
							}

							if conditions.Operator != nil {
								conditionsMap["operator"] = conditions.Operator
							}

							if conditions.MatchContent != nil {
								conditionsMap["match_content"] = conditions.MatchContent
							}

							conditionsList = append(conditionsList, conditionsMap)
						}

						userRulesMap["conditions"] = []interface{}{conditionsList}
					}

					if userRules.RulePriority != nil {
						userRulesMap["rule_priority"] = userRules.RulePriority
					}

					if userRules.UpdateTime != nil {
						userRulesMap["update_time"] = userRules.UpdateTime
					}

					userRulesList = append(userRulesList, userRulesMap)
				}

				rateLimitConfigMap["user_rules"] = []interface{}{userRulesList}
			}

			if securityPolicy.Config.RateLimitConfig.Template != nil {
				templateMap := map[string]interface{}{}

				if securityPolicy.Config.RateLimitConfig.Template.Mode != nil {
					templateMap["mode"] = securityPolicy.Config.RateLimitConfig.Template.Mode
				}

				if securityPolicy.Config.RateLimitConfig.Template.Detail != nil {
					detailMap := map[string]interface{}{}

					if securityPolicy.Config.RateLimitConfig.Template.Detail.Mode != nil {
						detailMap["mode"] = securityPolicy.Config.RateLimitConfig.Template.Detail.Mode
					}

					if securityPolicy.Config.RateLimitConfig.Template.Detail.ID != nil {
						detailMap["i_d"] = securityPolicy.Config.RateLimitConfig.Template.Detail.ID
					}

					if securityPolicy.Config.RateLimitConfig.Template.Detail.Action != nil {
						detailMap["action"] = securityPolicy.Config.RateLimitConfig.Template.Detail.Action
					}

					if securityPolicy.Config.RateLimitConfig.Template.Detail.PunishTime != nil {
						detailMap["punish_time"] = securityPolicy.Config.RateLimitConfig.Template.Detail.PunishTime
					}

					if securityPolicy.Config.RateLimitConfig.Template.Detail.Threshold != nil {
						detailMap["threshold"] = securityPolicy.Config.RateLimitConfig.Template.Detail.Threshold
					}

					if securityPolicy.Config.RateLimitConfig.Template.Detail.Period != nil {
						detailMap["period"] = securityPolicy.Config.RateLimitConfig.Template.Detail.Period
					}

					templateMap["detail"] = []interface{}{detailMap}
				}

				rateLimitConfigMap["template"] = []interface{}{templateMap}
			}

			if securityPolicy.Config.RateLimitConfig.Intelligence != nil {
				intelligenceMap := map[string]interface{}{}

				if securityPolicy.Config.RateLimitConfig.Intelligence.Switch != nil {
					intelligenceMap["switch"] = securityPolicy.Config.RateLimitConfig.Intelligence.Switch
				}

				if securityPolicy.Config.RateLimitConfig.Intelligence.Action != nil {
					intelligenceMap["action"] = securityPolicy.Config.RateLimitConfig.Intelligence.Action
				}

				rateLimitConfigMap["intelligence"] = []interface{}{intelligenceMap}
			}

			configMap["rate_limit_config"] = []interface{}{rateLimitConfigMap}
		}

		if securityPolicy.Config.AclConfig != nil {
			aclConfigMap := map[string]interface{}{}

			if securityPolicy.Config.AclConfig.Switch != nil {
				aclConfigMap["switch"] = securityPolicy.Config.AclConfig.Switch
			}

			if securityPolicy.Config.AclConfig.UserRules != nil {
				userRulesList := []interface{}{}
				for _, userRules := range securityPolicy.Config.AclConfig.UserRules {
					userRulesMap := map[string]interface{}{}

					if userRules.RuleID != nil {
						userRulesMap["rule_i_d"] = userRules.RuleID
					}

					if userRules.RuleName != nil {
						userRulesMap["rule_name"] = userRules.RuleName
					}

					if userRules.Action != nil {
						userRulesMap["action"] = userRules.Action
					}

					if userRules.RuleStatus != nil {
						userRulesMap["rule_status"] = userRules.RuleStatus
					}

					if userRules.Conditions != nil {
						conditionsList := []interface{}{}
						for _, conditions := range userRules.Conditions {
							conditionsMap := map[string]interface{}{}

							if conditions.MatchFrom != nil {
								conditionsMap["match_from"] = conditions.MatchFrom
							}

							if conditions.MatchParam != nil {
								conditionsMap["match_param"] = conditions.MatchParam
							}

							if conditions.Operator != nil {
								conditionsMap["operator"] = conditions.Operator
							}

							if conditions.MatchContent != nil {
								conditionsMap["match_content"] = conditions.MatchContent
							}

							conditionsList = append(conditionsList, conditionsMap)
						}

						userRulesMap["conditions"] = []interface{}{conditionsList}
					}

					if userRules.RulePriority != nil {
						userRulesMap["rule_priority"] = userRules.RulePriority
					}

					if userRules.UpdateTime != nil {
						userRulesMap["update_time"] = userRules.UpdateTime
					}

					if userRules.PunishTime != nil {
						userRulesMap["punish_time"] = userRules.PunishTime
					}

					if userRules.PunishTimeUnit != nil {
						userRulesMap["punish_time_unit"] = userRules.PunishTimeUnit
					}

					if userRules.Name != nil {
						userRulesMap["name"] = userRules.Name
					}

					if userRules.PageId != nil {
						userRulesMap["page_id"] = userRules.PageId
					}

					if userRules.RedirectUrl != nil {
						userRulesMap["redirect_url"] = userRules.RedirectUrl
					}

					if userRules.ResponseCode != nil {
						userRulesMap["response_code"] = userRules.ResponseCode
					}

					userRulesList = append(userRulesList, userRulesMap)
				}

				aclConfigMap["user_rules"] = []interface{}{userRulesList}
			}

			configMap["acl_config"] = []interface{}{aclConfigMap}
		}

		if securityPolicy.Config.BotConfig != nil {
			botConfigMap := map[string]interface{}{}

			if securityPolicy.Config.BotConfig.Switch != nil {
				botConfigMap["switch"] = securityPolicy.Config.BotConfig.Switch
			}

			if securityPolicy.Config.BotConfig.ManagedRule != nil {
				managedRuleMap := map[string]interface{}{}

				if securityPolicy.Config.BotConfig.ManagedRule.RuleID != nil {
					managedRuleMap["rule_i_d"] = securityPolicy.Config.BotConfig.ManagedRule.RuleID
				}

				if securityPolicy.Config.BotConfig.ManagedRule.Action != nil {
					managedRuleMap["action"] = securityPolicy.Config.BotConfig.ManagedRule.Action
				}

				if securityPolicy.Config.BotConfig.ManagedRule.PunishTime != nil {
					managedRuleMap["punish_time"] = securityPolicy.Config.BotConfig.ManagedRule.PunishTime
				}

				if securityPolicy.Config.BotConfig.ManagedRule.PunishTimeUnit != nil {
					managedRuleMap["punish_time_unit"] = securityPolicy.Config.BotConfig.ManagedRule.PunishTimeUnit
				}

				if securityPolicy.Config.BotConfig.ManagedRule.Name != nil {
					managedRuleMap["name"] = securityPolicy.Config.BotConfig.ManagedRule.Name
				}

				if securityPolicy.Config.BotConfig.ManagedRule.PageId != nil {
					managedRuleMap["page_id"] = securityPolicy.Config.BotConfig.ManagedRule.PageId
				}

				if securityPolicy.Config.BotConfig.ManagedRule.RedirectUrl != nil {
					managedRuleMap["redirect_url"] = securityPolicy.Config.BotConfig.ManagedRule.RedirectUrl
				}

				if securityPolicy.Config.BotConfig.ManagedRule.ResponseCode != nil {
					managedRuleMap["response_code"] = securityPolicy.Config.BotConfig.ManagedRule.ResponseCode
				}

				if securityPolicy.Config.BotConfig.ManagedRule.TransManagedIds != nil {
					managedRuleMap["trans_managed_ids"] = securityPolicy.Config.BotConfig.ManagedRule.TransManagedIds
				}

				if securityPolicy.Config.BotConfig.ManagedRule.AlgManagedIds != nil {
					managedRuleMap["alg_managed_ids"] = securityPolicy.Config.BotConfig.ManagedRule.AlgManagedIds
				}

				if securityPolicy.Config.BotConfig.ManagedRule.CapManagedIds != nil {
					managedRuleMap["cap_managed_ids"] = securityPolicy.Config.BotConfig.ManagedRule.CapManagedIds
				}

				if securityPolicy.Config.BotConfig.ManagedRule.MonManagedIds != nil {
					managedRuleMap["mon_managed_ids"] = securityPolicy.Config.BotConfig.ManagedRule.MonManagedIds
				}

				if securityPolicy.Config.BotConfig.ManagedRule.DropManagedIds != nil {
					managedRuleMap["drop_managed_ids"] = securityPolicy.Config.BotConfig.ManagedRule.DropManagedIds
				}

				botConfigMap["managed_rule"] = []interface{}{managedRuleMap}
			}

			if securityPolicy.Config.BotConfig.PortraitRule != nil {
				portraitRuleMap := map[string]interface{}{}

				if securityPolicy.Config.BotConfig.PortraitRule.RuleID != nil {
					portraitRuleMap["rule_i_d"] = securityPolicy.Config.BotConfig.PortraitRule.RuleID
				}

				if securityPolicy.Config.BotConfig.PortraitRule.AlgManagedIds != nil {
					portraitRuleMap["alg_managed_ids"] = securityPolicy.Config.BotConfig.PortraitRule.AlgManagedIds
				}

				if securityPolicy.Config.BotConfig.PortraitRule.CapManagedIds != nil {
					portraitRuleMap["cap_managed_ids"] = securityPolicy.Config.BotConfig.PortraitRule.CapManagedIds
				}

				if securityPolicy.Config.BotConfig.PortraitRule.MonManagedIds != nil {
					portraitRuleMap["mon_managed_ids"] = securityPolicy.Config.BotConfig.PortraitRule.MonManagedIds
				}

				if securityPolicy.Config.BotConfig.PortraitRule.DropManagedIds != nil {
					portraitRuleMap["drop_managed_ids"] = securityPolicy.Config.BotConfig.PortraitRule.DropManagedIds
				}

				if securityPolicy.Config.BotConfig.PortraitRule.Switch != nil {
					portraitRuleMap["switch"] = securityPolicy.Config.BotConfig.PortraitRule.Switch
				}

				botConfigMap["portrait_rule"] = []interface{}{portraitRuleMap}
			}

			if securityPolicy.Config.BotConfig.IntelligenceRule != nil {
				intelligenceRuleMap := map[string]interface{}{}

				if securityPolicy.Config.BotConfig.IntelligenceRule.Switch != nil {
					intelligenceRuleMap["switch"] = securityPolicy.Config.BotConfig.IntelligenceRule.Switch
				}

				if securityPolicy.Config.BotConfig.IntelligenceRule.Items != nil {
					itemsList := []interface{}{}
					for _, items := range securityPolicy.Config.BotConfig.IntelligenceRule.Items {
						itemsMap := map[string]interface{}{}

						if items.Label != nil {
							itemsMap["label"] = items.Label
						}

						if items.Action != nil {
							itemsMap["action"] = items.Action
						}

						itemsList = append(itemsList, itemsMap)
					}

					intelligenceRuleMap["items"] = []interface{}{itemsList}
				}

				botConfigMap["intelligence_rule"] = []interface{}{intelligenceRuleMap}
			}

			configMap["bot_config"] = []interface{}{botConfigMap}
		}

		if securityPolicy.Config.SwitchConfig != nil {
			switchConfigMap := map[string]interface{}{}

			if securityPolicy.Config.SwitchConfig.WebSwitch != nil {
				switchConfigMap["web_switch"] = securityPolicy.Config.SwitchConfig.WebSwitch
			}

			configMap["switch_config"] = []interface{}{switchConfigMap}
		}

		if securityPolicy.Config.IpTableConfig != nil {
			ipTableConfigMap := map[string]interface{}{}

			if securityPolicy.Config.IpTableConfig.Switch != nil {
				ipTableConfigMap["switch"] = securityPolicy.Config.IpTableConfig.Switch
			}

			if securityPolicy.Config.IpTableConfig.Rules != nil {
				rulesList := []interface{}{}
				for _, rules := range securityPolicy.Config.IpTableConfig.Rules {
					rulesMap := map[string]interface{}{}

					if rules.Action != nil {
						rulesMap["action"] = rules.Action
					}

					if rules.MatchFrom != nil {
						rulesMap["match_from"] = rules.MatchFrom
					}

					if rules.MatchContent != nil {
						rulesMap["match_content"] = rules.MatchContent
					}

					if rules.RuleID != nil {
						rulesMap["rule_i_d"] = rules.RuleID
					}

					if rules.UpdateTime != nil {
						rulesMap["update_time"] = rules.UpdateTime
					}

					rulesList = append(rulesList, rulesMap)
				}

				ipTableConfigMap["rules"] = []interface{}{rulesList}
			}

			configMap["ip_table_config"] = []interface{}{ipTableConfigMap}
		}

		if securityPolicy.Config.ExceptConfig != nil {
			exceptConfigMap := map[string]interface{}{}

			if securityPolicy.Config.ExceptConfig.Switch != nil {
				exceptConfigMap["switch"] = securityPolicy.Config.ExceptConfig.Switch
			}

			if securityPolicy.Config.ExceptConfig.ExceptUserRules != nil {
				exceptUserRulesList := []interface{}{}
				for _, exceptUserRules := range securityPolicy.Config.ExceptConfig.ExceptUserRules {
					exceptUserRulesMap := map[string]interface{}{}

					if exceptUserRules.RuleID != nil {
						exceptUserRulesMap["rule_i_d"] = exceptUserRules.RuleID
					}

					if exceptUserRules.RuleName != nil {
						exceptUserRulesMap["rule_name"] = exceptUserRules.RuleName
					}

					if exceptUserRules.Action != nil {
						exceptUserRulesMap["action"] = exceptUserRules.Action
					}

					if exceptUserRules.RuleStatus != nil {
						exceptUserRulesMap["rule_status"] = exceptUserRules.RuleStatus
					}

					if exceptUserRules.UpdateTime != nil {
						exceptUserRulesMap["update_time"] = exceptUserRules.UpdateTime
					}

					if exceptUserRules.RulePriority != nil {
						exceptUserRulesMap["rule_priority"] = exceptUserRules.RulePriority
					}

					if exceptUserRules.ExceptUserRuleConditions != nil {
						exceptUserRuleConditionsList := []interface{}{}
						for _, exceptUserRuleConditions := range exceptUserRules.ExceptUserRuleConditions {
							exceptUserRuleConditionsMap := map[string]interface{}{}

							if exceptUserRuleConditions.MatchFrom != nil {
								exceptUserRuleConditionsMap["match_from"] = exceptUserRuleConditions.MatchFrom
							}

							if exceptUserRuleConditions.MatchParam != nil {
								exceptUserRuleConditionsMap["match_param"] = exceptUserRuleConditions.MatchParam
							}

							if exceptUserRuleConditions.Operator != nil {
								exceptUserRuleConditionsMap["operator"] = exceptUserRuleConditions.Operator
							}

							if exceptUserRuleConditions.MatchContent != nil {
								exceptUserRuleConditionsMap["match_content"] = exceptUserRuleConditions.MatchContent
							}

							exceptUserRuleConditionsList = append(exceptUserRuleConditionsList, exceptUserRuleConditionsMap)
						}

						exceptUserRulesMap["except_user_rule_conditions"] = []interface{}{exceptUserRuleConditionsList}
					}

					if exceptUserRules.ExceptUserRuleScope != nil {
						exceptUserRuleScopeMap := map[string]interface{}{}

						if exceptUserRules.ExceptUserRuleScope.Modules != nil {
							exceptUserRuleScopeMap["modules"] = exceptUserRules.ExceptUserRuleScope.Modules
						}

						exceptUserRulesMap["except_user_rule_scope"] = []interface{}{exceptUserRuleScopeMap}
					}

					exceptUserRulesList = append(exceptUserRulesList, exceptUserRulesMap)
				}

				exceptConfigMap["except_user_rules"] = []interface{}{exceptUserRulesList}
			}

			configMap["except_config"] = []interface{}{exceptConfigMap}
		}

		if securityPolicy.Config.DropPageConfig != nil {
			dropPageConfigMap := map[string]interface{}{}

			if securityPolicy.Config.DropPageConfig.Switch != nil {
				dropPageConfigMap["switch"] = securityPolicy.Config.DropPageConfig.Switch
			}

			if securityPolicy.Config.DropPageConfig.WafDropPageDetail != nil {
				wafDropPageDetailMap := map[string]interface{}{}

				if securityPolicy.Config.DropPageConfig.WafDropPageDetail.PageId != nil {
					wafDropPageDetailMap["page_id"] = securityPolicy.Config.DropPageConfig.WafDropPageDetail.PageId
				}

				if securityPolicy.Config.DropPageConfig.WafDropPageDetail.StatusCode != nil {
					wafDropPageDetailMap["status_code"] = securityPolicy.Config.DropPageConfig.WafDropPageDetail.StatusCode
				}

				if securityPolicy.Config.DropPageConfig.WafDropPageDetail.Name != nil {
					wafDropPageDetailMap["name"] = securityPolicy.Config.DropPageConfig.WafDropPageDetail.Name
				}

				if securityPolicy.Config.DropPageConfig.WafDropPageDetail.Type != nil {
					wafDropPageDetailMap["type"] = securityPolicy.Config.DropPageConfig.WafDropPageDetail.Type
				}

				dropPageConfigMap["waf_drop_page_detail"] = []interface{}{wafDropPageDetailMap}
			}

			if securityPolicy.Config.DropPageConfig.AclDropPageDetail != nil {
				aclDropPageDetailMap := map[string]interface{}{}

				if securityPolicy.Config.DropPageConfig.AclDropPageDetail.PageId != nil {
					aclDropPageDetailMap["page_id"] = securityPolicy.Config.DropPageConfig.AclDropPageDetail.PageId
				}

				if securityPolicy.Config.DropPageConfig.AclDropPageDetail.StatusCode != nil {
					aclDropPageDetailMap["status_code"] = securityPolicy.Config.DropPageConfig.AclDropPageDetail.StatusCode
				}

				if securityPolicy.Config.DropPageConfig.AclDropPageDetail.Name != nil {
					aclDropPageDetailMap["name"] = securityPolicy.Config.DropPageConfig.AclDropPageDetail.Name
				}

				if securityPolicy.Config.DropPageConfig.AclDropPageDetail.Type != nil {
					aclDropPageDetailMap["type"] = securityPolicy.Config.DropPageConfig.AclDropPageDetail.Type
				}

				dropPageConfigMap["acl_drop_page_detail"] = []interface{}{aclDropPageDetailMap}
			}

			configMap["drop_page_config"] = []interface{}{dropPageConfigMap}
		}

		_ = d.Set("config", []interface{}{configMap})
	}

	return nil
}

func resourceTencentCloudTeoSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifySecurityPolicyRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	entity := idSplit[1]

	request.ZoneId = &zoneId
	request.Entity = &entity

	immutableArgs := []string{"zone_id", "entity", "config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
			securityConfig := teo.SecurityConfig{}
			if wafConfigMap, ok := helper.InterfaceToMap(dMap, "waf_config"); ok {
				wafConfig := teo.WafConfig{}
				if v, ok := wafConfigMap["switch"]; ok {
					wafConfig.Switch = helper.String(v.(string))
				}
				if v, ok := wafConfigMap["level"]; ok {
					wafConfig.Level = helper.String(v.(string))
				}
				if v, ok := wafConfigMap["mode"]; ok {
					wafConfig.Mode = helper.String(v.(string))
				}
				if wafRulesMap, ok := helper.InterfaceToMap(wafConfigMap, "waf_rules"); ok {
					wafRule := teo.WafRule{}
					if v, ok := wafRulesMap["switch"]; ok {
						wafRule.Switch = helper.String(v.(string))
					}
					if v, ok := wafRulesMap["block_rule_i_ds"]; ok {
						blockRuleIDsSet := v.(*schema.Set).List()
						for i := range blockRuleIDsSet {
							blockRuleIDs := blockRuleIDsSet[i].(int)
							wafRule.BlockRuleIDs = append(wafRule.BlockRuleIDs, helper.IntInt64(blockRuleIDs))
						}
					}
					if v, ok := wafRulesMap["observe_rule_i_ds"]; ok {
						observeRuleIDsSet := v.(*schema.Set).List()
						for i := range observeRuleIDsSet {
							observeRuleIDs := observeRuleIDsSet[i].(int)
							wafRule.ObserveRuleIDs = append(wafRule.ObserveRuleIDs, helper.IntInt64(observeRuleIDs))
						}
					}
					wafConfig.WafRules = &wafRule
				}
				if aiRuleMap, ok := helper.InterfaceToMap(wafConfigMap, "ai_rule"); ok {
					aiRule := teo.AiRule{}
					if v, ok := aiRuleMap["mode"]; ok {
						aiRule.Mode = helper.String(v.(string))
					}
					wafConfig.AiRule = &aiRule
				}
				securityConfig.WafConfig = &wafConfig
			}
			if rateLimitConfigMap, ok := helper.InterfaceToMap(dMap, "rate_limit_config"); ok {
				rateLimitConfig := teo.RateLimitConfig{}
				if v, ok := rateLimitConfigMap["switch"]; ok {
					rateLimitConfig.Switch = helper.String(v.(string))
				}
				if v, ok := rateLimitConfigMap["user_rules"]; ok {
					for _, item := range v.([]interface{}) {
						userRulesMap := item.(map[string]interface{})
						rateLimitUserRule := teo.RateLimitUserRule{}
						if v, ok := userRulesMap["rule_name"]; ok {
							rateLimitUserRule.RuleName = helper.String(v.(string))
						}
						if v, ok := userRulesMap["threshold"]; ok {
							rateLimitUserRule.Threshold = helper.IntInt64(v.(int))
						}
						if v, ok := userRulesMap["period"]; ok {
							rateLimitUserRule.Period = helper.IntInt64(v.(int))
						}
						if v, ok := userRulesMap["action"]; ok {
							rateLimitUserRule.Action = helper.String(v.(string))
						}
						if v, ok := userRulesMap["punish_time"]; ok {
							rateLimitUserRule.PunishTime = helper.IntInt64(v.(int))
						}
						if v, ok := userRulesMap["punish_time_unit"]; ok {
							rateLimitUserRule.PunishTimeUnit = helper.String(v.(string))
						}
						if v, ok := userRulesMap["rule_status"]; ok {
							rateLimitUserRule.RuleStatus = helper.String(v.(string))
						}
						if v, ok := userRulesMap["freq_fields"]; ok {
							freqFieldsSet := v.(*schema.Set).List()
							for i := range freqFieldsSet {
								freqFields := freqFieldsSet[i].(string)
								rateLimitUserRule.FreqFields = append(rateLimitUserRule.FreqFields, &freqFields)
							}
						}
						if v, ok := userRulesMap["conditions"]; ok {
							for _, item := range v.([]interface{}) {
								conditionsMap := item.(map[string]interface{})
								aCLCondition := teo.ACLCondition{}
								if v, ok := conditionsMap["match_from"]; ok {
									aCLCondition.MatchFrom = helper.String(v.(string))
								}
								if v, ok := conditionsMap["match_param"]; ok {
									aCLCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := conditionsMap["operator"]; ok {
									aCLCondition.Operator = helper.String(v.(string))
								}
								if v, ok := conditionsMap["match_content"]; ok {
									aCLCondition.MatchContent = helper.String(v.(string))
								}
								rateLimitUserRule.Conditions = append(rateLimitUserRule.Conditions, &aCLCondition)
							}
						}
						if v, ok := userRulesMap["rule_priority"]; ok {
							rateLimitUserRule.RulePriority = helper.IntInt64(v.(int))
						}
						rateLimitConfig.UserRules = append(rateLimitConfig.UserRules, &rateLimitUserRule)
					}
				}
				if templateMap, ok := helper.InterfaceToMap(rateLimitConfigMap, "template"); ok {
					rateLimitTemplate := teo.RateLimitTemplate{}
					if v, ok := templateMap["mode"]; ok {
						rateLimitTemplate.Mode = helper.String(v.(string))
					}
					if detailMap, ok := helper.InterfaceToMap(templateMap, "detail"); ok {
						rateLimitTemplateDetail := teo.RateLimitTemplateDetail{}
						if v, ok := detailMap["mode"]; ok {
							rateLimitTemplateDetail.Mode = helper.String(v.(string))
						}
						if v, ok := detailMap["i_d"]; ok {
							rateLimitTemplateDetail.ID = helper.IntInt64(v.(int))
						}
						if v, ok := detailMap["action"]; ok {
							rateLimitTemplateDetail.Action = helper.String(v.(string))
						}
						if v, ok := detailMap["punish_time"]; ok {
							rateLimitTemplateDetail.PunishTime = helper.IntInt64(v.(int))
						}
						if v, ok := detailMap["threshold"]; ok {
							rateLimitTemplateDetail.Threshold = helper.IntInt64(v.(int))
						}
						if v, ok := detailMap["period"]; ok {
							rateLimitTemplateDetail.Period = helper.IntInt64(v.(int))
						}
						rateLimitTemplate.Detail = &rateLimitTemplateDetail
					}
					rateLimitConfig.Template = &rateLimitTemplate
				}
				if intelligenceMap, ok := helper.InterfaceToMap(rateLimitConfigMap, "intelligence"); ok {
					rateLimitIntelligence := teo.RateLimitIntelligence{}
					if v, ok := intelligenceMap["switch"]; ok {
						rateLimitIntelligence.Switch = helper.String(v.(string))
					}
					if v, ok := intelligenceMap["action"]; ok {
						rateLimitIntelligence.Action = helper.String(v.(string))
					}
					rateLimitConfig.Intelligence = &rateLimitIntelligence
				}
				securityConfig.RateLimitConfig = &rateLimitConfig
			}
			if aclConfigMap, ok := helper.InterfaceToMap(dMap, "acl_config"); ok {
				aclConfig := teo.AclConfig{}
				if v, ok := aclConfigMap["switch"]; ok {
					aclConfig.Switch = helper.String(v.(string))
				}
				if v, ok := aclConfigMap["user_rules"]; ok {
					for _, item := range v.([]interface{}) {
						userRulesMap := item.(map[string]interface{})
						aCLUserRule := teo.ACLUserRule{}
						if v, ok := userRulesMap["rule_name"]; ok {
							aCLUserRule.RuleName = helper.String(v.(string))
						}
						if v, ok := userRulesMap["action"]; ok {
							aCLUserRule.Action = helper.String(v.(string))
						}
						if v, ok := userRulesMap["rule_status"]; ok {
							aCLUserRule.RuleStatus = helper.String(v.(string))
						}
						if v, ok := userRulesMap["conditions"]; ok {
							for _, item := range v.([]interface{}) {
								conditionsMap := item.(map[string]interface{})
								aCLCondition := teo.ACLCondition{}
								if v, ok := conditionsMap["match_from"]; ok {
									aCLCondition.MatchFrom = helper.String(v.(string))
								}
								if v, ok := conditionsMap["match_param"]; ok {
									aCLCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := conditionsMap["operator"]; ok {
									aCLCondition.Operator = helper.String(v.(string))
								}
								if v, ok := conditionsMap["match_content"]; ok {
									aCLCondition.MatchContent = helper.String(v.(string))
								}
								aCLUserRule.Conditions = append(aCLUserRule.Conditions, &aCLCondition)
							}
						}
						if v, ok := userRulesMap["rule_priority"]; ok {
							aCLUserRule.RulePriority = helper.IntInt64(v.(int))
						}
						if v, ok := userRulesMap["punish_time"]; ok {
							aCLUserRule.PunishTime = helper.IntInt64(v.(int))
						}
						if v, ok := userRulesMap["punish_time_unit"]; ok {
							aCLUserRule.PunishTimeUnit = helper.String(v.(string))
						}
						if v, ok := userRulesMap["name"]; ok {
							aCLUserRule.Name = helper.String(v.(string))
						}
						if v, ok := userRulesMap["page_id"]; ok {
							aCLUserRule.PageId = helper.IntInt64(v.(int))
						}
						if v, ok := userRulesMap["redirect_url"]; ok {
							aCLUserRule.RedirectUrl = helper.String(v.(string))
						}
						if v, ok := userRulesMap["response_code"]; ok {
							aCLUserRule.ResponseCode = helper.IntInt64(v.(int))
						}
						aclConfig.UserRules = append(aclConfig.UserRules, &aCLUserRule)
					}
				}
				securityConfig.AclConfig = &aclConfig
			}
			if botConfigMap, ok := helper.InterfaceToMap(dMap, "bot_config"); ok {
				botConfig := teo.BotConfig{}
				if v, ok := botConfigMap["switch"]; ok {
					botConfig.Switch = helper.String(v.(string))
				}
				if managedRuleMap, ok := helper.InterfaceToMap(botConfigMap, "managed_rule"); ok {
					botManagedRule := teo.BotManagedRule{}
					if v, ok := managedRuleMap["action"]; ok {
						botManagedRule.Action = helper.String(v.(string))
					}
					if v, ok := managedRuleMap["punish_time"]; ok {
						botManagedRule.PunishTime = helper.IntInt64(v.(int))
					}
					if v, ok := managedRuleMap["punish_time_unit"]; ok {
						botManagedRule.PunishTimeUnit = helper.String(v.(string))
					}
					if v, ok := managedRuleMap["name"]; ok {
						botManagedRule.Name = helper.String(v.(string))
					}
					if v, ok := managedRuleMap["page_id"]; ok {
						botManagedRule.PageId = helper.IntInt64(v.(int))
					}
					if v, ok := managedRuleMap["redirect_url"]; ok {
						botManagedRule.RedirectUrl = helper.String(v.(string))
					}
					if v, ok := managedRuleMap["response_code"]; ok {
						botManagedRule.ResponseCode = helper.IntInt64(v.(int))
					}
					if v, ok := managedRuleMap["trans_managed_ids"]; ok {
						transManagedIdsSet := v.(*schema.Set).List()
						for i := range transManagedIdsSet {
							transManagedIds := transManagedIdsSet[i].(int)
							botManagedRule.TransManagedIds = append(botManagedRule.TransManagedIds, helper.IntInt64(transManagedIds))
						}
					}
					if v, ok := managedRuleMap["alg_managed_ids"]; ok {
						algManagedIdsSet := v.(*schema.Set).List()
						for i := range algManagedIdsSet {
							algManagedIds := algManagedIdsSet[i].(int)
							botManagedRule.AlgManagedIds = append(botManagedRule.AlgManagedIds, helper.IntInt64(algManagedIds))
						}
					}
					if v, ok := managedRuleMap["cap_managed_ids"]; ok {
						capManagedIdsSet := v.(*schema.Set).List()
						for i := range capManagedIdsSet {
							capManagedIds := capManagedIdsSet[i].(int)
							botManagedRule.CapManagedIds = append(botManagedRule.CapManagedIds, helper.IntInt64(capManagedIds))
						}
					}
					if v, ok := managedRuleMap["mon_managed_ids"]; ok {
						monManagedIdsSet := v.(*schema.Set).List()
						for i := range monManagedIdsSet {
							monManagedIds := monManagedIdsSet[i].(int)
							botManagedRule.MonManagedIds = append(botManagedRule.MonManagedIds, helper.IntInt64(monManagedIds))
						}
					}
					if v, ok := managedRuleMap["drop_managed_ids"]; ok {
						dropManagedIdsSet := v.(*schema.Set).List()
						for i := range dropManagedIdsSet {
							dropManagedIds := dropManagedIdsSet[i].(int)
							botManagedRule.DropManagedIds = append(botManagedRule.DropManagedIds, helper.IntInt64(dropManagedIds))
						}
					}
					botConfig.ManagedRule = &botManagedRule
				}
				if portraitRuleMap, ok := helper.InterfaceToMap(botConfigMap, "portrait_rule"); ok {
					botPortraitRule := teo.BotPortraitRule{}
					if v, ok := portraitRuleMap["alg_managed_ids"]; ok {
						algManagedIdsSet := v.(*schema.Set).List()
						for i := range algManagedIdsSet {
							algManagedIds := algManagedIdsSet[i].(int)
							botPortraitRule.AlgManagedIds = append(botPortraitRule.AlgManagedIds, helper.IntInt64(algManagedIds))
						}
					}
					if v, ok := portraitRuleMap["cap_managed_ids"]; ok {
						capManagedIdsSet := v.(*schema.Set).List()
						for i := range capManagedIdsSet {
							capManagedIds := capManagedIdsSet[i].(int)
							botPortraitRule.CapManagedIds = append(botPortraitRule.CapManagedIds, helper.IntInt64(capManagedIds))
						}
					}
					if v, ok := portraitRuleMap["mon_managed_ids"]; ok {
						monManagedIdsSet := v.(*schema.Set).List()
						for i := range monManagedIdsSet {
							monManagedIds := monManagedIdsSet[i].(int)
							botPortraitRule.MonManagedIds = append(botPortraitRule.MonManagedIds, helper.IntInt64(monManagedIds))
						}
					}
					if v, ok := portraitRuleMap["drop_managed_ids"]; ok {
						dropManagedIdsSet := v.(*schema.Set).List()
						for i := range dropManagedIdsSet {
							dropManagedIds := dropManagedIdsSet[i].(int)
							botPortraitRule.DropManagedIds = append(botPortraitRule.DropManagedIds, helper.IntInt64(dropManagedIds))
						}
					}
					if v, ok := portraitRuleMap["switch"]; ok {
						botPortraitRule.Switch = helper.String(v.(string))
					}
					botConfig.PortraitRule = &botPortraitRule
				}
				if intelligenceRuleMap, ok := helper.InterfaceToMap(botConfigMap, "intelligence_rule"); ok {
					intelligenceRule := teo.IntelligenceRule{}
					if v, ok := intelligenceRuleMap["switch"]; ok {
						intelligenceRule.Switch = helper.String(v.(string))
					}
					if v, ok := intelligenceRuleMap["items"]; ok {
						for _, item := range v.([]interface{}) {
							itemsMap := item.(map[string]interface{})
							intelligenceRuleItem := teo.IntelligenceRuleItem{}
							if v, ok := itemsMap["label"]; ok {
								intelligenceRuleItem.Label = helper.String(v.(string))
							}
							if v, ok := itemsMap["action"]; ok {
								intelligenceRuleItem.Action = helper.String(v.(string))
							}
							intelligenceRule.Items = append(intelligenceRule.Items, &intelligenceRuleItem)
						}
					}
					botConfig.IntelligenceRule = &intelligenceRule
				}
				securityConfig.BotConfig = &botConfig
			}
			if switchConfigMap, ok := helper.InterfaceToMap(dMap, "switch_config"); ok {
				switchConfig := teo.SwitchConfig{}
				if v, ok := switchConfigMap["web_switch"]; ok {
					switchConfig.WebSwitch = helper.String(v.(string))
				}
				securityConfig.SwitchConfig = &switchConfig
			}
			if ipTableConfigMap, ok := helper.InterfaceToMap(dMap, "ip_table_config"); ok {
				ipTableConfig := teo.IpTableConfig{}
				if v, ok := ipTableConfigMap["switch"]; ok {
					ipTableConfig.Switch = helper.String(v.(string))
				}
				if v, ok := ipTableConfigMap["rules"]; ok {
					for _, item := range v.([]interface{}) {
						rulesMap := item.(map[string]interface{})
						ipTableRule := teo.IpTableRule{}
						if v, ok := rulesMap["action"]; ok {
							ipTableRule.Action = helper.String(v.(string))
						}
						if v, ok := rulesMap["match_from"]; ok {
							ipTableRule.MatchFrom = helper.String(v.(string))
						}
						if v, ok := rulesMap["match_content"]; ok {
							ipTableRule.MatchContent = helper.String(v.(string))
						}
						ipTableConfig.Rules = append(ipTableConfig.Rules, &ipTableRule)
					}
				}
				securityConfig.IpTableConfig = &ipTableConfig
			}
			if exceptConfigMap, ok := helper.InterfaceToMap(dMap, "except_config"); ok {
				exceptConfig := teo.ExceptConfig{}
				if v, ok := exceptConfigMap["switch"]; ok {
					exceptConfig.Switch = helper.String(v.(string))
				}
				if v, ok := exceptConfigMap["except_user_rules"]; ok {
					for _, item := range v.([]interface{}) {
						exceptUserRulesMap := item.(map[string]interface{})
						exceptUserRule := teo.ExceptUserRule{}
						if v, ok := exceptUserRulesMap["action"]; ok {
							exceptUserRule.Action = helper.String(v.(string))
						}
						if v, ok := exceptUserRulesMap["rule_status"]; ok {
							exceptUserRule.RuleStatus = helper.String(v.(string))
						}
						if v, ok := exceptUserRulesMap["rule_priority"]; ok {
							exceptUserRule.RulePriority = helper.IntInt64(v.(int))
						}
						if v, ok := exceptUserRulesMap["except_user_rule_conditions"]; ok {
							for _, item := range v.([]interface{}) {
								exceptUserRuleConditionsMap := item.(map[string]interface{})
								exceptUserRuleCondition := teo.ExceptUserRuleCondition{}
								if v, ok := exceptUserRuleConditionsMap["match_from"]; ok {
									exceptUserRuleCondition.MatchFrom = helper.String(v.(string))
								}
								if v, ok := exceptUserRuleConditionsMap["match_param"]; ok {
									exceptUserRuleCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := exceptUserRuleConditionsMap["operator"]; ok {
									exceptUserRuleCondition.Operator = helper.String(v.(string))
								}
								if v, ok := exceptUserRuleConditionsMap["match_content"]; ok {
									exceptUserRuleCondition.MatchContent = helper.String(v.(string))
								}
								exceptUserRule.ExceptUserRuleConditions = append(exceptUserRule.ExceptUserRuleConditions, &exceptUserRuleCondition)
							}
						}
						if exceptUserRuleScopeMap, ok := helper.InterfaceToMap(exceptUserRulesMap, "except_user_rule_scope"); ok {
							exceptUserRuleScope := teo.ExceptUserRuleScope{}
							if v, ok := exceptUserRuleScopeMap["modules"]; ok {
								modulesSet := v.(*schema.Set).List()
								for i := range modulesSet {
									modules := modulesSet[i].(string)
									exceptUserRuleScope.Modules = append(exceptUserRuleScope.Modules, &modules)
								}
							}
							exceptUserRule.ExceptUserRuleScope = &exceptUserRuleScope
						}
						exceptConfig.ExceptUserRules = append(exceptConfig.ExceptUserRules, &exceptUserRule)
					}
				}
				securityConfig.ExceptConfig = &exceptConfig
			}
			if dropPageConfigMap, ok := helper.InterfaceToMap(dMap, "drop_page_config"); ok {
				dropPageConfig := teo.DropPageConfig{}
				if v, ok := dropPageConfigMap["switch"]; ok {
					dropPageConfig.Switch = helper.String(v.(string))
				}
				if wafDropPageDetailMap, ok := helper.InterfaceToMap(dropPageConfigMap, "waf_drop_page_detail"); ok {
					dropPageDetail := teo.DropPageDetail{}
					if v, ok := wafDropPageDetailMap["page_id"]; ok {
						dropPageDetail.PageId = helper.IntInt64(v.(int))
					}
					if v, ok := wafDropPageDetailMap["status_code"]; ok {
						dropPageDetail.StatusCode = helper.IntInt64(v.(int))
					}
					if v, ok := wafDropPageDetailMap["name"]; ok {
						dropPageDetail.Name = helper.String(v.(string))
					}
					if v, ok := wafDropPageDetailMap["type"]; ok {
						dropPageDetail.Type = helper.String(v.(string))
					}
					dropPageConfig.WafDropPageDetail = &dropPageDetail
				}
				if aclDropPageDetailMap, ok := helper.InterfaceToMap(dropPageConfigMap, "acl_drop_page_detail"); ok {
					dropPageDetail := teo.DropPageDetail{}
					if v, ok := aclDropPageDetailMap["page_id"]; ok {
						dropPageDetail.PageId = helper.IntInt64(v.(int))
					}
					if v, ok := aclDropPageDetailMap["status_code"]; ok {
						dropPageDetail.StatusCode = helper.IntInt64(v.(int))
					}
					if v, ok := aclDropPageDetailMap["name"]; ok {
						dropPageDetail.Name = helper.String(v.(string))
					}
					if v, ok := aclDropPageDetailMap["type"]; ok {
						dropPageDetail.Type = helper.String(v.(string))
					}
					dropPageConfig.AclDropPageDetail = &dropPageDetail
				}
				securityConfig.DropPageConfig = &dropPageConfig
			}
			request.Config = &securityConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifySecurityPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo securityPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoSecurityPolicyRead(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
