Provides a resource to create a Config compliance pack.

~> **NOTE:** The current resource `tencentcloud_config_compliance_pack` does not support destroy, please contact the work order for processing.

Example Usage

```hcl
resource "tencentcloud_config_compliance_pack" "example" {
  compliance_pack_name = "tf-example"
  risk_level           = 2
  description          = "tf example compliance pack"

  config_rules {
    identifier = "cos-default-encryption-kms"
    rule_name  = "my-rule1"
    risk_level = 1
  }

  config_rules {
    identifier  = "cam-user-group-bound"
    rule_name   = "my-rule2"
    description = "rule description"
    risk_level  = 3
    input_parameter {
      parameter_key = "maxMemorySize"
      type          = "Require"
      value         = "512"
    }
  }
}
```

With status control

```hcl
resource "tencentcloud_config_compliance_pack" "example" {
  compliance_pack_name = "tf-example"
  risk_level           = 1
  description          = "high risk compliance pack"
  status               = "UN_ACTIVE"

  config_rules {
    identifier = "cam-user-mfa-check"
    risk_level = 1
  }
}
```

Import

Config compliance pack can be imported using the id, e.g.

```
terraform import tencentcloud_config_compliance_pack.example cp-33mA27YUlOJWG4sJ53Sx
```
