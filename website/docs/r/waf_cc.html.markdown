---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_cc"
sidebar_current: "docs-tencentcloud-resource-waf_cc"
description: |-
  Provides a resource to create a WAF cc
---

# tencentcloud_waf_cc

Provides a resource to create a WAF cc

## Example Usage

### If advance is 0(IP model)

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain      = "www.demo.com"
  name        = "tf-example"
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
  logical_op  = "and"
  options_arr = jsonencode(
    [
      {
        "key" : "URL",
        "args" : [
          "=cHJlZml4"
        ],
        "match" : "2",
        "encodeflag" : true
      },
      {
        "key" : "Method",
        "args" : [
          "=POST" # if encodeflag is false, parameter value needs to be prefixed with an = sign.
        ],
        "match" : "0",
        "encodeflag" : false
      },
      {
        "key" : "Post",
        "args" : [
          "S2V5=VmFsdWU"
        ],
        "match" : "0",
        "encodeflag" : true
      },
      {
        "key" : "Referer",
        "args" : [
          "="
        ],
        "match" : "12",
        "encodeflag" : true
      },
      {
        "key" : "Cookie",
        "args" : [
          "S2V5=VmFsdWU"
        ],
        "match" : "3",
        "encodeflag" : true
      },
      {
        "key" : "IPLocation",
        "args" : [
          "=eyJMYW5nIjoiY24iLCJBcmVhcyI6W3siQ291bnRyeSI6IuWbveWkliJ9XX0"
        ],
        "match" : "13",
        "encodeflag" : true
      }
    ]
  )
}
```

### If advance is 1(SESSION model)

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain          = "news.bots.icu"
  name            = "tf-example"
  status          = 1
  advance         = "1"
  limit           = "60"
  interval        = "60"
  url             = "/cc_demo"
  match_func      = 0
  action_type     = "22"
  priority        = 50
  valid_time      = 600
  edition         = "sparta-waf"
  type            = 1
  session_applied = [0]
  limit_method    = "only_limit"
  logical_op      = "or"
  cel_rule        = "(has(request.url) && request.url.startsWith('/prefix')) && (has(request.method) && request.method == 'POST')"
}
```

## Argument Reference

The following arguments are supported:

* `action_type` - (Required, String) Rule Action, 20 means observation, 21 means human-machine identification, 22 means interception, 23 means precise interception, 26 means precise human-machine identification, and 27 means JS verification.
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
* `cel_rule` - (Optional, String) Cel expression.
* `event_id` - (Optional, String) Event ID.
* `limit_method` - (Optional, String) Frequency limiting method.
* `logical_op` - (Optional, String) Logical operator of configuration mode, and/or.
* `options_arr` - (Optional, String) CC matching conditions JSON serialized string. Example: [{"key":"Method","args":["=R0VU"],"match":"0","encodeflag":true}]. 

Supported key types: URL, Method, Post, Referer, Cookie, User-Agent, CustomHeader, IPLocation, CaptchaRisk, CaptchaDeviceRisk, CaptchaScore. 

Match operators by key type:
- When Key is URL: 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains)
- When Key is Method: 0 (equal), 3 (not equal)
- When Key is Post: 0 (equal), 3 (not equal), 2 (contains), 7 (not contains)
- When Key is Cookie: 0 (equal), 3 (not equal), 2 (contains), 7 (not contains)
- When Key is Referer: 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains), 12 (exists), 5 (not exists), 4 (empty)
- When Key is User-Agent: 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains), 12 (exists), 5 (not exists), 4 (empty)
- When Key is CustomHeader: 0 (equal), 3 (not equal), 2 (contains), 7 (not contains), 4 (empty), 5 (not exists)
- When Key is IPLocation: 13 (belongs to), 14 (not belongs to)
- When Key is CaptchaRisk: 15 (numerically equal), 16 (numerically not equal), 13 (belongs to), 14 (not belongs to), 12 (exists), 5 (not exists)
- When Key is CaptchaDeviceRisk: 13 (belongs to), 14 (not belongs to), 12 (exists), 5 (not exists)
- When Key is CaptchaScore: 15 (numerically equal), 17 (numerically greater than), 18 (numerically less than), 19 (numerically greater than or equal), 20 (numerically less than or equal), 12 (exists), 5 (not exists)

Encoding rules: The args parameter requires encodeflag to be set to true. For Post, Cookie, or CustomHeader keys, Base64 encode both parameter name and value (remove trailing =), then concatenate with = sign (e.g., Base64(name)=Base64(value)). For Referer or User-Agent keys, Base64 encode the value (remove trailing =) and prefix with = sign (e.g., =Base64(value)).
* `session_applied` - (Optional, Set: [`Int`]) Session ID that needs to be enabled for the rule.
* `type` - (Optional, Int) Operate Type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


