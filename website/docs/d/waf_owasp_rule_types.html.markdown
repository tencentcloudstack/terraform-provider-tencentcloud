---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_owasp_rule_types"
sidebar_current: "docs-tencentcloud-datasource-waf_owasp_rule_types"
description: |-
  Use this data source to query detailed information of WAF owasp rule types
---

# tencentcloud_waf_owasp_rule_types

Use this data source to query detailed information of WAF owasp rule types

## Example Usage

```hcl
data "tencentcloud_waf_owasp_rule_types" "example" {
  domain = "demo.com"
  filters {
    name        = "RuleId"
    values      = ["10000001"]
    exact_match = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain names to be queried.
* `filters` - (Optional, List) Filter conditions. supports RuleId, CveID, and Desc.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `exact_match` - (Required, Bool) Exact search or not.
* `name` - (Required, String) Field name, used for filtering
Filter the sub-order number (value) by DealName.
* `values` - (Required, Set) Values after filtering.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Rule type list and information.
  * `action` - Protection mode of the rule type. valid values: 0 (observation), 1 (intercept).
  * `active_rule` - Indicates the total number of rules enabled under the rule type.
  * `classification` - Data type category.
  * `description` - Type description.
  * `level` - Protection level of the rule type. valid values: 100 (loose), 200 (normal), 300 (strict), 400 (ultra-strict).
  * `status` - The switch status of the rule type. valid values: 0 (disabled), 1 (enabled).
  * `total_rule` - Specifies all rules under the rule type. always.
  * `type_id` - Type ID.
  * `type_name` - Type name.


