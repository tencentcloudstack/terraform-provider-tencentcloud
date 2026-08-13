---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_api_sec_sensitive_white_rule"
sidebar_current: "docs-tencentcloud-resource-waf_api_sec_sensitive_white_rule"
description: |-
  Provides a resource to create a WAF api sec sensitive white rule
---

# tencentcloud_waf_api_sec_sensitive_white_rule

Provides a resource to create a WAF api sec sensitive white rule

## Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_white_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 1
  white_mode  = 1
  description = "tf example white rule"

  api_name_op {
    op    = "belong"
    value = ["/api/user/info"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_name` - (Required, String, ForceNew) White rule name.
* `status` - (Required, Int) Rule switch, 0: off, 1: on.
* `api_name_op` - (Optional, List) API match list.
* `description` - (Optional, String) Rule description.
* `white_fields` - (Optional, List) White field config list.
* `white_mode` - (Optional, Int) White mode. Enum values: 1: whitelist the whole API, 2: whitelist specified fields.

The `api_name_method` object of `api_name_op` supports the following:

* `api_name` - (Optional, String) API name.
* `method` - (Optional, String) API request method.

The `api_name_op` object supports the following:

* `api_name_method` - (Optional, List) When manually filtering, this structure should be passed.
* `op` - (Optional, String) Match method, such as belong and regex.
* `value` - (Optional, Set) Match value list.

The `white_fields` object supports the following:

* `field_name` - (Optional, String) Field name.
* `field_type` - (Optional, String) Field position.
* `sensitive_types` - (Optional, Set) Sensitive data type list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `update_time` - Update timestamp.


## Import

WAF api sec sensitive white rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_white_rule.example www.example.com#tf-example
```

