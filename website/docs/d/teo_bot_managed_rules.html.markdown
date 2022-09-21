---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_bot_managed_rules"
sidebar_current: "docs-tencentcloud-datasource-teo_bot_managed_rules"
description: |-
  Use this data source to query detailed information of teo botManagedRules
---

# tencentcloud_teo_bot_managed_rules

Use this data source to query detailed information of teo botManagedRules

## Example Usage

```hcl
data "tencentcloud_teo_bot_managed_rules" "botManagedRules" {
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

* `rules` - Managed rules list.
  * `description` - Description of the rule.
  * `rule_id` - Rule ID.
  * `rule_type_name` - Type of the rule.
  * `status` - Status of the rule.


