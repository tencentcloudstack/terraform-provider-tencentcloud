---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_param_templates"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_param_templates"
description: |-
  Use this data source to query detailed information of cynosdb param_templates
---

# tencentcloud_cynosdb_param_templates

Use this data source to query detailed information of cynosdb param_templates

## Example Usage

```hcl
data "tencentcloud_cynosdb_param_templates" "param_templates" {
}
```

## Argument Reference

The following arguments are supported:

* `db_modes` - (Optional, Set: [`String`]) Database mode, optional values: NORMAL, SERVERLESS.
* `engine_types` - (Optional, Set: [`String`]) Engine types.
* `engine_versions` - (Optional, Set: [`String`]) Database engine version number.
* `limit` - (Optional, Int) Query limit.
* `offset` - (Optional, Int) Page offset.
* `order_by` - (Optional, String) The sort field for the returned results.
* `order_direction` - (Optional, String) Sort by (asc, desc).
* `products` - (Optional, Set: [`String`]) The product type corresponding to the query template.
* `result_output_file` - (Optional, String) Used to save results.
* `template_ids` - (Optional, Set: [`Int`]) The id list of templates.
* `template_names` - (Optional, Set: [`String`]) The name list of templates.
* `template_types` - (Optional, Set: [`String`]) Template types.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Parameter Template Information.
  * `db_mode` - Database mode, optional values: NORMAL, SERVERLESS.
  * `engine_version` - Engine version.
  * `id` - The ID of template.
  * `param_info_set` - Parameter template details.Note: This field may return null, indicating that no valid value can be obtained.
    * `current_value` - Current value.
    * `default` - Default value.
    * `description` - The description of parameter.
    * `enum_value` - An optional set of value types when the parameter type is enum.Note: This field may return null, indicating that no valid value can be obtained.
    * `max` - The maximum value when the parameter type is float/integer.Note: This field may return null, indicating that no valid value can be obtained.
    * `min` - The minimum value when the parameter type is float/integer.Note: This field may return null, indicating that no valid value can be obtained.
    * `need_reboot` - Whether to reboot.
    * `param_name` - The name of parameter.
    * `param_type` - Parameter type: integer/float/string/enum.
  * `template_description` - The description of template.
  * `template_name` - The name of template.


