---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_bot_portrait_rules"
sidebar_current: "docs-tencentcloud-datasource-teo_bot_portrait_rules"
description: |-
  Use this data source to query detailed information of teo botPortraitRules
---

# tencentcloud_teo_bot_portrait_rules

Use this data source to query detailed information of teo botPortraitRules

## Example Usage

```hcl
data "tencentcloud_teo_bot_portrait_rules" "botPortraitRules" {
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

* `rules` - Portrait rules list.
  * `classification_id` - Classification of the rule. Note: This field may return null, indicating that no valid value can be obtained.
  * `description` - Description of the rule. Note: This field may return null, indicating that no valid value can be obtained.
  * `rule_id` - Rule ID.
  * `rule_type_name` - Type of the rule. Note: This field may return null, indicating that no valid value can be obtained.
  * `status` - Status of the rule. Note: This field may return null, indicating that no valid value can be obtained.


