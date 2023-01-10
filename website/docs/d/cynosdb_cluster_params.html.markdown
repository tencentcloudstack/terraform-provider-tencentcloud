---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_params"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_cluster_params"
description: |-
  Use this data source to query detailed information of cynosdb cluster_params
---

# tencentcloud_cynosdb_cluster_params

Use this data source to query detailed information of cynosdb cluster_params

## Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_params" "cluster_params" {
  cluster_id = "cynosdbmysql-bws8h88b"
  param_name = "innodb_checksum_algorithm"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The ID of cluster.
* `param_name` - (Optional, String) Parameter name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Instance parameter list. Note: This field may return null, indicating that no valid value can be obtained.
  * `current_value` - Current value.
  * `default` - Default value.
  * `description` - The description of parameter.
  * `enum_value` - When the parameter is enum/string/bool, the optional value list.Note: This field may return null, indicating that no valid value can be obtained.
  * `func` - Function.Note: This field may return null, indicating that no valid value can be obtained.
  * `is_func` - Is it a function.Note: This field may return null, indicating that no valid value can be obtained.
  * `is_global` - Is it a global parameter.Note: This field may return null, indicating that no valid value can be obtained.
  * `match_type` - Matching type, multiVal, regex is used when the parameter type is string.
  * `match_value` - Match the target value, when multiVal, each key is divided by `;`.
  * `max` - The maximum value when the parameter type is float/integer.
  * `min` - The minimum value when the parameter type is float/integer.
  * `need_reboot` - Whether to reboot.
  * `param_name` - The name of parameter.
  * `param_type` - Parameter type: integer/float/string/enum/bool.


