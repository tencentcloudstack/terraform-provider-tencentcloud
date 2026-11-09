---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results"
sidebar_current: "docs-tencentcloud-data_source-config_list_aggregate_config_rule_resource_evaluation_results"
description: |-
  Use this data source to query Config aggregate config rule resource evaluation results.
---

# tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results

Use this data source to query Config aggregate config rule resource evaluation results.

## Example Usage

### Query aggregate config rule resource evaluation results for a specific account group

```hcl
data "tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results" "example" {
  account_group_id = "ca-aggregator-abcdefg"
}
```

### Query aggregate config rule resource evaluation results for all account groups

```hcl
data "tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results" "example" {}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Optional, String) Account group ID. When omitted, the aggregate full results are queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `configuration` - Actual resource configuration, i.e. the non-compliant configuration. Up to 256 characters.
    * `desired_value` - Expected resource configuration, i.e. the compliant configuration. Up to 256 characters.
    * `operator` - Comparison operator between the actual and expected configuration. Up to 16 characters. This field may be empty for custom rules.
    * `property` - JSON path of the current configuration in the resource attribute structure. Up to 256 characters. This field may be empty for custom rules.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `annotation` - Compliance detail.
  * `resource_id` - Resource ID.
  * `resource_name` - Resource name.
  * `resource_region` - Resource region.
  * `resource_tags` - Resource tags.
  * `resource_type` - Resource type.
  * `rule_description` - Rule description.
  * `rule_id` - Rule ID.
  * `rule_identifier` - Rule identifier.
  * `rule_name` - Rule name.
  * `rule_owner_id` - Owner account ID of the rule.
  * `rule_risk_level` - Risk level of the rule.
* `items` - List of config rule resource evaluation results.


