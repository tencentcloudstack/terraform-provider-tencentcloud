---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_compliance_pack"
sidebar_current: "docs-tencentcloud-resource-config_compliance_pack"
description: |-
  Provides a resource to create a Config compliance pack.
---

# tencentcloud_config_compliance_pack

Provides a resource to create a Config compliance pack.

~> **NOTE:** The current resource `tencentcloud_config_compliance_pack` does not support destroy, please contact the work order for processing.

## Example Usage

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

### With status control

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

## Argument Reference

The following arguments are supported:

* `compliance_pack_name` - (Required, String) Compliance pack name.
* `config_rules` - (Required, Set) List of compliance pack rules.
* `risk_level` - (Required, Int) Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `description` - (Optional, String) Description of the compliance pack.
* `status` - (Optional, String) Compliance pack status. Valid values: ACTIVE, UN_ACTIVE.

The `config_rules` object supports the following:

* `identifier` - (Required, String) Rule identifier (managed rule name or custom rule cloud function ARN).
* `compliance_pack_id` - (Optional, String) Compliance pack ID that this rule belongs to.
* `config_rule_id` - (Optional, String) Config rule ID.
* `description` - (Optional, String) Rule description.
* `input_parameter` - (Optional, List) Rule input parameters.
* `managed_rule_identifier` - (Optional, String) Managed rule identifier (preset rule identity).
* `risk_level` - (Optional, Int) Rule risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `rule_name` - (Optional, String) Rule name.

The `input_parameter` object of `config_rules` supports the following:

* `parameter_key` - (Required, String) Parameter key.
* `type` - (Optional, String) Parameter type: Require or Optional.
* `value` - (Optional, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `compliance_pack_id` - Compliance pack ID.
* `create_time` - Creation time of the compliance pack.


## Import

Config compliance pack can be imported using the id, e.g.

```
terraform import tencentcloud_config_compliance_pack.example cp-33mA27YUlOJWG4sJ53Sx
```

