---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_web_security_template"
sidebar_current: "docs-tencentcloud-resource-teo_web_security_template"
description: |-
  Provides a resource to create a TEO web security template
---

# tencentcloud_teo_web_security_template

Provides a resource to create a TEO web security template

## Example Usage

```hcl
resource "tencentcloud_teo_web_security_template" "example" {
  zone_id       = "zone-3fkff38fyw8s"
  template_name = "example"
  security_policy {
    bot_management {
      enabled = "off"
      basic_bot_settings {
        bot_intelligence {
          enabled = "off"
          bot_ratings {
            high_risk_bot_requests_action {
              name = "Monitor"
            }

            human_requests_action {
              name = "Allow"
            }

            likely_bot_requests_action {
              name = "Monitor"
            }

            verified_bot_requests_action {
              name = "Monitor"
            }
          }
        }

        ip_reputation {
          enabled = "off"
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
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }

      bandwidth_abuse_defense {
        enabled = "off"
        action {
          name = "Monitor"
        }
      }

      client_filtering {
        enabled = "on"
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }

      slow_attack_defense {
        enabled = "off"
        action {
          name = "Deny"
        }

        minimal_request_body_transfer_rate {
          counting_period                     = "60s"
          enabled                             = "off"
          minimal_avg_transfer_rate_threshold = "80bps"
        }

        request_body_transfer_timeout {
          enabled      = "off"
          idle_timeout = "5s"
        }
      }
    }

    managed_rules {
      detection_only    = "on"
      enabled           = "on"
      semantic_analysis = "off"

      auto_update {
        auto_update_to_latest_version = "off"
      }

      managed_rule_groups {
        group_id          = "wafgroup-webshell-attacks"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-xss-attacks"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-xxe-attacks"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-vulnerability-scanners"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-non-compliant-protocol-usages"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-cms-vulnerabilities"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-file-upload-attacks"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-other-vulnerabilities"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-command-and-code-injections"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-sql-injections"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-shiro-vulnerabilities"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-file-accesses"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-ldap-injections"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-oa-vulnerabilities"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-ssrf-attacks"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-ssti-attacks"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-accesses"
        sensitivity_level = "strict"

        action {
          name = "Monitor"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) Policy template name. Composed of Chinese characters, letters, digits, and underscores. Cannot begin with an underscore and must be less than or equal to 32 characters.
* `zone_id` - (Required, String, ForceNew) Zone ID. Explicitly identifies the zone to which the policy template belongs for access control purposes.
* `security_policy` - (Optional, List) Web security policy template configuration. Generates default config if empty. Supported: Exception rules, custom rules, rate limiting rules, managed rules. Not supported: Bot management rules (under development).

The `action` object of `adaptive_frequency_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `bandwidth_abuse_defense` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `bot_management_action_overrides` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `client_filtering` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `frequent_scanning_protection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `managed_rule_groups` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `rule_actions` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `rules` supports the following:

* `bot_session_validation` - (Optional, List) Configures Cookie verification and session tracking.
* `client_behavior_detection` - (Optional, List) Configures client behavior validation.

The `action` object of `rules` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `rules` supports the following:

* `security_action` - (Optional, List) The handling method of the Bot custom rule. valid values: <li>Allow: pass, where AllowActionParameters supports MinDelayTime and MaxDelayTime configuration;</li> <li>Deny: block, where DenyActionParameters supports BlockIp, ReturnCustomPage, and Stall configuration;</li> <li>Monitor: observation;</li> <li>Challenge: Challenge, where ChallengeActionParameters.ChallengeOption supports JSChallenge and ManagedChallenge;</li> <li>Redirect: Redirect to URL.</li>.
* `weight` - (Optional, Int) The Weight of the current SecurityAction, only supported between 10 and 100 and must be a multiple of 10. the total of all Weight parameters must equal 100.

The `action` object of `slow_attack_defense` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `adaptive_frequency_control` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether adaptive frequency control is enabled. valid values: <li>on: enable;</li> <li>off: disable.</li>.
* `action` - (Optional, List) The handling method of adaptive frequency control. this field is required when Enabled is on. valid values for SecurityAction Name: <li>Monitor: observation;</li> <li>Deny: block;</li> <li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge.</li>.
* `sensitivity` - (Optional, String) The restriction level of adaptive frequency control. required when Enabled is on. valid values: <li>Loose: Loose</li><li>Moderate: Moderate</li><li>Strict: Strict</li>.

The `allow_action_parameters` object of `action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `base_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `bot_client_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `challenge_not_finished_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `challenge_timeout_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `high_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `high_risk_request_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `human_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `invalid_attestation_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `likely_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `low_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `medium_risk_request_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `mid_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `security_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `session_expired_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `session_invalid_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `allow_action_parameters` object of `verified_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: <li>s: seconds, value ranges from 5 to 10.</li>.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: <li>s: seconds, value ranges from 0 to 5.</li>.

The `auto_update` object of `managed_rules` supports the following:

* `auto_update_to_latest_version` - (Required, String) Enable automatic update to the latest version or not. Values: <li>`on`: enabled</li> <li>`off`: disabled</li>.

The `bandwidth_abuse_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether bandwidth abuse protection (applicable to chinese mainland only) is enabled. valid values: <li>on: enabled;</li> <li>off: disabled.</li>.
* `action` - (Optional, List) Bandwidth abuse protection (applicable to chinese mainland) handling method. required when Enabled is on. valid values for SecurityAction Name: <li>Monitor: observe;</li> <li>Deny: block;</li> <li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge.</li>.

The `base_action` object of `ip_reputation_group` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `base_action` object of `known_bot_categories` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `base_action` object of `search_engine_bots` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `base_action` object of `source_idc` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `basic_bot_settings` object of `bot_management` supports the following:

* `bot_intelligence` - (Optional, List) Specifies the configuration for Bot intelligent analysis.
* `ip_reputation` - (Optional, List) Threat intelligence database (originally client profile analysis) configuration, used for handling client ips with specific risk characteristics in recent access behavior.
* `known_bot_categories` - (Optional, List) Commercial or open-source tool UA feature configuration (original UA feature rule), used to handle access requests from known commercial or open-source tools. the User-Agent header of such requests complies with known commercial or open-source tool features.
* `search_engine_bots` - (Optional, List) Search engine crawler configuration, used to handle requests from search engine crawlers. the IP, User-Agent, or rDNS results of such requests match known search engine crawlers.
* `source_idc` - (Optional, List) Client IP source IDC configuration, used for handling access requests from client ips in idcs (data centers). such source requests are not directly accessed by mobile terminals or browser-side.

The `block_ip_action_parameters` object of `action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `base_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `bot_client_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `challenge_not_finished_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `challenge_timeout_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `high_rate_session_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `high_risk_request_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `human_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `invalid_attestation_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `likely_bot_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `low_rate_session_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `medium_risk_request_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `mid_rate_session_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `security_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `session_expired_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `session_invalid_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `block_ip_action_parameters` object of `verified_bot_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: <li>`s`: second, value range 1-120;</li> <li>`m`: minute, value range 1-120;</li> <li>`h`: hour, value range 1-48.</li>.

The `bot_client_action` object of `client_behavior_detection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `bot_intelligence` object of `basic_bot_settings` supports the following:

* `bot_ratings` - (Optional, List) Based on client and request features, divides request sources into human requests, legitimate Bot requests, suspected Bot requests, and high-risk Bot requests, and provides request handling options.
* `enabled` - (Optional, String) Specifies the switch for Bot intelligent analysis configuration. valid values:.

on: enabled.
off: disabled.

The `bot_management_action_overrides` object of `ip_reputation_group` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: <li>Deny: block;</li><li>Monitor: observe;</li><li>Disabled: Disabled, disable the specified rule;</li><li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li><li>Allow: pass (only for Bot basic feature management).</li>.
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management_action_overrides` object of `known_bot_categories` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: <li>Deny: block;</li><li>Monitor: observe;</li><li>Disabled: Disabled, disable the specified rule;</li><li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li><li>Allow: pass (only for Bot basic feature management).</li>.
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management_action_overrides` object of `search_engine_bots` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: <li>Deny: block;</li><li>Monitor: observe;</li><li>Disabled: Disabled, disable the specified rule;</li><li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li><li>Allow: pass (only for Bot basic feature management).</li>.
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management_action_overrides` object of `source_idc` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: <li>Deny: block;</li><li>Monitor: observe;</li><li>Disabled: Disabled, disable the specified rule;</li><li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li><li>Allow: pass (only for Bot basic feature management).</li>.
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management` object of `security_policy` supports the following:

* `basic_bot_settings` - (Optional, List) Bot management basic configuration. takes effect on all domains associated with the policy. can be customized through CustomRules.
* `browser_impersonation_detection` - (Optional, List) Configures browser spoofing identification rules (formerly active feature detection rule). sets the response page range for JavaScript injection, browser check options, and handling method for non-browser clients.
* `client_attestation_rules` - (Optional, List) Definition list of client authentication rules. this feature is in beta test. submit a ticket if you need to use it.
* `custom_rules` - (Optional, List) Bot management custom rule combines various crawlers and request behavior characteristics to accurately define bots and configure customized handling methods.
* `enabled` - (Optional, String) Whether Bot management is enabled. valid values: <li>on: enabled;</li><li>off: disabled.</li>.

The `bot_ratings` object of `bot_intelligence` supports the following:

* `high_risk_bot_requests_action` - (Optional, List) Execution action for malicious Bot requests. valid values for the Name parameter in SecurityAction: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Allow: pass;</li> <li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.</li>.
* `human_requests_action` - (Optional, List) Execution action for a normal Bot request. valid values for the Name parameter in SecurityAction: <li>Allow: pass.</li>.
* `likely_bot_requests_action` - (Optional, List) The execution action for suspected Bot requests. valid values for the Name parameter in SecurityAction: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Allow: pass;</li> <li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.</li>.
* `verified_bot_requests_action` - (Optional, List) Execution action for friendly Bot request. SecurityAction Name parameter supports: <li>Deny: block;</li><li>Monitor: observe;</li><li>Allow: pass;</li><li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.</li>.

The `bot_session_validation` object of `action` supports the following:

* `issue_new_bot_session_cookie` - (Optional, String) Whether to update Cookie and validate. valid values: <li>on: update Cookie and validate;</li> <li>off: verify only.</li>.
* `max_new_session_trigger_config` - (Optional, List) Specifies the trigger threshold for updating and validating cookies. valid only when IssueNewBotSessionCookie is set to on.
* `session_expired_action` - (Optional, List) Execution action when no Cookie is carried or the Cookie expired. valid values for the Name parameter in SecurityAction: <li>Deny: block, where Stall can be configured in DenyActionParameters;</li><li>Monitor: observe;</li><li>Allow: respond after wait, where MinDelayTime and MaxDelayTime must be configured in AllowActionParameters.</li>.
* `session_invalid_action` - (Optional, List) Execution action for invalid Cookie. valid values for the Name parameter in SecurityAction: <li>Deny: block, where the DenyActionParameters supports Stall configuration;</li><li>Monitor: observe;</li><li>Allow: respond after wait, where AllowActionParameters requires MinDelayTime and MaxDelayTime configuration.</li>.
* `session_rate_control` - (Optional, List) Specifies the session rate and periodic feature verification configuration.

The `browser_impersonation_detection` object of `bot_management` supports the following:

* `rules` - (Optional, List) List of browser spoofing identification Rules. when using ModifySecurityPolicy to modify Web protection configuration: <br> <li>if Rules parameter in SecurityPolicy.BotManagement.BrowserImpersonationDetection is not specified or parameter length is zero: clear all browser spoofing identification rule configurations.</li> <li>if BrowserImpersonationDetection parameter value is unspecified in SecurityPolicy.BotManagement parameters: keep existing browser spoofing identification rule configurations without modification.</li>.

The `challenge_action_parameters` object of `action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `base_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `bot_client_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `challenge_not_finished_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `challenge_timeout_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `high_rate_session_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `high_risk_request_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `human_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `invalid_attestation_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `likely_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `low_rate_session_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `medium_risk_request_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `mid_rate_session_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `security_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `session_expired_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `session_invalid_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_action_parameters` object of `verified_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: <li> InterstitialChallenge: interstitial challenge;</li> <li> InlineChallenge: embedded challenge;</li> <li> JSChallenge: JavaScript challenge;</li> <li> ManagedChallenge: managed challenge.</li>.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: <li>s: second, value ranges from 1 to 60;</li><li>m: minute, value ranges from 1 to 60;</li><li>h: hour, value ranges from 1 to 24.</li>.

The `challenge_not_finished_action` object of `client_behavior_detection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `challenge_timeout_action` object of `client_behavior_detection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `client_attestation_rules` object of `bot_management` supports the following:

* `rules` - (Optional, List) List of client authentication. when using ModifySecurityPolicy to modify Web protection configuration: <li> if Rules in SecurityPolicy.BotManagement.ClientAttestationRules is not specified or the parameter length of Rules is zero: clear all client authentication rule configuration. </li> <li> if ClientAttestationRules in SecurityPolicy.BotManagement parameters is unspecified: keep existing client authentication rule configuration and do not modify. </li>.

The `client_behavior_detection` object of `action` supports the following:

* `bot_client_action` - (Optional, List) The execution action of the Bot client. valid values for the Name parameter in SecurityAction: <li>Deny: block, where the Stall configuration is supported in DenyActionParameters;</li><li>Monitor: observation;</li><li>Allow: respond after wait, where MinDelayTime and MaxDelayTime configurations are required in AllowActionParameters.</li>.
* `challenge_not_finished_action` - (Optional, List) Execution action when client-side javascript is not enabled (test not completed). valid values for SecurityAction Name: <li>Deny: block, where Stall configuration is supported in DenyActionParameters;</li><li>Monitor: observe;</li><li>Allow: respond after waiting, where MinDelayTime and MaxDelayTime configuration is required in AllowActionParameters.</li>.
* `challenge_timeout_action` - (Optional, List) The execution action for client-side detection timeout. valid values for the Name parameter in SecurityAction: <li>Deny: block, where Stall can be configured in DenyActionParameters;</li> <li>Monitor: observe;</li> <li>Allow: respond after wait, where MinDelayTime and MaxDelayTime must be configured in AllowActionParameters.</li>.
* `crypto_challenge_delay_before` - (Optional, String) Specifies the execution mode for client behavior verification. valid values: <li>0ms: immediate execution;</li> <li>100ms: delay 100ms execution;</li> <li>200ms: delay 200ms execution;</li> <li>300ms: delay 300ms execution;</li> <li>400ms: delay 400ms execution;</li> <li>500ms: delay 500ms execution;</li> <li>600ms: delay 600ms execution;</li> <li>700ms: delay 700ms execution;</li> <li>800ms: delay 800ms execution;</li> <li>900ms: delay 900ms execution;</li> <li>1000ms: delay 1000ms execution.</li>.
* `crypto_challenge_intensity` - (Optional, String) Specifies the proof-of-work strength. valid values: <li>low: low;</li><li>medium: medium;</li><li>high: high.</li>.
* `max_challenge_count_interval` - (Optional, String) Time window for trigger threshold statistics. valid values: <li>5s: within 5 seconds;</li><li>10s: within 10 seconds;</li><li>15s: within 15 seconds;</li><li>30s: within 30 seconds;</li><li>60s: within 60 seconds;</li><li>5m: within 5 minutes;</li><li>10m: within 10 minutes;</li><li>30m: within 30 minutes;</li><li>60m: within 60 minutes.</li>.
* `max_challenge_count_threshold` - (Optional, Int) Trigger threshold cumulative count. value range: 1-100000000.

The `client_filtering` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether intelligent client filtering is enabled. valid values: <li>on: enable;</li> <li>off: disable.</li>.
* `action` - (Optional, List) The handling method of intelligent client filtering. when Enabled is on, this field is required. the Name parameter of SecurityAction supports: <li>Monitor: observation;</li> <li>Deny: block;</li> <li>Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge.</li>.

The `custom_rules` object of `bot_management` supports the following:

* `rules` - (Optional, List) List of Bot custom Rules. when using ModifySecurityPolicy to modify Web protection configuration: <br> <li> if Rules in SecurityPolicy.BotManagement.CustomRules is not specified or parameter length of Rules is zero: clear all Bot custom rule configurations.</li> <li> if CustomRules in SecurityPolicy.BotManagement parameters is unspecified: keep existing Bot custom rule configurations and do not modify them.</li>.

The `custom_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) The custom rule. <br>when modifying the Web protection configuration using ModifySecurityPolicy: <br> - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations. <br> - if the Rules parameter is not specified: keep the existing custom rule configuration without modification.

The `deny_action_parameters` object of `action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `base_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `bot_client_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `challenge_not_finished_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `challenge_timeout_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `high_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `high_risk_request_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `human_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `invalid_attestation_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `likely_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `low_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `medium_risk_request_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `mid_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `security_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `session_expired_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `session_invalid_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `verified_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.
Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.
Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:.
<li>`on`: Enable;</li>

<li>off: Disable.</li>

Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources.
Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `device_profiles` object of `rules` supports the following:

* `client_type` - (Required, String) Client device type. valid values: <li>iOS;</li> <li>Android;</li> <li>WebView.</li>.
* `high_risk_min_score` - (Optional, Int) The minimum value to determine a request as high-risk ranges from 1-99. the larger the value, the higher the request risk, and the closer it resembles a request initiated by a Bot client. the default value is 50, corresponding to high-risk for values 51-100.
* `high_risk_request_action` - (Optional, List) Handling method for high-risk requests. valid values for SecurityAction Name: <li>Deny: block;</li> <li>Monitor: observation;</li> <li>Redirect: redirection;</li> <li>Challenge: Challenge.</li> default value: Monitor.
* `medium_risk_min_score` - (Optional, Int) Specifies the minimum value to determine a request as medium-risk. value range: 1-99. the larger the value, the higher the request risk, resembling requests initiated by a Bot client. default value: 15, corresponding to medium-risk for values 16-50.
* `medium_risk_request_action` - (Optional, List) Handling method for medium-risk requests. SecurityAction Name parameter supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Redirect: Redirect;</li> <li>Challenge: Challenge.</li> default value is Monitor.

The `exception_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) Definition list of exception Rules. when using ModifySecurityPolicy to modify Web protection configuration: <li>if the Rules parameter is not specified or the parameter length is zero: clear all exception rule configurations.</li><li>if the ExceptionRules parameter value is not specified in SecurityPolicy: keep existing exception rule configurations without modification.</li>.

The `frequent_scanning_protection` object of `managed_rules` supports the following:

* `action_duration` - (Optional, String) This parameter specifies the duration of the handling Action set by the high frequency scan protection Action parameter. value range: 60 to 86400. measurement unit: seconds (s) only, for example 60s. this field is required when Enabled is on.
* `action` - (Optional, List) The handling action for high-frequency scan protection. required when Enabled is on. valid values for SecurityAction Name: <li>Deny: block and respond with an interception page;</li> <li>Monitor: observe without processing requests, log security events in logs;</li> <li>JSChallenge: respond with a JavaScript challenge page.</li>.
* `block_threshold` - (Optional, Int) This parameter specifies the threshold for high-frequency scan protection, which is the intercept count of managed rules set to interception within the time range set by CountingPeriod. value range: 1 to 4294967294, for example 100. when exceeding this statistical value, subsequent requests will trigger the handling Action set by Action. required when Enabled is on.
* `count_by` - (Optional, String) The match mode for request statistics. required when Enabled is on. valid values: <li>http.request.xff_header_ip: client ip (priority match xff header);</li><li>http.request.ip: client ip.</li>.
* `counting_period` - (Optional, String) This parameter specifies the statistical time window for high-frequency scan protection, which is the time window for counting requests that hit managed rules configured as block. valid values: 5-1800. measurement unit: seconds (s) only, such as 5s. this field is required when Enabled is on.
* `enabled` - (Optional, String) Whether the high-frequency scan protection rule is enabled. valid values: <li>on: enable. the high-frequency scan protection rule takes effect.</li><li>off: disable. the high-frequency scan protection rule does not take effect.</li>.

The `high_rate_session_action` object of `session_rate_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `high_risk_bot_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `high_risk_request_action` object of `device_profiles` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `http_ddos_protection` object of `security_policy` supports the following:

* `adaptive_frequency_control` - (Optional, List) Specifies the specific configuration of adaptive frequency control.
* `bandwidth_abuse_defense` - (Optional, List) Specifies the specific configuration for bandwidth abuse protection.
* `client_filtering` - (Optional, List) Specifies the intelligent client filter configuration.
* `slow_attack_defense` - (Optional, List) Specifies the configuration of slow attack protection.

The `human_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `invalid_attestation_action` object of `rules` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `ip_reputation_group` object of `ip_reputation` supports the following:

* `base_action` - (Optional, List) Execution action of the IP intelligence library (formerly client profile analysis). SecurityAction Name parameter supports: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Disabled: not enabled, disable specified rule;</li> <li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.</li>.
* `bot_management_action_overrides` - (Optional, List) The specific configuration of the IP intelligence library (originally client profile analysis), used to override the default configuration in BaseAction. among them, the Ids in BotManagementActionOverrides can be filled with: <li>IPREP_WEB_AND_DDOS_ATTACKERS_LOW: network attack - general confidence;</li> <li>IPREP_WEB_AND_DDOS_ATTACKERS_MID: network attack - medium confidence;</li> <li>IPREP_WEB_AND_DDOS_ATTACKERS_HIGH: network attack - HIGH confidence;</li> <li>IPREP_PROXIES_AND_ANONYMIZERS_LOW: network proxy - general confidence;</li> <li>IPREP_PROXIES_AND_ANONYMIZERS_MID: network proxy - medium confidence;</li> <li>IPREP_PROXIES_AND_ANONYMIZERS_HIGH: network proxy - HIGH confidence;</li> <li>IPREP_SCANNING_TOOLS_LOW: scanner - general confidence;</li> <li>IPREP_SCANNING_TOOLS_MID: scanner - medium confidence;</li> <li>IPREP_SCANNING_TOOLS_HIGH: scanner - HIGH confidence;</li> <li>IPREP_ATO_ATTACKERS_LOW: account takeover attack - general confidence;</li> <li>IPREP_ATO_ATTACKERS_MID: account takeover attack - medium confidence;</li> <li>IPREP_ATO_ATTACKERS_HIGH: account takeover attack - HIGH confidence;</li> <li>IPREP_WEB_SCRAPERS_AND_TRAFFIC_BOTS_LOW: malicious BOT - general confidence;</li> <li>IPREP_WEB_SCRAPERS_AND_TRAFFIC_BOTS_MID: malicious BOT - medium confidence;</li> <li>IPREP_WEB_SCRAPERS_AND_TRAFFIC_BOTS_HIGH: malicious BOT - HIGH confidence.</li>.

The `ip_reputation` object of `basic_bot_settings` supports the following:

* `enabled` - (Optional, String) IP intelligence library (formerly client profile analysis). valid values: <li>on: enable;</li> <li>off: disable.</li>.
* `ip_reputation_group` - (Optional, List) IP intelligence library (formerly client profile analysis) configuration content.

The `known_bot_categories` object of `basic_bot_settings` supports the following:

* `base_action` - (Optional, List) Handling method for access requests from known commercial tools or open-source tools. specifies the Name parameter value of SecurityAction: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Disabled: not enabled, disable specified rule;</li> <li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li> <li>Allow: pass (to be deprecated).</li>.
* `bot_management_action_overrides` - (Optional, List) Specifies the handling method for access requests from known commercial tools or open-source tools.

The `likely_bot_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `low_rate_session_action` object of `session_rate_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `managed_rule_groups` object of `managed_rules` supports the following:

* `action` - (Required, List) Action for ManagedRuleGroup. the Name parameter value of SecurityAction supports: <li>`Deny`: block and respond with a block page;</li> <li>`Monitor`: observe, do not process requests and record security events in logs;</li> <li>`Disabled`: not enabled, do not scan requests and skip this rule.</li>.
* `group_id` - (Required, String) Name of the managed rule group, if the configuration for the rule group is not specified, it will be processed by default, refer to product documentation for the specific value of GroupId.
* `sensitivity_level` - (Required, String) Protection level of the managed rule group. Values: <li>`loose`: lenient, only contain ultra-high risk rules, at this point, Action parameter needs configured instead of RuleActions parameter;</li> <li>`normal`: normal, contain ultra-high risk and high-risk rules, at this point,Action parameter needs configured instead of RuleActions parameter;</li> <li>`strict`: strict, contains ultra-high risk, high-risk and medium-risk rules, at this point, Action parameter needs configured instead of RuleActions parameter;</li> <li>`extreme`: super strict, contains ultra-high risk, high-risk, medium-risk and low-risk rules, at this point, Action parameter needs configured instead of RuleActions parameter;</li> <li>`custom`: custom, refined strategy, configure the RuleActions parameter for each individual rule, at this point, the Action field is invalid, use RuleActions to configure the refined strategy for each individual rule.</li>.
* `rule_actions` - (Optional, List) Specific configuration of rule items under the managed rule group, valid only when SensitivityLevel is custom.

The `managed_rules` object of `security_policy` supports the following:

* `detection_only` - (Required, String) Evaluation mode is enabled or not, it is valid only when the `Enabled` parameter is set to `on`. Values: <li>`on`: enabled, all managed rules take effect in `observe` mode.</li> <li>off: disabled, all managed rules take effect according to the specified configuration.</li>.
* `enabled` - (Required, String) The managed rule status. Values: <li>`on`: enabled, all managed rules take effect as configured;</li> <li>`off`: disabled, all managed rules do not take effect.</li>.
* `auto_update` - (Optional, List) Managed rule automatic update option.
* `frequent_scanning_protection` - (Optional, List) High-Frequency scan protection configuration option. when a visitor's frequent requests hit the managed rule configured as block within a period of time, all requests from that visitor are blocked.
* `managed_rule_groups` - (Optional, List) Configuration of the managed rule group. If this structure is passed as an empty array or the GroupId is not included in the array, it will be processed based by default.
* `semantic_analysis` - (Optional, String) Managed rule semantic analysis is enabled or not, it is valid only when the `Enabled` parameter is `on`. Values: <li>`on`: enabled, perform semantic analysis  before processing requests;</li> <li>`off`: disabled, process requests directly without semantic analysis.</li> <br/>The default value is `off`.

The `max_new_session_trigger_config` object of `bot_session_validation` supports the following:

* `max_new_session_count_interval` - (Optional, String) Time window for trigger threshold statistics. valid values: <li>5s: within 5 seconds;</li><li>10s: within 10 seconds;</li><li>15s: within 15 seconds;</li><li>30s: within 30 seconds;</li><li>60s: within 60 seconds;</li><li>5m: within 5 minutes;</li><li>10m: within 10 minutes;</li><li>30m: within 30 minutes;</li><li>60m: within 60 minutes.</li>.
* `max_new_session_count_threshold` - (Optional, Int) Trigger threshold cumulative count. value range: 1-100000000.

The `medium_risk_request_action` object of `device_profiles` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `meta_data` object of `managed_rule_groups` supports the following:


The `mid_rate_session_action` object of `session_rate_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `minimal_request_body_transfer_rate` object of `slow_attack_defense` supports the following:

* `counting_period` - (Required, String) Minimum body transfer rate statistical time range, valid values: <li>10s: 10 seconds;</li> <li>30s: 30 seconds;</li> <li>60s: 60 seconds;</li> <li>120s: 120 seconds.</li>.
* `enabled` - (Required, String) Specifies whether the minimum body transfer rate threshold is enabled. valid values: <li>on: enable;</li> <li>off: disable.</li>.
* `minimal_avg_transfer_rate_threshold` - (Required, String) Minimum body transfer rate threshold, the measurement unit is only supported in bps.

The `rate_limiting_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) Definition list of precise rate limiting. when using ModifySecurityPolicy to modify the Web protection configuration: <br> <li> if the Rules parameter is not specified or its length is zero: clear all precision rate limiting configurations.</li> <li> if the RateLimitingRules parameter value is unspecified in the SecurityPolicy parameter: retain the existing custom rule configuration without modification.</li>.

The `redirect_action_parameters` object of `action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `base_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `bot_client_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `challenge_not_finished_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `challenge_timeout_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `high_rate_session_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `high_risk_request_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `human_requests_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `invalid_attestation_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `likely_bot_requests_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `low_rate_session_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `medium_risk_request_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `mid_rate_session_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `security_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `session_expired_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `session_invalid_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `redirect_action_parameters` object of `verified_bot_requests_action` supports the following:

* `url` - (Required, String) Redirect URL.

The `request_body_transfer_timeout` object of `slow_attack_defense` supports the following:

* `enabled` - (Required, String) Whether body transfer timeout is enabled. valid values: <li>`on`: enable</li> <li>`off`: disable</li>.
* `idle_timeout` - (Required, String) Body transfer timeout duration. valid values: 5-120. measurement unit: seconds (s) only.

The `request_fields_for_exception` object of `rules` supports the following:

* `condition` - (Required, String) Skip specific field expression must comply with expression grammar.
Condition supports expression configuration syntax: <li> write according to the matching conditional expression syntax of rules, with support for referencing key and value.</li> <li> supports in, like operators, and logical combination with and.</li>.
For example: <li>${key} in ['x-trace-id']: the parameter name equals x-trace-id.</li> <li>${key} in ['x-trace-id'] and ${value} like ['Bearer *']: the parameter name equals x-trace-id and the parameter value wildcard matches Bearer *.</li>.
* `scope` - (Required, String) Skip specific field. supported values:.
<li>body.json: parameter content in json requests. at this point, Condition supports key and value, TargetField supports key and value, for example { "Scope": "body.json", "Condition": "", "TargetField": "key" }, which means all parameters in json requests skip WAF scan.</li>.
<li style="margin-top:5px">cookie: cookie; at this point Condition supports key, value, TargetField supports key, value, for example { "Scope": "cookie", "Condition": "${key} in ['account-id'] and ${value} like ['prefix-*']", "TargetField": "value" }, which means the cookie parameter name equals account-id and the parameter value wildcard matches prefix-* to skip WAF scan;</li>.
<li style="margin-top:5px">header: HTTP header parameters. at this point, Condition supports key and value, TargetField supports key and value, for example { "Scope": "header", "Condition": "${key} like ['x-auth-*']", "TargetField": "value" }, which means header parameter name wildcard match x-auth-* skips WAF scan.</li>.
<li style="margin-top:5px">uri.query: URL encoding content/query parameter. at this point, Condition supports key and value, TargetField supports key and value. example: { "Scope": "uri.query", "Condition": "${key} in ['action'] and ${value} in ['upload', 'delete']", "TargetField": "value" }. indicates URL encoding content/query parameter name equal to action and parameter value equal to upload or delete skips WAF scan.</li>.
<li style="margin-top:5px">uri: specifies the request path uri. at this point, Condition must be empty. TargetField supports query, path, fullpath, such as {"Scope": "uri", "Condition": "", "TargetField": "query"}, indicates the request path uri skips WAF scan for query parameters.</li>.
<li style="margin-top:5px">body: request body content. at this point Condition must be empty, TargetField supports fullbody, multipart, such as { "Scope": "body", "Condition": "", "TargetField": "fullbody" }, which means the request body content skips WAF scan as a full request.</li>.
* `target_field` - (Required, String) The Scope parameter takes different values. the TargetField expression supports the following values:.
<Li> body.json: supports key, value.</li>.
<li>cookie: supports key and value.</li>.
<li>header: supports key, value</li>.
<Li> uri.query: supports key and value</li>.
<li>uri. specifies path, query, or fullpath.</li>.
<Li>Body: supports fullbody and multipart.</li>.

The `return_custom_page_action_parameters` object of `action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `base_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `bot_client_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `challenge_not_finished_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `challenge_timeout_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `high_rate_session_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `high_risk_request_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `human_requests_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `invalid_attestation_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `likely_bot_requests_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `low_rate_session_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `medium_risk_request_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `mid_rate_session_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `security_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `session_expired_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `session_invalid_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `return_custom_page_action_parameters` object of `verified_bot_requests_action` supports the following:

* `error_page_id` - (Required, String) Response custom page ID.
* `response_code` - (Required, String) Response custom status code.

The `rule_actions` object of `managed_rule_groups` supports the following:

* `action` - (Required, List) Action for the managed rule item specified by RuleId, the SecurityAction Name parameter supports: <li>`Deny`: block and respond with an block page;</li> <li>`Monitor`: observe, do not process the request and record the security event in logs;</li> <li>`Disabled`: disabled, do not scan the request and skip this rule.</li>.
* `rule_id` - (Required, String) Specific items under ManagedRuleGroup, used to rewrite the configuration of this individual rule item, refer to product documentation for details.

The `rule_details` object of `meta_data` supports the following:


The `rules` object of `browser_impersonation_detection` supports the following:

* `action` - (Optional, List) Describes the handling method for browser spoofing identification rules, including Cookie verification, session tracking configuration, and client behavior validation configuration.
* `condition` - (Optional, String) Specifies the specific content of browser spoofing identification rules, which only support configuration of request Method (Method), request Path (Path), and request URL, and must comply with expression grammar. for detailed specifications, please refer to the product document.
* `enabled` - (Optional, String) Whether browser spoofing detection is enabled. valid values: <li>on: enabled;</li><li>off: disabled.</li>.
* `id` - (Optional, String) Browser spoofing identification rule ID. rule ID supports different rule configuration operations: <li> <b>add</b> a new rule: ID is empty or without specifying the ID parameter;</li> <li> <b>modify</b> an existing rule: specify the rule ID that needs to be updated/modified;</li> <li> <b>delete</b> an existing rule: existing Rules not included in the Rules list of the BrowserImpersonationDetection parameter will be deleted.</li>.
* `name` - (Optional, String) Specifies the name of the browser spoofing identification rule.

The `rules` object of `client_attestation_rules` supports the following:

* `attester_id` - (Optional, String) Specifies the client authentication option ID.
* `condition` - (Optional, String) The rule content must comply with expression grammar. for details, see the product document.
* `device_profiles` - (Optional, List) Client device configuration. if the DeviceProfiles parameter value is not specified in the ClientAttestationRules parameter, keep the existing client device configuration and do not modify it.
* `enabled` - (Optional, String) Whether the rule is enabled. valid values: <li>`on`: enable</li> <li>`off`: disable</li>.
* `id` - (Optional, String) Client authentication rule ID. supported rule configuration operations by rule ID: <li> <b>add</b> a new rule: leave the ID empty or do not specify the ID parameter.</li> <li> <b>modify</b> an existing rule: specify the rule ID that needs to be updated/modified.</li> <li> <b>delete</b> an existing rule: existing rules not included in the ClientAttestationRule list under BotManagement parameters will be deleted.</li>.
* `invalid_attestation_action` - (Optional, List) Handling method for failed client authentication. valid values for SecurityAction Name: <li>Deny: block;</li> <li>Monitor: observation;</li> <li>Redirect: redirection;</li> <li>Challenge: Challenge.</li> default value: Monitor.
* `name` - (Optional, String) Specifies the name of the client authentication rule.
* `priority` - (Optional, Int) Priority of rules. a smaller value indicates higher priority execution. value range: 0-100. default value: 0.

The `rules` object of `custom_rules` supports the following:

* `action` - (Optional, List) The handling method for Bot custom rules. valid values: <li>Monitor: observation;</li><li>Deny: block, where DenyActionParameters.Name supports Deny and ReturnCustomPage;</li><li>Challenge: Challenge, where ChallengeActionParameters.Name supports JSChallenge and ManagedChallenge;</li><li>Redirect: Redirect to URL.</li>.
* `condition` - (Optional, String) Specifies the specific content of the Bot custom rule, which must comply with expression grammar. for detailed specifications, refer to the product document.
* `enabled` - (Optional, String) Whether the custom Bot rule is enabled. valid values: <li>on: enabled;</li><li>off: disabled.</li>.
* `id` - (Optional, String) The ID of a Bot custom rule. different rule configuration operations are supported by rule ID: <li><b>add</b> a new rule: leave the ID empty or do not specify the ID parameter.</li> <li><b>modify</b> an existing rule: specify the rule ID that needs to be updated/modified.</li> <li><b>delete</b> an existing rule: existing Rules not included in the Rules list under the BotManagementCustomRules parameter will be deleted.</li>.
* `name` - (Optional, String) Specifies the name of the Bot custom rule.
* `priority` - (Optional, Int) Priority of custom Bot rules. value range: 1-100. default value is 50.

The `rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Action for custom rules. The Name parameter of SecurityAction supports: <li>`Deny`: block;</li> <li>`Monitor`: observe;</li> <li>`ReturnCustomPage`: block with customized page;</li> <li>`Redirect`: Redirect to URL;</li> <li>`BlockIP`: IP blocking;</li> <li>`JSChallenge`: JavaScript challenge;</li> <li>`ManagedChallenge`: managed challenge;</li> <li>`Allow`: Allow.</li>.
* `condition` - (Required, String) The specifics of the custom rule, must comply with the expression grammar, please refer to product documentation for details.
* `enabled` - (Required, String) The custom rule status. Values: <li>`on`: enabled</li> <li>`off`: disabled</li>.
* `name` - (Required, String) The custom rule name.
* `id` - (Optional, String) Custom rule ID. <br>Different rule configuration operations are supported by rule ID: <br> Add a new rule: ID is empty or the ID parameter is not specified; <br> Modify an existing rule: specify the rule ID that needs to be updated/modified; <br> Delete an existing rule: existing rules not included in the Rules parameter will be deleted.
* `priority` - (Optional, Int) Customize the priority of custom rule. Range: 0-100, the default value is 0, this parameter only supports PreciseMatchRule.

The `rules` object of `exception_rules` supports the following:

* `condition` - (Optional, String) Describes the specific content of the exception rule, which must comply with the expression grammar. for details, please refer to the product document.
* `enabled` - (Optional, String) Whether the exception rule is enabled. valid values: <li>`on`: enable</li> <li>`off`: disable</li>.
* `id` - (Optional, String) The ID of the exception rule. different rule configuration operations are supported by rule ID: <li> <b>add</b> a new rule: leave the ID empty or do not specify the ID parameter.</li> <li> <b>modify</b> an existing rule: specify the rule ID that needs to be updated/modified.</li> <li> <b>delete</b> an existing rule: existing Rules not included in the Rules list under the ExceptionRules parameter will be deleted.</li>.
* `managed_rule_groups_for_exception` - (Optional, Set) A managed rule group with designated exception rules is valid only when SkipScope is ManagedRules, and at this point you cannot specify ManagedRulesForException.
* `managed_rules_for_exception` - (Optional, Set) Specifies the managed rule for the exception rule. valid only when SkipScope is ManagedRules. cannot specify ManagedRuleGroupsForException at this time.
* `name` - (Optional, String) The name of the exception rule.
* `request_fields_for_exception` - (Optional, List) Specify exception rules to skip request fields. valid only when SkipScope is ManagedRules and SkipOption is SkipOnSpecifiedRequestFields.
* `skip_option` - (Optional, String) Skip the specific type of request. valid values: <li>SkipOnAllRequestFields: skip all requests;</li> <li>SkipOnSpecifiedRequestFields: skip specified request fields.</li> valid only when SkipScope is ManagedRules.
* `skip_scope` - (Optional, String) Exception rule execution options, valid values: <li>WebSecurityModules: designate the security protection module for the exception rule.</li> <li>ManagedRules: designate the managed rule.</li>.
* `web_security_modules_for_exception` - (Optional, Set) Specifies the security protection module for exception rules. valid only when SkipScope is WebSecurityModules. valid values: <li>websec-mod-managed-rules: managed rule.</li><li>websec-mod-rate-limiting: rate limit.</li><li>websec-mod-custom-rules: custom rule.</li><li>websec-mod-adaptive-control: adaptive frequency control, intelligent client filtering, slow attack protection, traffic theft protection.</li><li>websec-mod-bot: bot management.</li>.

The `rules` object of `rate_limiting_rules` supports the following:

* `action_duration` - (Optional, String) The duration of an Action is only supported in the following units: <li>s: seconds, value range 1-120;</li> <li>m: minutes, value range 1-120;</li> <li>h: hours, value range 1-48;</li> <li>d: days, value range 1-30.</li>.
* `action` - (Optional, List) Precision rate limiting handling methods. valid values: <li>Monitor: Monitor;</li> <li>Deny: block, where DenyActionParameters.Name supports Deny and ReturnCustomPage;</li> <li>Challenge: Challenge, where ChallengeActionParameters.Name supports JSChallenge and ManagedChallenge;</li> <li>Redirect: Redirect to URL;</li>.
* `condition` - (Optional, String) The specific content of precise speed limit shall comply with the expression syntax. for detailed specifications, see the product documentation.
* `count_by` - (Optional, Set) Rate threshold request feature match mode. this field is required when Enabled is on.  when there are multiple conditions, composite multiple conditions will perform statistics count. the maximum number of conditions must not exceed 5. valid values: <li><b>http.request.ip</b>: client ip;</li> <li><b>http.request.xff_header_ip</b>: client ip (priority match xff header);</li> <li><b>http.request.uri.path</b>: request access path;</li> <li><b>http.request.cookies['session']</b>: Cookie named session, where session can be replaced with your own specified parameter;</li> <li><b>http.request.headers['user-agent']</b>: http header named user-agent, where user-agent can be replaced with your own specified parameter;</li> <li><b>http.request.ja3</b>: request ja3 fingerprint;</li> <li><b>http.request.uri.query['test']</b>: URL query parameter named test, where test can be replaced with your own specified parameter.</li>.
* `counting_period` - (Optional, String) Specifies the time window for statistics. valid values: <li>1s: 1 second;</li><li>5s: 5 seconds;</li><li>10s: 10 seconds;</li><li>20s: 20 seconds;</li><li>30s: 30 seconds;</li><li>40s: 40 seconds;</li><li>50s: 50 seconds;</li><li>1m: 1 minute;</li><li>2m: 2 minutes;</li><li>5m: 5 minutes;</li><li>10m: 10 minutes;</li><li>1h: 1 hour.</li>.
* `enabled` - (Optional, String) Whether the precise rate limiting rule is enabled. valid values: <li>on: enabled;</li> <li>off: disabled.</li>.
* `id` - (Optional, String) The ID of precise rate limiting. rule ID supports different rule configuration operations: <li><b>add</b> a new rule: leave the ID empty or do not specify the ID parameter.</li> <li><b>modify</b> an existing rule: specify the rule ID that needs to be updated/modified.</li> <li><b>delete</b> an existing rule: existing Rules not included in the Rules list under the RateLimitingRules parameter will be deleted.</li>.
* `max_request_threshold` - (Optional, Int) Precision rate limiting specifies the cumulative number of interceptions within the time range. value ranges from 1 to 100000.
* `name` - (Optional, String) Specifies the name of the precise rate limit.
* `priority` - (Optional, Int) Precision rate limiting specifies the priority. value range is 0 to 100. default is 0.

The `search_engine_bots` object of `basic_bot_settings` supports the following:

* `base_action` - (Optional, List) Specifies the action for requests from search engine crawlers. valid values for SecurityAction Name: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Disabled: not enabled, disable specified rule;</li> <li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li> <li>Allow: pass (to be deprecated).</li>.
* `bot_management_action_overrides` - (Optional, List) Specifies the handling method for search engine crawler requests.

The `security_action` object of `action` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `security_policy` object supports the following:

* `bot_management` - (Optional, List) Bot management configuration.
* `custom_rules` - (Optional, List) Custom rules. If the parameter is null or not filled, the configuration last set will be used by default.
Note: This field may return null, indicating that no valid value can be obtained.
* `exception_rules` - (Optional, List) Exception rule configuration.
* `http_ddos_protection` - (Optional, List) HTTP DDOS protection configuration.
* `managed_rules` - (Optional, List) Managed. If the parameter is null or not filled, the configuration last set will be used by default.
Note: This field may return null, indicating that no valid value can be obtained.
* `rate_limiting_rules` - (Optional, List) Configures the rate limiting rule.

The `session_expired_action` object of `bot_session_validation` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `session_invalid_action` object of `bot_session_validation` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `session_rate_control` object of `bot_session_validation` supports the following:

* `enabled` - (Optional, String) Specifies whether session rate and periodic feature verification are enabled. valid values: <li>on: enable</li><li>off: disable</li>.
* `high_rate_session_action` - (Optional, List) Session rate and periodic feature verification high-risk execution actions. SecurityAction Name valid values: <li>Deny: block, where Stall configuration is supported in DenyActionParameters;</li> <li>Monitor: observation;</li> <li>Allow: respond after wait, where MinDelayTime and MaxDelayTime configuration is required in AllowActionParameters.</li>.
* `low_rate_session_action` - (Optional, List) Session rate and periodic feature verification low risk execution action. SecurityAction Name parameter supports: <li>Deny: block, where DenyActionParameters supports Stall configuration;</li><li>Monitor: observe;</li><li>Allow: respond after wait, where AllowActionParameters requires MinDelayTime and MaxDelayTime configuration.</li>.
* `mid_rate_session_action` - (Optional, List) Session rate and periodic feature verification medium-risk execution action. SecurityAction Name parameter supports: <li>Deny: block, where DenyActionParameters supports Stall configuration;</li><li>Monitor: observe;</li><li>Allow: respond after wait, where AllowActionParameters requires MinDelayTime and MaxDelayTime configuration.</li>.

The `slow_attack_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether slow attack protection is enabled. valid values: <li>on: enabled;</li> <li>off: disabled.</li>.
* `action` - (Optional, List) Slow attack protection handling method. required when Enabled is on. valid values for SecurityAction Name: <li>Monitor: observation;</li> <li>Deny: block;</li>.
* `minimal_request_body_transfer_rate` - (Optional, List) The specific configuration of the minimum body transfer rate threshold is required when Enabled is on.
* `request_body_transfer_timeout` - (Optional, List) Specifies the specific configuration of body transfer timeout duration. required when Enabled is on.

The `source_idc` object of `basic_bot_settings` supports the following:

* `base_action` - (Optional, List) Handling method for requests from the specified IDC. valid values for SecurityAction Name: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>Disabled: not enabled, disable specified rule;</li> <li>Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;</li> <li>Allow: pass (to be deprecated).</li>.
* `bot_management_action_overrides` - (Optional, List) Specifies the handling method for the specified id request.

The `verified_bot_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:.
<Li>Deny. specifies to block requests from accessing site resources.</li>.
<Li>Monitor: observation, only record logs.</li>.
<li>Redirect: Redirect to URL.</li>.
<Li>Disabled: specifies that the rule is not enabled.</li>.
<Li>Allow: specifies whether to allow access with delayed processing of requests.</li>.
<Li>Challenge: specifies the challenge content to respond to.</li>.
<Li>Trans: pass and allow requests to directly access site resources.</li>.
<Li>BlockIP: to be deprecated. ip block.</li>.
<Li>ReturnCustomPage: to be deprecated. use specified page for interception.</li>.
<li>JSChallenge: to be deprecated, JavaScript challenge;</li>.
<Li>ManagedChallenge: to be deprecated. managed challenge.</li>.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `template_id` - Template ID.


## Import

TEO web security template can be imported using the zoneId#templateId, e.g.

```
terraform import tencentcloud_teo_web_security_template.example zone-3fkff38fyw8s#temp-p3p973nu
```

