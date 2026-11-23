---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_rule"
sidebar_current: "docs-tencentcloud-resource-config_rule"
description: |-
  Provides a resource to create a Config rule.
---

# tencentcloud_config_rule

Provides a resource to create a Config rule.

## Example Usage

### System preset rule with scheduled trigger

```hcl
resource "tencentcloud_config_rule" "example" {
  identifier      = "cam-user-group-bound"
  identifier_type = "SYSTEM"
  rule_name       = "CAM访问管理子用户必须关联用户组"
  resource_type   = ["QCS::CAM::User"]
  risk_level      = 3
  description     = "帐号访问管理中用户至少关联一个用户组，则符合规则。"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  status = "ACTIVE"
}
```

### Rule with input parameters

```hcl
resource "tencentcloud_config_rule" "example_with_params" {
  identifier      = "cam-user-mfa-check"
  identifier_type = "SYSTEM"
  rule_name       = "CAM子用户开启MFA"
  resource_type   = ["QCS::CAM::User"]
  risk_level      = 2

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  input_parameter {
    parameter_key = "maxMemorySize"
    type          = "Require"
    value         = "512"
  }

  status = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

* `identifier_type` - (Required, String, ForceNew) Rule template type. Valid values: SYSTEM (system preset), CUSTOMIZE (custom).
* `identifier` - (Required, String, ForceNew) Rule template identifier. For system preset rules use the identifier name; for custom rules use the cloud function ARN (region:functionName).
* `resource_type` - (Required, List: [`String`], ForceNew) Supported resource type list (e.g. QCS::CAM::User).
* `risk_level` - (Required, Int) Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `rule_name` - (Required, String) Rule name.
* `trigger_type` - (Required, List) Trigger type list, up to 2 entries.
* `description` - (Optional, String) Rule description (0~1024 characters).
* `exclude_resource_ids_scope` - (Optional, List: [`String`]) Resource ID list excluded from rule evaluation.
* `input_parameter` - (Optional, List) Rule input parameter list.
* `regions_scope` - (Optional, List: [`String`]) Region scope list; rule only applies to resources in the specified regions.
* `status` - (Optional, String) Rule status. Valid values: ACTIVE (enabled), UN_ACTIVE (disabled).
* `tags_scope` - (Optional, List) Tag scope list; rule only applies to resources with the specified tags.

The `input_parameter` object supports the following:

* `parameter_key` - (Required, String) Parameter key.
* `type` - (Optional, String) Parameter type: Require or Optional.
* `value` - (Optional, String) Parameter value.

The `tags_scope` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

The `trigger_type` object supports the following:

* `message_type` - (Required, String) Trigger message type. Valid values: ScheduledNotification, ConfigurationItemChangeNotification.
* `maximum_execution_frequency` - (Optional, String) Maximum execution frequency (only for ScheduledNotification). e.g. TwentyFour_Hours.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `compliance_result` - Compliance result.
* `config_rule_id` - Config rule ID.
* `config_rule_invoked_time` - Last rule evaluation time.
* `create_time` - Creation time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [AddConfigRule](https://www.tencentcloud.com/document/api/1164/84976)
  - [CloseConfigRule](https://www.tencentcloud.com/document/api/1164/84974)
  - [OpenConfigRule](https://www.tencentcloud.com/document/api/1164/84960)
* Read: [DescribeConfigRule](https://www.tencentcloud.com/document/api/1164/84968)
* Update:
  - [CloseConfigRule](https://www.tencentcloud.com/document/api/1164/84974)
  - [OpenConfigRule](https://www.tencentcloud.com/document/api/1164/84960)
  - [UpdateConfigRule](https://www.tencentcloud.com/document/api/1164/84955)
* Delete:
  - [CloseConfigRule](https://www.tencentcloud.com/document/api/1164/84974)
  - [DeleteConfigRule](https://www.tencentcloud.com/document/api/1164/84971)
  - [OpenConfigRule](https://www.tencentcloud.com/document/api/1164/84960)

## Import

Config rule can be imported using the ruleId, e.g.

```
terraform import tencentcloud_config_rule.example cr-3xhsd76j603v0a8ma0i73
```

