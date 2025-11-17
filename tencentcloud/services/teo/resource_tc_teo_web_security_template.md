Provides a resource to create a TEO web security template

Example Usage

```hcl
resource "tencentcloud_teo_web_security_template" "example" {
  zone_id       = "zone-3fkff38fyw8s"
  template_name = "example"
  security_policy {
    exception_rules {
      rules {
        name                               = "test"
        condition                          = "$${http.request.host} in ['1.1.1.1']"
        skip_scope                         = "WebSecurityModules"
        skip_option                        = "SkipOnAllRequestFields"
        web_security_modules_for_exception = ["websec-mod-managed-rules"]
        enabled                            = "on"
      }
    }

    custom_rules {
      rules {
        name      = "test"
        condition = "$${http.request.ip} in ['1.1.1.1']"
        enabled   = "on"
        rule_type = "BasicAccessRule"
        action {
          name = "Deny"
        }
      }
    }

    rate_limiting_rules {
      rules {
        name                  = "单 IP 请求速率限制"
        condition             = "$${http.request.uri.path} contain ['/checkout/submit']"
        count_by              = ["http.request.ip"]
        max_request_threshold = 300
        counting_period       = "60s"
        action_duration       = "30m"
        priority              = 50
        enabled               = "on"
        action {
          name = "Challenge"
          challenge_action_parameters {
            challenge_option = "JSChallenge"
          }
        }
      }
    }

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

Import

TEO web security template can be imported using the zoneId#templateId, e.g.

```
terraform import tencentcloud_teo_web_security_template.example zone-3fkff38fyw8s#temp-p3p973nu
```
