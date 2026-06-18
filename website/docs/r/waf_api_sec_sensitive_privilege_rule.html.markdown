---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_privilege_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_privilege_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive privilege rule
---

# tencentcloud_waf_api_sec_sensitive_privilege_rule

Provides a resource to create a WAF api sec sensitive privilege rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_privilege_rule" "example" {
  domain         = "www.example.com"
  rule_name      = "tf-example"
  status         = 1
  api_name       = ["/api/user/info"]
  position       = "QUERY"
  parameter_list = ["parameter"]
  option         = 1
  source         = "custom"

  api_name_op {
    op    = "belong"
    value = ["/api/user/info"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) Rule name, cannot be repeated.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `api_name_op` - (Optional, List) API match list.
* `api_name` - (Optional, Set: [`String`]) Up to 20 APIs can be entered.
* `option` - (Optional, Int) Application object value, 1 means manual filling, 2 means obtaining from API assets.
* `parameter_list` - (Optional, Set: [`String`]) Authentication parameter list.
* `position` - (Optional, String) Authentication position.
* `source` - (Optional, String) Rule source.

The `api_name_method` object of `api_name_op` supports the following:

* `api_name` - (Optional, String) API name.
* `method` - (Optional, String) API request method.

The `api_name_op` object supports the following:

* `api_name_method` - (Optional, List) When manually filtering, this structure should be passed.
* `op` - (Optional, String) Match method, such as belong and regex.
* `value` - (Optional, Set) Match value list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Update timestamp.


## Import

WAF api sec sensitive privilege rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_privilege_rule.example www.example.com#tf-example
```

