---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_aggregate_config_rule"
sidebar_current: "docs-tencentcloud-resource-config_aggregate_config_rule"
description: |-
  Provides a resource to create a TencentCloud Config aggregate config rule under an account group.
---

# tencentcloud_config_aggregate_config_rule

Provides a resource to create a TencentCloud Config aggregate config rule under an account group.

## Example Usage

### System preset rule with scheduled trigger

```hcl
resource "tencentcloud_config_aggregate_config_rule" "example" {
  identifier       = "cam-user-group-bound"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-aggregate-rule"
  resource_type    = ["QCS::CAM::User"]
  risk_level       = 3
  account_group_id = "ca-xxxxxxxx"
  description      = "tf example aggregate rule"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }
}
```

### Rule with input parameters and scope filters

```hcl
resource "tencentcloud_config_aggregate_config_rule" "example_with_scope" {
  identifier       = "cam-user-mfa-check"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-rule-with-scope"
  resource_type    = ["QCS::CAM::User"]
  risk_level       = 2
  account_group_id = "ca-xxxxxxxx"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  input_parameter {
    parameter_key = "maxMemorySize"
    type          = "Require"
    value         = "512"
  }

  region_scope               = ["ap-guangzhou"]
  exclude_resource_ids_scope = ["ins-keep"]

  tags_scope {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String) Account group ID. The aggregate rule is created under this account group.
* `identifier_type` - (Required, String) Rule template type. Valid values: `SYSTEM` (system preset), `CUSTOMIZE` (custom).
* `identifier` - (Required, String) Rule template identifier. For system preset rules use the identifier name; for custom rules use the cloud function ARN (region:functionName).
* `resource_type` - (Required, Set) List of supported resource types (e.g. QCS::CAM::User). The rule only applies to resources of the specified types.
* `risk_level` - (Required, Int) Risk level. Valid values: `1` (high risk), `2` (medium risk), `3` (low risk).
* `rule_name` - (Required, String) Rule name.
* `trigger_type` - (Required, List) Trigger type list, up to 2 entries.
* `description` - (Optional, String) Rule description.
* `exclude_resource_ids_scope` - (Optional, Set) Resource ID list excluded from rule evaluation.
* `input_parameter` - (Optional, List) Rule input parameter list.
* `region_scope` - (Optional, Set) Region scope list; the rule only applies to resources in the specified regions.
* `tags_scope` - (Optional, List) Tag scope list; the rule only applies to resources with the specified tags.

The `input_parameter` object supports the following:

* `parameter_key` - (Required, String) Parameter key.
* `type` - (Optional, String) Parameter type: `Require` or `Optional`.
* `value` - (Optional, String) Parameter value.

The `tags_scope` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `trigger_type` object supports the following:

* `message_type` - (Required, String) Trigger message type. Valid values: `ScheduledNotification` (periodic), `ConfigurationItemChangeNotification` (change).
* `maximum_execution_frequency` - (Optional, String) Maximum execution frequency, only for `ScheduledNotification`. e.g. `TwentyFour_Hours`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `configuration` - Current resource configuration.
    * `default_value` - Default value.
    * `description` - Description.
    * `desired_value` - Desired value.
    * `desired_value` - Expected resource configuration.
    * `empty_as` - Condition when empty.
    * `maximum_execution_frequency` - Trigger time period.
    * `message_type` - Trigger type.
    * `operator` - Comparison operator.
    * `operator` - Operator.
    * `parameter_key` - Parameter key.
    * `parameter_key` - Parameter key.
    * `property` - JSON path of the configuration property.
    * `required` - Whether required.
    * `select_path` - Configuration path.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
    * `type` - Parameter type.
    * `type` - Parameter type.
    * `value_type` - Value type.
    * `value` - Parameter value.
  * `account_group_id` - Account group ID.
  * `account_group_name` - Account group name.
  * `annotation` - Non-compliance annotation.
  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `compliance_result` - Compliance result.
  * `config_rule_id` - Rule ID.
  * `config_rule_invoked_time` - Last rule evaluation time.
  * `create_time` - Creation time.
  * `description` - Rule description.
  * `exclude_resource_ids_scope` - Excluded resource IDs.
  * `identifier_type` - CUSTOMIZE: custom rule, SYSTEM: managed rule.
  * `identifier` - Rule identifier.
  * `input_parameter` - Rule parameters.
  * `labels` - Rule labels.
  * `manage_input_parameter` - Managed parameter detail.
  * `manage_trigger_type` - Trigger types supported by preset rules.
  * `regions_scope` - Region scope.
  * `resource_type` - Supported resource types.
  * `risk_level` - Risk level. 1:low, 2:medium, 3:high.
  * `rule_name` - Rule name.
  * `rule_owner_id` - Rule owner user ID.
  * `service_function` - Corresponding cloud function.
  * `source_condition` - Rule trigger conditions.
  * `status` - Rule status. ACTIVE: enabled, NO_ACTIVE: stopped.
  * `tags_scope` - Tag scope.
  * `trigger_type` - Trigger types.
* `config_rule` - Rule detail returned by the cloud.
* `id` - Resource identifier in the form `<account_group_id>#<rule_id>`.
* `rule_id` - Config rule ID. Returned by the cloud after creation and refreshed by Read.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [AddAggregateConfigRule](https://www.tencentcloud.com/document/api/1164/84977)
  - [DescribeAggregateConfigRule](https://www.tencentcloud.com/document/api/1164/84969)
  - [ListAggregateConfigRules](https://www.tencentcloud.com/document/api/1164/66768)
* Read: [DescribeAggregateConfigRule](https://www.tencentcloud.com/document/api/1164/84969)
* Update:
  - [DescribeAggregateConfigRule](https://www.tencentcloud.com/document/api/1164/84969)
  - [UpdateAggregateConfigRule](https://www.tencentcloud.com/document/api/1164/84956)
* Delete: [DeleteAggregateConfigRule](https://www.tencentcloud.com/document/api/1164/84972)

## Import

Config aggregate config rule can be imported using the composite id `<account_group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_config_aggregate_config_rule.example ca-xxxxxxxx#cr-yyyyyyyy
```

