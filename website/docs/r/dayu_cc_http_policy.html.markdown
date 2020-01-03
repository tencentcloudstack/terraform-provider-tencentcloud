---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_cc_http_policy"
sidebar_current: "docs-tencentcloud-resource-dayu_cc_http_policy"
description: |-
  Use this resource to create a dayu CC self-define http policy
---

# tencentcloud_dayu_cc_http_policy

Use this resource to create a dayu CC self-define http policy

## Example Usage

```hcl
resource "tencentcloud_dayu_cc_http_policy" "test_bgpip" {
  resource_type = "bgpip"
  resource_id   = "bgpip-00000294"
  name          = "policy_match"
  smode         = "matching"
  action        = "drop"
  switch        = true
  rule_list {
    skey     = "host"
    operator = "include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_net" {
  resource_type = "net"
  resource_id   = "net-0000007e"
  name          = "policy_match"
  smode         = "matching"
  action        = "drop"
  switch        = true
  rule_list {
    skey     = "cgi"
    operator = "equal"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_bgpmultip" {
  resource_type = "bgp-multip"
  resource_id   = "bgp-0000008o"
  name          = "policy_match"
  smode         = "matching"
  action        = "alg"
  switch        = true
  ip            = "111.230.178.25"

  rule_list {
    skey     = "referer"
    operator = "not_include"
    value    = "123"
  }
}

resource "tencentcloud_dayu_cc_http_policy" "test_bgp" {
  resource_type = "bgp"
  resource_id   = "bgp-000006mq"
  name          = "policy_match"
  smode         = "matching"
  action        = "alg"
  switch        = true

  rule_list {
    skey     = "ua"
    operator = "not_include"
    value    = "123"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the CC self-define http policy. Length should between 1 and 20.
* `resource_id` - (Required, ForceNew) ID of the resource that the CC self-define http policy works for.
* `resource_type` - (Required, ForceNew) Type of the resource that the CC self-define http policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.
* `action` - (Optional) Action mode, only valid when `smode` is `matching`. Valid values are `alg` and `drop`.
* `frequency` - (Optional) Max frequency per minute, only valid when `smode` is `speedlimit`, the valid value ranges from 1 to 10000.
* `ip` - (Optional) Ip of the CC self-define http policy, only valid when `resource_type` is `bgp-multip`. The num of list items can only be set one.
* `rule_list` - (Optional) Rule list of the CC self-define http policy,  only valid when `smode` is `matching`.
* `smode` - (Optional) Match mode, and valid values are `matching`, `speedlimit`. Note: the speed limit type CC self-define policy can only set one.
* `switch` - (Optional) Indicate the CC self-define http policy takes effect or not.

The `rule_list` object supports the following:

* `operator` - (Optional) Operator of the rule, valid values are `include`, `not_include`, `equal`.
* `skey` - (Optional) Key of the rule, valid values are `host`, `cgi`, `ua`, `referer`.
* `value` - (Optional) Rule value, then length should be less than 31 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the CC self-define http policy.
* `policy_id` - Id of the CC self-define http policy.


