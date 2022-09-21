---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_waf_rule_groups"
sidebar_current: "docs-tencentcloud-datasource-teo_waf_rule_groups"
description: |-
  Use this data source to query detailed information of teo wafRuleGroups
---

# tencentcloud_teo_waf_rule_groups

Use this data source to query detailed information of teo wafRuleGroups

## Example Usage

```hcl
data "tencentcloud_teo_waf_rule_groups" "wafRuleGroups" {
  zone_id = ""
  entity  = ""
}
```

## Argument Reference

The following arguments are supported:

* `entity` - (Required, String) Subdomain or application name.
* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `waf_rule_groups` - List of WAF rule groups.
  * `rule_type_desc` - Description of rule type in this group.
  * `rule_type_id` - Type id of rules in this group.
  * `rule_type_name` - Type name of rules in this group.
  * `rules` - Rules detail.
    * `description` - Description of the rule.
    * `rule_id` - WAF managed rule id.
    * `rule_level_desc` - System default level of the rule.
    * `rule_tags` - Tags of the rule. Note: This field may return null, indicating that no valid value can be obtained.


