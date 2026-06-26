---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_custom_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_custom_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive custom rule
---

# tencentcloud_waf_api_sec_sensitive_custom_rule

Provides a resource to create a WAF api sec sensitive custom rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 1
  position    = ["headers"]
  match_key   = "key_match"
  match_value = ["admin", "cookie"]
  level       = "100"
  match_cond  = ["and"]
  is_pan      = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) Rule name.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `is_pan` - (Optional, Int) Whether the rule is generalized, default 0 means not generalized.
* `level` - (Optional, String) Risk level.
* `match_cond` - (Optional, Set: [`String`]) Match symbol, pass this value when the match condition is keyword match or character match, multiple values can be passed.
* `match_key` - (Optional, String) Match condition.
* `match_value` - (Optional, Set: [`String`]) Match value.
* `position` - (Optional, Set: [`String`]) Parameter position.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WAF api sec sensitive custom rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_rule.example www.example.com#tf-example
```

