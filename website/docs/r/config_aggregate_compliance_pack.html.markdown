---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_aggregate_compliance_pack"
sidebar_current: "docs-tencentcloud-resource-config_aggregate_compliance_pack"
description: |-
  Provides a resource to create a Config aggregate (account group) compliance pack.
---

# tencentcloud_config_aggregate_compliance_pack

Provides a resource to create a Config aggregate (account group) compliance pack.

## Example Usage

```hcl
resource "tencentcloud_config_aggregate_compliance_pack" "example" {
  account_group_id     = "ca-Zp9EiU7Pa4Vy0lFFm8XM"
  compliance_pack_name = "tf-example"
  risk_level           = 1
  description          = "remark."
  status               = "ACTIVE"

  config_rules = [
    {
      rule_name               = "CAM用户密码策略符合要求"
      risk_level              = 1
      identifier              = "cam-password-policy-check"
      managed_rule_identifier = "cam-password-policy-check"
      description             = "CAM用户密码策略中各项配置满足参数设置的值，视为“合规”。CAM用户密码策略中各项配置不满足参数设置的值，视为“不合规”。"
      input_parameter = [
        { parameter_key = "PasswordExpirationInvalidation", type = "Optional", value = "true" },
        { parameter_key = "ForcePasswordChange", type = "Optional", value = "365" },
        { parameter_key = "RetryPasswordLimit", type = "Optional", value = "10" },
        { parameter_key = "MinimumLength", type = "Optional", value = "8" },
        { parameter_key = "ReusePasswordLimit", type = "Optional", value = "1" },
        { parameter_key = "requireLowercaseCharacters", type = "Optional", value = "true" },
        { parameter_key = "requireNumbers", type = "Optional", value = "true" },
        { parameter_key = "requireSymbols", type = "Optional", value = "true" },
        { parameter_key = "requireUppercaseCharacters", type = "Optional", value = "true" },
      ]
    },
    {
      rule_name               = "操作审计开启全量日志跟踪"
      risk_level              = 1
      identifier              = "actiontrail-trail-intact-enabled"
      managed_rule_identifier = "actiontrail-trail-intact-enabled"
      description             = "操作审计中存在开启状态的跟踪，且跟踪全部地域和全部事件类型，视为“合规”，否则视为不合规。"
      input_parameter         = []
    },
    {
      rule_name               = "使用主机安全、容器安全专业版或以上版本"
      risk_level              = 2
      identifier              = "cwp-and-tcss-version-check"
      managed_rule_identifier = "cwp-and-tcss-version-check"
      description             = "主机安全为轻量版、专业版或以上版本，同时容器安全服务为专业版或以上版本，视为“合规”；否则视为“不合规”。"
      input_parameter         = []
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String) Account group id. Immutable: changing it forces resource replacement.
* `compliance_pack_name` - (Required, String) Compliance pack name.
* `config_rules` - (Required, List) List of compliance pack rules.
* `risk_level` - (Required, Int) Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `description` - (Optional, String) Description of the compliance pack.
* `status` - (Optional, String) Compliance pack status. Valid values: `ACTIVE` (enabled) and `UN_ACTIVE` (disabled). A pack is created as `ACTIVE`; set `UN_ACTIVE` to disable it.

The `config_rules` object supports the following:

* `input_parameter` - (Required, List) Rule input parameters.
* `risk_level` - (Required, Int) Rule risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `rule_name` - (Required, String) Rule name.
* `compliance_pack_id` - (Optional, String) Compliance pack id that this rule belongs to.
* `config_rule_id` - (Optional, String) Config rule id.
* `description` - (Optional, String) Rule description.
* `identifier` - (Optional, String) Rule identifier (managed rule name or custom rule cloud function ARN).
* `managed_rule_identifier` - (Optional, String) Managed rule identifier (preset rule identity).

The `input_parameter` object of `config_rules` supports the following:

* `parameter_key` - (Required, String) Parameter key.
* `type` - (Optional, String) Parameter type. Required: `Require`, optional: `Optional`.
* `value` - (Optional, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `compliance_pack_id` - Compliance pack id.
* `create_time` - Creation time of the compliance pack.
* `id` - Resource identifier in the form `<account_group_id>#<compliance_pack_id>`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85010)
* Update: [UpdateAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/84998)
* Delete: [DeleteAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85008)

## Import

Config aggregate compliance pack can be imported using the composite id, e.g. `account_group_id#compliance_pack_id`.

```
terraform import tencentcloud_config_aggregate_compliance_pack.example ca-Zp9EiU7Pa4Vy0lFFm8XM#cp-VNsYhja10ewnInqtc0qp
```

-> `account_group_id` is immutable: changing it forces resource replacement. The `managed_rule_identifier` and per-rule `compliance_pack_id` fields are not echoed back by the Describe API, so drift detection for those fields is limited; their values are preserved from the Terraform configuration on read.

