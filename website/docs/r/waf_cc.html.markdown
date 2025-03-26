---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_cc"
sidebar_current: "docs-tencentcloud-resource-waf_cc"
description: |-
  Provides a resource to create a waf cc
---

# tencentcloud_waf_cc

Provides a resource to create a waf cc

## Example Usage

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain      = "www.demo.com"
  name        = "terraform"
  status      = 1
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
  options_arr = jsonencode(
    [
      {
        "key" : "Method",
        "args" : ["=R0VU"],
        "match" : "0",
        "encodeflag" : true
      }
    ]
  )
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, String) Rule Action, 20 log, 21 captcha, 22 deny, 23 accurate deny.
* `advance` - (Required, String) Advanced mode (whether to use session detection). 0(disabled) 1(enabled).
* `domain` - (Required, String) Domain.
* `edition` - (Required, String) WAF edition. clb-waf means clb-waf, sparta-waf means saas-waf.
* `interval` - (Required, String) CC detection cycle.
* `limit` - (Required, String) CC detection threshold.
* `match_func` - (Required, Int) Match method, 0(equal), 1(prefix), 2(contains), 3(not equal), 6(suffix), 7(not contains).
* `name` - (Required, String) Rule Name.
* `priority` - (Required, Int) Rule Priority.
* `status` - (Required, Int) Rule Status, 0 rule close, 1 rule open.
* `url` - (Required, String) Detection URL.
* `valid_time` - (Required, Int) Action ValidTime, minute unit. Min: 60, Max: 604800.
* `event_id` - (Optional, String) Event ID.
* `options_arr` - (Optional, String) JSON serialized string of CC matching conditions, example:[{\"key\":\"Method\",\"args\":[\"=R0VU\"],\"match\":\"0\",\"encodeflag\":true}]
        Key optional values are Method, Post, Referer, Cookie, User-Agent, CustomHeader
        Match optional values are, when Key is Method, optional values are 0 (equal), 3 (not equal).
        When the key is Post, the optional values are 0 (equal to), 3 (not equal to), when the key is Cookie, the optional values are 0 (equal to), 2 (included), 3 (not equal to), 7 (not included),
        When the key is Referer, the optional values are 0 (equal to), 3 (not equal to), 1 (prefix match), 6 (suffix match), 2 (included), 7 (not included), 12 (exists), 5 (not included), 4 (content is empty),
        When the key is Cookie, the optional values are 0 (equal to), 3 (not equal to), 2 (included), 7 (not included),
        When the key is User-Agent, the optional values are 0 (equal to), 3 (not equal to), 1 (prefix match), 6 (suffix match), 2 (included), 7 (not included), 12 (exists), 5 (not included), 4 (content is empty),
        When the key is CustomHeader, the optional values are 0 (equal to), 3 (not equal to), 2 (included), 7 (not included), 12 (exists), 5 (not included), 4 (content is empty).
        args is used to indicate the matching content. You need to set encodeflag to true. When the Key is Post, Cookie, or CustomHeader, use the equal sign = to concatenate the Key and Value, and encode them using Base64, similar to YWJj=YWJj. When the Key is Referer or User-Agent, use the equal sign = to concatenate the Value, similar to =YWJj.
* `session_applied` - (Optional, Set: [`Int`]) Session ID that needs to be enabled for the rule.
* `type` - (Optional, Int) Operate Type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


