---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_owasp_rules"
sidebar_current: "docs-tencentcloud-datasource-waf_owasp_rules"
description: |-
  Use this data source to query detailed information of WAF owasp rules
---

# tencentcloud_waf_owasp_rules

Use this data source to query detailed information of WAF owasp rules

## Example Usage

```hcl
data "tencentcloud_waf_owasp_rules" "example" {
  domain = "example.qcloud.com"
  by     = "RuleId"
  order  = "desc"
  filters {
    name        = "RuleId"
    values      = ["106251141"]
    exact_match = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain to be queried.
* `by` - (Optional, String) Specifies the field used to sort. valid values: RuleId, ModifyTime.
* `filters` - (Optional, List) Specifies the criteria, support RuleId, TypeId, Desc, CveID, Status, and VulLevel.
* `order` - (Optional, String) Sorting method. supports asc, desc.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `exact_match` - (Required, Bool) Exact search or not.
* `name` - (Required, String) Field name, used for filtering
Filter the sub-order number (value) by DealName.
* `values` - (Required, Set) Values after filtering.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of rules.
  * `create_time` - Creation time.
  * `cve_id` - CVE ID.
  * `description` - Rule description.
  * `level` - Protection level of the rule. valid values: 100 (loose), 200 (normal), 300 (strict), 400 (ultra-strict).
  * `locked` - Whether the user is locked.
  * `modify_time` - Update time.
  * `reason` - Reason for modification

0: none (compatibility records are empty).
1: avoid false positives due to business characteristics.
2: reporting of rule-based false positives.
3: gray release of core business rules.
4: others.
  * `rule_id` - Rule ID.
  * `status` - Rule switch. valid values: 0 (disabled), 1 (enabled), 2 (observation only).
  * `type_id` - Specifies the rule type ID.
  * `vul_level` - Threat level. valid values: 0 (unknown), 100 (low risk), 200 (medium risk), 300 (high risk), 400 (critical).


