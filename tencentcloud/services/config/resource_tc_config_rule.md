Provides a resource to create a Config rule.

Example Usage

System preset rule with scheduled trigger

```hcl
resource "tencentcloud_config_rule" "example" {
  identifier      = "cam-user-group-bound"
  identifier_type = "SYSTEM"
  rule_name       = "CAM访问管理子用户必须关联用户组"
  resource_type   = ["QCS::CAM::User"]
  risk_level      = 3
  description     = "帐号访问管理中用户至少关联一个用户组，则符合规则。"

  trigger_type {
    message_type               = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  status = "ACTIVE"
}
```

Rule with input parameters

```hcl
resource "tencentcloud_config_rule" "example_with_params" {
  identifier      = "cam-user-mfa-check"
  identifier_type = "SYSTEM"
  rule_name       = "CAM子用户开启MFA"
  resource_type   = ["QCS::CAM::User"]
  risk_level      = 2

  trigger_type {
    message_type               = "ScheduledNotification"
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

Import

Config rule can be imported using the ruleId, e.g.

```
terraform import tencentcloud_config_rule.example cr-3xhsd76j603v0a8ma0i73
```
