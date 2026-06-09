---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_rules"
sidebar_current: "docs-tencentcloud-datasource-config_rules"
description: |-
  Use this data source to query detailed information of Config rules.
---

# tencentcloud_config_rules

Use this data source to query detailed information of Config rules.

## Example Usage

### Query all config rules

```hcl
data "tencentcloud_config_rules" "example" {}
```

### Query config rules by name

```hcl
data "tencentcloud_config_rules" "example" {
  rule_name = "cam-user-mfa-check"
}
```

### Query config rules by filters

```hcl
data "tencentcloud_config_rules" "example" {
  risk_level        = [1, 2]
  state             = "ACTIVE"
  compliance_result = ["NON_COMPLIANT"]
  order_type        = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `compliance_result` - (Optional, List: [`String`]) Compliance result list for filtering. Valid values: COMPLIANT, NON_COMPLIANT.
* `order_type` - (Optional, String) Sort type by rule name. Valid values: desc (descending), asc (ascending).
* `result_output_file` - (Optional, String) Used to save results.
* `risk_level` - (Optional, List: [`Int`]) Risk level list for filtering. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
* `rule_name` - (Optional, String) Rule name for filtering.
* `state` - (Optional, String) Rule state for filtering. Valid values: ACTIVE, UN_ACTIVE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rule_list` - Config rule list.
  * `compliance_pack_id` - Compliance pack ID.
  * `compliance_pack_name` - Compliance pack name.
  * `compliance_result` - Compliance result. Valid values: COMPLIANT, NON_COMPLIANT, NOT_APPLICABLE.
  * `config_rule_id` - Config rule ID.
  * `config_rule_invoked_time` - Rule evaluation time.
  * `create_time` - Creation time.
  * `description` - Rule description.
  * `identifier_type` - Rule type. Valid values: CUSTOMIZE (custom rule), SYSTEM (managed rule).
  * `identifier` - Rule identifier.
  * `labels` - Rule label list.
  * `resource_type` - Supported resource type list.
  * `risk_level` - Risk level. Valid values: 1 (low risk), 2 (medium risk), 3 (high risk).
  * `rule_name` - Rule name.
  * `service_function` - Corresponding service function.
  * `status` - Rule status. Valid values: ACTIVE, NO_ACTIVE.


