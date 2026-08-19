---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_scene_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_scene_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive scene rule
---

# tencentcloud_waf_api_sec_sensitive_scene_rule

Provides a resource to create a WAF api sec sensitive scene rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_scene_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 1
  source    = "custom"

  rule_list {
    key     = "api"
    operate = "equal"
    value   = ["/login", "/user"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) Scene name.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `rule_list` - (Optional, List) Rule list.
* `source` - (Optional, String) Rule source, system built-in: OS, customer custom: custom.

The `rule_list` object supports the following:

* `key` - (Optional, String) Match field.
* `name` - (Optional, String) When the match field is get parameter value, post parameter value, cookie parameter value, header parameter value or rsp parameter value, this field can be filled.
* `operate` - (Optional, String) Operator.
* `value` - (Optional, Set) Match value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Update timestamp.


## Import

WAF api sec sensitive scene rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_scene_rule.example www.example.com#tf-example
```

