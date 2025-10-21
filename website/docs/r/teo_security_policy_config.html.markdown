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
      precise_match_rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['test']"
        enabled   = "on"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      basic_access_rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
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

    http_ddos_protection {
      adaptive_frequency_control {
        enabled     = "on"
        sensitivity = "Loose"
        action {
          name = "Challenge"
          challenge_action_parameters {
            challenge_option = "JSChallenge"
          }
        }
      }

      client_filtering {
        enabled = "on"
        action {
          name = "Challenge"
          challenge_action_parameters {
            challenge_option = "JSChallenge"
          }
        }
      }

      bandwidth_abuse_defense {
        enabled = "on"
        action {
          name = "Deny"
        }
      }

      slow_attack_defense {
        enabled = "on"
        action {
          name = "Deny"
        }

        minimal_request_body_transfer_rate {
          minimal_avg_transfer_rate_threshold = "80bps"
          counting_period                     = "60s"
          enabled                             = "on"
        }

        request_body_transfer_timeout {
          idle_timeout = "5s"
          enabled      = "on"
        }
      }
    }

    rate_limiting_rules {
      rules {
        name                  = "Single IP request rate limit"
        condition             = "$${http.request.uri.path} contain ['/checkout/submit']"
        count_by              = ["http.request.ip"]
        max_request_threshold = 300
        counting_period       = "60s"
        action_duration       = "30m"
        action {
          name = "Challenge"
          challenge_action_parameters {
            challenge_option = "JSChallenge"
          }
        }
        priority = 50
        enabled  = "on"
      }
    }

    exception_rules {
      rules {
        name                               = "High-frequency API bypasses rate limits"
        condition                          = "$${http.request.method} in ['POST'] and $${http.request.uri.path} in ['/api/EventLogUpload']"
        skip_scope                         = "WebSecurityModules"
        skip_option                        = "SkipOnAllRequestFields"
        web_security_modules_for_exception = ["websec-mod-adaptive-control"]
        enabled                            = "off"
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
      precise_match_rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      basic_access_rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
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
      precise_match_rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      basic_access_rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
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

The `action` object of `adaptive_frequency_control` supports the following:

* `name` - (Required, String) The specific action of security execution. The values are:
<li>Deny: intercept, block the request to access site resources;</li>
<li>Monitor: observe, only record logs;</li>
<li>Redirect: redirect to URL;</li>
<li>Disabled: disabled, do not enable the specified rule;</li>
<li>Allow: allow access, but delay processing requests;</li>
<li>Challenge: challenge, respond to challenge content;</li>
<li>BlockIP: to be abandoned, IP ban;</li>
<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>
<li>JSChallenge: to be abandoned, JavaScript challenge;</li>
<li>ManagedChallenge: to be abandoned, managed challenge.</li>.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `bandwidth_abuse_defense` supports the following:

* `name` - (Required, String) The specific action of security execution. The values are:
<li>Deny: intercept, block the request to access site resources;</li>
<li>Monitor: observe, only record logs;</li>
<li>Redirect: redirect to URL;</li>
<li>Disabled: disabled, do not enable the specified rule;</li>
<li>Allow: allow access, but delay processing requests;</li>
<li>Challenge: challenge, respond to challenge content;</li>
<li>BlockIP: to be abandoned, IP ban;</li>
<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>
<li>JSChallenge: to be abandoned, JavaScript challenge;</li>
<li>ManagedChallenge: to be abandoned, managed challenge.</li>.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `basic_access_rules` supports the following:

* `name` - (Required, String) Specific actions for safe execution. valid values:.
<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.
* `block_ip_action_parameters` - (Optional, List) Additional parameter when Name is BlockIP.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) Additional parameter when Name is ReturnCustomPage.

The `action` object of `client_filtering` supports the following:

* `name` - (Required, String) The specific action of security execution. The values are:
<li>Deny: intercept, block the request to access site resources;</li>
<li>Monitor: observe, only record logs;</li>
<li>Redirect: redirect to URL;</li>
<li>Disabled: disabled, do not enable the specified rule;</li>
<li>Allow: allow access, but delay processing requests;</li>
<li>Challenge: challenge, respond to challenge content;</li>
<li>BlockIP: to be abandoned, IP ban;</li>
<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>
<li>JSChallenge: to be abandoned, JavaScript challenge;</li>
<li>ManagedChallenge: to be abandoned, managed challenge.</li>.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `managed_rule_groups` supports the following:

* `name` - (Required, String) Specific actions for safe execution. valid values:.
<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: use specified page to block</li> <li>Redirect: Redirect to URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>Disabled: Disabled</li> <li>Allow: Allow</li>.
* `block_ip_action_parameters` - (Optional, List) Additional parameter when Name is BlockIP.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) Additional parameter when Name is ReturnCustomPage.

The `action` object of `precise_match_rules` supports the following:

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

The `action` object of `rules` supports the following:

* `name` - (Required, String) The specific action of security execution. The values are:
<li>Deny: intercept, block the request to access site resources;</li>
<li>Monitor: observe, only record logs;</li>
<li>Redirect: redirect to URL;</li>
<li>Disabled: disabled, do not enable the specified rule;</li>
<li>Allow: allow access, but delay processing requests;</li>
<li>Challenge: challenge, respond to challenge content;</li>
<li>BlockIP: to be abandoned, IP ban;</li>
<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>
<li>JSChallenge: to be abandoned, JavaScript challenge;</li>
<li>ManagedChallenge: to be abandoned, managed challenge.</li>.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `slow_attack_defense` supports the following:

* `name` - (Required, String) The specific action of security execution. The values are:
<li>Deny: intercept, block the request to access site resources;</li>
<li>Monitor: observe, only record logs;</li>
<li>Redirect: redirect to URL;</li>
<li>Disabled: disabled, do not enable the specified rule;</li>
<li>Allow: allow access, but delay processing requests;</li>
<li>Challenge: challenge, respond to challenge content;</li>
<li>BlockIP: to be abandoned, IP ban;</li>
<li>ReturnCustomPage: to be abandoned, use the specified page to intercept;</li>
<li>JSChallenge: to be abandoned, JavaScript challenge;</li>
<li>ManagedChallenge: to be abandoned, managed challenge.</li>.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `adaptive_frequency_control` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether adaptive frequency control is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The handling method of adaptive frequency control. When Enabled is on, this field is required. SecurityAction's Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.
* `sensitivity` - (Optional, String) The restriction level of adaptive frequency control. When Enabled is on, this field is required. The values are: <li>Loose: loose; </li><li>Moderate: moderate; </li><li>Strict: strict. </li>.

The `auto_update` object of `managed_rules` supports the following:

* `auto_update_to_latest_version` - (Required, String) Indicates whether to enable automatic update to the latest version. valid values: <li>on: enabled</li> <li>off: disabled</li>.

The `bandwidth_abuse_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether the anti-theft feature (only applicable to mainland China) is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The method for preventing traffic fraud (only applicable to mainland China). When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.

The `basic_access_rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.
* `condition` - (Required, String) The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.
* `enabled` - (Required, String) Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `name` - (Required, String) The name of the custom rule.
* `priority` - (Optional, Int) Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports `rule_type` is `PreciseMatchRule`.

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.

The `challenge_action_parameters` object of `action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.
* `attester_id` - (Optional, String) Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.

The `client_filtering` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether smart client filtering is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The method of intelligent client filtering. When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.

The `custom_rules` object of `security_policy` supports the following:

* `basic_access_rules` - (Optional, List) List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.
* `precise_match_rules` - (Optional, List) List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.
* `rules` - (Optional, List, **Deprecated**) It has been deprecated from version 1.81.184. Please use `precise_match_rules` or `basic_access_rules` instead. List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.

The `deny_action_parameters` object of `action` supports the following:

* `block_ip_duration` - (Optional, String) When BlockIP is on, the IP blocking duration.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. The possible values are:
<li>on: on;</li>
<li>off: off.</li>
When enabled, the client IP that triggers the rule will be blocked continuously. When this option is enabled, the BlockIpDuration parameter must be specified at the same time.
Note: This option cannot be enabled at the same time as the ReturnCustomPage or Stall options.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. The possible values are:
<li>on: on;</li>
<li>off: off.</li>
After enabling, use custom page content to intercept (respond to) requests. When enabling this option, you must specify the ResponseCode and ErrorPageId parameters at the same time.
Note: This option cannot be enabled at the same time as the BlockIp or Stall options.
* `stall` - (Optional, String) Whether to ignore the request source suspension. The value is:
<li>on: Enable;</li>
<li>off: Disable.</li>
After enabling, it will no longer respond to requests in the current connection session and will not actively disconnect. It is used to fight against crawlers and consume client connection resources.
Note: This option cannot be enabled at the same time as the BlockIp or ReturnCustomPage options.

The `exception_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) Definition list of exception rules. When using ModifySecurityPolicy to modify the Web protection configuration: <li>If the Rules parameter is not specified, or the length of the Rules parameter is zero: clear all exception rule configurations. </li>.<li>If the ExceptionRules parameter value is not specified in the SecurityPolicy parameter: keep the existing exception rule configurations and do not modify them. </li>.

The `http_ddos_protection` object of `security_policy` supports the following:

* `adaptive_frequency_control` - (Optional, List) Specific configuration of adaptive frequency control.
* `bandwidth_abuse_defense` - (Optional, List) Specific configuration of traffic fraud prevention.
* `client_filtering` - (Optional, List) Specific configuration of intelligent client filtering.
* `slow_attack_defense` - (Optional, List) Specific configuration of slow attack protection.

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


The `minimal_request_body_transfer_rate` object of `slow_attack_defense` supports the following:

* `counting_period` - (Required, String) The minimum text transmission rate statistics time range, the possible values are: <li>10s: 10 seconds; </li><li>30s: 30 seconds; </li><li>60s: 60 seconds; </li><li>120s: 120 seconds. </li>.
* `enabled` - (Required, String) Whether the text transmission minimum rate threshold is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `minimal_avg_transfer_rate_threshold` - (Required, String) Minimum text transmission rate threshold. The unit only supports bps.

The `precise_match_rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.
* `condition` - (Required, String) The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.
* `enabled` - (Required, String) Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `name` - (Required, String) The name of the custom rule.
* `priority` - (Optional, Int) Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports `rule_type` is `PreciseMatchRule`.

The `rate_limiting_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) A list of precise rate limiting definitions. When using ModifySecurityPolicy to modify the Web protection configuration: <br> <li> If the Rules parameter is not specified, or the Rules parameter length is zero: clear all precise rate limiting configurations. </li>. <li> If the RateLimitingRules parameter value is not specified in the SecurityPolicy parameter: keep the existing custom rule configuration and do not modify it. </li>.

The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `request_body_transfer_timeout` object of `slow_attack_defense` supports the following:

* `enabled` - (Required, String) Whether the text transmission timeout is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `idle_timeout` - (Required, String) The text transmission timeout period is between 5 and 120, and the unit only supports seconds (s).

The `request_fields_for_exception` object of `rules` supports the following:

* `condition` - (Required, String) The expression of the specific field to be skipped must conform to the expression syntax. <br />
Condition supports expression configuration syntax: <li> Written according to the matching condition expression syntax of the rule, supporting references to key and value. </li>.<li> Supports in, like operators, and and logical combinations. </li>.
For example: <li>${key} in ['x-trace-id']: parameter name is equal to x-trace-id. </li>.<li>${key} in ['x-trace-id'] and ${value} like ['Bearer *']: parameter name is equal to x-trace-id and the parameter value wildcard matches Bearer *. </li>.
* `scope` - (Required, String) Specific fields to skip. Supported values:<br/>
<li>body.json: JSON request content; in this case, Condition supports key and value, and TargetField supports key and value, for example, { "Scope": "body.json", "Condition": "", "TargetField": "key" }, which means that all parameters of JSON request content skip WAF scanning;</li>
<li style="margin-top:5px">cookie: Cookie; in this case, Condition supports key and value, and TargetField supports key and value, for example, { "Scope": "cookie", "Condition": "${key} in ['account-id'] and ${value} like ['prefix-*']", "TargetField": "value" }, which means that the Cookie parameter name is equal to account-id and the parameter value wildcard matches prefix-* to skip WAF scanning;</li>
<li style="margin-top:5px">header: HTTP header parameter; Condition supports key and value, TargetField supports key and value, for example { "Scope": "header", "Condition": "${key} like ['x-auth-*']", "TargetField": "value" }, which means that the header parameter name wildcard matches x-auth-* and skips WAF scanning; </li>
<li style="margin-top:5px">uri.query: URL encoded content/query parameter; Condition supports key and value, TargetField supports key and value, for example { "Scope": "uri.query", "Condition": "${key} in ['action'] and ${value} in ['upload', 'delete']", "TargetField": "value" }, which means that the parameter name of the URL encoded content/query parameter is equal to action And the parameter value is equal to upload or delete to skip WAF scanning;</li>
<li style="margin-top:5px">uri: request path URI; in this case, Condition must be empty, TargetField supports query, path, fullpath, for example, { "Scope": "uri", "Condition": "", "TargetField": "query" }, indicating that the request path URI only query parameters skip WAF scanning;</li>
<li style="margin-top:5px">body: request body content. In this case, Condition must be empty, TargetField supports fullbody and multipart, for example, { "Scope": "body", "Condition": "", "TargetField": "fullbody" }, indicating that the request body content is the complete request body and skips WAF scanning;</li>.
* `target_field` - (Required, String) When the Scope parameter uses different values, the supported values in the TargetField expression are as follows:
<li> body.json: supports key and value</li>
<li> cookie: supports key and value</li>
<li> header: supports key and value</li>
<li> uri.query: supports key and value</li>
<li> uri: supports path, query and fullpath</li>
<li> body: supports fullbody and multipart</li>.

The `return_custom_page_action_parameters` object of `action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response status code.

The `return_custom_page_action_parameters` object of `action` supports the following:

* `error_page_id` - (Required, String) The custom page ID of the response.
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

The `rules` object of `exception_rules` supports the following:

* `condition` - (Optional, String) The specific content of the exception rule must comply with the expression syntax. For detailed specifications, see the product documentation.
* `enabled` - (Optional, String) Whether the exception rule is enabled. The values are: <li>on: enabled</li><li>off: disabled</li>.
* `id` - (Optional, String) The ID of the exception rule. <br>The rule ID can support different rule configuration operations: <br> <li> <b>Add</b> a new rule: the ID is empty or the ID parameter is not specified; </li><li> <b>Modify</b> an existing rule: specify the rule ID to be updated/modified; </li><li> <b>Delete</b> an existing rule: in the ExceptionRules parameter, the existing rules not included in the Rules list will be deleted. </li>.
* `managed_rule_groups_for_exception` - (Optional, Set) Specifies the managed rule group for the exception rule. This is only valid when SkipScope is ManagedRules and ManagedRulesForException cannot be specified.
* `managed_rules_for_exception` - (Optional, Set) Specifies the specific managed rule for the exception rule. This is only valid when SkipScope is ManagedRules and ManagedRuleGroupsForException cannot be specified.
* `name` - (Optional, String) The name of the exception rule.
* `request_fields_for_exception` - (Optional, List) Specifies the specific configuration of the exception rule to skip the specified request field. This is only valid when SkipScope is ManagedRules and SkipOption is SkipOnSpecifiedRequestFields.
* `skip_option` - (Optional, String) The specific type of the skipped request. The possible values are: <li>SkipOnAllRequestFields: skip all requests; </li><li>SkipOnSpecifiedRequestFields: skip specified request fields. </li>. This option is only valid when SkipScope is ManagedRules.
* `skip_scope` - (Optional, String) Exception rule execution options, the values are: <li>WebSecurityModules: Specifies the security protection module for the exception rule. </li>.<li>ManagedRules: Specifies the managed rules. </li>.
* `web_security_modules_for_exception` - (Optional, Set) Specifies the security protection module for the exception rule. It is valid only when SkipScope is WebSecurityModules. The possible values are: <li>websec-mod-managed-rules: managed rules; </li><li>websec-mod-rate-limiting: rate limiting; </li><li>websec-mod-custom-rules: custom rules; </li><li>websec-mod-adaptive-control: adaptive frequency control, intelligent client filtering, slow attack protection, traffic theft protection; </li><li>websec-mod-bot: Bot management. </li>.

The `rules` object of `rate_limiting_rules` supports the following:

* `action_duration` - (Optional, String) Action The duration of the action. The supported units are: <li>s: seconds, with a value of 1 to 120; </li><li>m: minutes, with a value of 1 to 120; </li><li>h: hours, with a value of 1 to 48; </li><li>d: days, with a value of 1 to 30. </li>.
* `action` - (Optional, List) The precise rate limit handling method. The values are: <li>Monitor: Observe; </li><li>Deny: Intercept, where DenyActionParameters.Name supports Deny and ReturnCustomPage; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name supports JSChallenge and ManagedChallenge; </li><li>Redirect: Redirect to URL; </li>.
* `condition` - (Optional, String) The specific content of the precise rate limit must conform to the expression syntax. For detailed specifications, see the product documentation.
* `count_by` - (Optional, Set) The matching method of the rate threshold request feature. When Enabled is on, this field is required. <br /><br />When there are multiple conditions, multiple conditions will be combined for statistical calculation. The number of conditions cannot exceed 5. The possible values are: <br/><li><b>http.request.ip</b>: client IP; </li><li><b>http.request.xff_header_ip</b>: client IP (matching XFF header first); </li><li><b>http.request.uri.path</b>: requested access path; </li><li><b>http.request.cookies['session']</b>: cookie named session, where session can be replaced by the parameter you specify; </li><li><b>http.request.headers['user-agent']</b>: HTTP header named user-agent, where user-agent can be replaced by the parameter you specify; </li><li><b>http.request.ja3</b>: requested JA3 fingerprint; </li><li><b>http.request.uri.query['test']</b>: URL query parameter named test, where test can be replaced by the parameter you specify. </li>.
* `counting_period` - (Optional, String) The statistical time window, the possible values are: <li>1s: 1 second; </li><li>5s: 5 seconds; </li><li>10s: 10 seconds; </li><li>20s: 20 seconds; </li><li>30s: 30 seconds; </li><li>40s: 40 seconds; </li><li>50s: 50 seconds; </li><li>1m: 1 minute; </li><li>2m: 2 minutes; </li><li>5m: 5 minutes; </li><li>10m: 10 minutes; </li><li>1h: 1 hour. </li>.
* `enabled` - (Optional, String) Whether the precise rate limit rule is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `id` - (Optional, String) The ID of the precise rate limit. <br>The rule ID can support different rule configuration operations: <br> <li> <b>Add</b> a new rule: the ID is empty or the ID parameter is not specified; </li><li> <b>Modify</b> an existing rule: specify the rule ID to be updated/modified; </li><li> <b>Delete</b> an existing rule: in the RateLimitingRules parameter, the existing rules not included in the Rules list will be deleted. </li>.
* `max_request_threshold` - (Optional, Int) The cumulative number of interceptions within the time range of the precise rate limit, ranging from 1 to 100000.
* `name` - (Optional, String) The name of the precise rate limit.
* `priority` - (Optional, Int) The priority of precise rate limiting ranges from 0 to 100, and the default is 0.

The `security_policy` object supports the following:

* `custom_rules` - (Optional, List) Custom rule configuration.
* `exception_rules` - (Optional, List) Exception rule configuration.
* `http_ddos_protection` - (Optional, List) HTTP DDOS protection configuration.
* `managed_rules` - (Optional, List) Managed rule configuration.
* `rate_limiting_rules` - (Optional, List) Rate limiting rule configuration.

The `slow_attack_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether slow attack protection is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The handling method of slow attack protection. When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li>.
* `minimal_request_body_transfer_rate` - (Optional, List) Specific configuration of the minimum rate threshold for text transmission. This field is required when Enabled is on.
* `request_body_transfer_timeout` - (Optional, List) Specific configuration of the text transmission timeout. When Enabled is on, this field is required.

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

