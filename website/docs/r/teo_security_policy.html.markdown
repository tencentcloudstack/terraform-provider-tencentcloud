---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_policy"
sidebar_current: "docs-tencentcloud-resource-teo_security_policy"
description: |-
  Provides a resource to create a teo security_policy
---

# tencentcloud_teo_security_policy

Provides a resource to create a teo security_policy

## Example Usage

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
        block_rule_ids = [
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

## Argument Reference

The following arguments are supported:

* `entity` - (Required, String) Subdomain.
* `zone_id` - (Required, String) Site ID.
* `config` - (Optional, List) Security policy configuration.

The `acl_config` object supports the following:

* `switch` - (Required, String) - `on`: Enable.- `off`: Disable.
* `user_rules` - (Optional, List) Custom configuration.

The `acl_drop_page_detail` object supports the following:

* `name` - (Optional, String) File name or URL.
* `page_id` - (Optional, Int) ID of the custom error page. when set to 0, use system default error page.
* `status_code` - (Optional, Int) HTTP status code to use. Valid range: 100-600.
* `type` - (Optional, String) Type of the custom error page. Valid values: `file`, `url`.

The `ai_rule` object supports the following:

* `mode` - (Optional, String) Valid values:- `smart_status_close`: disabled.- `smart_status_open`: blocked.- `smart_status_observe`: observed.

The `bot_config` object supports the following:

* `intelligence_rule` - (Optional, List) Bot intelligent rule configuration.
* `managed_rule` - (Optional, List) Preset rules.
* `portrait_rule` - (Optional, List) Portrait rule.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `conditions` object supports the following:

* `match_content` - (Required, String) Content to match.
* `match_from` - (Required, String) Items to match. Valid values:- `host`: Host of the request.- `sip`: Client IP.- `ua`: User-Agent.- `cookie`: Session cookie.- `cgi`: CGI script.- `xff`: XFF extension header.- `url`: URL of the request.- `accept`: Accept encoding of the request.- `method`: HTTP method of the request.- `header`: HTTP header of the request.- `sip_proto`: Network protocol of the request.
* `match_param` - (Required, String) Parameter for match item. For example, when match from header, match parameter can be set to a header key.
* `operator` - (Required, String) Valid values:- `equal`: string equal.- `not_equal`: string not equal.- `include`: string include.- `not_include`: string not include.- `match`: ip match.- `not_match`: ip not match.- `include_area`: area include.- `is_empty`: field existed but empty.- `not_exists`: field is not existed.- `regexp`: regex match.- `len_gt`: value greater than.- `len_lt`: value less than.- `len_eq`: value equal.- `match_prefix`: string prefix match.- `match_suffix`: string suffix match.- `wildcard`: wildcard match.

The `config` object supports the following:

* `acl_config` - (Optional, List) ACL configuration.
* `bot_config` - (Optional, List) Bot Configuration.
* `drop_page_config` - (Optional, List) Custom drop page configuration.
* `except_config` - (Optional, List) Exception rule configuration.
* `ip_table_config` - (Optional, List) Basic access control.
* `rate_limit_config` - (Optional, List) RateLimit Configuration.
* `switch_config` - (Optional, List) Main switch of 7-layer security.
* `waf_config` - (Optional, List) WAF (Web Application Firewall) Configuration.

The `detail` object supports the following:

* `action` - (Optional, String) Action to take.
* `id` - (Optional, Int) Template ID. Note: This field may return null, indicating that no valid value can be obtained.
* `mode` - (Optional, String) Template Name. Note: This field may return null, indicating that no valid value can be obtained.
* `period` - (Optional, Int) Period.
* `punish_time` - (Optional, Int) Punish time.
* `threshold` - (Optional, Int) Threshold.

The `drop_page_config` object supports the following:

* `acl_drop_page_detail` - (Optional, List) Custom error page of ACL rules.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.
* `waf_drop_page_detail` - (Optional, List) Custom error page of WAF rules.

The `except_config` object supports the following:

* `except_user_rules` - (Optional, List) Exception rules.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `except_user_rule_conditions` object supports the following:

* `match_content` - (Optional, String) Content to match.
* `match_from` - (Optional, String) Items to match. Valid values:- `host`: Host of the request.- `sip`: Client IP.- `ua`: User-Agent.- `cookie`: Session cookie.- `cgi`: CGI script.- `xff`: XFF extension header.- `url`: URL of the request.- `accept`: Accept encoding of the request.- `method`: HTTP method of the request.- `header`: HTTP header of the request.- `sip_proto`: Network protocol of the request.
* `match_param` - (Optional, String) Parameter for match item. For example, when match from header, match parameter can be set to a header key.
* `operator` - (Optional, String) Valid values:- `equal`: string equal.- `not_equal`: string not equal.- `include`: string include.- `not_include`: string not include.- `match`: ip match.- `not_match`: ip not match.- `include_area`: area include.- `is_empty`: field existed but empty.- `not_exists`: field is not existed.- `regexp`: regex match.- `len_gt`: value greater than.- `len_lt`: value less than.- `len_eq`: value equal.- `match_prefix`: string prefix match.- `match_suffix`: string suffix match.- `wildcard`: wildcard match.

The `except_user_rule_scope` object supports the following:

* `modules` - (Optional, Set) Modules in which the rule take effect. Valid values: `waf`.

The `except_user_rules` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `skip`.
* `except_user_rule_conditions` - (Optional, List) Conditions of the rule.
* `except_user_rule_scope` - (Optional, List) Scope of the rule in effect.
* `rule_priority` - (Optional, Int) Priority of the rule. Valid value range: 0-100.
* `rule_status` - (Optional, String) Status of the rule. Valid values:- `on`: Enabled.- `off`: Disabled.

The `intelligence_rule` object supports the following:

* `items` - (Optional, List) Configuration detail.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `intelligence` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `monitor`, `alg`.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `ip_table_config` object supports the following:

* `rules` - (Optional, List) Rules list.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `items` object supports the following:

* `action` - (Optional, String) Action to take. Valid values: `trans`, `monitor`, `alg`, `captcha`, `drop`.
* `label` - (Optional, String) Bot label, valid values: `evil_bot`, `suspect_bot`, `good_bot`, `normal`.

The `managed_rule` object supports the following:

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
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `rate_limit_config` object supports the following:

* `intelligence` - (Optional, List) Intelligent client filter.
* `switch` - (Optional, String) - `on`: Enable.- `off`: Disable.
* `template` - (Optional, List) Default Template. Note: This field may return null, indicating that no valid value can be obtained.
* `user_rules` - (Optional, List) Custom configuration.

The `rules` object supports the following:

* `action` - (Optional, String) Actions to take. Valid values: `drop`, `trans`, `monitor`.
* `match_content` - (Optional, String) Matching content.
* `match_from` - (Optional, String) Matching type. Valid values: `ip`, `area`.

The `switch_config` object supports the following:

* `web_switch` - (Optional, String) - `on`: Enable.- `off`: Disable.

The `template` object supports the following:

* `detail` - (Optional, List) Detail of the template.
* `mode` - (Optional, String) Template Name. Note: This field may return null, indicating that no valid value can be obtained.

The `user_rules` object supports the following:

* `action` - (Required, String) Action to take. Valid values: `trans`, `drop`, `monitor`, `ban`, `redirect`, `page`, `alg`.
* `conditions` - (Required, List) Conditions of the rule.
* `rule_name` - (Required, String) Rule name.
* `rule_priority` - (Required, Int) Priority of the rule. Valid value range: 0-100.
* `rule_status` - (Required, String) Status of the rule. Valid values: `on`, `off`.
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
* `punish_time_unit` - (Required, String) Time unit of the punish time. Valid values: `second`, `minutes`, `hour`.
* `punish_time` - (Required, Int) Punish time, Valid value range: 0-2 days.
* `rule_name` - (Required, String) Rule Name.
* `rule_priority` - (Required, Int) Priority of the rule. Valid value range: 1-100.
* `threshold` - (Required, Int) Threshold of the rate limit. Valid value range: 0-4294967294.
* `freq_fields` - (Optional, Set) Filter words.
* `rule_status` - (Optional, String) Status of the rule. Valid values: `on`, `off`, `hour`.

The `waf_config` object supports the following:

* `level` - (Required, String) Protection level. Valid values: `loose`, `normal`, `strict`, `stricter`, `custom`.
* `mode` - (Required, String) Protection mode. Valid values:- `block`: use block mode globally, you still can set a group of rules to use observe mode.- `observe`: use observe mode globally.
* `switch` - (Required, String) Whether to enable WAF rules. Valid values:- `on`: Enable.- `off`: Disable.
* `waf_rules` - (Required, List) WAF Rules Configuration.
* `ai_rule` - (Optional, List) AI based rules configuration.

The `waf_drop_page_detail` object supports the following:

* `name` - (Optional, String) File name or URL.
* `page_id` - (Optional, Int) ID of the custom error page. when set to 0, use system default error page.
* `status_code` - (Optional, Int) HTTP status code to use. Valid range: 100-600.
* `type` - (Optional, String) Type of the custom error page. Valid values: `file`, `url`.

The `waf_rules` object supports the following:

* `block_rule_ids` - (Required, Set) Block mode rules list. See details in data source `waf_managed_rules`.
* `switch` - (Required, String) Whether to host the rules&#39; configuration.- `on`: Enable.- `off`: Disable.
* `observe_rule_ids` - (Optional, Set) Observe rules list. See details in data source `waf_managed_rules`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo security_policy can be imported using the zoneId#entity, e.g.
```
$ terraform import tencentcloud_teo_security_policy.security_policy zone-2983wizgxqvm#aaa.sfurnace.work
```

