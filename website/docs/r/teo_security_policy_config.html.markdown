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

    bot_management_lite {
      captcha_page_challenge {
        enabled = "on"
      }

      ai_crawler_detection {
        enabled = "on"
        action {
          name = "Deny"
          deny_action_parameters {
            block_ip          = "on"
            block_ip_duration = "120s"
          }
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
* `security_config` - (Optional, List) Security configuration. Classic web protection settings. Note: the DescribeSecurityPolicy API does not return SecurityConfig, so this field is write-only for state consistency. For each sub-configuration, if not specified, the existing API configuration is kept.
* `security_policy` - (Optional, List) Security policy configuration. it is recommended to use for custom policies and managed rule configurations of Web protection. it supports configuring security policies with expression grammar.
* `template_id` - (Optional, String, ForceNew) Specify the policy Template ID. use this parameter to specify the ID of the policy Template when the Entity parameter value is Template.

The `acl_conditions` object of `acl_user_rules` supports the following:

* `match_content` - (Optional, String) Match content.
* `match_from` - (Optional, String) Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`).
* `match_param` - (Optional, String) Match parameter. For `header` MatchFrom, the header key.
* `operator` - (Optional, String) Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`).

The `acl_conditions` object of `bot_user_rules` supports the following:

* `match_content` - (Optional, String) Match content.
* `match_from` - (Optional, String) Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`).
* `match_param` - (Optional, String) Match parameter. For `header` MatchFrom, the header key.
* `operator` - (Optional, String) Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`).

The `acl_conditions` object of `customizes` supports the following:

* `match_content` - (Optional, String) Match content.
* `match_from` - (Optional, String) Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`).
* `match_param` - (Optional, String) Match parameter. For `header` MatchFrom, the header key.
* `operator` - (Optional, String) Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`).

The `acl_conditions` object of `rate_limit_customizes` supports the following:

* `match_content` - (Optional, String) Match content.
* `match_from` - (Optional, String) Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`).
* `match_param` - (Optional, String) Match parameter. For `header` MatchFrom, the header key.
* `operator` - (Optional, String) Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`).

The `acl_conditions` object of `rate_limit_user_rules` supports the following:

* `match_content` - (Optional, String) Match content.
* `match_from` - (Optional, String) Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`).
* `match_param` - (Optional, String) Match parameter. For `header` MatchFrom, the header key.
* `operator` - (Optional, String) Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`).

The `acl_config` object of `security_config` supports the following:

* `acl_user_rules` - (Optional, List) User-defined ACL rules.
* `customizes` - (Optional, List) Managed customized ACL rules.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `acl_drop_page_detail` object of `drop_page_config` supports the following:

* `custom_response_id` - (Optional, String) Custom response ID.
* `name` - (Optional, String) Block page file name or URL.
* `page_id` - (Optional, Int) The unique ID of the block page. The system includes a built-in block page with ID 0.
* `status_code` - (Optional, Int) HTTP status code for the block page. Range: 100-600, excluding 3xx.
* `type` - (Optional, String) Page type. Valid values: `page`.

The `acl_user_rules` object of `acl_config` supports the following:

* `acl_conditions` - (Optional, List) Rule ACL conditions.
* `action` - (Optional, String) Action. Valid values: `trans`, `drop`, `monitor`, `ban`, `redirect`, `page`, `alg`.
* `custom_response_id` - (Optional, String) Custom response ID.
* `name` - (Optional, String) Custom response page name.
* `page_id` - (Optional, Int) Custom page instance ID. Deprecated.
* `punish_time_unit` - (Optional, String) Penalty time unit. Valid values: `second`, `minutes`, `hour`.
* `punish_time` - (Optional, Int) IP ban penalty time.
* `redirect_url` - (Optional, String) Redirect URL.
* `response_code` - (Optional, Int) Custom response code.
* `rule_name` - (Optional, String) Rule name.
* `rule_priority` - (Optional, Int) Rule priority (0-100).
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.

The `action_overrides` object of `ip_reputation_group` supports the following:

* `rule_id` - (Required, String) Rule ID or category ID for action override.
* `action` - (Optional, List) Action override configuration.

The `action_overrides` object of `known_bot_categories` supports the following:

* `rule_id` - (Required, String) Rule ID or category ID for action override.
* `action` - (Optional, List) Action override configuration.

The `action_overrides` object of `search_engine_bots` supports the following:

* `rule_id` - (Required, String) Rule ID or category ID for action override.
* `action` - (Optional, List) Action override configuration.

The `action_overrides` object of `source_idc` supports the following:

* `rule_id` - (Required, String) Rule ID or category ID for action override.
* `action` - (Optional, List) Action override configuration.

The `action` object of `action_overrides` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `action` object of `action` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

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

The `action` object of `ai_crawler_detection` supports the following:

* `name` - (Required, String) The security action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.

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

The `action` object of `browser_impersonation_detection` supports the following:

* `bot_session_validation` - (Optional, List) Cookie validation and session tracking.
* `client_behavior_detection` - (Optional, List) Client behavior detection.

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

The `action` object of `custom_rules` supports the following:

* `action` - (Required, List) Security action configuration.
* `weight` - (Required, Int) Action weight (10-100, must be multiples of 10). Sum of all weights must equal 100.

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

The `ai_crawler_detection` object of `bot_management_lite` supports the following:

* `enabled` - (Required, String) Whether AI crawler detection is enabled. Valid values: `on`, `off`.
* `action` - (Optional, List) Execution action when Enabled is on. When Enabled is on, this field is required. SecurityAction Name value supports: Deny, Monitor, Allow, Challenge.

The `ai_rule` object of `waf_config` supports the following:

* `mode` - (Optional, String) AI rule mode. Valid values: `smart_status_close`, `smart_status_open`, `smart_status_observe`.

The `alg_conditions` object of `alg_detect_rule` supports the following:

* `match_content` - (Optional, String) Match content.
* `match_from` - (Optional, String) Match field. See product doc for valid values (e.g. `host`, `sip`, `ua`, `cookie`, `cgi`, `xff`, `url`, `accept`, `method`, `header`, `app_proto`, `sip_proto`).
* `match_param` - (Optional, String) Match parameter. For `header` MatchFrom, the header key.
* `operator` - (Optional, String) Match operator (e.g. `equal`, `not_equal`, `include`, `regexp`, `match_prefix`, `wildcard`).

The `alg_detect_js` object of `alg_detect_rule` supports the following:

* `execute_mode` - (Optional, Int) JS execution delay in ms (0-1000, default 500).
* `invalid_stat_time` - (Optional, Int) Statistical period for invalid JS (5-3600s, default 10).
* `invalid_threshold` - (Optional, Int) Threshold for invalid JS (1-100000000, default 300).
* `name` - (Optional, String) Operation name.
* `work_level` - (Optional, String) Proof-of-work strength. Valid values: `low`, `middle`, `high` (default `low`).

The `alg_detect_rule` object of `bot_config` supports the following:

* `alg_conditions` - (Optional, List) Custom conditions.
* `alg_detect_js` - (Optional, List) Client behavior detection.
* `alg_detect_session` - (Optional, List) Cookie validation and session behavior analysis.
* `rule_name` - (Optional, String) Rule name.
* `switch` - (Optional, String) Rule switch.

The `alg_detect_session` object of `alg_detect_rule` supports the following:

* `detect_mode` - (Optional, String) Detection mode. Valid values: `detect`, `update_detect`.
* `invalid_stat_time` - (Optional, Int) Statistical period for missing/expired cookie (5-3600s, default 10).
* `invalid_threshold` - (Optional, Int) Trigger threshold for missing/expired cookie (1-100000000, default 300).
* `name` - (Optional, String) Operation name.
* `session_analyze_switch` - (Optional, String) Session behavior analysis switch. Valid values: `off`, `on`.

The `allow_action_parameters` object of `action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `base_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `bot_client_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `challenge_not_finished_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `challenge_timeout_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `high_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `human_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `invalid_attestation_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `likely_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `low_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `mid_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `session_expired_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `session_invalid_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `allow_action_parameters` object of `verified_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delay response time. Supported unit: seconds, range 5-10.
* `min_delay_time` - (Optional, String) Minimum delay response time. Supported unit: seconds, range 0-5.

The `auto_update` object of `managed_rules` supports the following:

* `auto_update_to_latest_version` - (Required, String) Indicates whether to enable automatic update to the latest version. valid values: <li>on: enabled</li> <li>off: disabled</li>.

The `bandwidth_abuse_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether the anti-theft feature (only applicable to mainland China) is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The method for preventing traffic fraud (only applicable to mainland China). When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.

The `base_action` object of `ip_reputation_group` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `base_action` object of `known_bot_categories` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `base_action` object of `search_engine_bots` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `base_action` object of `source_idc` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `basic_access_rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.
* `condition` - (Required, String) The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.
* `enabled` - (Required, String) Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `name` - (Required, String) The name of the custom rule.
* `priority` - (Optional, Int) Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports `rule_type` is `PreciseMatchRule`.

The `basic_bot_settings` object of `bot_management` supports the following:

* `bot_intelligence` - (Optional, List) Bot intelligence configuration.
* `ip_reputation` - (Optional, List) IP reputation configuration.
* `known_bot_categories` - (Optional, List) Known bot categories configuration.
* `search_engine_bots` - (Optional, List) Search engine bots configuration.
* `source_idc` - (Optional, List) Source IDC configuration.

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) Penalty duration for blocking ips. supported units: <li>s: second, value range 1-120;</li> <li>m: minute, value range 1-120;</li> <li>h: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) The penalty duration for banning an IP. Supported units are: <li>s: seconds, value range 1 to 120; </li><li>m: minutes, value range 1 to 120; </li><li>h: hours, value range 1 to 48. </li>.

The `bot_client_action` object of `client_behavior_detection` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `bot_config` object of `security_config` supports the following:

* `alg_detect_rule` - (Optional, List) Bot active feature detection rules.
* `bot_managed_rule` - (Optional, List) Generic bot managed rules.
* `bot_portrait_rule` - (Optional, List) User portrait rule.
* `bot_user_rules` - (Optional, List) Bot user-defined rules.
* `customizes` - (Optional, List) Bot managed customized rules.
* `intelligence_rule` - (Optional, List) Bot intelligence rule.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `bot_intelligence` object of `basic_bot_settings` supports the following:

* `bot_ratings` - (Optional, List) Bot ratings configuration.
* `enabled` - (Optional, String) Whether bot intelligence is enabled. Valid values: `on`, `off`.

The `bot_managed_rule` object of `bot_config` supports the following:

* `action` - (Optional, String) Action. Valid values: `drop`, `trans`, `alg`, `monitor`.
* `alg_managed_ids` - (Optional, List) Rule IDs with JS challenge.
* `cap_managed_ids` - (Optional, List) Rule IDs with CAPTCHA.
* `drop_managed_ids` - (Optional, List) Rule IDs to drop.
* `mon_managed_ids` - (Optional, List) Rule IDs in monitor mode.
* `trans_managed_ids` - (Optional, List) Rule IDs to allow.

The `bot_management_lite` object of `security_policy` supports the following:

* `ai_crawler_detection` - (Optional, List) AI crawler detection configuration.
* `captcha_page_challenge` - (Optional, List) CAPTCHA page challenge configuration.

The `bot_management` object of `security_policy` supports the following:

* `basic_bot_settings` - (Optional, List) Basic bot settings.
* `browser_impersonation_detection` - (Optional, List) Browser impersonation detection rules.
* `client_attestation_rules` - (Optional, List) Client attestation rules (beta feature).
* `custom_rules` - (Optional, List) Bot management custom rules.
* `enabled` - (Optional, String) Whether bot management is enabled. Valid values: `on`, `off`.

The `bot_portrait_rule` object of `bot_config` supports the following:

* `alg_managed_ids` - (Optional, List) Rule IDs with JS challenge.
* `cap_managed_ids` - (Optional, List) Rule IDs with CAPTCHA.
* `drop_managed_ids` - (Optional, List) Rule IDs to drop.
* `mon_managed_ids` - (Optional, List) Rule IDs in monitor mode.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `bot_ratings` object of `bot_intelligence` supports the following:

* `high_risk_bot_requests_action` - (Optional, List) Action for high risk bot requests.
* `human_requests_action` - (Optional, List) Action for human requests.
* `likely_bot_requests_action` - (Optional, List) Action for likely bot requests.
* `verified_bot_requests_action` - (Optional, List) Action for verified bot requests.

The `bot_session_validation` object of `action` supports the following:

* `issue_new_bot_session_cookie` - (Optional, String) Whether to issue new session cookie. Valid values: `on`, `off`.
* `max_new_session_trigger_config` - (Optional, List) Trigger config for new session.
* `session_expired_action` - (Optional, List) Action when session is expired.
* `session_invalid_action` - (Optional, List) Action when session is invalid.
* `session_rate_control` - (Optional, List) Session rate control.

The `bot_user_rules` object of `bot_config` supports the following:

* `acl_conditions` - (Optional, List) Rule ACL conditions.
* `action` - (Optional, String) Action. Valid values: `drop`, `monitor`, `trans`, `redirect`, `page`, `alg`, `captcha`, `random`, `silence`, `shortdelay`, `longdelay`.
* `custom_response_id` - (Optional, String) Custom response ID.
* `extend_actions` - (Optional, List) Random action weighted distribution.
* `freq_fields` - (Optional, List) Filter fields.
* `freq_scope` - (Optional, List) Statistical scope.
* `name` - (Optional, String) Custom response page name.
* `redirect_url` - (Optional, String) Redirect URL.
* `response_code` - (Optional, Int) Custom response code.
* `rule_name` - (Optional, String) Rule name.
* `rule_priority` - (Optional, Int) Rule priority (0-100).
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.

The `browser_impersonation_detection` object of `bot_management` supports the following:

* `condition` - (Required, String) Rule condition in expression syntax.
* `enabled` - (Required, String) Whether the rule is enabled. Valid values: `on`, `off`.
* `name` - (Required, String) Rule name.
* `action` - (Optional, List) Action configuration.
* `id` - (Optional, String) Rule ID.

The `captcha_page_challenge` object of `bot_management_lite` supports the following:

* `enabled` - (Required, String) Whether CAPTCHA page challenge is enabled. Valid values: `on`, `off`.

The `challenge_action_parameters` object of `action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. The possible values are: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>.
* `attester_id` - (Optional, String) Client authentication method ID. This field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) The time interval for repeating the challenge. When Name is InterstitialChallenge/InlineChallenge, this field is required. The default value is 300s. Supported units are: <li>s: seconds, value range 1 to 60; </li><li>m: minutes, value range 1 to 60; </li><li>h: hours, value range 1 to 24. </li>.

The `challenge_action_parameters` object of `action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `base_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `bot_client_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `challenge_not_finished_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `challenge_timeout_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `high_rate_session_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `human_requests_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `invalid_attestation_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `likely_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `low_rate_session_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `mid_rate_session_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `session_expired_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `session_invalid_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_action_parameters` object of `verified_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) The specific challenge action to be executed safely. Valid values: `InterstitialChallenge`, `InlineChallenge`, `JSChallenge`, `ManagedChallenge`.
* `attester_id` - (Optional, String) Client authentication method ID.
* `interval` - (Optional, String) The time interval for repeating the challenge.

The `challenge_not_finished_action` object of `client_behavior_detection` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `challenge_timeout_action` object of `client_behavior_detection` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `client_attestation_rules` object of `bot_management` supports the following:

* `condition` - (Required, String) Rule condition in expression syntax.
* `enabled` - (Required, String) Whether the rule is enabled. Valid values: `on`, `off`.
* `name` - (Required, String) Rule name.
* `attester_id` - (Optional, String) Client authentication method ID.
* `id` - (Optional, String) Rule ID.
* `invalid_attestation_action` - (Optional, List) Action when attestation is invalid.
* `priority` - (Optional, Int) Rule priority (0-100).

The `client_behavior_detection` object of `action` supports the following:

* `bot_client_action` - (Optional, List) Action for bot client.
* `challenge_not_finished_action` - (Optional, List) Action when challenge not finished.
* `challenge_timeout_action` - (Optional, List) Action when challenge timeout.
* `crypto_challenge_delay_before` - (Optional, String) Challenge delay before execution. Valid values: `0ms`, `100ms`, `200ms`, `300ms`, `400ms`, `500ms`, `600ms`, `700ms`, `800ms`, `900ms`, `1000ms`.
* `crypto_challenge_intensity` - (Optional, String) Proof-of-work challenge intensity. Valid values: `low`, `medium`, `high`.
* `max_challenge_count_interval` - (Optional, String) Statistics time window for trigger threshold. Valid values: `5s`, `10s`, `15s`, `30s`, `60s`, `5m`, `10m`, `30m`, `60m`.
* `max_challenge_count_threshold` - (Optional, Int) Cumulative count for trigger threshold. Range: 1-100000000.

The `client_filtering` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether smart client filtering is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The method of intelligent client filtering. When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge. </li>.

The `conditions` object of `detect_length_limit_rules` supports the following:


The `custom_rules` object of `bot_management` supports the following:

* `condition` - (Required, String) Rule condition in expression syntax.
* `enabled` - (Required, String) Whether the rule is enabled. Valid values: `on`, `off`.
* `name` - (Required, String) Rule name.
* `action` - (Optional, List) Rule action with weight.
* `id` - (Optional, String) Rule ID. If not specified, a new rule will be created. If specified, the existing rule will be updated or deleted.
* `priority` - (Optional, Int) Rule priority (0-100).

The `custom_rules` object of `security_policy` supports the following:

* `basic_access_rules` - (Optional, List) List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.
* `precise_match_rules` - (Optional, List) List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.
* `rules` - (Optional, List, **Deprecated**) It has been deprecated from version 1.81.184. Please use `precise_match_rules` or `basic_access_rules` instead. List of custom rule definitions. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the parameter value of CustomRules in the SecurityPolicy parameter is not specified: keep the existing custom rule configuration without modification.

The `customizes` object of `acl_config` supports the following:

* `acl_conditions` - (Optional, List) Rule ACL conditions.
* `action` - (Optional, String) Action. Valid values: `trans`, `drop`, `monitor`, `ban`, `redirect`, `page`, `alg`.
* `custom_response_id` - (Optional, String) Custom response ID.
* `name` - (Optional, String) Custom response page name.
* `page_id` - (Optional, Int) Custom page instance ID. Deprecated.
* `punish_time_unit` - (Optional, String) Penalty time unit. Valid values: `second`, `minutes`, `hour`.
* `punish_time` - (Optional, Int) IP ban penalty time.
* `redirect_url` - (Optional, String) Redirect URL.
* `response_code` - (Optional, Int) Custom response code.
* `rule_name` - (Optional, String) Rule name.
* `rule_priority` - (Optional, Int) Rule priority (0-100).
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.

The `customizes` object of `bot_config` supports the following:

* `acl_conditions` - (Optional, List) Rule ACL conditions.
* `action` - (Optional, String) Action. Valid values: `drop`, `monitor`, `trans`, `redirect`, `page`, `alg`, `captcha`, `random`, `silence`, `shortdelay`, `longdelay`.
* `custom_response_id` - (Optional, String) Custom response ID.
* `extend_actions` - (Optional, List) Random action weighted distribution.
* `freq_fields` - (Optional, List) Filter fields.
* `freq_scope` - (Optional, List) Statistical scope.
* `name` - (Optional, String) Custom response page name.
* `redirect_url` - (Optional, String) Redirect URL.
* `response_code` - (Optional, Int) Custom response code.
* `rule_name` - (Optional, String) Rule name.
* `rule_priority` - (Optional, Int) Rule priority (0-100).
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.

The `default_deny_security_action_parameters` object of `security_policy` supports the following:

* `managed_rules` - (Optional, List) Managed rules default deny action configuration. Supported parameters: `return_custom_page`, `response_code`, `error_page_id`.
* `other_modules` - (Optional, List) Default deny action configuration for security protection rules other than managed rules (custom rules, rate limiting and Bot management). Supported parameters: `return_custom_page`, `response_code`, `error_page_id`.

The `deny_action_parameters` object of `action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

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

The `deny_action_parameters` object of `action` supports the following:

* `block_ip_duration` - (Optional, String) When BlockIP is on, the IP blocking duration.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `base_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `bot_client_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `challenge_not_finished_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `challenge_timeout_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `high_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `human_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `invalid_attestation_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `likely_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `low_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `mid_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `session_expired_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `session_invalid_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `deny_action_parameters` object of `verified_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) IP blocking duration when BlockIP is on.
* `block_ip` - (Optional, String) Whether to extend the blocking of source IP. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) The PageId of the custom page.
* `response_code` - (Optional, String) Customize the status code of the page.
* `return_custom_page` - (Optional, String) Whether to use custom pages. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to ignore the request source suspension. Valid values: `on`, `off`.

The `detect_length_limit_config` object of `security_config` supports the following:


The `detect_length_limit_rules` object of `detect_length_limit_config` supports the following:


The `drop_page_config` object of `security_config` supports the following:

* `acl_drop_page_detail` - (Optional, List) Custom rule drop page.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.
* `waf_drop_page_detail` - (Optional, List) Managed rule drop page.

The `except_config` object of `security_config` supports the following:

* `except_user_rules` - (Optional, List) Exception rules detail.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `except_user_rule_conditions` object of `except_user_rules` supports the following:

* `match_content` - (Optional, String) Match value.
* `match_from` - (Optional, String) Match field.
* `match_param` - (Optional, String) Match parameter (e.g. header key when MatchFrom=header).
* `operator` - (Optional, String) Operator.

The `except_user_rule_scope` object of `except_user_rules` supports the following:

* `modules` - (Optional, List) Effective modules. Valid values: `waf`, `rate`, `acl`, `cc`, `bot`.
* `partial_modules` - (Optional, List) Partial rule ID exceptions.
* `skip_conditions` - (Optional, List) Conditions to skip.
* `type` - (Optional, String) Scope type. Valid values: `complete`, `partial`.

The `except_user_rules` object of `except_config` supports the following:

* `action` - (Optional, String) Rule action. Only `skip` is supported.
* `except_user_rule_conditions` - (Optional, List) Match conditions.
* `except_user_rule_scope` - (Optional, List) Rule effective scope.
* `rule_name` - (Optional, String) Rule name (no Chinese characters).
* `rule_priority` - (Optional, Int) Priority (0-100). Default 0.
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.

The `exception_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) Definition list of exception rules. When using ModifySecurityPolicy to modify the Web protection configuration: <li>If the Rules parameter is not specified, or the length of the Rules parameter is zero: clear all exception rule configurations. </li>.<li>If the ExceptionRules parameter value is not specified in the SecurityPolicy parameter: keep the existing exception rule configurations and do not modify them. </li>.

The `extend_actions` object of `bot_user_rules` supports the following:

* `action` - (Optional, String) Action. Valid values: `monitor`, `alg`, `captcha`, `random`, `silence`, `shortdelay`, `longdelay`.
* `percent` - (Optional, Int) Action probability (0-100).

The `extend_actions` object of `customizes` supports the following:

* `action` - (Optional, String) Action. Valid values: `monitor`, `alg`, `captcha`, `random`, `silence`, `shortdelay`, `longdelay`.
* `percent` - (Optional, Int) Action probability (0-100).

The `first_part_config` object of `slow_post_config` supports the following:

* `stat_time` - (Optional, Int) First segment statistical duration in seconds (default 5).
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `high_rate_session_action` object of `session_rate_control` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `high_risk_bot_requests_action` object of `bot_ratings` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `http_ddos_protection` object of `security_policy` supports the following:

* `adaptive_frequency_control` - (Optional, List) Specific configuration of adaptive frequency control.
* `bandwidth_abuse_defense` - (Optional, List) Specific configuration of traffic fraud prevention.
* `client_filtering` - (Optional, List) Specific configuration of intelligent client filtering.
* `slow_attack_defense` - (Optional, List) Specific configuration of slow attack protection.

The `human_requests_action` object of `bot_ratings` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `intelligence_rule_items` object of `intelligence_rule` supports the following:

* `action` - (Optional, String) Action. Valid values: `drop`, `trans`, `alg`, `captcha`, `monitor`.
* `label` - (Optional, String) Intelligence label. Valid values: `evil_bot`, `suspect_bot`, `good_bot`, `normal`.

The `intelligence_rule` object of `bot_config` supports the following:

* `intelligence_rule_items` - (Optional, List) Intelligence rule items.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `invalid_attestation_action` object of `client_attestation_rules` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `ip_reputation_group` object of `ip_reputation` supports the following:

* `action_overrides` - (Optional, List) Action overrides for specific IP reputation types.
* `base_action` - (Optional, List) Base action for IP reputation.

The `ip_reputation` object of `basic_bot_settings` supports the following:

* `enabled` - (Optional, String) Whether IP reputation is enabled. Valid values: `on`, `off`.
* `ip_reputation_group` - (Optional, List) IP reputation group configuration.

The `ip_table_config` object of `security_config` supports the following:

* `ip_table_rules` - (Optional, List) IP table rules.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `ip_table_rules` object of `ip_table_config` supports the following:

* `action` - (Optional, String) Action. Valid values: `drop`, `trans`, `monitor`.
* `match_content` - (Optional, String) Match content. Comma-separated for multi-values.
* `match_from` - (Optional, String) Match field. Valid values: `ip`, `area`, `asn`, `referer`, `ua`, `url`.
* `operator` - (Optional, String) Operator. Valid values include `match`, `not_match`, `include_area`, `not_include_area`, `asn_match`, `asn_not_match`, `equal`, `not_equal`, `include`, `not_include`, `is_emty`, `not_exists`.
* `rule_name` - (Optional, String) Rule name.
* `status` - (Optional, String) Rule status. Valid values: `on`, `off`. Default: `on`.

The `known_bot_categories` object of `basic_bot_settings` supports the following:

* `action_overrides` - (Optional, List) Action overrides for specific bot types.
* `base_action` - (Optional, List) Base action for known bot categories.

The `likely_bot_requests_action` object of `bot_ratings` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `low_rate_session_action` object of `session_rate_control` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `managed_rule_groups` object of `managed_rules` supports the following:

* `action` - (Required, List) Handling actions for managed rule groups. the Name parameter value of SecurityAction supports: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe, do not process requests and record security events in logs;</li> <li>Disabled: not enabled, do not scan requests and skip this rule.</li>.
* `group_id` - (Required, String) Group name of the managed rule. if the rule group for the configuration is not specified, it will be processed based on the default configuration. refer to product documentation for the specific value of GroupId.
* `sensitivity_level` - (Required, String) Protection level of the managed rule group. valid values: <li>loose: lenient, only contains ultra-high risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>normal: normal, contains ultra-high risk and high-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>strict: strict, contains ultra-high risk, high-risk and medium-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>extreme: super strict, contains ultra-high risk, high-risk, medium-risk and low-risk rules. at this point, configure Action, and RuleActions configuration is invalid;</li> <li>custom: custom, refined strategy. configure the disposal method for each individual rule. at this point, the Action field is invalid. use RuleActions to configure the refined strategy for each individual rule.</li>.
* `rule_actions` - (Optional, List) Specific configuration of rule items under the managed rule group. the configuration is effective only when SensitivityLevel is custom.

The `managed_rules` object of `default_deny_security_action_parameters` supports the following:

* `block_ip_duration` - (Optional, String) IP block duration when `block_ip` is `on`.
* `block_ip` - (Optional, String) Whether to extend the source IP block. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) PageId of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Whether to use a custom page. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to suspend the request source. Valid values: `on`, `off`.

The `managed_rules` object of `security_policy` supports the following:

* `detection_only` - (Required, String) Indicates whether the evaluation mode is Enabled. it is valid only when the Enabled parameter is set to on. valid values: <li>on: Enabled. all managed rules take effect in observation mode.</li> <li>off: disabled. all managed rules take effect according to the actual configuration.</li>.
* `enabled` - (Required, String) Indicates whether the managed rule is enabled. valid values: <li>on: enabled. all managed rules take effect as configured;</li> <li>off: disabled. all managed rules do not take effect.</li>.
* `auto_update` - (Optional, List) Managed rule automatic update option.
* `managed_rule_groups` - (Optional, Set) Configuration of the managed rule group. if this structure is passed as an empty array or the GroupId is not included in the list, it will be processed based on the default method.
* `semantic_analysis` - (Optional, String) Whether the managed rule semantic analysis option is Enabled is valid only when the Enabled parameter is on. valid values: <li>on: enable. perform semantic analysis on requests before processing them;</li> <li>off: disable. process requests directly without semantic analysis.</li> <br/>default off.

The `max_new_session_trigger_config` object of `bot_session_validation` supports the following:

* `max_new_session_count_interval` - (Optional, String) Statistics time window for trigger threshold. Valid values: `5s`, `10s`, `15s`, `30s`, `60s`, `5m`, `10m`, `30m`, `60m`.
* `max_new_session_count_threshold` - (Optional, Int) Cumulative count for trigger threshold. Range: 1-100000000.

The `meta_data` object of `managed_rule_groups` supports the following:


The `mid_rate_session_action` object of `session_rate_control` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `minimal_request_body_transfer_rate` object of `slow_attack_defense` supports the following:

* `counting_period` - (Required, String) The minimum text transmission rate statistics time range, the possible values are: <li>10s: 10 seconds; </li><li>30s: 30 seconds; </li><li>60s: 60 seconds; </li><li>120s: 120 seconds. </li>.
* `enabled` - (Required, String) Whether the text transmission minimum rate threshold is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `minimal_avg_transfer_rate_threshold` - (Required, String) Minimum text transmission rate threshold. The unit only supports bps.

The `other_modules` object of `default_deny_security_action_parameters` supports the following:

* `block_ip_duration` - (Optional, String) IP block duration when `block_ip` is `on`.
* `block_ip` - (Optional, String) Whether to extend the source IP block. Valid values: `on`, `off`.
* `error_page_id` - (Optional, String) PageId of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Whether to use a custom page. Valid values: `on`, `off`.
* `stall` - (Optional, String) Whether to suspend the request source. Valid values: `on`, `off`.

The `partial_modules` object of `except_user_rule_scope` supports the following:

* `include` - (Optional, List) Rule IDs to include.
* `module` - (Optional, String) Module name. Valid values: `managed-rule`, `managed-group`, `waf` (deprecated).

The `precise_match_rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Execution actions for custom rules. the Name parameter value of SecurityAction supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block using a specified page;</li> <li>Redirect: Redirect to URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>.
* `condition` - (Required, String) The specific content of the custom rule must comply with the expression grammar. please refer to the product document for detailed specifications.
* `enabled` - (Required, String) Indicates whether the custom rule is enabled. valid values: <li>on: enabled</li> <li>off: disabled</li>.
* `name` - (Required, String) The name of the custom rule.
* `priority` - (Optional, Int) Customizes the priority of rules. value range: 0-100. it defaults to 0. only supports `rule_type` is `PreciseMatchRule`.

The `rate_limit_config` object of `security_config` supports the following:

* `rate_limit_customizes` - (Optional, List) Managed customized rate limit rules.
* `rate_limit_intelligence` - (Optional, List) Intelligent client filtering.
* `rate_limit_template` - (Optional, List) Rate limit template.
* `rate_limit_user_rules` - (Optional, List) User-defined rate limit rules.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `rate_limit_customizes` object of `rate_limit_config` supports the following:

* `acl_conditions` - (Optional, List) Rule ACL conditions.
* `action` - (Optional, String) Action. Valid values: `monitor`, `drop`, `redirect`, `page`, `alg`.
* `custom_response_id` - (Optional, String) Custom response ID.
* `freq_fields` - (Optional, List) Filter fields. Valid values: `sip`.
* `freq_scope` - (Optional, List) Statistical scope. Valid values: `source_to_eo`, `client_to_eo`.
* `name` - (Optional, String) Custom response page name. Required when Action is `page`.
* `period` - (Optional, Int) Rate limit statistical period in seconds (10/20/30/40/50/60).
* `punish_time_unit` - (Optional, String) Penalty duration unit. Valid values: `second`, `minutes`, `hour`.
* `punish_time` - (Optional, Int) Penalty duration (0-2 days).
* `redirect_url` - (Optional, String) Redirect URL. Required when Action is `redirect`.
* `response_code` - (Optional, Int) Custom response code (100-600, excl. 3xx).
* `rule_name` - (Optional, String) Rule name.
* `rule_priority` - (Optional, Int) Rule priority (0-100).
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.
* `threshold` - (Optional, Int) Rate limit threshold in count. Range 0-4294967294.

The `rate_limit_intelligence` object of `rate_limit_config` supports the following:

* `action` - (Optional, String) Action. Valid values: `monitor`, `alg`.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `rate_limit_template_detail` object of `rate_limit_template` supports the following:


The `rate_limit_template` object of `rate_limit_config` supports the following:

* `action` - (Optional, String) Template action, e.g. `alg`.
* `mode` - (Optional, String) Template level. Valid values: `sup_loose`, `loose`, `emergency`, `normal`, `strict`, `close`.

The `rate_limit_user_rules` object of `rate_limit_config` supports the following:

* `acl_conditions` - (Optional, List) Rule ACL conditions.
* `action` - (Optional, String) Action. Valid values: `monitor`, `drop`, `redirect`, `page`, `alg`.
* `custom_response_id` - (Optional, String) Custom response ID.
* `freq_fields` - (Optional, List) Filter fields. Valid values: `sip`.
* `freq_scope` - (Optional, List) Statistical scope. Valid values: `source_to_eo`, `client_to_eo`.
* `name` - (Optional, String) Custom response page name. Required when Action is `page`.
* `period` - (Optional, Int) Rate limit statistical period in seconds (10/20/30/40/50/60).
* `punish_time_unit` - (Optional, String) Penalty duration unit. Valid values: `second`, `minutes`, `hour`.
* `punish_time` - (Optional, Int) Penalty duration (0-2 days).
* `redirect_url` - (Optional, String) Redirect URL. Required when Action is `redirect`.
* `response_code` - (Optional, Int) Custom response code (100-600, excl. 3xx).
* `rule_name` - (Optional, String) Rule name.
* `rule_priority` - (Optional, Int) Rule priority (0-100).
* `rule_status` - (Optional, String) Rule status. Valid values: `on`, `off`.
* `threshold` - (Optional, Int) Rate limit threshold in count. Range 0-4294967294.

The `rate_limiting_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) A list of precise rate limiting definitions. When using ModifySecurityPolicy to modify the Web protection configuration: <br> <li> If the Rules parameter is not specified, or the Rules parameter length is zero: clear all precise rate limiting configurations. </li>. <li> If the RateLimitingRules parameter value is not specified in the SecurityPolicy parameter: keep the existing custom rule configuration and do not modify it. </li>.

The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `base_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `bot_client_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `challenge_not_finished_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `challenge_timeout_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `high_rate_session_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `human_requests_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `invalid_attestation_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `likely_bot_requests_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `low_rate_session_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `mid_rate_session_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `session_expired_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `session_invalid_action` supports the following:

* `url` - (Required, String) The URL to redirect.

The `redirect_action_parameters` object of `verified_bot_requests_action` supports the following:

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

The `search_engine_bots` object of `basic_bot_settings` supports the following:

* `action_overrides` - (Optional, List) Action overrides for specific bot types.
* `base_action` - (Optional, List) Base action for search engine bots.

The `security_config` object supports the following:

* `acl_config` - (Optional, List) Custom rule configuration.
* `bot_config` - (Optional, List) Bot configuration.
* `drop_page_config` - (Optional, List) Drop page configuration.
* `except_config` - (Optional, List) Exception rules configuration.
* `ip_table_config` - (Optional, List) Basic access control.
* `rate_limit_config` - (Optional, List) Rate limit configuration.
* `slow_post_config` - (Optional, List) Slow attack configuration.
* `switch_config` - (Optional, List) Layer-7 protection master switch.
* `waf_config` - (Optional, List) Managed rules configuration.

The `security_policy` object supports the following:

* `bot_management_lite` - (Optional, List) Basic Bot management configuration.
* `bot_management` - (Optional, List) Bot management configuration.
* `custom_rules` - (Optional, List) Custom rule configuration.
* `default_deny_security_action_parameters` - (Optional, List) Default deny action configuration. If not specified, the existing configuration is kept.
* `exception_rules` - (Optional, List) Exception rule configuration.
* `http_ddos_protection` - (Optional, List) HTTP DDOS protection configuration.
* `managed_rules` - (Optional, List) Managed rule configuration.
* `rate_limiting_rules` - (Optional, List) Rate limiting rule configuration.

The `session_expired_action` object of `bot_session_validation` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `session_invalid_action` object of `bot_session_validation` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `session_rate_control` object of `bot_session_validation` supports the following:

* `enabled` - (Optional, String) Whether session rate control is enabled. Valid values: `on`, `off`.
* `high_rate_session_action` - (Optional, List) Action for high-rate session.
* `low_rate_session_action` - (Optional, List) Action for low-rate session.
* `mid_rate_session_action` - (Optional, List) Action for mid-rate session.

The `skip_conditions` object of `except_user_rule_scope` supports the following:

* `match_content_type` - (Optional, String) Content match type. Valid values: `equal`, `wildcard`.
* `match_content` - (Optional, List) Content values.
* `match_from_type` - (Optional, String) Key match type. Valid values: `equal`, `wildcard`.
* `match_from` - (Optional, List) Key values.
* `selector` - (Optional, String) Selector. Valid values: `args`, `path`, `full`, `upload_filename`, `keys`, `values`, `key_value`.
* `type` - (Optional, String) Skip type. Valid values: `header_fields`, `cookie`, `query_string`, `uri`, `body_raw`, `body_json`.

The `slow_attack_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether slow attack protection is enabled. The possible values are: <li>on: enabled; </li><li>off: disabled. </li>.
* `action` - (Optional, List) The handling method of slow attack protection. When Enabled is on, this field is required. SecurityAction Name value supports: <li>Monitor: Observe; </li><li>Deny: Intercept; </li>.
* `minimal_request_body_transfer_rate` - (Optional, List) Specific configuration of the minimum rate threshold for text transmission. This field is required when Enabled is on.
* `request_body_transfer_timeout` - (Optional, List) Specific configuration of the text transmission timeout. When Enabled is on, this field is required.

The `slow_post_config` object of `security_config` supports the following:

* `action` - (Optional, String) Action. Valid values: `monitor`, `drop`.
* `first_part_config` - (Optional, List) First packet configuration.
* `slow_rate_config` - (Optional, List) Slow rate configuration.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

The `slow_rate_config` object of `slow_post_config` supports the following:

* `interval` - (Optional, Int) Statistical interval in seconds.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.
* `threshold` - (Optional, Int) Rate threshold in bps.

The `source_idc` object of `basic_bot_settings` supports the following:

* `action_overrides` - (Optional, List) Action overrides for specific bot types.
* `base_action` - (Optional, List) Base action for IDC requests.

The `switch_config` object of `security_config` supports the following:

* `web_switch` - (Optional, String) Web master switch. Valid values: `on`, `off`. Does not affect DDoS or Bot switches.

The `template_config` object of `security_config` supports the following:


The `verified_bot_requests_action` object of `bot_ratings` supports the following:

* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `challenge_action_parameters` - (Optional, List) Additional parameters when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `name` - (Optional, String) Action name. Valid values: `Deny`, `Monitor`, `Allow`, `Challenge`, `Disabled`.
* `redirect_action_parameters` - (Optional, List) Additional parameters when Name is Redirect.

The `waf_config` object of `security_config` supports the following:

* `ai_rule` - (Optional, List) AI rule engine configuration.
* `level` - (Optional, String) Protection level. Valid values: `loose`, `normal`, `strict`, `stricter`, `custom`.
* `mode` - (Optional, String) Global WAF mode. Valid values: `block`, `observe`.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.
* `waf_rule` - (Optional, List) Managed rule detail configuration.

The `waf_drop_page_detail` object of `drop_page_config` supports the following:

* `custom_response_id` - (Optional, String) Custom response ID.
* `name` - (Optional, String) Block page file name or URL.
* `page_id` - (Optional, Int) The unique ID of the block page. The system includes a built-in block page with ID 0.
* `status_code` - (Optional, Int) HTTP status code for the block page. Range: 100-600, excluding 3xx.
* `type` - (Optional, String) Page type. Valid values: `page`.

The `waf_rule` object of `waf_config` supports the following:

* `block_rule_ids` - (Optional, List) Rule IDs to block (disable).
* `observe_rule_ids` - (Optional, List) Rule IDs in observe mode.
* `switch` - (Optional, String) Switch. Valid values: `on`, `off`.

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

