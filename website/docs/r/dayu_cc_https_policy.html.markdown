---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_cc_https_policy"
sidebar_current: "docs-tencentcloud-resource-dayu_cc_https_policy"
description: |-
  Use this resource to create a dayu CC self-define https policy
---

# tencentcloud_dayu_cc_https_policy

Use this resource to create a dayu CC self-define https policy

~> **NOTE:** creating CC self-define https policy need a valid resource `tencentcloud_dayu_l7_rule`; The resource only support Anti-DDoS of resource type `bgpip`.

## Example Usage

```hcl
resource "tencentcloud_dayu_cc_https_policy" "test_policy" {
  resource_type = tencentcloud_dayu_l7_rule.test_rule.resource_type
  resource_id   = tencentcloud_dayu_l7_rule.test_rule.resource_id
  rule_id       = tencentcloud_dayu_l7_rule.test_rule.rule_id
  domain        = tencentcloud_dayu_l7_rule.test_rule.domain
  name          = "policy_test"
  exe_mode      = "drop"
  switch        = true

  rule_list {
    skey     = "cgi"
    operator = "include"
    value    = "123"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) Domain that the CC self-define https policy works for, only valid when `protocol` is `https`.
* `name` - (Required, ForceNew) Name of the CC self-define https policy. Length should between 1 and 20.
* `resource_id` - (Required, ForceNew) ID of the resource that the CC self-define https policy works for.
* `resource_type` - (Required, ForceNew) Type of the resource that the CC self-define https policy works for, valid value is `bgpip`.
* `rule_id` - (Required, ForceNew) Rule id of the domain that the CC self-define https policy works for, only valid when `protocol` is `https`.
* `rule_list` - (Required) Rule list of the CC self-define https policy.
* `exe_mode` - (Optional) Execute mode. Valid values are `alg` and `drop`.
* `switch` - (Optional) Indicate the CC self-define https policy takes effect or not.

The `rule_list` object supports the following:

* `operator` - (Required) Operator of the rule, valid values are `include` and `equal`.
* `skey` - (Required) Key of the rule, valid values are `cgi`, `ua` and `referer`.
* `value` - (Required) Rule value, then length should be less than 31 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the CC self-define https policy.
* `ip_list` - Ip of the CC self-define https policy.
* `policy_id` - Id of the CC self-define https policy.


