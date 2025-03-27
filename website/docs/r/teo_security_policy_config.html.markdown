---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_policy_config"
sidebar_current: "docs-tencentcloud-resource-teo_security_policy_config"
description: |-
  Provides a resource to create a teo security policy
---

# tencentcloud_teo_security_policy_config

Provides a resource to create a teo security policy

~> **NOTE:** If the user's EO version is the personal version, `managed_rule_groups` needs to set one; If the user's EO version is a non personal version, `managed_rule_groups` needs to set 17. If the user does not set the `managed_rule_groups` parameter, the system will generate it by default.

## Example Usage

### If entity is ZoneDefaultPolicy

```hcl
resource "tencentcloud_teo_security_policy_config" "example" {
  zone_id = "zone-37u62pwxfo8s"
  entity  = "ZoneDefaultPolicy"
  security_policy {
    custom_rules {
      rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        rule_type = "PreciseMatchRule"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
        id        = "2182252647"
        rule_type = "BasicAccessRule"
        action {
          name = "Deny"
        }
      }
    }

    managed_rules {
      enabled           = "on"
      detection_only    = "off"
      semantic_analysis = "off"
      auto_update {
        auto_update_to_latest_version = "off"
      }

      managed_rule_groups {
        group_id          = "wafgroup-webshell-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xxe-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-non-compliant-protocol-usages"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-file-upload-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-command-and-code-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ldap-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssrf-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xss-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-vulnerability-scanners"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-cms-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-other-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-sql-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-file-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-oa-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssti-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-shiro-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }
    }
  }
}
```

### If entity is Host

```hcl
resource "tencentcloud_teo_security_policy_config" "example" {
  zone_id = "zone-37u62pwxfo8s"
  entity  = "Host"
  host    = "www.example.com"
  security_policy {
    custom_rules {
      rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        rule_type = "PreciseMatchRule"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
        id        = "2182252647"
        rule_type = "BasicAccessRule"
        action {
          name = "Deny"
        }
      }
    }

    managed_rules {
      enabled           = "on"
      detection_only    = "off"
      semantic_analysis = "off"
      auto_update {
        auto_update_to_latest_version = "off"
      }

      managed_rule_groups {
        group_id          = "wafgroup-webshell-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xxe-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-non-compliant-protocol-usages"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-file-upload-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-command-and-code-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ldap-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssrf-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xss-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-vulnerability-scanners"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-cms-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-other-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-sql-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-file-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-oa-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssti-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-shiro-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }
    }
  }
}
```

### If entity is Template

```hcl
resource "tencentcloud_teo_security_policy_config" "example" {
  zone_id     = "zone-37u62pwxfo8s"
  entity      = "Template"
  template_id = "temp-05dtxkyw"
  security_policy {
    custom_rules {
      rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        rule_type = "PreciseMatchRule"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
        id        = "2182252647"
        rule_type = "BasicAccessRule"
        action {
          name = "Deny"
        }
      }
    }

    managed_rules {
      enabled           = "on"
      detection_only    = "off"
      semantic_analysis = "off"
      auto_update {
        auto_update_to_latest_version = "off"
      }

      managed_rule_groups {
        group_id          = "wafgroup-webshell-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xxe-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-non-compliant-protocol-usages"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-file-upload-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-command-and-code-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ldap-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssrf-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xss-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-vulnerability-scanners"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-cms-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-other-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-sql-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-file-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-oa-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssti-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-shiro-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
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

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.

The `custom_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.

The `managed_rule_groups` object of `managed_rules` supports the following:

* `action` - (Required, List) Handling actions for managed rule groups. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process requests and record security events in logs;</li> <li>Disabled: not enabled, do not scan requests and skip this rule.</li>.
* `group_id` - (Required, String) Group name of the managed rule. if the rule group for the configuration is not specified, it will be processed based on the default configuration. refer to product documentation for the specific value of GroupId.
* `sensitivity_level` - (Required, String) Protection level of the managed rule group. valid values: <li>loose: lenient, only contains ultra-high risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>normal: normal, contains ultra-high risk and high-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>strict: strict, contains ultra-high risk, high-risk and medium-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>extreme: super strict, contains ultra-high risk, high-risk, medium-risk and low-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>custom: custom, refined strategy. configure the disposal method for each individual rule. at this point, the Action field is invalid. use RuleActions to configure the refined strategy for each individual rule.</li>.
* `rule_actions` - (Optional, List) Specific configuration of rule items under the managed rule group. the configuration is effective only when SensitivityLevel is custom.

The `managed_rules` object of `security_policy` supports the following:

* `detection_only` - (Required, String) Indicates whether the evaluation mode is Enabled. it is valid only when the Enabled parameter is set to on. valid values: <li>on: Enabled. all managed rules take effect in observation mode.</li> <li>off: disabled. all managed rules take effect according to the actual configuration.</li>.
* `enabled` - (Required, String) Indicates whether the managed rule is enabled. valid values: <li>on: enabled. all managed rules take effect as configured;</li> <li>off: disabled. all managed rules do not take effect.</li>.
* `auto_update` - (Optional, List) Managed rule automatic update option.
* `managed_rule_groups` - (Optional, Set) Configuration of the managed rule group. if this structure is passed as an empty array or the GroupId is not included in the list, it will be processed based on the default method.
* `semantic_analysis` - (Optional, String) Whether the managed rule semantic analysis option is Enabled is valid only when the Enabled parameter is on. valid values: <li>on: enable. perform semantic analysis on requests before processing them;</li> <li>off: disable. process requests directly without semantic analysis.</li> <br/>default off.

The `meta_data` object of `managed_rule_groups` supports the following:


The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) Redirect URL.

The `return_custom_page_action_parameters` object of `action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response status code.

The `rule_actions` object of `managed_rule_groups` supports the following:

* `action` - (Required, List) Specify the handling action for the managed rule item in RuleId. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process the request and record the security event in logs;</li> <li>Disabled: Disabled, do not scan the request and skip this rule.</li>.
* `rule_id` - (Required, String) Specific items under the managed rule group, which are used to rewrite the configuration content of this individual rule item. refer to product documentation for details.

The `rule_details` object of `meta_data` supports the following:


The `rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.
* `condition` - (Required, String) The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.
* `enabled` - (Required, String) Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `name` - (Required, String) The name of the custom rule.
* `id` - (Optional, String) The ID of a custom rule. <br> the rule ID supports different rule configuration operations: <br> - add a new rule: ID is empty or the ID parameter is not specified; <br> - modify an existing rule: specify the rule ID that needs to be updated/modified; <br> - delete an existing rule: existing Rules not included in the Rules list of the CustomRules parameter will be deleted.
* `priority` - (Optional, Int) Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports `rule_type` is `PreciseMatchRule`.
* `rule_type` - (Optional, String) Type of custom rule. valid values: <li>BasicAccessRule: basic access control;</li> <li>PreciseMatchRule: exact matching rule, default;</li> <li>ManagedAccessRule: expert customized rule, for output only.</li> the default value is PreciseMatchRule.

The `security_policy` object supports the following:

* `custom_rules` - (Optional, List) Custom rule configuration.
* `managed_rules` - (Optional, List) Managed rule configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo security policy can be imported using the id, e.g.

```
# If entity is ZoneDefaultPolicy 
terraform import tencentcloud_teo_security_policy_config.example zone-37u62pwxfo8s#ZoneDefaultPolicy
# If entity is Host
terraform import tencentcloud_teo_security_policy_config.example zone-37u62pwxfo8s#Host#www.example.com
# If entity is Template
terraform import tencentcloud_teo_security_policy_config.example zone-37u62pwxfo8s#Template#temp-05dtxkyw
```

