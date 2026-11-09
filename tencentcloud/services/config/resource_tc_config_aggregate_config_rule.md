Provides a resource to create a TencentCloud Config aggregate config rule under an account group.

Example Usage

System preset rule with scheduled trigger

```hcl
resource "tencentcloud_config_aggregate_config_rule" "example" {
  identifier        = "cam-user-group-bound"
  identifier_type   = "SYSTEM"
  rule_name         = "tf-example-aggregate-rule"
  resource_type     = ["QCS::CAM::User"]
  risk_level        = 3
  account_group_id = "ca-xxxxxxxx"
  description       = "tf example aggregate rule"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }
}
```

Rule with input parameters and scope filters

```hcl
resource "tencentcloud_config_aggregate_config_rule" "example_with_scope" {
  identifier        = "cam-user-mfa-check"
  identifier_type   = "SYSTEM"
  rule_name         = "tf-example-rule-with-scope"
  resource_type     = ["QCS::CAM::User"]
  risk_level        = 2
  account_group_id = "ca-xxxxxxxx"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  input_parameter {
    parameter_key = "maxMemorySize"
    type            = "Require"
    value           = "512"
  }

  region_scope               = ["ap-guangzhou"]
  exclude_resource_ids_scope = ["ins-keep"]

  tags_scope {
    tag_key   = "env"
    tag_value = "prod"
  }
}
```

Import

Config aggregate config rule can be imported using the composite id `<account_group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_config_aggregate_config_rule.example ca-xxxxxxxx#cr-yyyyyyyy
```