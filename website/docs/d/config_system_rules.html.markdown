---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_system_rules"
sidebar_current: "docs-tencentcloud-datasource-config_system_rules"
description: |-
  Use this data source to query detailed information of Config system preset rules.
---

# tencentcloud_config_system_rules

Use this data source to query detailed information of Config system preset rules.

## Example Usage

### Query all system preset rules

```hcl
data "tencentcloud_config_system_rules" "example" {}
```

### Query system rules by keyword

```hcl
data "tencentcloud_config_system_rules" "example" {
  keyword = "cam"
}
```

### Query system rules by risk level

```hcl
data "tencentcloud_config_system_rules" "example" {
  risk_level = 1
}
```

## Argument Reference

The following arguments are supported:

* `keyword` - (Optional, String) Search keyword. Supports identifier/name/label/description search.
* `result_output_file` - (Optional, String) Used to save results.
* `risk_level` - (Optional, Int) Risk level for filtering. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rule_list` - System preset rule list.
  * `create_time` - Creation time.
  * `description` - Rule description.
  * `identifier_type` - Rule type.
  * `identifier` - Rule unique identifier.
  * `label` - Rule label list.
  * `reference_count` - Number of times this rule is referenced.
  * `resource_type` - Supported resource type list.
  * `risk_level` - Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).
  * `rule_name` - Rule name.
  * `service_function` - Corresponding service function.
  * `trigger_type` - Trigger type list.
  * `update_time` - Last update time.


