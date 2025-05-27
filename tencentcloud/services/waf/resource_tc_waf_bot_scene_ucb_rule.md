Provides a resource to create a WAF bot scene ucb rule

Example Usage

The rules are permanently effective

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

The rules take effect on a scheduled basis

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

    on_off    = "on"
    rule_type = 0
    prior     = 100
    label     = "正常流量"
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

The rules take effect on a weekly basis

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

Import

WAF bot scene ucb rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_bot_scene_ucb_rule.example examle.com#3000000791#3000003489
```
