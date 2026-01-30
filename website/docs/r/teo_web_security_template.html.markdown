---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_web_security_template"
sidebar_current: "docs-tencentcloud-resource-teo_web_security_template"
description: |-
  Provides a resource to create a teo web security template
---

# tencentcloud_teo_web_security_template

Provides a resource to create a teo web security template

~> **NOTE:** The current resources do not support managed_rule_groups.

## Example Usage

### Basic usage

```hcl
resource "tencentcloud_teo_web_security_template" "web_security_template" {
  template_name = "tf-test"
  zone_id       = "zone-3fkff38fyw8s"
  security_policy {
    bot_management {
      enabled = "on"
      basic_bot_settings {
        bot_intelligence {
          enabled = "on"
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
          enabled = "on"
          ip_reputation_group {
          }
        }
        known_bot_categories {
          bot_management_action_overrides {
            ids = ["9395241960"]
            action {
              name = "Allow"
            }
          }
        }
        search_engine_bots {
          bot_management_action_overrides {
            ids = ["9126905504"]
            action {
              name = "Deny"
            }
          }
        }
        source_idc {
          bot_management_action_overrides {
            ids = ["8868370049", "8868370048"]
            action {
              name = "Deny"
            }
          }
        }
      }
      browser_impersonation_detection {
        rules {
          condition = "$${http.request.uri.path} like ['/*'] and $${http.request.method} in ['get']"
          enabled   = "on"
          name      = "Block Non-Browser Crawler Access"
          action {
            bot_session_validation {
              issue_new_bot_session_cookie = "on"
              max_new_session_trigger_config {
                max_new_session_count_interval  = "10s"
                max_new_session_count_threshold = 300
              }
              session_expired_action {
                name = "Deny"
              }
              session_invalid_action {
                name = "Deny"
                deny_action_parameters {
                  block_ip           = null
                  block_ip_duration  = null
                  error_page_id      = null
                  response_code      = null
                  return_custom_page = null
                  stall              = "on"
                }
              }
              session_rate_control {
                enabled = "off"
              }
            }
          }
        }
      }
      client_attestation_rules {
      }
      custom_rules {
        rules {
          condition = "$${http.request.ip} in ['222.22.22.0/24'] and $${http.request.headers['user-agent']} contain ['cURL']"
          enabled   = "on"
          name      = "Login API Request Surge Protection"
          priority  = 50
          action {
            weight = 100
            security_action {
              name = "Deny"
              deny_action_parameters {
                block_ip           = null
                block_ip_duration  = null
                error_page_id      = null
                response_code      = null
                return_custom_page = null
                stall              = "on"
              }
            }
          }
        }
      }
    }
    custom_rules {
      rules {
        condition = "$${http.request.headers['user-agent']} contain ['curl/','Wget/','ApacheBench/']"
        enabled   = "on"
        name      = "Malicious User-Agent Blacklist"
        priority  = 50
        rule_type = "PreciseMatchRule"
        action {
          name = "JSChallenge"
        }
      }
      rules {
        condition = "$${http.request.ip} in ['36']"
        enabled   = "on"
        name      = "Custom Rule"
        priority  = 0
        rule_type = "BasicAccessRule"
        action {
          name = "Monitor"
        }
      }
    }
    exception_rules {
      rules {
        condition                          = "$${http.request.method} in ['post'] and $${http.request.uri.path} in ['/api/EventLogUpload']"
        enabled                            = "on"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "High Frequency API Skip Rate Limit 1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control"]
      }
      rules {
        condition                          = "$${http.request.ip} in ['123.123.123.0/24']"
        enabled                            = "on"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "IP Whitelist 1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control", "websec-mod-bot", "websec-mod-custom-rules", "websec-mod-managed-rules", "websec-mod-rate-limiting"]
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
        enabled = "on"
        action {
          name = "Deny"
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
        enabled = "on"
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
    rate_limiting_rules {
      rules {
        action_duration       = "30m"
        condition             = "$${http.request.uri.path} contain ['/checkout/submit']"
        count_by              = ["http.request.ip"]
        counting_period       = "60s"
        enabled               = "on"
        max_request_threshold = 300
        name                  = "Single IP Request Rate Limit 1"
        priority              = 50
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) Policy template name. Composed of Chinese characters, letters, digits, and underscores. Cannot begin with an underscore and must be less than or equal to 32 characters.
* `zone_id` - (Required, String) Zone ID. Explicitly identifies the zone to which the policy template belongs for access control purposes.
* `security_policy` - (Optional, List) Web security policy template configuration. Generates default config if empty. Supported: Exception rules, custom rules, rate limiting rules, managed rules. Not supported: Bot management rules (under development).

The `action` object of `adaptive_frequency_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `bandwidth_abuse_defense` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `bot_management_action_overrides` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `client_filtering` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `frequent_scanning_protection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
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

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `action` object of `rules` supports the following:

* `security_action` - (Optional, List) The handling method of the Bot custom rule. valid values: - Allow: pass, where AllowActionParameters supports MinDelayTime and MaxDelayTime configuration; - Deny: block, where DenyActionParameters supports BlockIp, ReturnCustomPage, and Stall configuration; - Monitor: observation; - Challenge: Challenge, where ChallengeActionParameters.ChallengeOption supports JSChallenge and ManagedChallenge; - Redirect: Redirect to URL.
* `weight` - (Optional, Int) The Weight of the current SecurityAction, only supported between 10 and 100 and must be a multiple of 10. the total of all Weight parameters must equal 100.

The `action` object of `slow_attack_defense` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `adaptive_frequency_control` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether adaptive frequency control is enabled. valid values: - on: enable; - off: disable.
* `action` - (Optional, List) The handling method of adaptive frequency control. this field is required when Enabled is on. valid values for SecurityAction Name: - Monitor: observation; - Deny: block; - Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge.
* `sensitivity` - (Optional, String) The restriction level of adaptive frequency control. required when Enabled is on. valid values: - Loose: Loose- Moderate: Moderate- Strict: Strict.

The `allow_action_parameters` object of `action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `base_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `bot_client_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `challenge_not_finished_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `challenge_timeout_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `high_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `high_risk_request_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `human_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `invalid_attestation_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `likely_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `low_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `medium_risk_request_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `mid_rate_session_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `security_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `session_expired_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `session_invalid_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `allow_action_parameters` object of `verified_bot_requests_action` supports the following:

* `max_delay_time` - (Optional, String) Maximum delayed response time. supported units: - s: seconds, value ranges from 5 to 10.
* `min_delay_time` - (Optional, String) Minimum latency response time. when configured as 0s, it means no delay for direct response. supported units: - s: seconds, value ranges from 0 to 5.

The `auto_update` object of `managed_rules` supports the following:

* `auto_update_to_latest_version` - (Required, String) Enable automatic update to the latest version or not. Values: - `on`: enabled - `off`: disabled.
* `ruleset_version` - (Optional, String) Current version, compliant with ISO 8601 standard format, such as 2023-12-21T12:00:32Z, empty by default, output parameter only.

The `bandwidth_abuse_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether bandwidth abuse protection (applicable to chinese mainland only) is enabled. valid values: - on: enabled; - off: disabled.
* `action` - (Optional, List) Bandwidth abuse protection (applicable to chinese mainland) handling method. required when Enabled is on. valid values for SecurityAction Name: - Monitor: observe; - Deny: block; - Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge.

The `base_action` object of `ip_reputation_group` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `base_action` object of `known_bot_categories` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `base_action` object of `search_engine_bots` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `base_action` object of `source_idc` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
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

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `base_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `bot_client_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `challenge_not_finished_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `challenge_timeout_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `high_rate_session_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `high_risk_request_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `human_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `invalid_attestation_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `likely_bot_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `low_rate_session_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `medium_risk_request_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `mid_rate_session_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `security_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `session_expired_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `session_invalid_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `block_ip_action_parameters` object of `verified_bot_requests_action` supports the following:

* `duration` - (Required, String) Penalty duration for `BlockIP`. Units: - `s`: second, value range 1-120; - `m`: minute, value range 1-120; - `h`: hour, value range 1-48.

The `bot_client_action` object of `client_behavior_detection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `bot_intelligence` object of `basic_bot_settings` supports the following:

* `bot_ratings` - (Optional, List) Based on client and request features, divides request sources into human requests, legitimate Bot requests, suspected Bot requests, and high-risk Bot requests, and provides request handling options.
* `enabled` - (Optional, String) Specifies the switch for Bot intelligent analysis configuration. valid values:.  on: enabled. off: disabled.

The `bot_management_action_overrides` object of `ip_reputation_group` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: - Deny: block;- Monitor: observe;- Disabled: Disabled, disable the specified rule;- Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;- Allow: pass (only for Bot basic feature management).
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management_action_overrides` object of `known_bot_categories` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: - Deny: block;- Monitor: observe;- Disabled: Disabled, disable the specified rule;- Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;- Allow: pass (only for Bot basic feature management).
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management_action_overrides` object of `search_engine_bots` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: - Deny: block;- Monitor: observe;- Disabled: Disabled, disable the specified rule;- Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;- Allow: pass (only for Bot basic feature management).
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management_action_overrides` object of `source_idc` supports the following:

* `action` - (Optional, List) Specifies the handling action for Bot rule items in Ids. valid values for the Name parameter in SecurityAction: - Deny: block;- Monitor: observe;- Disabled: Disabled, disable the specified rule;- Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge;- Allow: pass (only for Bot basic feature management).
* `ids` - (Optional, Set) Specific item under Bot rules used to rewrite the configuration content of this single rule. refer to the returned message from the DescribeBotManagedRules API for detailed information corresponding to Ids.

The `bot_management` object of `security_policy` supports the following:

* `basic_bot_settings` - (Optional, List) Bot management basic configuration. takes effect on all domains associated with the policy. can be customized through CustomRules.
* `browser_impersonation_detection` - (Optional, List) Configures browser spoofing identification rules (formerly active feature detection rule). sets the response page range for JavaScript injection, browser check options, and handling method for non-browser clients.
* `client_attestation_rules` - (Optional, List) Definition list of client authentication rules. this feature is in beta test. submit a ticket if you need to use it.
* `custom_rules` - (Optional, List) Bot management custom rule combines various crawlers and request behavior characteristics to accurately define bots and configure customized handling methods.
* `enabled` - (Optional, String) Whether Bot management is enabled. valid values: - on: enabled;- off: disabled.

The `bot_ratings` object of `bot_intelligence` supports the following:

* `high_risk_bot_requests_action` - (Optional, List) Execution action for malicious Bot requests. valid values for the Name parameter in SecurityAction: - Deny: block; - Monitor: observe; - Allow: pass; - Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.
* `human_requests_action` - (Optional, List) Execution action for a normal Bot request. valid values for the Name parameter in SecurityAction: - Allow: pass.
* `likely_bot_requests_action` - (Optional, List) The execution action for suspected Bot requests. valid values for the Name parameter in SecurityAction: - Deny: block; - Monitor: observe; - Allow: pass; - Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.
* `verified_bot_requests_action` - (Optional, List) Execution action for friendly Bot request. SecurityAction Name parameter supports: - Deny: block;- Monitor: observe;- Allow: pass;- Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.

The `bot_session_validation` object of `action` supports the following:

* `issue_new_bot_session_cookie` - (Optional, String) Whether to update Cookie and validate. valid values: - on: update Cookie and validate; - off: verify only.
* `max_new_session_trigger_config` - (Optional, List) Specifies the trigger threshold for updating and validating cookies. valid only when IssueNewBotSessionCookie is set to on.
* `session_expired_action` - (Optional, List) Execution action when no Cookie is carried or the Cookie expired. valid values for the Name parameter in SecurityAction: - Deny: block, where Stall can be configured in DenyActionParameters;- Monitor: observe;- Allow: respond after wait, where MinDelayTime and MaxDelayTime must be configured in AllowActionParameters.
* `session_invalid_action` - (Optional, List) Execution action for invalid Cookie. valid values for the Name parameter in SecurityAction: - Deny: block, where the DenyActionParameters supports Stall configuration;- Monitor: observe;- Allow: respond after wait, where AllowActionParameters requires MinDelayTime and MaxDelayTime configuration.
* `session_rate_control` - (Optional, List) Specifies the session rate and periodic feature verification configuration.

The `browser_impersonation_detection` object of `bot_management` supports the following:

* `rules` - (Optional, List) List of browser spoofing identification Rules. When using ModifySecurityPolicy to modify Web protection configuration: - if Rules parameter in SecurityPolicy.BotManagement.BrowserImpersonationDetection is not specified or parameter length is zero: clear all browser spoofing identification rule configurations; - if BrowserImpersonationDetection parameter value is unspecified in SecurityPolicy.BotManagement parameters: keep existing browser spoofing identification rule configurations without modification.

The `challenge_action_parameters` object of `action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `base_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `bot_client_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `challenge_not_finished_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `challenge_timeout_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `high_rate_session_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `high_risk_request_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `human_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `invalid_attestation_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `likely_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `low_rate_session_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `medium_risk_request_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `mid_rate_session_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `security_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `session_expired_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `session_invalid_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_action_parameters` object of `verified_bot_requests_action` supports the following:

* `challenge_option` - (Required, String) Safe execution challenge actions. valid values: -  InterstitialChallenge: interstitial challenge; -  InlineChallenge: embedded challenge; -  JSChallenge: JavaScript challenge; -  ManagedChallenge: managed challenge.
* `attester_id` - (Optional, String) Client authentication method ID. this field is required when Name is InterstitialChallenge/InlineChallenge.
* `interval` - (Optional, String) Specifies the time interval for challenge repetition. this field is required when Name is InterstitialChallenge/InlineChallenge. default value is 300s. supported units: - s: second, value ranges from 1 to 60;- m: minute, value ranges from 1 to 60;- h: hour, value ranges from 1 to 24.

The `challenge_not_finished_action` object of `client_behavior_detection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `challenge_timeout_action` object of `client_behavior_detection` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `client_attestation_rules` object of `bot_management` supports the following:

* `rules` - (Optional, List) List of client authentication. when using ModifySecurityPolicy to modify Web protection configuration: -  if Rules in SecurityPolicy.BotManagement.ClientAttestationRules is not specified or the parameter length of Rules is zero: clear all client authentication rule configuration.  -  if ClientAttestationRules in SecurityPolicy.BotManagement parameters is unspecified: keep existing client authentication rule configuration and do not modify..

The `client_behavior_detection` object of `action` supports the following:

* `bot_client_action` - (Optional, List) The execution action of the Bot client. valid values for the Name parameter in SecurityAction: - Deny: block, where the Stall configuration is supported in DenyActionParameters;- Monitor: observation;- Allow: respond after wait, where MinDelayTime and MaxDelayTime configurations are required in AllowActionParameters.
* `challenge_not_finished_action` - (Optional, List) Execution action when client-side javascript is not enabled (test not completed). valid values for SecurityAction Name: - Deny: block, where Stall configuration is supported in DenyActionParameters;- Monitor: observe;- Allow: respond after waiting, where MinDelayTime and MaxDelayTime configuration is required in AllowActionParameters.
* `challenge_timeout_action` - (Optional, List) The execution action for client-side detection timeout. valid values for the Name parameter in SecurityAction: - Deny: block, where Stall can be configured in DenyActionParameters; - Monitor: observe; - Allow: respond after wait, where MinDelayTime and MaxDelayTime must be configured in AllowActionParameters.
* `crypto_challenge_delay_before` - (Optional, String) Specifies the execution mode for client behavior verification. valid values: - 0ms: immediate execution; - 100ms: delay 100ms execution; - 200ms: delay 200ms execution; - 300ms: delay 300ms execution; - 400ms: delay 400ms execution; - 500ms: delay 500ms execution; - 600ms: delay 600ms execution; - 700ms: delay 700ms execution; - 800ms: delay 800ms execution; - 900ms: delay 900ms execution; - 1000ms: delay 1000ms execution.
* `crypto_challenge_intensity` - (Optional, String) Specifies the proof-of-work strength. valid values: - low: low;- medium: medium;- high: high.
* `max_challenge_count_interval` - (Optional, String) Time window for trigger threshold statistics. valid values: - 5s: within 5 seconds;- 10s: within 10 seconds;- 15s: within 15 seconds;- 30s: within 30 seconds;- 60s: within 60 seconds;- 5m: within 5 minutes;- 10m: within 10 minutes;- 30m: within 30 minutes;- 60m: within 60 minutes.
* `max_challenge_count_threshold` - (Optional, Int) Trigger threshold cumulative count. value range: 1-100000000.

The `client_filtering` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether intelligent client filtering is enabled. valid values: - on: enable; - off: disable.
* `action` - (Optional, List) The handling method of intelligent client filtering. when Enabled is on, this field is required. the Name parameter of SecurityAction supports: - Monitor: observation; - Deny: block; - Challenge: Challenge, where ChallengeActionParameters.Name only supports JSChallenge.

The `custom_rules` object of `bot_management` supports the following:

* `rules` - (Optional, List) List of Bot custom Rules. When using ModifySecurityPolicy to modify Web protection configuration: - if Rules in SecurityPolicy.BotManagement.CustomRules is not specified or parameter length of Rules is zero: clear all Bot custom rule configurations; - if CustomRules in SecurityPolicy.BotManagement parameters is unspecified: keep existing Bot custom rule configurations and do not modify them.

The `custom_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) The custom rule. When modifying the Web protection configuration using ModifySecurityPolicy: - if the Rules parameter is not specified or the parameter length of Rules is zero: clear all custom rule configurations; - if the Rules parameter is not specified: keep the existing custom rule configuration without modification.

The `deny_action_parameters` object of `action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated.Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable; - off: Disable. Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously.Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `base_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `bot_client_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `challenge_not_finished_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `challenge_timeout_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `high_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `high_risk_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `high_risk_request_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `human_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `invalid_attestation_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `likely_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `low_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `medium_risk_request_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `mid_rate_session_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `security_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `session_expired_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `session_invalid_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `deny_action_parameters` object of `verified_bot_requests_action` supports the following:

* `block_ip_duration` - (Optional, String) The ban duration when BlockIP is on.
* `block_ip` - (Optional, String) Specifies whether to extend the ban on the source IP. valid values. - `on`: Enable;  - off: Disable.  After enabled, continuously blocks client ips that trigger the rule. when this option is enabled, the BlockIpDuration parameter must be simultaneously designated. Note: this option cannot intersect with ReturnCustomPage or Stall.
* `error_page_id` - (Optional, String) Specifies the page id of the custom page.
* `response_code` - (Optional, String) Status code of the custom page.
* `return_custom_page` - (Optional, String) Specifies whether to use a custom page. valid values:. - `on`: Enable;  - off: Disable.  Enabled, use custom page content to intercept requests. when this option is enabled, ResponseCode and ErrorPageId parameters must be specified simultaneously. Note: this option cannot intersect with the BlockIp or Stall option.
* `stall` - (Optional, String) Specifies whether to suspend the request source without processing. valid values:. - `on`: Enable;  - off: Disable.  Enabled, no longer responds to requests in the current connection session and does not actively disconnect. used for crawler combat to consume client connection resources. Note: this option cannot intersect with BlockIp or ReturnCustomPage options.

The `device_profiles` object of `rules` supports the following:

* `client_type` - (Required, String) Client device type. valid values: - iOS; - Android; - WebView.
* `high_risk_min_score` - (Optional, Int) The minimum value to determine a request as high-risk ranges from 1-99. the larger the value, the higher the request risk, and the closer it resembles a request initiated by a Bot client. the default value is 50, corresponding to high-risk for values 51-100.
* `high_risk_request_action` - (Optional, List) Handling method for high-risk requests. valid values for SecurityAction Name: - Deny: block; - Monitor: observation; - Redirect: redirection; - Challenge: Challenge. default value: Monitor.
* `medium_risk_min_score` - (Optional, Int) Specifies the minimum value to determine a request as medium-risk. value range: 1-99. the larger the value, the higher the request risk, resembling requests initiated by a Bot client. default value: 15, corresponding to medium-risk for values 16-50.
* `medium_risk_request_action` - (Optional, List) Handling method for medium-risk requests. SecurityAction Name parameter supports: - Deny: block; - Monitor: observe; - Redirect: Redirect; - Challenge: Challenge. default value is Monitor.

The `exception_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) Definition list of exception Rules. when using ModifySecurityPolicy to modify Web protection configuration: - if the Rules parameter is not specified or the parameter length is zero: clear all exception rule configurations.- if the ExceptionRules parameter value is not specified in SecurityPolicy: keep existing exception rule configurations without modification.

The `frequent_scanning_protection` object of `managed_rules` supports the following:

* `action_duration` - (Optional, String) This parameter specifies the duration of the handling Action set by the high frequency scan protection Action parameter. value range: 60 to 86400. measurement unit: seconds (s) only, for example 60s. this field is required when Enabled is on.
* `action` - (Optional, List) The handling action for high-frequency scan protection. required when Enabled is on. valid values for SecurityAction Name: - Deny: block and respond with an interception page; - Monitor: observe without processing requests, log security events in logs; - JSChallenge: respond with a JavaScript challenge page.
* `block_threshold` - (Optional, Int) This parameter specifies the threshold for high-frequency scan protection, which is the intercept count of managed rules set to interception within the time range set by CountingPeriod. value range: 1 to 4294967294, for example 100. when exceeding this statistical value, subsequent requests will trigger the handling Action set by Action. required when Enabled is on.
* `count_by` - (Optional, String) The match mode for request statistics. required when Enabled is on. valid values: - http.request.xff_header_ip: client ip (priority match xff header);- http.request.ip: client ip.
* `counting_period` - (Optional, String) This parameter specifies the statistical time window for high-frequency scan protection, which is the time window for counting requests that hit managed rules configured as block. valid values: 5-1800. measurement unit: seconds (s) only, such as 5s. this field is required when Enabled is on.
* `enabled` - (Optional, String) Whether the high-frequency scan protection rule is enabled. valid values: - on: enable. the high-frequency scan protection rule takes effect.- off: disable. the high-frequency scan protection rule does not take effect.

The `high_rate_session_action` object of `session_rate_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `high_risk_bot_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `high_risk_request_action` object of `device_profiles` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
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

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `invalid_attestation_action` object of `rules` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `ip_reputation_group` object of `ip_reputation` supports the following:

* `base_action` - (Optional, List) Execution action of the IP intelligence library (formerly client profile analysis). SecurityAction Name parameter supports: - Deny: block; - Monitor: observe; - Disabled: not enabled, disable specified rule; - Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge.
* `bot_management_action_overrides` - (Optional, List) The specific configuration of the IP intelligence library (originally client profile analysis), used to override the default configuration in BaseAction. among them, the Ids in BotManagementActionOverrides can be filled with: - IPREP_WEB_AND_DDOS_ATTACKERS_LOW: network attack - general confidence; - IPREP_WEB_AND_DDOS_ATTACKERS_MID: network attack - medium confidence; - IPREP_WEB_AND_DDOS_ATTACKERS_HIGH: network attack - HIGH confidence; - IPREP_PROXIES_AND_ANONYMIZERS_LOW: network proxy - general confidence; - IPREP_PROXIES_AND_ANONYMIZERS_MID: network proxy - medium confidence; - IPREP_PROXIES_AND_ANONYMIZERS_HIGH: network proxy - HIGH confidence; - IPREP_SCANNING_TOOLS_LOW: scanner - general confidence; - IPREP_SCANNING_TOOLS_MID: scanner - medium confidence; - IPREP_SCANNING_TOOLS_HIGH: scanner - HIGH confidence; - IPREP_ATO_ATTACKERS_LOW: account takeover attack - general confidence; - IPREP_ATO_ATTACKERS_MID: account takeover attack - medium confidence; - IPREP_ATO_ATTACKERS_HIGH: account takeover attack - HIGH confidence; - IPREP_WEB_SCRAPERS_AND_TRAFFIC_BOTS_LOW: malicious BOT - general confidence; - IPREP_WEB_SCRAPERS_AND_TRAFFIC_BOTS_MID: malicious BOT - medium confidence; - IPREP_WEB_SCRAPERS_AND_TRAFFIC_BOTS_HIGH: malicious BOT - HIGH confidence.

The `ip_reputation` object of `basic_bot_settings` supports the following:

* `enabled` - (Optional, String) IP intelligence library (formerly client profile analysis). valid values: - on: enable; - off: disable.
* `ip_reputation_group` - (Optional, List) IP intelligence library (formerly client profile analysis) configuration content.

The `known_bot_categories` object of `basic_bot_settings` supports the following:

* `base_action` - (Optional, List) Handling method for access requests from known commercial tools or open-source tools. specifies the Name parameter value of SecurityAction: - Deny: block; - Monitor: observe; - Disabled: not enabled, disable specified rule; - Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge; - Allow: pass (to be deprecated).
* `bot_management_action_overrides` - (Optional, List) Specifies the handling method for access requests from known commercial tools or open-source tools.

The `likely_bot_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `low_rate_session_action` object of `session_rate_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `managed_rules` object of `security_policy` supports the following:

* `detection_only` - (Required, String) Evaluation mode is enabled or not, it is valid only when the `Enabled` parameter is set to `on`. Values: - `on`: enabled, all managed rules take effect in `observe` mode. - off: disabled, all managed rules take effect according to the specified configuration.
* `enabled` - (Required, String) The managed rule status. Values: - `on`: enabled, all managed rules take effect as configured; - `off`: disabled, all managed rules do not take effect.
* `auto_update` - (Optional, List) Managed rule automatic update option.
* `frequent_scanning_protection` - (Optional, List) High-Frequency scan protection configuration option. when a visitor's frequent requests hit the managed rule configured as block within a period of time, all requests from that visitor are blocked.
* `semantic_analysis` - (Optional, String) Managed rule semantic analysis is enabled or not, it is valid only when the `Enabled` parameter is `on`. Values: `on`: enabled, perform semantic analysis before processing requests; `off`: disabled, process requests directly without semantic analysis. The default value is `off`.

The `max_new_session_trigger_config` object of `bot_session_validation` supports the following:

* `max_new_session_count_interval` - (Optional, String) Time window for trigger threshold statistics. valid values: - 5s: within 5 seconds;- 10s: within 10 seconds;- 15s: within 15 seconds;- 30s: within 30 seconds;- 60s: within 60 seconds;- 5m: within 5 minutes;- 10m: within 10 minutes;- 30m: within 30 minutes;- 60m: within 60 minutes.
* `max_new_session_count_threshold` - (Optional, Int) Trigger threshold cumulative count. value range: 1-100000000.

The `medium_risk_request_action` object of `device_profiles` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `mid_rate_session_action` object of `session_rate_control` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `minimal_request_body_transfer_rate` object of `slow_attack_defense` supports the following:

* `counting_period` - (Required, String) Minimum body transfer rate statistical time range, valid values: - 10s: 10 seconds; - 30s: 30 seconds; - 60s: 60 seconds; - 120s: 120 seconds.
* `enabled` - (Required, String) Specifies whether the minimum body transfer rate threshold is enabled. valid values: - on: enable; - off: disable.
* `minimal_avg_transfer_rate_threshold` - (Required, String) Minimum body transfer rate threshold, the measurement unit is only supported in bps.

The `rate_limiting_rules` object of `security_policy` supports the following:

* `rules` - (Optional, List) Definition list of precise rate limiting. When using ModifySecurityPolicy to modify the Web protection configuration: - if the Rules parameter is not specified or its length is zero: clear all precision rate limiting configurations; - if the RateLimitingRules parameter value is unspecified in the SecurityPolicy parameter: retain the existing custom rule configuration without modification.

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

* `enabled` - (Required, String) Whether body transfer timeout is enabled. valid values: - `on`: enable - `off`: disable.
* `idle_timeout` - (Required, String) Body transfer timeout duration. valid values: 5-120. measurement unit: seconds (s) only.

The `request_fields_for_exception` object of `rules` supports the following:

* `condition` - (Required, String) Skip specific field expression must comply with expression grammar. Condition supports expression configuration syntax: -  write according to the matching conditional expression syntax of rules, with support for referencing key and value. -  supports in, like operators, and logical combination with and. For example: - ${key} in ['x-trace-id']: the parameter name equals x-trace-id. - ${key} in ['x-trace-id'] and ${value} like ['Bearer *']: the parameter name equals x-trace-id and the parameter value wildcard matches Bearer *.
* `scope` - (Required, String) Skip specific field. supported values:. - body.json: parameter content in json requests. at this point, Condition supports key and value, TargetField supports key and value, for example { "Scope": "body.json", "Condition": "", "TargetField": "key" }, which means all parameters in json requests skip WAF scan. cookie: cookie; at this point Condition supports key, value, TargetField supports key, value, for example { "Scope": "cookie", "Condition": "${key} in ['account-id'] and ${value} like ['prefix-*']", "TargetField": "value" }, which means the cookie parameter name equals account-id and the parameter value wildcard matches prefix-* to skip WAF scan;. header: HTTP header parameters. at this point, Condition supports key and value, TargetField supports key and value, for example { "Scope": "header", "Condition": "${key} like ['x-auth-*']", "TargetField": "value" }, which means header parameter name wildcard match x-auth-* skips WAF scan. uri.query: URL encoding content/query parameter. at this point, Condition supports key and value, TargetField supports key and value. example: { "Scope": "uri.query", "Condition": "${key} in ['action'] and ${value} in ['upload', 'delete']", "TargetField": "value" }. indicates URL encoding content/query parameter name equal to action and parameter value equal to upload or delete skips WAF scan. uri: specifies the request path uri. at this point, Condition must be empty. TargetField supports query, path, fullpath, such as {"Scope": "uri", "Condition": "", "TargetField": "query"}, indicates the request path uri skips WAF scan for query parameters. body: request body content. at this point Condition must be empty, TargetField supports fullbody, multipart, such as { "Scope": "body", "Condition": "", "TargetField": "fullbody" }, which means the request body content skips WAF scan as a full request.
* `target_field` - (Required, String) The Scope parameter takes different values. the TargetField expression supports the following values:. -  body.json: supports key, value. - cookie: supports key and value. - header: supports key, value. -  uri.query: supports key and value. - uri. specifies path, query, or fullpath. - Body: supports fullbody and multipart.

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

The `rules` object of `browser_impersonation_detection` supports the following:

* `action` - (Optional, List) Describes the handling method for browser spoofing identification rules, including Cookie verification, session tracking configuration, and client behavior validation configuration.
* `condition` - (Optional, String) Specifies the specific content of browser spoofing identification rules, which only support configuration of request Method (Method), request Path (Path), and request URL, and must comply with expression grammar. for detailed specifications, please refer to the product document.
* `enabled` - (Optional, String) Whether browser spoofing detection is enabled. valid values: - on: enabled;- off: disabled.
* `name` - (Optional, String) Specifies the name of the browser spoofing identification rule.

The `rules` object of `client_attestation_rules` supports the following:

* `attester_id` - (Optional, String) Specifies the client authentication option ID.
* `condition` - (Optional, String) The rule content must comply with expression grammar. for details, see the product document.
* `device_profiles` - (Optional, List) Client device configuration. if the DeviceProfiles parameter value is not specified in the ClientAttestationRules parameter, keep the existing client device configuration and do not modify it.
* `enabled` - (Optional, String) Whether the rule is enabled. valid values: - `on`: enable - `off`: disable.
* `invalid_attestation_action` - (Optional, List) Handling method for failed client authentication. valid values for SecurityAction Name: - Deny: block; - Monitor: observation; - Redirect: redirection; - Challenge: Challenge. default value: Monitor.
* `name` - (Optional, String) Specifies the name of the client authentication rule.
* `priority` - (Optional, Int) Priority of rules. a smaller value indicates higher priority execution. value range: 0-100. default value: 0.

The `rules` object of `custom_rules` supports the following:

* `action` - (Optional, List) The handling method for Bot custom rules. valid values: - Monitor: observation;- Deny: block, where DenyActionParameters.Name supports Deny and ReturnCustomPage;- Challenge: Challenge, where ChallengeActionParameters.Name supports JSChallenge and ManagedChallenge;- Redirect: Redirect to URL.
* `condition` - (Optional, String) Specifies the specific content of the Bot custom rule, which must comply with expression grammar. for detailed specifications, refer to the product document.
* `enabled` - (Optional, String) Whether the custom Bot rule is enabled. valid values: - on: enabled;- off: disabled.
* `name` - (Optional, String) Specifies the name of the Bot custom rule.
* `priority` - (Optional, Int) Priority of custom Bot rules. value range: 1-100. default value is 50.

The `rules` object of `custom_rules` supports the following:

* `action` - (Required, List) Action for custom rules. The Name parameter of SecurityAction supports: - `Deny`: block; - `Monitor`: observe; - `ReturnCustomPage`: block with customized page; - `Redirect`: Redirect to URL; - `BlockIP`: IP blocking; - `JSChallenge`: JavaScript challenge; - `ManagedChallenge`: managed challenge; - `Allow`: Allow.
* `condition` - (Required, String) The specifics of the custom rule, must comply with the expression grammar, please refer to product documentation for details.
* `enabled` - (Required, String) The custom rule status. Values: - `on`: enabled - `off`: disabled.
* `name` - (Required, String) The custom rule name.
* `priority` - (Optional, Int) Customize the priority of custom rule. Range: 0-100, the default value is 0, this parameter only supports PreciseMatchRule.
* `rule_type` - (Optional, String) Type of custom rule. Values: - `BasicAccessRule`: basic access control; - `PreciseMatchRule`: exact custom rule, default; - `ManagedAccessRule`: expert customized rule, output parameter only.The default value is PreciseMatchRule.

The `rules` object of `exception_rules` supports the following:

* `condition` - (Optional, String) Describes the specific content of the exception rule, which must comply with the expression grammar. for details, please refer to the product document.
* `enabled` - (Optional, String) Whether the exception rule is enabled. valid values: - `on`: enable - `off`: disable.
* `managed_rule_groups_for_exception` - (Optional, Set) A managed rule group with designated exception rules is valid only when SkipScope is ManagedRules, and at this point you cannot specify ManagedRulesForException.
* `managed_rules_for_exception` - (Optional, Set) Specifies the managed rule for the exception rule. valid only when SkipScope is ManagedRules. cannot specify ManagedRuleGroupsForException at this time.
* `name` - (Optional, String) The name of the exception rule.
* `request_fields_for_exception` - (Optional, List) Specify exception rules to skip request fields. valid only when SkipScope is ManagedRules and SkipOption is SkipOnSpecifiedRequestFields.
* `skip_option` - (Optional, String) Skip the specific type of request. valid values: - SkipOnAllRequestFields: skip all requests; - SkipOnSpecifiedRequestFields: skip specified request fields. valid only when SkipScope is ManagedRules.
* `skip_scope` - (Optional, String) Exception rule execution options, valid values: - WebSecurityModules: designate the security protection module for the exception rule. - ManagedRules: designate the managed rule.
* `web_security_modules_for_exception` - (Optional, Set) Specifies the security protection module for exception rules. valid only when SkipScope is WebSecurityModules. valid values: - websec-mod-managed-rules: managed rule.- websec-mod-rate-limiting: rate limit.- websec-mod-custom-rules: custom rule.- websec-mod-adaptive-control: adaptive frequency control, intelligent client filtering, slow attack protection, traffic theft protection.- websec-mod-bot: bot management.

The `rules` object of `rate_limiting_rules` supports the following:

* `action_duration` - (Optional, String) The duration of an Action is only supported in the following units: - s: seconds, value range 1-120; - m: minutes, value range 1-120; - h: hours, value range 1-48; - d: days, value range 1-30.
* `action` - (Optional, List) Precision rate limiting handling methods. valid values: - Monitor: Monitor; - Deny: block, where DenyActionParameters.Name supports Deny and ReturnCustomPage; - Challenge: Challenge, where ChallengeActionParameters.Name supports JSChallenge and ManagedChallenge; - Redirect: Redirect to URL;.
* `condition` - (Optional, String) The specific content of precise speed limit shall comply with the expression syntax. for detailed specifications, see the product documentation.
* `count_by` - (Optional, Set) Rate threshold request feature match mode. this field is required when Enabled is on.  when there are multiple conditions, composite multiple conditions will perform statistics count. the maximum number of conditions must not exceed 5. valid values: - http.request.ip: client ip; - http.request.xff_header_ip: client ip (priority match xff header); - http.request.uri.path: request access path; - http.request.cookies['session']: Cookie named session, where session can be replaced with your own specified parameter; - http.request.headers['user-agent']: http header named user-agent, where user-agent can be replaced with your own specified parameter; - http.request.ja3: request ja3 fingerprint; - http.request.uri.query['test']: URL query parameter named test, where test can be replaced with your own specified parameter.
* `counting_period` - (Optional, String) Specifies the time window for statistics. valid values: - 1s: 1 second;- 5s: 5 seconds;- 10s: 10 seconds;- 20s: 20 seconds;- 30s: 30 seconds;- 40s: 40 seconds;- 50s: 50 seconds;- 1m: 1 minute;- 2m: 2 minutes;- 5m: 5 minutes;- 10m: 10 minutes;- 1h: 1 hour.
* `enabled` - (Optional, String) Whether the precise rate limiting rule is enabled. valid values: - on: enabled; - off: disabled(No other fields are required when closing).
* `max_request_threshold` - (Optional, Int) Precision rate limiting specifies the cumulative number of interceptions within the time range. value ranges from 1 to 100000.
* `name` - (Optional, String) Specifies the name of the precise rate limit.
* `priority` - (Optional, Int) Precision rate limiting specifies the priority. value range is 0 to 100. default is 0.

The `search_engine_bots` object of `basic_bot_settings` supports the following:

* `base_action` - (Optional, List) Specifies the action for requests from search engine crawlers. valid values for SecurityAction Name: - Deny: block; - Monitor: observe; - Disabled: not enabled, disable specified rule; - Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge; - Allow: pass (to be deprecated).
* `bot_management_action_overrides` - (Optional, List) Specifies the handling method for search engine crawler requests.

The `security_action` object of `action` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `security_policy` object supports the following:

* `bot_management` - (Optional, List) Bot management configuration.
* `custom_rules` - (Optional, List) Custom rules. If the parameter is null or not filled, the configuration last set will be used by default. Note: This field may return null, indicating that no valid value can be obtained.
* `exception_rules` - (Optional, List) Exception rule configuration.
* `http_ddos_protection` - (Optional, List) HTTP DDOS protection configuration.
* `managed_rules` - (Optional, List) Managed. If the parameter is null or not filled, the configuration last set will be used by default. Note: This field may return null, indicating that no valid value can be obtained.
* `rate_limiting_rules` - (Optional, List) Configures the rate limiting rule.

The `session_expired_action` object of `bot_session_validation` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `session_invalid_action` object of `bot_session_validation` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

The `session_rate_control` object of `bot_session_validation` supports the following:

* `enabled` - (Optional, String) Specifies whether session rate and periodic feature verification are enabled. valid values: - on: enable- off: disable.
* `high_rate_session_action` - (Optional, List) Session rate and periodic feature verification high-risk execution actions. SecurityAction Name valid values: - Deny: block, where Stall configuration is supported in DenyActionParameters; - Monitor: observation; - Allow: respond after wait, where MinDelayTime and MaxDelayTime configuration is required in AllowActionParameters.
* `low_rate_session_action` - (Optional, List) Session rate and periodic feature verification low risk execution action. SecurityAction Name parameter supports: - Deny: block, where DenyActionParameters supports Stall configuration;- Monitor: observe;- Allow: respond after wait, where AllowActionParameters requires MinDelayTime and MaxDelayTime configuration.
* `mid_rate_session_action` - (Optional, List) Session rate and periodic feature verification medium-risk execution action. SecurityAction Name parameter supports: - Deny: block, where DenyActionParameters supports Stall configuration;- Monitor: observe;- Allow: respond after wait, where AllowActionParameters requires MinDelayTime and MaxDelayTime configuration.

The `slow_attack_defense` object of `http_ddos_protection` supports the following:

* `enabled` - (Required, String) Whether slow attack protection is enabled. valid values: - on: enabled; - off: disabled.
* `action` - (Optional, List) Slow attack protection handling method. required when Enabled is on. valid values for SecurityAction Name: - Monitor: observation; - Deny: block;.
* `minimal_request_body_transfer_rate` - (Optional, List) The specific configuration of the minimum body transfer rate threshold is required when Enabled is on.
* `request_body_transfer_timeout` - (Optional, List) Specifies the specific configuration of body transfer timeout duration. required when Enabled is on.

The `source_idc` object of `basic_bot_settings` supports the following:

* `base_action` - (Optional, List) Handling method for requests from the specified IDC. valid values for SecurityAction Name: - Deny: block; - Monitor: observe; - Disabled: not enabled, disable specified rule; - Challenge: Challenge, where ChallengeOption in ChallengeActionParameters supports JSChallenge and ManagedChallenge; - Allow: pass (to be deprecated).
* `bot_management_action_overrides` - (Optional, List) Specifies the handling method for the specified id request.

The `verified_bot_requests_action` object of `bot_ratings` supports the following:

* `name` - (Required, String) Specifies the specific actions for safe execution. valid values:. - Deny. specifies to block requests from accessing site resources. - Monitor: observation, only record logs. - Redirect: Redirect to URL. - Disabled: specifies that the rule is not enabled. - Allow: specifies whether to allow access with delayed processing of requests. - Challenge: specifies the challenge content to respond to. - Trans: pass and allow requests to directly access site resources. - BlockIP: to be deprecated. ip block. - ReturnCustomPage: to be deprecated. use specified page for interception. - JSChallenge: to be deprecated, JavaScript challenge;. - ManagedChallenge: to be deprecated. managed challenge.
* `allow_action_parameters` - (Optional, List) Additional parameters when Name is Allow.
* `block_ip_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is BlockIP.
* `challenge_action_parameters` - (Optional, List) Additional parameter when Name is Challenge.
* `deny_action_parameters` - (Optional, List) Additional parameters when Name is Deny.
* `redirect_action_parameters` - (Optional, List) Additional parameter when Name is Redirect.
* `return_custom_page_action_parameters` - (Optional, List) To be deprecated, additional parameter when Name is ReturnCustomPage.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo web security template can be imported using the id, e.g.

```
terraform import tencentcloud_teo_web_security_template.example zone-37u62pwxfo8s#temp-05dtxkyw
```

