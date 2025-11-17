---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_owasp_rule_type_config"
sidebar_current: "docs-tencentcloud-resource-waf_owasp_rule_type_config"
description: |-
  Provides a resource to create a WAF owasp rule type config
---

# tencentcloud_waf_owasp_rule_type_config

Provides a resource to create a WAF owasp rule type config

## Example Usage

```hcl
resource "tencentcloud_waf_owasp_rule_type_config" "example" {
  domain           = "demo.com"
  type_id          = "30000000"
  rule_type_status = 1
  rule_type_action = 1
  rule_type_level  = 200
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `type_id` - (Required, String, ForceNew) Rule type ID.
* `rule_type_action` - (Optional, Int) Protection mode of the rule type. valid values: 0 (observation), 1 (intercept).
* `rule_type_level` - (Optional, Int) Protection level of the rule. valid values: 100 (loose), 200 (normal), 300 (strict), 400 (ultra-strict).
* `rule_type_status` - (Optional, Int) The switch status of the rule type. valid values: 0 (disabled), 1 (enabled).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `active_rule` - Indicates the total number of rules enabled under the rule type.
* `classification` - Data type category.
* `description` - Rule type description.
* `rule_type_name` - Rule type name.
* `total_rule` - Specifies all rules under the rule type. always.


## Import

WAF owasp rule type config can be imported using the domain#typeId, e.g.

```
terraform import tencentcloud_waf_owasp_rule_type_config.example demo.com#30000000
```

