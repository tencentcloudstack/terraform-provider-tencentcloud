---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_policy"
sidebar_current: "docs-tencentcloud-resource-teo_security_policy"
description: |-
  Provides a resource to create a teo teo_security_policy
---

# tencentcloud_teo_security_policy

Provides a resource to create a teo teo_security_policy

## Example Usage

```hcl
resource "tencentcloud_teo_security_policy" "teo_security_policy" {
  security_config = {
    waf_config = {
      waf_rule = {
      }
      ai_rule = {
      }
    }
    rate_limit_config = {
      rate_limit_user_rules = {
        acl_conditions = {
        }
      }
      rate_limit_template = {
        rate_limit_template_detail = {
        }
      }
      rate_limit_intelligence = {
      }
      rate_limit_customizes = {
        acl_conditions = {
        }
      }
    }
    acl_config = {
      acl_user_rules = {
        acl_conditions = {
        }
      }
      customizes = {
        acl_conditions = {
        }
      }
    }
    bot_config = {
      bot_managed_rule = {
      }
      bot_portrait_rule = {
      }
      intelligence_rule = {
        intelligence_rule_items = {
        }
      }
      bot_user_rules = {
        extend_actions = {
        }
        acl_conditions = {
        }
      }
      alg_detect_rule = {
        alg_conditions = {
        }
        alg_detect_session = {
          alg_detect_results = {
          }
          session_behaviors = {
          }
        }
        alg_detect_js = {
          alg_detect_results = {
          }
        }
      }
      customizes = {
        extend_actions = {
        }
        acl_conditions = {
        }
      }
    }
    switch_config = {
    }
    ip_table_config = {
      ip_table_rules = {
      }
    }
    except_config = {
      except_user_rules = {
        except_user_rule_conditions = {
        }
        except_user_rule_scope = {
          partial_modules = {
          }
          skip_conditions = {
          }
        }
      }
    }
    drop_page_config = {
      waf_drop_page_detail = {
      }
      acl_drop_page_detail = {
      }
    }
    template_config = {
    }
    slow_post_config = {
      first_part_config = {
      }
      slow_rate_config = {
      }
    }
    detect_length_limit_config = {
      detect_length_limit_rules = {
        conditions = {
        }
      }
    }
  }
  security_policy = {
    custom_rules = {
      rules = {
        action = {
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
          redirect_action_parameters = {
          }
        }
      }
    }
    managed_rules = {
      auto_update = {
      }
      managed_rule_groups = {
        action = {
          block_ip_action_parameters = {
          }
          return_custom_page_action_parameters = {
          }
          redirect_action_parameters = {
          }
        }
        rule_actions = {
          action = {
            block_ip_action_parameters = {
            }
            return_custom_page_action_parameters = {
            }
            redirect_action_parameters = {
            }
          }
        }
        meta_data = {
          rule_details = {
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String, ForceNew) Zone ID.
* `entity` - (Optional, String, ForceNew) Security policy type. the following parameter values can be used: <li>ZoneDefaultPolicy: used to specify a site-level policy;</li> <li>Template: used to specify a policy Template. you need to simultaneously specify the TemplateId parameter;</li> <li>Host: used to specify a domain-level policy (note: when using a domain name to specify a dns service policy, only dns services or policy templates that have applied a domain-level policy are supported).</li>.
* `host` - (Optional, String, ForceNew) Specifies the specified domain. when the Entity parameter value is Host, use the domain-level policy specified by this parameter. for example: use www.example.com to configure the domain-level policy of the domain.
* `security_policy` - (Optional, List) Security policy configuration. it is recommended to use for custom policies and managed rule configurations of Web protection. it supports configuring security policies with expression grammar.
* `template_id` - (Optional, String, ForceNew) Specify the policy Template ID. use this parameter to specify the ID of the policy Template when the Entity parameter value is Template.

The `action` object of `managed_rule_groups` supports the following:

* `name` - (Required, String) Specific actions for safe execution. valid values:.
<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.
* `block_ip_action_parameters` - (Optional, List) Additional parameter when Name is BlockIP.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) Additional parameter when Name is ReturnCustomPage.

The `action` object of `rule_actions` supports the following:

* `name` - (Required, String) Specific actions for safe execution. valid values:.
<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.
* `block_ip_action_parameters` - (Optional, List) Additional parameter when Name is BlockIP.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) Additional parameter when Name is ReturnCustomPage.

The `action` object of `rules` supports the following:

* `name` - (Required, String) Specific actions for safe execution. valid values:.
<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.
* `block_ip_action_parameters` - (Optional, List) Additional parameter when Name is BlockIP.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) Additional parameter when Name is ReturnCustomPage.

The `auto_update` object of `managed_rules` supports the following:

* `auto_update_to_latest_version` - (Required, String) Indicates whether to enable automatic update to the latest version. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `ruleset_version` - (Optional, String) The currently used version, in the format compliant with ISO 8601 standard, such as 2023-12-21T12:00:32Z. it is empty by default and is only an output parameter.

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.

The `custom_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.

The `managed_rule_groups` object of `managed_rules` supports the following:

* `action` - (Required, List) Handling actions for managed rule groups. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process requests and record security events in logs;</li> <li>Disabled: not enabled, do not scan requests and skip this rule.</li>.
* `group_id` - (Required, String) Group name of the managed rule. if the rule group for the configuration is not specified, it will be processed based on the default configuration. refer to product documentation for the specific value of GroupId.
* `sensitivity_level` - (Required, String) Protection level of the managed rule group. valid values: <li>loose: lenient, only contains ultra-high risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>normal: normal, contains ultra-high risk and high-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>strict: strict, contains ultra-high risk, high-risk and medium-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>extreme: super strict, contains ultra-high risk, high-risk, medium-risk and low-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>custom: custom, refined strategy. configure the disposal method for each individual rule. at this point, the Action field is invalid. use RuleActions to configure the refined strategy for each individual rule.</li>.
* `meta_data` - (Optional, List) Managed rule group information, for output only.
* `rule_actions` - (Optional, List) Specific configuration of rule items under the managed rule group. the configuration is effective only when SensitivityLevel is custom.

The `managed_rules` object of `security_policy` supports the following:

* `detection_only` - (Required, String) Indicates whether the evaluation mode is Enabled. it is valid only when the Enabled parameter is set to on. valid values: <li>on: Enabled. all managed rules take effect in observation mode.</li> <li>off: disabled. all managed rules take effect according to the actual configuration.</li>.
* `enabled` - (Required, String) Indicates whether the managed rule is enabled. valid values: <li>on: enabled. all managed rules take effect as configured;</li> <li>off: disabled. all managed rules do not take effect.</li>.
* `auto_update` - (Optional, List) Managed rule automatic update option.
* `managed_rule_groups` - (Optional, List) Configuration of the managed rule group. if this structure is passed as an empty array or the GroupId is not included in the list, it will be processed based on the default method.
* `semantic_analysis` - (Optional, String) Whether the managed rule semantic analysis option is Enabled is valid only when the Enabled parameter is on. valid values: <li>on: enable. perform semantic analysis on requests before processing them;</li> <li>off: disable. process requests directly without semantic analysis.</li> <br/>default off.

The `meta_data` object of `managed_rule_groups` supports the following:

* `group_detail` - (Optional, String) Managed rule group description, for output only.
* `group_name` - (Optional, String) Managed rule group name, for output only.
* `rule_details` - (Optional, List) All sub-rule information under the current managed rule group, for output only.

The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) Redirect URL.

The `return_custom_page_action_parameters` object of `action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response status code.

The `rule_actions` object of `managed_rule_groups` supports the following:

* `action` - (Required, List) Specify the handling action for the managed rule item in RuleId. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process the request and record the security event in logs;</li> <li>Disabled: Disabled, do not scan the request and skip this rule.</li>.
* `rule_id` - (Required, String) Specific items under the managed rule group, which are used to rewrite the configuration content of this individual rule item. refer to product documentation for details.

The `rule_details` object of `meta_data` supports the following:

* `description` - (Optional, String) Rule description.
* `risk_level` - (Optional, String) Protection level of managed rules. valid values: <li>low: low risk. this rule has a relatively low risk and is applicable to access scenarios in a very strict control environment. this level of rule may generate considerable false alarms.</li> <li>medium: medium risk. this means the risk of this rule is normal and it is suitable for protection scenarios with stricter requirements.</li> <li>high: high risk. this indicates that the risk of this rule is relatively high and it will not generate false alarms in most scenarios.</li> <li>extreme: ultra-high risk. this represents that the risk of this rule is extremely high and it will not generate false alarms basically.</li>.
* `rule_id` - (Optional, String) Managed rule Id.
* `rule_version` - (Optional, String) Rule ownership version.
* `tags` - (Optional, Set) Rule tag. some types of rules do not have tags.

The `rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.
* `condition` - (Required, String) The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.
* `enabled` - (Required, String) Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `name` - (Required, String) The name of the custom rule.
* `id` - (Optional, String) The ID of a custom rule. <br> the rule ID supports different rule configuration operations: <br> - add a new rule: ID is empty or the ID parameter is not specified; <br> - modify an existing rule: specify the rule ID that needs to be updated/modified; <br> - delete an existing rule: existing Rules not included in the Rules list of the CustomRules parameter will be deleted.
* `priority` - (Optional, Int) Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports PreciseMatchRule.
* `rule_type` - (Optional, String) Type of custom rule. valid values: <li>BasicAccessRule: basic access control;</li> <li>PreciseMatchRule: exact matching rule, default;</li> <li>ManagedAccessRule: expert customized rule, for output only.</li> the default value is PreciseMatchRule.

The `security_policy` object supports the following:

* `custom_rules` - (Optional, List) Custom rule configuration.
* `managed_rules` - (Optional, List) Managed rule configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo teo_security_policy can be imported using the id, e.g.

```
terraform import tencentcloud_teo_security_policy.teo_security_policy teo_security_policy_id
```

