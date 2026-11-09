---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_list_aggregate_config_rules"
sidebar_current: "docs-tencentcloud-data_source-config_list_aggregate_config_rules"
description: |-
  Use this data source to query the list of Config aggregate config rules under an account group.
---

# tencentcloud_config_list_aggregate_config_rules

Use this data source to query the list of Config aggregate config rules under an account group.

## Example Usage

### Query aggregate config rules by account group id

```hcl
data "tencentcloud_config_list_aggregate_config_rules" "example" {
  account_group_id = "ca-aggregator-abcdefg"
}
```

### Query aggregate config rules with filters

```hcl
data "tencentcloud_config_list_aggregate_config_rules" "example" {
  account_group_id  = "ca-aggregator-abcdefg"
  risk_level        = [1, 2]
  state             = "ACTIVE"
  compliance_result = ["NON_COMPLIANT"]
  order_type        = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String) Account group ID.
* `compliance_result` - (Optional, List) Compliance result list for filtering.
* `order_type` - (Optional, String) Sort order of the rule list. Valid values: `desc` (descending), `asc` (ascending).
* `risk_level` - (Optional, List) Risk level list for filtering. Valid values: `1` (high risk), `2` (medium risk), `3` (low risk).
* `rule_name` - (Optional, String) Rule name for filtering.
* `rule_owner_id` - (Optional, Int) Rule owner account ID.
* `state` - (Optional, String) Rule state for filtering.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `configuration` - Actual configuration of the resource, i.e. the non-compliant part. Up to 256 characters.
    * `default_value` - Default value.
    * `description` - Description.
    * `desired_value` - Expected configuration of the resource. Up to 256 characters.
    * `desired_value` - Expected value.
    * `empty_as` - Result when the condition is empty. `COMPLIANT` for compliant, `NON_COMPLIANT` for non-compliant, `NOT_APPLICABLE` for not applicable.
    * `maximum_execution_frequency` - Trigger period.
    * `message_type` - Trigger type.
    * `operator` - Comparison operator between the actual and expected configuration. Up to 16 characters. May be empty for custom rules.
    * `operator` - Comparison operator.
    * `parameter_key` - Parameter key.
    * `parameter_key` - Parameter name.
    * `property` - JSON path of the current configuration in the resource attribute structure. Up to 256 characters. May be empty for custom rules.
    * `required` - Whether the condition is required.
    * `select_path` - Configuration path.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
    * `type` - Parameter type. `Require` for required, `Optional` for optional.
    * `type` - Parameter type. `Require` for required, `Optional` for optional.
    * `value_type` - Value type. `Integer` for number, `String` for string.
    * `value` - Parameter value.
  * `account_group_id` - Account group ID.
  * `account_group_name` - Account group name.
  * `annotation` - Annotation of the non-compliant resource.
  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `compliance_result` - Compliance result. `COMPLIANT` for compliant, `NON_COMPLIANT` for non-compliant, `NOT_APPLICABLE` for not applicable.
  * `config_rule_id` - Rule ID.
  * `config_rule_invoked_time` - Rule evaluation time in the format YYYY-MM-DD h:i:s.
  * `create_time` - Creation time in the format YYYY-MM-DD h:i:s.
  * `description` - Rule description.
  * `exclude_resource_ids_scope` - Resource IDs excluded from evaluation.
  * `identifier_type` - Rule type. `CUSTOMIZE` for custom rule, `SYSTEM` for managed rule.
  * `identifier` - Rule identifier.
  * `input_parameter` - Rule parameters.
  * `labels` - Rule labels.
  * `manage_input_parameter` - Parameter details.
  * `manage_trigger_type` - Trigger types supported by preset rules. `ScheduledNotification` for periodic, `ConfigurationItemChangeNotification` for change.
  * `regions_scope` - Associated regions.
  * `resource_type` - Supported resource types. The rule only applies to the specified resource types.
  * `risk_level` - Rule risk level. `1` for low risk, `2` for medium risk, `3` for high risk.
  * `rule_name` - Rule name.
  * `rule_owner_id` - Rule owner user ID.
  * `service_function` - Corresponding service function.
  * `source_condition` - Rule trigger conditions.
  * `status` - Rule status. `ACTIVE` for enabled, `NO_ACTIVE` for stopped.
  * `tags_scope` - Associated tags.
  * `trigger_type` - Trigger types. `ScheduledNotification` for periodic trigger, `ConfigurationItemChangeNotification` for change trigger.
* `items` - List of matched aggregate config rules.
* `total` - Total number of matched rules.


