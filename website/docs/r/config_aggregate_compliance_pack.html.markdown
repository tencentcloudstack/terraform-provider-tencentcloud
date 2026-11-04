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
  compliance_pack_name = "tf-example-aggregate-pack"
  account_group_id     = "ca-xxxxxxxx"
  risk_level           = 1
  description          = "tf example aggregate compliance pack"

  config_rules {
    rule_name               = "rule-1"
    risk_level              = 1
    identifier              = "conformance-pack-1"
    managed_rule_identifier = "qcloud::cvm::check-instance-vpc-subnet"
    description             = "rule description"

    input_parameter {
      parameter_key = "param1"
      type          = "Require"
      value         = "value1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_group_id` - (Required, String) Account group id. Immutable: changing it forces resource replacement.
* `compliance_pack_name` - (Required, String) Compliance pack name.
* `config_rules` - (Required, List) List of compliance pack rules.
* `risk_level` - (Required, Int) Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `description` - (Optional, String) Description of the compliance pack.

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
* `status` - Compliance pack status.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [AddAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85010)
  - [DescribeAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85006)
* Read: [DescribeAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85006)
* Update:
  - [DescribeAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85006)
  - [UpdateAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/84998)
* Delete: [DeleteAggregateCompliancePack](https://www.tencentcloud.com/document/api/1164/85008)

## Import

Config aggregate compliance pack can be imported using the composite id, e.g. `account_group_id#compliance_pack_id`.

```
terraform import tencentcloud_config_aggregate_compliance_pack.example ca-xxxxxxxx#cp-xxxxxxxx
```

-> `account_group_id` is immutable: changing it forces resource replacement. The `managed_rule_identifier` and per-rule `compliance_pack_id` fields are not echoed back by the Describe API, so drift detection for those fields is limited; their values are preserved from the Terraform configuration on read.

