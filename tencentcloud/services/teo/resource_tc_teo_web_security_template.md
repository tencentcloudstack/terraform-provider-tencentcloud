Provides a resource to create a teo web security template

~> **NOTE:** The current resources do not support managed_rule_groups.

Example Usage

Basic usage

```hcl

resource "tencentcloud_teo_web_security_template" "web_security_template" {
  template_name = "tf测试"
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
          name      = "拦截非浏览器爬虫访问"
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
          name      = "登录接口请求突增防护"
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
        name      = "盗刷 User-Agent 黑名单"
        priority  = 50
        rule_type = "PreciseMatchRule"
        action {
          name = "JSChallenge"
        }
      }
      rules {
        condition = "$${http.request.ip} in ['36']"
        enabled   = "on"
        name      = "自定义规则"
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
        name                               = "高频 API 跳过速率限制1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control"]
      }
      rules {
        condition                          = "$${http.request.ip} in ['123.123.123.0/24']"
        enabled                            = "on"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "IP 白名单1"
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
        name                  = "单 IP 请求速率限制1"
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

Import

teo web security template can be imported using the id, e.g.

```
terraform import tencentcloud_teo_web_security_template.example zone-37u62pwxfo8s#temp-05dtxkyw
```
