---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_policy"
sidebar_current: "docs-tencentcloud-resource-teo_security_policy"
description: |-
  Provides a resource to create a teo securityPolicy
---

# tencentcloud_teo_security_policy

Provides a resource to create a teo securityPolicy

## Example Usage

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
        switch           = ""
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
        rule_priority = ""
      }
      template {
        mode = ""
        detail {
          mode        = ""
          id          = ""
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
        rule_name   = ""
        action      = ""
        rule_status = ""
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
        rule_id           = ""
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
        rule_id          = ""
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
        rule_id       = ""
      }
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required, List) Security policy configuration.
* `entity` - (Required, String) Subdomain.
* `zone_id` - (Required, String) Site ID.

The `acl_config` object supports the following:

* `switch` - (Required, String) - on: Enable.- off: Disable.
* `user_rules` - (Required, List) Custom configuration.

The `ai_rule` object supports the following:

* `mode` - (Optional, String) Valid values:- smart_status_close: disabled.- smart_status_open: blocked.- smart_status_observe: observed.

The `bot_config` object supports the following:

* `intelligence_rule` - (Optional, List) Bot intelligent rule configuration.
* `managed_rule` - (Optional, List) Preset rules.
* `portrait_rule` - (Optional, List) Portrait rule.
* `switch` - (Optional, String) - on: Enable.- off: Disable.

The `conditions` object supports the following:

* `match_content` - (Required, String) Matching content.
* `match_from` - (Required, String) Matching field.
* `match_param` - (Required, String) Matching string.
* `operator` - (Required, String) Matching operator.

The `config` object supports the following:

* `acl_config` - (Optional, List) ACL configuration.
* `bot_config` - (Optional, List) Bot Configuration.
* `ip_table_config` - (Optional, List) Basic access control.
* `rate_limit_config` - (Optional, List) RateLimit Configuration.
* `switch_config` - (Optional, List) Main switch of 7-layer security.
* `waf_config` - (Optional, List) WAF (Web Application Firewall) Configuration.

The `detail` object supports the following:

* `action` - (Optional, String) Action to take.
* `id` - (Optional, Int) Template ID.Note: This field may return null, indicating that no valid value can be obtained.
* `mode` - (Optional, String) Template Name.Note: This field may return null, indicating that no valid value can be obtained.
* `period` - (Optional, Int) Period.
* `punish_time` - (Optional, Int) Punish time.
* `threshold` - (Optional, Int) Threshold.

The `intelligence_rule` object supports the following:

* `items` - (Optional, List) Configuration detail.
* `switch` - (Optional, String) - on: Enable.- off: Disable.

The `intelligence` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `monitor`, `alg`.
* `switch` - (Optional, String) - on: Enable.- off: Disable.

The `ip_table_config` object supports the following:

* `rules` - (Optional, List) Rules list.
* `switch` - (Optional, String) - on: Enable.- off: Disable.

The `items` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `trans`, `monitor`, `alg`, `captcha`, `drop`.
* `label` - (Optional, String) Bot label, valid values: `evil_bot`, `suspect_bot`, `good_bot`, `normal`.

The `managed_rule` object supports the following:

* `rule_id` - (Required, Int) Rule ID.
* `action` - (Optional, String) Action to take. Valid values: `drop`, `trans`, `monitor`, `alg`.
* `alg_managed_ids` - (Optional, Set) Rules to enable when action is `alg`. See details in data source `bot_managed_rules`.
* `cap_managed_ids` - (Optional, Set) Rules to enable when action is `captcha`. See details in data source `bot_managed_rules`.
* `drop_managed_ids` - (Optional, Set) Rules to enable when action is `drop`. See details in data source `bot_managed_rules`.
* `mon_managed_ids` - (Optional, Set) Rules to enable when action is `monitor`. See details in data source `bot_managed_rules`.
* `name` - (Optional, String) Name of the custom response page.
* `page_id` - (Optional, Int) ID of the custom response page.
* `punish_time_unit` - (Optional, String) Time unit of the punish time.
* `punish_time` - (Optional, Int) Punish time.
* `redirect_url` - (Optional, String) Redirect target URL, must be an sub-domain from one of the account&#39;s site.
* `response_code` - (Optional, Int) Response code to use when redirecting.
* `trans_managed_ids` - (Optional, Set) Rules to enable when action is `trans`. See details in data source `bot_managed_rules`.

The `portrait_rule` object supports the following:

* `alg_managed_ids` - (Optional, Set) Rules to enable when action is `alg`. See details in data source `bot_portrait_rules`.
* `cap_managed_ids` - (Optional, Set) Rules to enable when action is `captcha`. See details in data source `bot_portrait_rules`.
* `drop_managed_ids` - (Optional, Set) Rules to enable when action is `drop`. See details in data source `bot_portrait_rules`.
* `mon_managed_ids` - (Optional, Set) Rules to enable when action is `monitor`. See details in data source `bot_portrait_rules`.
* `rule_id` - (Optional, Int) Rule ID.
* `switch` - (Optional, String) - on: Enable.- off: Disable.

The `rate_limit_config` object supports the following:

* `switch` - (Required, String) - on: Enable.- off: Disable.
* `user_rules` - (Required, List) Custom configuration.
* `intelligence` - (Optional, List) Intelligent client filter.
* `template` - (Optional, List) Default Template.Note: This field may return null, indicating that no valid value can be obtained.

The `rules` object supports the following:

* `action` - (Optional, String) Actions to take. Valid values: `drop`, `trans`, `monitor`.
* `match_content` - (Optional, String) Matching content.
* `match_from` - (Optional, String) Matching type. Valid values: `ip`, `area`.
* `rule_id` - (Optional, Int) Rule ID.

The `switch_config` object supports the following:

* `web_switch` - (Optional, String) - on: Enable.- off: Disable.

The `template` object supports the following:

* `detail` - (Optional, List) Detail of the template.
* `mode` - (Optional, String) Template Name.Note: This field may return null, indicating that no valid value can be obtained.

The `user_rules` object supports the following:

* `action` - (Required, String) Action to take. Valid values: `trans`, `drop`, `monitor`, `ban`, `redirect`, `page`, `alg`.
* `conditions` - (Required, List) Conditions of the rule.
* `rule_name` - (Required, String) Rule name.
* `rule_priority` - (Required, Int) Priority of the rule. Valid value range: 0-100.
* `rule_status` - (Required, String) Status of the rule.
* `name` - (Optional, String) Name of the custom response page.
* `page_id` - (Optional, Int) ID of the custom response page.
* `punish_time_unit` - (Optional, String) Time unit of the punish time. Valid values: `second`, `minutes`, `hour`.
* `punish_time` - (Optional, Int) Punish time, Valid value range: 0-2 days.
* `redirect_url` - (Optional, String) Redirect target URL, must be an sub-domain from one of the account&#39;s site.
* `response_code` - (Optional, Int) Response code to use when redirecting.

The `user_rules` object supports the following:

* `action` - (Required, String) Valid values: `monitor`, `drop`.
* `conditions` - (Required, List) Conditions of the rule.
* `period` - (Required, Int) Period of the rate limit. Valid values: 10, 20, 30, 40, 50, 60 (in seconds).
* `punish_time_unit` - (Required, String) Time unit of the punish time. Valid values: `second`,`minutes`,`hour`.
* `punish_time` - (Required, Int) Punish time, Valid value range: 0-2 days.
* `rule_name` - (Required, String) Rule Name.
* `rule_priority` - (Required, Int) Priority of the rule.
* `rule_status` - (Required, String) Status of the rule.
* `threshold` - (Required, Int) Threshold of the rate limit. Valid value range: 0-4294967294.
* `freq_fields` - (Optional, Set) Filter words.

The `waf_config` object supports the following:

* `level` - (Required, String) Protection level. Valid values: `loose`, `normal`, `strict`, `stricter`, `custom`.
* `mode` - (Required, String) Protection mode. Valid values:- block: use block mode globally, you still can set a group of rules to use observe mode.- observe: use observe mode globally.
* `switch` - (Required, String) Whether to enable WAF rules. Valid values:- on: Enable.- off: Disable.
* `waf_rules` - (Required, List) WAF Rules Configuration.
* `ai_rule` - (Optional, List) AI based rules configuration.

The `waf_rules` object supports the following:

* `block_rule_ids` - (Required, Set) Block mode rules list. See details in data source `waf_managed_rules`.
* `switch` - (Required, String) Whether to host the rules&#39; configuration.- on: Enable.- off: Disable.
* `observe_rule_ids` - (Optional, Set) Observe rules list. See details in data source `waf_managed_rules`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo securityPolicy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_security_policy.securityPolicy securityPolicy_id#entity
```

