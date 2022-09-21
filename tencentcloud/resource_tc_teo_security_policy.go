/*
Provides a resource to create a teo security_policy

Example Usage

```hcl
resource "tencentcloud_teo_security_policy" "security_policy" {
  entity  = "aaa.sfurnace.work"
  zone_id = "zone-2983wizgxqvm"

  config {
    acl_config {
      switch = "off"
    }

    bot_config {
      switch = "off"

      intelligence_rule {
        switch = "off"

        items {
          action = "drop"
          label  = "evil_bot"
        }
        items {
          action = "alg"
          label  = "suspect_bot"
        }
        items {
          action = "monitor"
          label  = "good_bot"
        }
        items {
          action = "trans"
          label  = "normal"
        }
      }

      managed_rule {
        action            = "monitor"
        alg_managed_ids   = []
        cap_managed_ids   = []
        drop_managed_ids  = []
        mon_managed_ids   = []
        page_id           = 0
        punish_time       = 0
        response_code     = 0
        rule_id           = 0
        trans_managed_ids = []
      }

      portrait_rule {
        alg_managed_ids  = []
        cap_managed_ids  = []
        drop_managed_ids = []
        mon_managed_ids  = []
        rule_id          = -1
        switch           = "off"
      }
    }

    drop_page_config {
      switch = "on"

      acl_drop_page_detail {
        name        = "-"
        page_id     = 0
        status_code = 569
        type        = "default"
      }

      waf_drop_page_detail {
        name        = "-"
        page_id     = 0
        status_code = 566
        type        = "default"
      }
    }

    except_config {
      switch = "on"
    }

    ip_table_config {
      switch = "off"
    }

    rate_limit_config {
      switch = "on"

      intelligence {
        action = "monitor"
        switch = "off"
      }

      template {
        mode = "sup_loose"

        detail {
          action      = "alg"
          id          = 831807989
          mode        = "sup_loose"
          period      = 1
          punish_time = 0
          threshold   = 2000
        }
      }
    }

    switch_config {
      web_switch = "on"
    }

    waf_config {
      level  = "strict"
      mode   = "block"
      switch = "on"

      ai_rule {
        mode = "smart_status_close"
      }

      waf_rules {
        block_rule_ids   = [
          22,
          84214562,
          106246133,
          106246507,
          106246508,
          106246523,
          106246524,
          106246679,
          106247029,
          106247048,
          106247140,
          106247356,
          106247357,
          106247358,
          106247378,
          106247389,
          106247392,
          106247394,
          106247405,
          106247409,
          106247413,
          106247558,
          106247795,
          106247819,
          106248021,
        ]
        observe_rule_ids = []
        switch           = "off"
      }
    }
  }
}

```
Import

teo security_policy can be imported using the zoneId#entity, e.g.
```
$ terraform import tencentcloud_teo_security_policy.security_policy zone-2983wizgxqvm#aaa.sfurnace.work
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
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
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
															"id": {
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
										Optional:    true,
										Computed:    true,
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
													Computed:    true,
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
												"rule_id": {
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
												"rule_id": {
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
												"rule_id": {
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
												"rule_id": {
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

	var (
		logId   = getLogId(contextNil)
		service = TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		zoneId  string
		entity  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
	}

	var ddosPolicy *teo.DescribeZoneDDoSPolicyResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := service.DescribeTeoZoneDDoSPolicyByFilter(ctx, map[string]interface{}{
			"zone_id": zoneId,
		})
		if e != nil {
			return retryError(e)
		}
		ddosPolicy = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo planInfo failed, reason:%+v", logId, err)
		return err
	}

	if len(ddosPolicy.ShieldAreas) > 0 {
	outer:
		for _, areas := range ddosPolicy.ShieldAreas {
			for _, host := range areas.DDoSHosts {
				if host.Host != nil && *host.Host == entity && host.SecurityType != nil && *host.SecurityType != "on" {
					request := teo.NewModifyDDoSPolicyHostRequest()
					request.ZoneId = &zoneId
					request.Host = &entity
					request.PolicyId = areas.PolicyId
					request.SecurityType = helper.String("on")
					request.AccelerateType = helper.String("on")

					err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
						result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDDoSPolicyHost(request)
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
					break outer
				}
			}
		}
	}

	d.SetId(zoneId + FILED_SP + entity)
	err = resourceTencentCloudTeoSecurityPolicyUpdate(d, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create teo ddosPolicy failed, reason:%+v", logId, err)
		return err
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
					wafRulesMap["block_rule_ids"] = securityPolicy.SecurityConfig.WafConfig.WafRule.BlockRuleIDs
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
						detailMap["id"] = securityPolicy.SecurityConfig.RateLimitConfig.RateLimitTemplate.RateLimitTemplateDetail.ID
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
					managedRuleMap["rule_id"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.RuleID
				}
				if securityPolicy.SecurityConfig.BotConfig.BotManagedRule.Action != nil {
					managedRuleMap["action"] = securityPolicy.SecurityConfig.BotConfig.BotManagedRule.Action
				}
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
					portraitRuleMap["rule_id"] = securityPolicy.SecurityConfig.BotConfig.BotPortraitRule.RuleID
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
		if securityPolicy.SecurityConfig.ExceptConfig != nil {
			exceptConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.ExceptConfig.Switch != nil {
				exceptConfigMap["switch"] = securityPolicy.SecurityConfig.ExceptConfig.Switch
			}
			if securityPolicy.SecurityConfig.ExceptConfig.ExceptUserRules != nil {
				exceptUserRulesList := []interface{}{}
				for _, exceptUserRules := range securityPolicy.SecurityConfig.ExceptConfig.ExceptUserRules {
					exceptUserRulesMap := map[string]interface{}{}
					if exceptUserRules.RuleID != nil {
						exceptUserRulesMap["rule_id"] = exceptUserRules.RuleID
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
						exceptUserRulesMap["except_user_rule_conditions"] = exceptUserRuleConditionsList
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
				exceptConfigMap["except_user_rules"] = exceptUserRulesList
			}

			configMap["except_config"] = []interface{}{exceptConfigMap}
		}
		if securityPolicy.SecurityConfig.DropPageConfig != nil {
			dropPageConfigMap := map[string]interface{}{}
			if securityPolicy.SecurityConfig.DropPageConfig.Switch != nil {
				dropPageConfigMap["switch"] = securityPolicy.SecurityConfig.DropPageConfig.Switch
			}
			if securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail != nil {
				wafDropPageDetailMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.PageId != nil {
					wafDropPageDetailMap["page_id"] = securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.PageId
				}
				if securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.StatusCode != nil {
					wafDropPageDetailMap["status_code"] = securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.StatusCode
				}
				if securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.Name != nil {
					wafDropPageDetailMap["name"] = securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.Name
				}
				if securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.Type != nil {
					wafDropPageDetailMap["type"] = securityPolicy.SecurityConfig.DropPageConfig.WafDropPageDetail.Type
				}

				dropPageConfigMap["waf_drop_page_detail"] = []interface{}{wafDropPageDetailMap}
			}
			if securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail != nil {
				aclDropPageDetailMap := map[string]interface{}{}
				if securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.PageId != nil {
					aclDropPageDetailMap["page_id"] = securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.PageId
				}
				if securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.StatusCode != nil {
					aclDropPageDetailMap["status_code"] = securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.StatusCode
				}
				if securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.Name != nil {
					aclDropPageDetailMap["name"] = securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.Name
				}
				if securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.Type != nil {
					aclDropPageDetailMap["type"] = securityPolicy.SecurityConfig.DropPageConfig.AclDropPageDetail.Type
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

	if d.HasChange("zone_id") {
		if old, _ := d.GetChange("zone_id"); old.(string) != "" {
			return fmt.Errorf("`zone_id` do not support change now.")
		}
	}

	if d.HasChange("entity") {
		if old, _ := d.GetChange("entity"); old.(string) != "" {
			return fmt.Errorf("`entity` do not support change now.")
		}
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
					if v, ok := ManagedRuleMap["rule_id"]; ok && v.(int) != 0 {
						botManagedRule.RuleID = helper.IntInt64(v.(int))
					}
					if v, ok := ManagedRuleMap["action"]; ok {
						botManagedRule.Action = helper.String(v.(string))
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
					botConfig.BotManagedRule = &botManagedRule
				}
				if PortraitRuleMap, ok := helper.InterfaceToMap(BotConfigMap, "portrait_rule"); ok {
					botPortraitRule := teo.BotPortraitRule{}
					if v, ok := PortraitRuleMap["rule_id"]; ok && v.(int) != 0 {
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
						if v, ok := RulesMap["rule_id"]; ok && v.(int) != 0 {
							ipTableRule.RuleID = helper.IntInt64(v.(int))
						}
						ipTableConfig.IpTableRules = append(ipTableConfig.IpTableRules, &ipTableRule)
					}
				}
				securityConfig.IpTableConfig = &ipTableConfig
			}
			if ExceptConfigMap, ok := helper.InterfaceToMap(dMap, "except_config"); ok {
				exceptConfig := teo.ExceptConfig{}
				if v, ok := ExceptConfigMap["switch"]; ok {
					exceptConfig.Switch = helper.String(v.(string))
				}
				if v, ok := ExceptConfigMap["except_user_rules"]; ok {
					for _, item := range v.([]interface{}) {
						ExceptUserRulesMap := item.(map[string]interface{})
						exceptUserRule := teo.ExceptUserRule{}
						if v, ok := ExceptUserRulesMap["action"]; ok {
							exceptUserRule.Action = helper.String(v.(string))
						}
						if v, ok := ExceptUserRulesMap["rule_status"]; ok {
							exceptUserRule.RuleStatus = helper.String(v.(string))
						}
						if v, ok := ExceptUserRulesMap["rule_priority"]; ok {
							exceptUserRule.RulePriority = helper.IntInt64(v.(int))
						}
						if v, ok := ExceptUserRulesMap["except_user_rule_conditions"]; ok {
							for _, item := range v.([]interface{}) {
								ExceptUserRuleConditionsMap := item.(map[string]interface{})
								exceptUserRuleCondition := teo.ExceptUserRuleCondition{}
								if v, ok := ExceptUserRuleConditionsMap["match_from"]; ok {
									exceptUserRuleCondition.MatchFrom = helper.String(v.(string))
								}
								if v, ok := ExceptUserRuleConditionsMap["match_param"]; ok {
									exceptUserRuleCondition.MatchParam = helper.String(v.(string))
								}
								if v, ok := ExceptUserRuleConditionsMap["operator"]; ok {
									exceptUserRuleCondition.Operator = helper.String(v.(string))
								}
								if v, ok := ExceptUserRuleConditionsMap["match_content"]; ok {
									exceptUserRuleCondition.MatchContent = helper.String(v.(string))
								}
								exceptUserRule.ExceptUserRuleConditions = append(exceptUserRule.ExceptUserRuleConditions, &exceptUserRuleCondition)
							}
						}
						if ExceptUserRuleScopeMap, ok := helper.InterfaceToMap(ExceptUserRulesMap, "except_user_rule_scope"); ok {
							exceptUserRuleScope := teo.ExceptUserRuleScope{}
							if v, ok := ExceptUserRuleScopeMap["modules"]; ok {
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
			if DropPageConfigMap, ok := helper.InterfaceToMap(dMap, "drop_page_config"); ok {
				dropPageConfig := teo.DropPageConfig{}
				if v, ok := DropPageConfigMap["switch"]; ok {
					dropPageConfig.Switch = helper.String(v.(string))
				}
				if WafDropPageDetailMap, ok := helper.InterfaceToMap(DropPageConfigMap, "waf_drop_page_detail"); ok {
					dropPageDetail := teo.DropPageDetail{}
					if v, ok := WafDropPageDetailMap["page_id"]; ok {
						dropPageDetail.PageId = helper.IntInt64(v.(int))
					}
					if v, ok := WafDropPageDetailMap["status_code"]; ok {
						dropPageDetail.StatusCode = helper.IntInt64(v.(int))
					}
					if v, ok := WafDropPageDetailMap["name"]; ok {
						dropPageDetail.Name = helper.String(v.(string))
					}
					if v, ok := WafDropPageDetailMap["type"]; ok {
						dropPageDetail.Type = helper.String(v.(string))
					}
					dropPageConfig.WafDropPageDetail = &dropPageDetail
				}
				if AclDropPageDetailMap, ok := helper.InterfaceToMap(DropPageConfigMap, "acl_drop_page_detail"); ok {
					dropPageDetail := teo.DropPageDetail{}
					if v, ok := AclDropPageDetailMap["page_id"]; ok {
						dropPageDetail.PageId = helper.IntInt64(v.(int))
					}
					if v, ok := AclDropPageDetailMap["status_code"]; ok {
						dropPageDetail.StatusCode = helper.IntInt64(v.(int))
					}
					if v, ok := AclDropPageDetailMap["name"]; ok {
						dropPageDetail.Name = helper.String(v.(string))
					}
					if v, ok := AclDropPageDetailMap["type"]; ok {
						dropPageDetail.Type = helper.String(v.(string))
					}
					dropPageConfig.AclDropPageDetail = &dropPageDetail
				}
				securityConfig.DropPageConfig = &dropPageConfig
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

	return resourceTencentCloudTeoSecurityPolicyRead(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_security_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
