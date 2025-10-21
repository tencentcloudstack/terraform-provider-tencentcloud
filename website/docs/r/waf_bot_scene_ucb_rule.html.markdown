---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_bot_scene_ucb_rule"
sidebar_current: "docs-tencentcloud-resource-waf_bot_scene_ucb_rule"
description: |-
  Provides a resource to create a WAF bot scene ucb rule
---

# tencentcloud_waf_bot_scene_ucb_rule

Provides a resource to create a WAF bot scene ucb rule

## Example Usage

### The rules are permanently effective

```hcl
resource "tencentcloud_waf_bot_scene_ucb_rule" "example" {
  domain   = "examle.com"
  scene_id = "3000000791"
  rule {
    domain = "examle.com"
    name   = "tf-example"
    rule {
      key  = "ip_scope"
      op   = "belong"
      lang = "cn"
      value {
        belong_value = ["1.1.1.1"]
      }
    }

    rule {
      key  = "url"
      op   = "rematch"
      lang = "cn"
      value {
        multi_value = [
          "/prefix",
          "/startwith"
        ]
      }
    }

    action       = "monitor"
    on_off       = "on"
    rule_type    = 0
    prior        = 100
    label        = "疑似BOT"
    appid        = 1276513791
    addition_arg = "none"
    desc         = "rule desc."
    pre_define   = true
    job_type     = "forever"
    job_date_time {
      timed {
        start_date_time = 0
        end_date_time   = 0
      }

      time_t_zone = "UTC+8"
    }
  }
}
```

### The rules take effect on a scheduled basis

```hcl
resource "tencentcloud_waf_bot_scene_ucb_rule" "example" {
  domain   = "examle.com"
  scene_id = "3000000791"
  rule {
    domain = "examle.com"
    name   = "tf-example"
    rule {
      key  = "header_value"
      op   = "logic"
      name = "token"
      lang = "cn"
      value {
        logic_value = true
      }
    }

    action = "multi_action"
    action_list {
      action     = "monitor"
      proportion = 0.3
    }

    action_list {
      action     = "intercept"
      proportion = 0.3
    }

    action_list {
      action     = "captcha"
      proportion = 0.4
    }

    on_off       = "on"
    rule_type    = 0
    prior        = 100
    label        = "正常流量"
    appid        = 1256704386
    addition_arg = "none"
    desc         = "rule desc."
    pre_define   = true
    job_type     = "timed_job"
    job_date_time {
      timed {
        start_date_time = 1747324800
        end_date_time   = 1747152000
      }

      time_t_zone = "UTC+8"
    }
  }
}
```

### The rules take effect on a weekly basis

```hcl
resource "tencentcloud_waf_bot_scene_ucb_rule" "example" {
  domain   = "examle.com"
  scene_id = "3000000791"
  rule {
    domain = "examle.com"
    name   = "tf-example"
    rule {
      key  = "post_value"
      op   = "prefix"
      lang = "cn"
      value {
        multi_value = [
          "terraform",
          "provider"
        ]
      }
    }

    action        = "intercept"
    on_off        = "on"
    rule_type     = 0
    prior         = 100
    label         = "恶意BOT"
    appid         = 1256704386
    addition_arg  = "none"
    desc          = "rule desc."
    pre_define    = true
    block_page_id = 71
    job_type      = "cron_week"
    job_date_time {
      cron {
        w_days     = [1, 2, 3, 4, 5]
        start_time = "00:00:00"
        end_time   = "23:59:59"
      }

      time_t_zone = "UTC+8"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `scene_id` - (Required, String, ForceNew) When calling at the BOT global whitelist, pass `global`; When configuring BOT scenarios, transmit the specific scenario ID.
* `rule` - (Optional, List) Rule content, add encoding SceneId information. When calling at the BOT global whitelist, SceneId is set to `global` and RuleType is passed as 10, Action is `permit`; When configuring BOT scenarios, SceneId is the scenario ID.

The `action_list` object of `rule` supports the following:

* `action` - (Optional, String) Action.
* `proportion` - (Optional, Float64) Proportion.

The `areas` object of `rule` supports the following:

* `country` - (Required, String) In addition to standard countries, the country also supports two special identifiers: domestic and foreign.
* `city` - (Optional, String) City.
* `region` - (Optional, String) Province.

The `cron` object of `job_date_time` supports the following:

* `days` - (Optional, Set) On what day of each month is it executed.
* `end_time` - (Optional, String) End time.
* `start_time` - (Optional, String) Start time.
* `w_days` - (Optional, Set) What day of the week is executed each week.

The `job_date_time` object of `rule` supports the following:

* `cron` - (Optional, List) Time parameter for cycle execution.
* `time_t_zone` - (Optional, String) Time zone.
* `timed` - (Optional, List) Time parameter for timed execution.

The `rule` object of `rule` supports the following:

* `areas` - (Optional, List) Regional selection.
* `key` - (Optional, String) Key.
* `lang` - (Optional, String) Language environment.
* `name` - (Optional, String) When using header parameter values.
* `op_arg` - (Optional, Set) Optional supplementary parameters.
* `op_op` - (Optional, String) Optional Supplementary Operators.
* `op_value` - (Optional, Float64) Optional supplementary values.
* `op` - (Optional, String) Operator.
* `value` - (Optional, List) Value.

The `rule` object supports the following:

* `action` - (Required, String) Disposal action.
* `domain` - (Required, String) Domain.
* `label` - (Required, String) Label.
* `name` - (Required, String) Rule name.
* `on_off` - (Required, String) Rule switch.
* `prior` - (Required, Int) Rule priority.
* `rule_type` - (Required, Int) Rule type.
* `rule` - (Required, List) Specific rule items of UCB.
* `action_list` - (Optional, List) When Action=intercept, this field is mandatory.
* `addition_arg` - (Optional, String) Additional parameters.
* `appid` - (Optional, Int) Appid.
* `block_page_id` - (Optional, Int) Customize interception page ID.
* `desc` - (Optional, String) Rule description.
* `expire_time` - (Optional, Int) Effective deadline.
* `id` - (Optional, String) Entry ID.
* `job_date_time` - (Optional, List) Scheduled task configuration.
* `job_type` - (Optional, String) Scheduled task type.
* `pre_define` - (Optional, Bool) True - System preset rules False - Custom rules.
* `scene_id` - (Optional, String) Scene ID.
* `valid_status` - (Optional, Int) Effective -1, Invalid -0.
* `valid_time` - (Optional, Int) Valid time.

The `timed` object of `job_date_time` supports the following:

* `end_date_time` - (Optional, Int) End timestamp, in seconds.
* `start_date_time` - (Optional, Int) Start timestamp, in seconds.

The `value` object of `rule` supports the following:

* `basic_value` - (Optional, String) String type value.
* `belong_value` - (Optional, Set) String array type value.
* `logic_value` - (Optional, Bool) Bool type value.
* `multi_value` - (Optional, Set) String array type value.
* `valid_key` - (Optional, String) Indicate valid fields.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

WAF bot scene ucb rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_bot_scene_ucb_rule.example examle.com#3000000791#3000003489
```

