---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_rule"
sidebar_current: "docs-tencentcloud-datasource-teo_l7_acc_rule"
description: |-
  Use this data source to query detailed information of TEO L7 access control rules
---

# tencentcloud_teo_l7_acc_rule

Use this data source to query detailed information of TEO L7 access control rules

## Example Usage

### Query all L7 access rules by zone ID

```hcl
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

### Query specific L7 access rule

```hcl
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
  rule_id = "rule-xxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the site ID.
* `rule_id` - (Optional, String) Specifies the rule ID. If not specified, return all rules in the zone.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - (Int) Total count of rules matching the query criteria.
* `rules` - (List) L7 access control rules.
  * `rule_id` - (String) Rule ID. Unique identifier of the rule.
  * `rule_name` - (String) Rule name.
  * `description` - (List) Rule annotation.
  * `rule_priority` - (Int) Rule priority.
  * `branches` - (List) Sub-Rule branch.
