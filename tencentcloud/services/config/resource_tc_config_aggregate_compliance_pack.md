Provides a resource to create a Config aggregate (account group) compliance pack.

Example Usage

```hcl
resource "tencentcloud_config_aggregate_compliance_pack" "example" {
  compliance_pack_name = "tf-example-aggregate-pack"
  account_group_id     = "ca-xxxxxxxx"
  risk_level           = 1
  description           = "tf example aggregate compliance pack"

  config_rules {
    rule_name          = "rule-1"
    risk_level         = 1
    identifier          = "conformance-pack-1"
    managed_rule_identifier = "qcloud::cvm::check-instance-vpc-subnet"
    description         = "rule description"

    input_parameter {
      parameter_key = "param1"
      type           = "Require"
      value          = "value1"
    }
  }
}
```

Import

Config aggregate compliance pack can be imported using the composite id, e.g. `account_group_id#compliance_pack_id`.

```
terraform import tencentcloud_config_aggregate_compliance_pack.example ca-xxxxxxxx#cp-xxxxxxxx
```

-> `account_group_id` is immutable: changing it forces resource replacement. The `managed_rule_identifier` and per-rule `compliance_pack_id` fields are not echoed back by the Describe API, so drift detection for those fields is limited; their values are preserved from the Terraform configuration on read.