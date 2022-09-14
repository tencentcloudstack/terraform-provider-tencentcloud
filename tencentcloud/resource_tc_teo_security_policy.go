/*
Provides a resource to create a teo security_policy

Example Usage

```hcl
resource "tencentcloud_teo_security_policy" "security_policy" {
  zone_id = ""
  entity = ""
  config {
		waf_config {
				switch = ""
				level = ""
				mode = ""
			waf_rules {
					switch = ""
					block_rule_i_ds = ""
					observe_rule_ids = ""
			}
			ai_rule {
					mode = ""
			}
		}
		rate_limit_config {
				switch = ""
			user_rules {
					rule_name = ""
					threshold = ""
					period = ""
					action = ""
					punish_time = ""
					punish_time_unit = ""
					rule_status = ""
					freq_fields = ""
				conditions {
						match_from = ""
						match_param = ""
						operator = ""
						match_content = ""
				}
					rule_priority = ""
			}
			template {
					mode = ""
				detail {
						mode = ""
						i_d = ""
						action = ""
						punish_time = ""
						threshold = ""
						period = ""
				}
			}
			intelligence {
					switch = ""
					action = ""
			}
		}
		acl_config {
				switch = ""
			user_rules {
					rule_name = ""
					action = ""
					rule_status = ""
				conditions {
						match_from = ""
						match_param = ""
						operator = ""
						match_content = ""
				}
					rule_priority = ""
					punish_time = ""
					punish_time_unit = ""
					name = ""
					page_id = ""
					redirect_url = ""
					response_code = ""
			}
		}
		bot_config {
				switch = ""
			managed_rule {
					rule_i_d = ""
					action = ""
					punish_time = ""
					punish_time_unit = ""
					name = ""
					page_id = ""
					redirect_url = ""
					response_code = ""
					trans_managed_ids = ""
					alg_managed_ids = ""
					cap_managed_ids = ""
					mon_managed_ids = ""
					drop_managed_ids = ""
			}
			portrait_rule {
					rule_i_d = ""
					alg_managed_ids = ""
					cap_managed_ids = ""
					mon_managed_ids = ""
					drop_managed_ids = ""
					switch = ""
			}
			intelligence_rule {
					switch = ""
				items {
						label = ""
						action = ""
				}
			}
		}
		switch_config {
				web_switch = ""
		}
		ip_table_config {
				switch = ""
			rules {
					action = ""
					match_from = ""
					match_content = ""
					rule_i_d = ""
			}
		}

  }
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo security_policy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_security_policy.security_policy securityPolicy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoSecurityPolicyRead,
		Create: resourceTencentCloudTeoSecurityPolicyCreate,
		Update: resourceTencentCloudTeoSecurityPolicyUpdate,
		Delete: resourceTencentCloudTeoSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain.",
			},

			"config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
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
													Description: "Whether to host the rules&#39; configuration.- `on`: Enable.- `off`: Disable.",
												},
												"block_rule_i_ds": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Required:    true,
													Description: "Block mode rules list. See details in data source `waf_managed_rules`.",
												},
												"observe_rule_ids": {
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
																Description: "Matching field.",
															},
															"match_param": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching string.",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching operator.",
															},
															"match_content": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching content.",
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
													Description: "Status of the rule.",
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
																Description: "Matching field.",
															},
															"match_param": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching string.",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching operator.",
															},
															"match_content": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Matching content.",
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
													Description: "Redirect target URL, must be an sub-domain from one of the account&#39;s site.",
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
													Required:    true,
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
													Description: "Redirect target URL, must be an sub-domain from one of the account&#39;s site.",
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
													Optional:    true,
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
													Optional:    true,
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
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTeoSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		zoneId string
		entity string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
	}

	err := resourceTencentCloudTeoSecurityPolicyUpdate(d, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create teo securityPolicy failed, reason:%+v", logId, err)
		return err
	}

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(60*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoSecurityPolicy(ctx, zoneId, entity)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.SecurityConfig.SwitchConfig.WebSwitch == "on" {
			return nil
		}
		if *instance.SecurityConfig.SwitchConfig.WebSwitch == "off" {
			return resource.NonRetryableError(fmt.Errorf("securityPolicy status is %v, operate failed.", *instance.SecurityConfig.SwitchConfig.WebSwitch))
		}
		return resource.RetryableError(fmt.Errorf("securityPolicy status is %v, retry...", *instance.SecurityConfig.SwitchConfig.WebSwitch))
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId + FILED_SP + entity)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, entity)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudTeoSecurityPolicyRead(d, meta)
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

	securityPolicy, err := service.DescribeTeoSecurityPolicy(ctx, zoneId, entity)

	if err != nil {
		return err
	}

	if securityPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityPolicy` %s does not exist", entity)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("entity", entity)

	if securityPolicy.SecurityConfig != nil {
		configMap := map[string]interface{}{}
		if securityPolicy.SecurityConfig.WafConfig != nil {
			wafConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.WafConfig.Switch != nil {
				wafConfigMap["switch"] = securityPolicy.SecurityConfig.WafConfig.Switch
			}
			if securityPolicy.SecurityConfig.WafConfig.Level != nil {
				wafConfigMap["level"] = securityPolicy.SecurityConfig.WafConfig.Level
			}
			if securityPolicy.SecurityConfig.WafConfig.Mode != nil {
				wafConfigMap["mode"] = securityPolicy.SecurityConfig.WafConfig.Mode
			}
			if securityPolicy.SecurityConfig.WafConfig.WafRule != nil {
				wafRulesMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.WafConfig.WafRule.Switch != nil {
					wafRulesMap["switch"] = securityPolicy.SecurityConfig.WafConfig.WafRule.Switch
				}
				if securityPolicy.SecurityConfig.WafConfig.WafRule.BlockRuleIDs != nil {
					wafRulesMap["block_rule_i_ds"] = securityPolicy.SecurityConfig.WafConfig.WafRule.BlockRuleIDs
				}
				if securityPolicy.SecurityConfig.WafConfig.WafRule.ObserveRuleIDs != nil {
					wafRulesMap["observe_rule_ids"] = securityPolicy.SecurityConfig.WafConfig.WafRule.ObserveRuleIDs
				}

				wafConfigMap["waf_rules"] = []interface{}{wafRulesMap}
			}
			if securityPolicy.SecurityConfig.WafConfig.AiRule != nil {
				aiRuleMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.WafConfig.AiRule.Mode != nil {
					aiRuleMap["mode"] = securityPolicy.SecurityConfig.WafConfig.AiRule.Mode
				}

				wafConfigMap["ai_rule"] = []interface{}{aiRuleMap}
			}

			configMap["waf_config"] = []interface{}{wafConfigMap}
		}
		if securityPolicy.SecurityConfig.RateLimitConfig != nil {
			rateLimitConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.RateLimitConfig.Switch != nil {
				rateLimitConfigMap["switch"] = securityPolicy.SecurityConfig.RateLimitConfig.Switch
			}
			if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitUserRules != nil {
				userRulesList := []interface{}{}
				for _, userRules := range securityPolicy.SecurityConfig.RateLimitConfig.RateLimitUserRules {
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
					if userRules.AclConditions != nil {
						conditionsList := []interface{}{}
						for _, conditions := range userRules.AclConditions {
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
						userRulesMap["conditions"] = conditionsList
					}
					if userRules.RulePriority != nil {
						userRulesMap["rule_priority"] = userRules.RulePriority
					}
					if userRules.UpdateTime != nil {
						userRulesMap["update_time"] = userRules.UpdateTime
					}

					userRulesList = append(userRulesList, userRulesMap)
				}
				rateLimitConfigMap["user_rules"] = userRulesList
			}
			if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate != nil {
				templateMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.Mode != nil {
					templateMap["mode"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.Mode
				}
				if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail != nil {
					detailMap := map[string]interface{}{}
					if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail != nil {
						detailMap["mode"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Mode
					}
					if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.ID != nil {
						detailMap["i_d"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.ID
					}
					if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Action != nil {
						detailMap["action"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Action
					}
					if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.PunishTime != nil {
						detailMap["punish_time"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.PunishTime
					}
					if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Threshold != nil {
						detailMap["threshold"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Threshold
					}
					if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Period != nil {
						detailMap["period"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.Period
					}

					templateMap["detail"] = []interface{}{detailMap}
				}

				rateLimitConfigMap["template"] = []interface{}{templateMap}
			}
			if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitIntelligence != nil {
				intelligenceMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitIntelligence.Switch != nil {
					intelligenceMap["switch"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitIntelligence.Switch
				}
				if securityPolicy.SecurityConfig.RateLimitConfig.RateLimitIntelligence.Action != nil {
					intelligenceMap["action"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitIntelligence.Action
				}

				rateLimitConfigMap["intelligence"] = []interface{}{intelligenceMap}
			}

			configMap["rate_limit_config"] = []interface{}{rateLimitConfigMap}
		}
		if securityPolicy.SecurityConfig.AclConfig != nil {
			aclConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.AclConfig.Switch != nil {
				aclConfigMap["switch"] = securityPolicy.SecurityConfig.AclConfig.Switch
			}
			if securityPolicy.SecurityConfig.AclConfig.AclUserRules != nil {
				userRulesList := []interface{}{}
				for _, userRules := range securityPolicy.SecurityConfig.AclConfig.AclUserRules {
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
					if userRules.AclConditions != nil {
						conditionsList := []interface{}{}
						for _, conditions := range userRules.AclConditions {
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
						userRulesMap["conditions"] = conditionsList
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
				aclConfigMap["user_rules"] = userRulesList
			}

			configMap["acl_config"] = []interface{}{aclConfigMap}
		}
		if securityPolicy.SecurityConfig.BotConfig != nil {
			botConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.BotConfig.Switch != nil {
				botConfigMap["switch"] = securityPolicy.SecurityConfig.BotConfig.Switch
			}
			if securityPolicy.SecurityConfig.BotConfig.BotManagedRule != nil {
				managedRuleMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.RuleID != nil {
					managedRuleMap["rule_i_d"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.RuleID
				}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.Action != nil {
					managedRuleMap["action"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.Action
				}
				//if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.PunishTime != nil {
				//	managedRuleMap["punish_time"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.PunishTime
				//}
				//if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.PunishTimeUnit != nil {
				//	managedRuleMap["punish_time_unit"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.PunishTimeUnit
				//}
				//if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.Name != nil {
				//	managedRuleMap["name"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.Name
				//}
				//if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.PageId != nil {
				//	managedRuleMap["page_id"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.PageId
				//}
				//if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.RedirectUrl != nil {
				//	managedRuleMap["redirect_url"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.RedirectUrl
				//}
				//if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.ResponseCode != nil {
				//	managedRuleMap["response_code"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.ResponseCode
				//}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.TransManagedIds != nil {
					managedRuleMap["trans_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.TransManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.AlgManagedIds != nil {
					managedRuleMap["alg_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.AlgManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.CapManagedIds != nil {
					managedRuleMap["cap_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.CapManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.MonManagedIds != nil {
					managedRuleMap["mon_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.MonManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.DropManagedIds != nil {
					managedRuleMap["drop_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.DropManagedIds
				}

				botConfigMap["managed_rule"] = []interface{}{managedRuleMap}
			}
			if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule != nil {
				portraitRuleMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.RuleID != nil {
					portraitRuleMap["rule_i_d"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.RuleID
				}
				if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.AlgManagedIds != nil {
					portraitRuleMap["alg_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.AlgManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.CapManagedIds != nil {
					portraitRuleMap["cap_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.CapManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.MonManagedIds != nil {
					portraitRuleMap["mon_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.MonManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.DropManagedIds != nil {
					portraitRuleMap["drop_managed_ids"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.DropManagedIds
				}
				if securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.Switch != nil {
					portraitRuleMap["switch"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.Switch
				}

				botConfigMap["portrait_rule"] = []interface{}{portraitRuleMap}
			}
			if securityPolicy.SecurityConfig.BotConfig.IntelligenceRule != nil {
				intelligenceRuleMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.BotConfig.IntelligenceRule.Switch != nil {
					intelligenceRuleMap["switch"] = securityPolicy.SecurityConfig.BotConfig.IntelligenceRule.Switch
				}
				if securityPolicy.SecurityConfig.BotConfig.IntelligenceRule.IntelligenceRuleItems != nil {
					itemsList := []interface{}{}
					for _, items := range securityPolicy.SecurityConfig.BotConfig.IntelligenceRule.IntelligenceRuleItems {
						itemsMap := map[string]interface{}{}
						if items.Label != nil {
							itemsMap["label"] = items.Label
						}
						if items.Action != nil {
							itemsMap["action"] = items.Action
						}

						itemsList = append(itemsList, itemsMap)
					}
					intelligenceRuleMap["items"] = itemsList
				}

				botConfigMap["intelligence_rule"] = []interface{}{intelligenceRuleMap}
			}

			configMap["bot_config"] = []interface{}{botConfigMap}
		}
		if securityPolicy.SecurityConfig.SwitchConfig != nil {
			switchConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.SwitchConfig.WebSwitch != nil {
				switchConfigMap["web_switch"] = securityPolicy.SecurityConfig.SwitchConfig.WebSwitch
			}

			configMap["switch_config"] = []interface{}{switchConfigMap}
		}
		if securityPolicy.SecurityConfig.IpTableConfig != nil {
			ipTableConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.IpTableConfig.Switch != nil {
				ipTableConfigMap["switch"] = securityPolicy.SecurityConfig.IpTableConfig.Switch
			}
			if securityPolicy.SecurityConfig.IpTableConfig.IpTableRules != nil {
				rulesList := []interface{}{}
				for _, rules := range securityPolicy.SecurityConfig.IpTableConfig.IpTableRules {
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
				ipTableConfigMap["rules"] = rulesList
			}

			configMap["ip_table_config"] = []interface{}{ipTableConfigMap}
		}

		_ = d.Set("config", []interface{}{configMap})
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTeoSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifySecurityPolicyRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	entity := idSplit[1]

	request.ZoneId = &zoneId
	request.Entity = &entity

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("entity") {
		return fmt.Errorf("`entity` do not support change now.")
	}

	if d.HasChange("config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
			securityConfig := teo.SecurityConfig{}
			if WafConfigMap, ok := helper.InterfaceToMap(dMap, "waf_config"); ok {
				wafConfig := teo.WafConfig{}
				if v, ok := WafConfigMap["switch"]; ok {
					wafConfig.Switch = helper.String(v.(string))
				}
				if v, ok := WafConfigMap["level"]; ok {
					wafConfig.Level = helper.String(v.(string))
				}
				if v, ok := WafConfigMap["mode"]; ok {
					wafConfig.Mode = helper.String(v.(string))
				}
				if WafRulesMap, ok := helper.InterfaceToMap(WafConfigMap, "waf_rules"); ok {
					wafRule := teo.WafRule{}
					if v, ok := WafRulesMap["switch"]; ok {
						wafRule.Switch = helper.String(v.(string))
					}
					if v, ok := WafRulesMap["block_rule_i_ds"]; ok {
						blockRuleIDsSet := v.(*schema.Set).List()
						for i := range blockRuleIDsSet {
							blockRuleIDs := blockRuleIDsSet[i].(int)
							wafRule.BlockRuleIDs = append(wafRule.BlockRuleIDs, helper.IntInt64(blockRuleIDs))
						}
					}
					if v, ok := WafRulesMap["observe_rule_ids"]; ok {
						observeRuleIDsSet := v.(*schema.Set).List()
						for i := range observeRuleIDsSet {
							observeRuleIDs := observeRuleIDsSet[i].(int)
							wafRule.ObserveRuleIDs = append(wafRule.ObserveRuleIDs, helper.IntInt64(observeRuleIDs))
						}
					}
					wafConfig.WafRule = &wafRule
				}
				if AiRuleMap, ok := helper.InterfaceToMap(WafConfigMap, "ai_rule"); ok {
					aiRule := teo.AiRule{}
					if v, ok := AiRuleMap["mode"]; ok {
						aiRule.Mode = helper.String(v.(string))
					}
					wafConfig.AiRule = &aiRule
				}
				securityConfig.WafConfig = &wafConfig
			}
			if RateLimitConfigMap, ok := helper.InterfaceToMap(dMap, "rate_limit_config"); ok {
				rateLimitConfig := teo.RateLimitConfig{}
				if v, ok := RateLimitConfigMap["switch"]; ok {
					rateLimitConfig.Switch = helper.String(v.(string))
				}
				if v, ok := RateLimitConfigMap["user_rules"]; ok {
					for _, item := range v.([]interface{}) {
						UserRulesMap := item.(map[string]interface{})
						rateLimitUserRule := teo.RateLimitUserRule{}
						if v, ok := UserRulesMap["rule_name"]; ok {
							rateLimitUserRule.RuleName = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["threshold"]; ok {
							rateLimitUserRule.Threshold = helper.IntInt64(v.(int))
						}
						if v, ok := UserRulesMap["period"]; ok {
							rateLimitUserRule.Period = helper.IntInt64(v.(int))
						}
						if v, ok := UserRulesMap["action"]; ok {
							rateLimitUserRule.Action = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["punish_time"]; ok {
							rateLimitUserRule.PunishTime = helper.IntInt64(v.(int))
						}
						if v, ok := UserRulesMap["punish_time_unit"]; ok {
							rateLimitUserRule.PunishTimeUnit = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["rule_status"]; ok {
							rateLimitUserRule.RuleStatus = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["freq_fields"]; ok {
							freqFieldsSet := v.(*schema.Set).List()
							for i := range freqFieldsSet {
								freqFields := freqFieldsSet[i].(string)
								rateLimitUserRule.FreqFields = append(rateLimitUserRule.FreqFields, &freqFields)
							}
						}
						if v, ok := UserRulesMap["conditions"]; ok {
							for _, item := range v.([]interface{}) {
								ConditionsMap := item.(map[string]interface{})
								aCLCondition := teo.AclCondition{}
								if v, ok := ConditionsMap["match_from"]; ok {
									aCLCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := ConditionsMap["match_param"]; ok {
									aCLCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := ConditionsMap["operator"]; ok {
									aCLCondition.Operator = helper.String(v.(string))
								}
								if v, ok := ConditionsMap["match_content"]; ok {
									aCLCondition.MatchContent = helper.String(v.(string))
								}
								rateLimitUserRule.AclConditions = append(rateLimitUserRule.AclConditions, &aCLCondition)
							}
						}
						if v, ok := UserRulesMap["rule_priority"]; ok {
							rateLimitUserRule.RulePriority = helper.IntInt64(v.(int))
						}
						rateLimitConfig.RateLimitUserRules = append(rateLimitConfig.RateLimitUserRules, &rateLimitUserRule)
					}
				}
				if TemplateMap, ok := helper.InterfaceToMap(RateLimitConfigMap, "template"); ok {
					rateLimitTemplate := teo.RateLimitTemplate{}
					if v, ok := TemplateMap["mode"]; ok {
						rateLimitTemplate.Mode = helper.String(v.(string))
					}
					if DetailMap, ok := helper.InterfaceToMap(TemplateMap, "detail"); ok {
						rateLimitTemplateDetail := teo.RateLimitTemplateDetail{}
						if v, ok := DetailMap["mode"]; ok {
							rateLimitTemplateDetail.Mode = helper.String(v.(string))
						}
						if v, ok := DetailMap["i_d"]; ok {
							rateLimitTemplateDetail.ID = helper.IntInt64(v.(int))
						}
						if v, ok := DetailMap["action"]; ok {
							rateLimitTemplateDetail.Action = helper.String(v.(string))
						}
						if v, ok := DetailMap["punish_time"]; ok {
							rateLimitTemplateDetail.PunishTime = helper.IntInt64(v.(int))
						}
						if v, ok := DetailMap["threshold"]; ok {
							rateLimitTemplateDetail.Threshold = helper.IntInt64(v.(int))
						}
						if v, ok := DetailMap["period"]; ok {
							rateLimitTemplateDetail.Period = helper.IntInt64(v.(int))
						}
						rateLimitTemplate.RateLimitTemplateDetail = &rateLimitTemplateDetail
					}
					rateLimitConfig.RateLimitTemplate = &rateLimitTemplate
				}
				if IntelligenceMap, ok := helper.InterfaceToMap(RateLimitConfigMap, "intelligence"); ok {
					rateLimitIntelligence := teo.RateLimitIntelligence{}
					if v, ok := IntelligenceMap["switch"]; ok {
						rateLimitIntelligence.Switch = helper.String(v.(string))
					}
					if v, ok := IntelligenceMap["action"]; ok {
						rateLimitIntelligence.Action = helper.String(v.(string))
					}
					rateLimitConfig.RateLimitIntelligence = &rateLimitIntelligence
				}
				securityConfig.RateLimitConfig = &rateLimitConfig
			}
			if AclConfigMap, ok := helper.InterfaceToMap(dMap, "acl_config"); ok {
				aclConfig := teo.AclConfig{}
				if v, ok := AclConfigMap["switch"]; ok {
					aclConfig.Switch = helper.String(v.(string))
				}
				if v, ok := AclConfigMap["user_rules"]; ok {
					for _, item := range v.([]interface{}) {
						UserRulesMap := item.(map[string]interface{})
						aCLUserRule := teo.AclUserRule{}
						if v, ok := UserRulesMap["rule_name"]; ok {
							aCLUserRule.RuleName = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["action"]; ok {
							aCLUserRule.Action = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["rule_status"]; ok {
							aCLUserRule.RuleStatus = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["conditions"]; ok {
							for _, item := range v.([]interface{}) {
								ConditionsMap := item.(map[string]interface{})
								aCLCondition := teo.AclCondition{}
								if v, ok := ConditionsMap["match_from"]; ok {
									aCLCondition.MatchFrom = helper.String(v.(string))
								}
								if v, ok := ConditionsMap["match_param"]; ok {
									aCLCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := ConditionsMap["operator"]; ok {
									aCLCondition.Operator = helper.String(v.(string))
								}
								if v, ok := ConditionsMap["match_content"]; ok {
									aCLCondition.MatchContent = helper.String(v.(string))
								}
								aCLUserRule.AclConditions = append(aCLUserRule.AclConditions, &aCLCondition)
							}
						}
						if v, ok := UserRulesMap["rule_priority"]; ok {
							aCLUserRule.RulePriority = helper.IntInt64(v.(int))
						}
						if v, ok := UserRulesMap["punish_time"]; ok {
							aCLUserRule.PunishTime = helper.IntInt64(v.(int))
						}
						if v, ok := UserRulesMap["punish_time_unit"]; ok {
							aCLUserRule.PunishTimeUnit = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["name"]; ok {
							aCLUserRule.Name = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["page_id"]; ok {
							aCLUserRule.PageId = helper.IntInt64(v.(int))
						}
						if v, ok := UserRulesMap["redirect_url"]; ok {
							aCLUserRule.RedirectUrl = helper.String(v.(string))
						}
						if v, ok := UserRulesMap["response_code"]; ok {
							aCLUserRule.ResponseCode = helper.IntInt64(v.(int))
						}
						aclConfig.AclUserRules = append(aclConfig.AclUserRules, &aCLUserRule)
					}
				}
				securityConfig.AclConfig = &aclConfig
			}
			if BotConfigMap, ok := helper.InterfaceToMap(dMap, "bot_config"); ok {
				botConfig := teo.BotConfig{}
				if v, ok := BotConfigMap["switch"]; ok {
					botConfig.Switch = helper.String(v.(string))
				}
				if ManagedRuleMap, ok := helper.InterfaceToMap(BotConfigMap, "managed_rule"); ok {
					botManagedRule := teo.BotManagedRule{}
					if v, ok := ManagedRuleMap["rule_i_d"]; ok {
						botManagedRule.RuleID = helper.IntInt64(v.(int))
					}
					if v, ok := ManagedRuleMap["action"]; ok {
						botManagedRule.Action = helper.String(v.(string))
					}
					//if v, ok := ManagedRuleMap["punish_time"]; ok {
					//	botManagedRule.PunishTime = helper.IntInt64(v.(int))
					//}
					//if v, ok := ManagedRuleMap["punish_time_unit"]; ok {
					//	botManagedRule.PunishTimeUnit = helper.String(v.(string))
					//}
					//if v, ok := ManagedRuleMap["name"]; ok {
					//	botManagedRule.Name = helper.String(v.(string))
					//}
					//if v, ok := ManagedRuleMap["page_id"]; ok {
					//	botManagedRule.PageId = helper.IntInt64(v.(int))
					//}
					//if v, ok := ManagedRuleMap["redirect_url"]; ok {
					//	botManagedRule.RedirectUrl = helper.String(v.(string))
					//}
					//if v, ok := ManagedRuleMap["response_code"]; ok {
					//	botManagedRule.ResponseCode = helper.IntInt64(v.(int))
					//}
					if v, ok := ManagedRuleMap["trans_managed_ids"]; ok {
						transManagedIdsSet := v.(*schema.Set).List()
						for i := range transManagedIdsSet {
							transManagedIds := transManagedIdsSet[i].(int)
							botManagedRule.TransManagedIds = append(botManagedRule.TransManagedIds, helper.IntInt64(transManagedIds))
						}
					}
					if v, ok := ManagedRuleMap["alg_managed_ids"]; ok {
						algManagedIdsSet := v.(*schema.Set).List()
						for i := range algManagedIdsSet {
							algManagedIds := algManagedIdsSet[i].(int)
							botManagedRule.AlgManagedIds = append(botManagedRule.AlgManagedIds, helper.IntInt64(algManagedIds))
						}
					}
					if v, ok := ManagedRuleMap["cap_managed_ids"]; ok {
						capManagedIdsSet := v.(*schema.Set).List()
						for i := range capManagedIdsSet {
							capManagedIds := capManagedIdsSet[i].(int)
							botManagedRule.CapManagedIds = append(botManagedRule.CapManagedIds, helper.IntInt64(capManagedIds))
						}
					}
					if v, ok := ManagedRuleMap["mon_managed_ids"]; ok {
						monManagedIdsSet := v.(*schema.Set).List()
						for i := range monManagedIdsSet {
							monManagedIds := monManagedIdsSet[i].(int)
							botManagedRule.MonManagedIds = append(botManagedRule.MonManagedIds, helper.IntInt64(monManagedIds))
						}
					}
					if v, ok := ManagedRuleMap["drop_managed_ids"]; ok {
						dropManagedIdsSet := v.(*schema.Set).List()
						for i := range dropManagedIdsSet {
							dropManagedIds := dropManagedIdsSet[i].(int)
							botManagedRule.DropManagedIds = append(botManagedRule.DropManagedIds, helper.IntInt64(dropManagedIds))
						}
					}
					botConfig.BotManagedRule = &botManagedRule
				}
				if PortraitRuleMap, ok := helper.InterfaceToMap(BotConfigMap, "portrait_rule"); ok {
					botPortraitRule := teo.BotPortraitRule{}
					if v, ok := PortraitRuleMap["rule_i_d"]; ok {
						botPortraitRule.RuleID = helper.IntInt64(v.(int))
					}
					if v, ok := PortraitRuleMap["alg_managed_ids"]; ok {
						algManagedIdsSet := v.(*schema.Set).List()
						for i := range algManagedIdsSet {
							algManagedIds := algManagedIdsSet[i].(int)
							botPortraitRule.AlgManagedIds = append(botPortraitRule.AlgManagedIds, helper.IntInt64(algManagedIds))
						}
					}
					if v, ok := PortraitRuleMap["cap_managed_ids"]; ok {
						capManagedIdsSet := v.(*schema.Set).List()
						for i := range capManagedIdsSet {
							capManagedIds := capManagedIdsSet[i].(int)
							botPortraitRule.CapManagedIds = append(botPortraitRule.CapManagedIds, helper.IntInt64(capManagedIds))
						}
					}
					if v, ok := PortraitRuleMap["mon_managed_ids"]; ok {
						monManagedIdsSet := v.(*schema.Set).List()
						for i := range monManagedIdsSet {
							monManagedIds := monManagedIdsSet[i].(int)
							botPortraitRule.MonManagedIds = append(botPortraitRule.MonManagedIds, helper.IntInt64(monManagedIds))
						}
					}
					if v, ok := PortraitRuleMap["drop_managed_ids"]; ok {
						dropManagedIdsSet := v.(*schema.Set).List()
						for i := range dropManagedIdsSet {
							dropManagedIds := dropManagedIdsSet[i].(int)
							botPortraitRule.DropManagedIds = append(botPortraitRule.DropManagedIds, helper.IntInt64(dropManagedIds))
						}
					}
					if v, ok := PortraitRuleMap["switch"]; ok {
						botPortraitRule.Switch = helper.String(v.(string))
					}
					botConfig.BotPortraitRule = &botPortraitRule
				}
				if IntelligenceRuleMap, ok := helper.InterfaceToMap(BotConfigMap, "intelligence_rule"); ok {
					intelligenceRule := teo.IntelligenceRule{}
					if v, ok := IntelligenceRuleMap["switch"]; ok {
						intelligenceRule.Switch = helper.String(v.(string))
					}
					if v, ok := IntelligenceRuleMap["items"]; ok {
						for _, item := range v.([]interface{}) {
							ItemsMap := item.(map[string]interface{})
							intelligenceRuleItem := teo.IntelligenceRuleItem{}
							if v, ok := ItemsMap["label"]; ok {
								intelligenceRuleItem.Label = helper.String(v.(string))
							}
							if v, ok := ItemsMap["action"]; ok {
								intelligenceRuleItem.Action = helper.String(v.(string))
							}
							intelligenceRule.IntelligenceRuleItems = append(intelligenceRule.IntelligenceRuleItems, &intelligenceRuleItem)
						}
					}
					botConfig.IntelligenceRule = &intelligenceRule
				}
				securityConfig.BotConfig = &botConfig
			}
			if SwitchConfigMap, ok := helper.InterfaceToMap(dMap, "switch_config"); ok {
				switchConfig := teo.SwitchConfig{}
				if v, ok := SwitchConfigMap["web_switch"]; ok {
					switchConfig.WebSwitch = helper.String(v.(string))
				}
				securityConfig.SwitchConfig = &switchConfig
			}
			if IpTableConfigMap, ok := helper.InterfaceToMap(dMap, "ip_table_config"); ok {
				ipTableConfig := teo.IpTableConfig{}
				if v, ok := IpTableConfigMap["switch"]; ok {
					ipTableConfig.Switch = helper.String(v.(string))
				}
				if v, ok := IpTableConfigMap["rules"]; ok {
					for _, item := range v.([]interface{}) {
						RulesMap := item.(map[string]interface{})
						ipTableRule := teo.IpTableRule{}
						if v, ok := RulesMap["action"]; ok {
							ipTableRule.Action = helper.String(v.(string))
						}
						if v, ok := RulesMap["match_from"]; ok {
							ipTableRule.MatchFrom = helper.String(v.(string))
						}
						if v, ok := RulesMap["match_content"]; ok {
							ipTableRule.MatchContent = helper.String(v.(string))
						}
						if v, ok := RulesMap["rule_i_d"]; ok {
							ipTableRule.RuleID = helper.IntInt64(v.(int))
						}
						ipTableConfig.IpTableRules = append(ipTableConfig.IpTableRules, &ipTableRule)
					}
				}
				securityConfig.IpTableConfig = &ipTableConfig
			}
			request.SecurityConfig = &securityConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifySecurityPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo securityPolicy failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("teo", "zone", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoSecurityPolicyRead(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
