/*
Provides a resource to create a teo securityPolicy

Example Usage

```hcl
resource "tencentcloud_teo_security_policy" "securityPolicy" {
  zone_id = ""
  entity  = ""
  config {
    waf_config {
      switch = ""
      level  = ""
      mode   = ""
      waf_rules {
        switch            = ""
        block_rule_ids   = ""
        observe_rule_ids = ""
      }
      ai_rule {
        mode = ""
      }
    }
    rate_limit_config {
      switch = ""
      user_rules {
        rule_name        = ""
        threshold        = ""
        period           = ""
        action           = ""
        punish_time      = ""
        punish_time_unit = ""
        rule_status      = ""
        freq_fields      = ""
        conditions {
          match_from    = ""
          match_param   = ""
          operator      = ""
          match_content = ""
        }
        rule_priority    = ""
      }
      template {
        mode = ""
        detail {
          mode        = ""
          id         = ""
          action      = ""
          punish_time = ""
          threshold   = ""
          period      = ""
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
        rule_name        = ""
        action           = ""
        rule_status      = ""
        conditions {
          match_from    = ""
          match_param   = ""
          operator      = ""
          match_content = ""
        }
        rule_priority    = ""
        punish_time      = ""
        punish_time_unit = ""
        name             = ""
        page_id          = ""
        redirect_url     = ""
        response_code    = ""
      }
    }
    bot_config {
      switch = ""
      managed_rule {
        rule_id          = ""
        action            = ""
        punish_time       = ""
        punish_time_unit  = ""
        name              = ""
        page_id           = ""
        redirect_url      = ""
        response_code     = ""
        trans_managed_ids = ""
        alg_managed_ids   = ""
        cap_managed_ids   = ""
        mon_managed_ids   = ""
        drop_managed_ids  = ""
      }
      portrait_rule {
        rule_id         = ""
        alg_managed_ids  = ""
        cap_managed_ids  = ""
        mon_managed_ids  = ""
        drop_managed_ids = ""
        switch           = ""
      }
      intelligence_rule {
        switch = ""
        items {
          label  = ""
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
        action        = ""
        match_from    = ""
        match_content = ""
        rule_id      = ""
      }
    }

  }
}

```
Import

teo securityPolicy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_security_policy.securityPolicy securityPolicy_id#entity
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
				Required:    true,
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
										Description: "Whether to enable WAF rules. Valid values:- on: Enable.- off: Disable.",
									},
									"level": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protection level. Valid values: `loose`, `normal`, `strict`, `stricter`, `custom`.",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protection mode. Valid values:- block: use block mode globally, you still can set a group of rules to use observe mode.- observe: use observe mode globally.",
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
													Description: "Whether to host the rules&#39; configuration.- on: Enable.- off: Disable.",
												},
												"block_rule_ids": {
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
													Description: "Valid values:- smart_status_close: disabled.- smart_status_open: blocked.- smart_status_observe: observed.",
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
										Required:    true,
										Description: "- on: Enable.- off: Disable.",
									},
									"user_rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Custom configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
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
													Description: "Time unit of the punish time. Valid values: `second`,`minutes`,`hour`.",
												},
												"rule_status": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Status of the rule.",
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
													Description: "Priority of the rule.",
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
										Description: "Default Template.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template Name.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"detail": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Detail of the template.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Template Name.Note: This field may return null, indicating that no valid value can be obtained.",
															},
															"id": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Template ID.Note: This field may return null, indicating that no valid value can be obtained.",
															},
															"action": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Action to take.",
															},
															"punish_time": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Punish time.",
															},
															"threshold": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Threshold.",
															},
															"period": {
																Type:        schema.TypeInt,
																Optional:    true,
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
													Description: "- on: Enable.- off: Disable.",
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
										Description: "- on: Enable.- off: Disable.",
									},
									"user_rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Custom configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
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
										Description: "- on: Enable.- off: Disable.",
									},
									"managed_rule": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Preset rules.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
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
										Description: "Portrait rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule_id": {
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
													Description: "- on: Enable.- off: Disable.",
												},
											},
										},
									},
									"intelligence_rule": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Bot intelligent rule configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "- on: Enable.- off: Disable.",
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
										Description: "- on: Enable.- off: Disable.",
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
										Description: "- on: Enable.- off: Disable.",
									},
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
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
												"rule_id": {
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
		},
	}
}

func resourceTencentCloudTeoSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.create")()
	defer inconsistentCheck(d, meta)()

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

	d.SetId(zoneId + FILED_SP + entity)
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

	securityPolicy, err := service.DescribeTeoSecurityPolicy(ctx, zoneId, entity)

	if err != nil {
		return err
	}

	if securityPolicy == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityPolicy` %s does not exist", d.Id())
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
					wafRulesMap["block_rule_ids"] = securityPolicy.Config.WafConfig.WafRules.BlockRuleIDs
				}
				if securityPolicy.Config.WafConfig.WafRules.ObserveRuleIDs != nil {
					wafRulesMap["observe_rule_ids"] = securityPolicy.Config.WafConfig.WafRules.ObserveRuleIDs
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
						userRulesMap["rule_id"] = userRules.RuleID
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
						detailMap["id"] = securityPolicy.Config.RateLimitConfig.Template.Detail.ID
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
						userRulesMap["rule_id"] = userRules.RuleID
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
		if securityPolicy.Config.BotConfig != nil {
			botConfigMap := map[string]interface{}{}
			if securityPolicy.Config.BotConfig.Switch != nil {
				botConfigMap["switch"] = securityPolicy.Config.BotConfig.Switch
			}
			if securityPolicy.Config.BotConfig.ManagedRule != nil {
				managedRuleMap := map[string]interface{}{}
				if securityPolicy.Config.BotConfig.ManagedRule.RuleID != nil {
					managedRuleMap["rule_id"] = securityPolicy.Config.BotConfig.ManagedRule.RuleID
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
					portraitRuleMap["rule_id"] = securityPolicy.Config.BotConfig.PortraitRule.RuleID
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
					intelligenceRuleMap["items"] = itemsList
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
						rulesMap["rule_id"] = rules.RuleID
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
				if v, ok := WafRulesMap["block_rule_ids"]; ok {
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
				wafConfig.WafRules = &wafRule
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
							aCLCondition := teo.ACLCondition{}
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
							rateLimitUserRule.Conditions = append(rateLimitUserRule.Conditions, &aCLCondition)
						}
					}
					if v, ok := UserRulesMap["rule_priority"]; ok {
						rateLimitUserRule.RulePriority = helper.IntInt64(v.(int))
					}
					rateLimitConfig.UserRules = append(rateLimitConfig.UserRules, &rateLimitUserRule)
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
					if v, ok := DetailMap["id"]; ok {
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
					rateLimitTemplate.Detail = &rateLimitTemplateDetail
				}
				rateLimitConfig.Template = &rateLimitTemplate
			}
			if IntelligenceMap, ok := helper.InterfaceToMap(RateLimitConfigMap, "intelligence"); ok {
				rateLimitIntelligence := teo.RateLimitIntelligence{}
				if v, ok := IntelligenceMap["switch"]; ok {
					rateLimitIntelligence.Switch = helper.String(v.(string))
				}
				if v, ok := IntelligenceMap["action"]; ok {
					rateLimitIntelligence.Action = helper.String(v.(string))
				}
				rateLimitConfig.Intelligence = &rateLimitIntelligence
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
					aCLUserRule := teo.ACLUserRule{}
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
							aCLCondition := teo.ACLCondition{}
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
							aCLUserRule.Conditions = append(aCLUserRule.Conditions, &aCLCondition)
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
					aclConfig.UserRules = append(aclConfig.UserRules, &aCLUserRule)
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
				if v, ok := ManagedRuleMap["rule_id"]; ok {
					botManagedRule.RuleID = helper.IntInt64(v.(int))
				}
				if v, ok := ManagedRuleMap["action"]; ok {
					botManagedRule.Action = helper.String(v.(string))
				}
				if v, ok := ManagedRuleMap["punish_time"]; ok {
					botManagedRule.PunishTime = helper.IntInt64(v.(int))
				}
				if v, ok := ManagedRuleMap["punish_time_unit"]; ok {
					botManagedRule.PunishTimeUnit = helper.String(v.(string))
				}
				if v, ok := ManagedRuleMap["name"]; ok {
					botManagedRule.Name = helper.String(v.(string))
				}
				if v, ok := ManagedRuleMap["page_id"]; ok {
					botManagedRule.PageId = helper.IntInt64(v.(int))
				}
				if v, ok := ManagedRuleMap["redirect_url"]; ok {
					botManagedRule.RedirectUrl = helper.String(v.(string))
				}
				if v, ok := ManagedRuleMap["response_code"]; ok {
					botManagedRule.ResponseCode = helper.IntInt64(v.(int))
				}
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
				botConfig.ManagedRule = &botManagedRule
			}
			if PortraitRuleMap, ok := helper.InterfaceToMap(BotConfigMap, "portrait_rule"); ok {
				botPortraitRule := teo.BotPortraitRule{}
				if v, ok := PortraitRuleMap["rule_id"]; ok {
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
				botConfig.PortraitRule = &botPortraitRule
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
						intelligenceRule.Items = append(intelligenceRule.Items, &intelligenceRuleItem)
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
					if v, ok := RulesMap["rule_id"]; ok {
						ipTableRule.RuleID = helper.IntInt64(v.(int))
					}
					ipTableConfig.Rules = append(ipTableConfig.Rules, &ipTableRule)
				}
			}
			securityConfig.IpTableConfig = &ipTableConfig
		}

		request.Config = &securityConfig
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

	return resourceTencentCloudTeoSecurityPolicyRead(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
