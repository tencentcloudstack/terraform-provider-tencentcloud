---
subcategory: "WeData"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_rule_templates"
sidebar_current: "docs-tencentcloud-datasource-wedata_rule_templates"
description: |-
  Use this data source to query detailed information of wedata rule templates
---

# tencentcloud_wedata_rule_templates

Use this data source to query detailed information of wedata rule templates

## Example Usage

```hcl
data "tencentcloud_wedata_rule_templates" "rule_templates" {
  type                = 2
  source_object_type  = 2
  project_id          = "1840731346428280832"
  source_engine_types = [2, 4, 16]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `source_engine_types` - (Optional, Set: [`Int`]) Applicable type of source data.
* `source_object_type` - (Optional, Int) Source data object type. `1`: Constant, `2`: Offline table level, `3`: Offline field level.
* `type` - (Optional, Int) Template type. `1` means System template, `2` means Custom template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - rule template list.
  * `citation_count` - Citations.
  * `compare_type` - The type of comparison method supported by the rule (1: fixed value comparison, greater than, less than, greater than or equal to, etc. 2: fluctuating value comparison, absolute value, rise, fall).
  * `description` - Description of rule template.
  * `multi_source_flag` - Whether to associate other library tables.
  * `name` - Name of rule template.
  * `quality_dim` - Quality inspection dimensions. `1`: Accuracy, `2`: Uniqueness, `3`: Completeness, `4`: Consistency, `5`: Timeliness, `6`: Effectiveness.
  * `rule_template_id` - ID of rule template.
  * `source_content` - Content of rule template.
  * `source_engine_types` - Applicable type of source data.
  * `source_object_data_type` - Source data object type. `1`: value, `2`: string.
  * `source_object_type` - Source object type. `1`: Constant, `2`: Offline table level, `3`: Offline field level.
  * `sql_expression` - Sql Expression.
  * `sub_quality_dim` - Sub Quality inspection dimensions. `1`: Accuracy, `2`: Uniqueness, `3`: Completeness, `4`: Consistency, `5`: Timeliness, `6`: Effectiveness.
  * `type` - Template type. `1` means System template, `2` means Custom template.
  * `update_time` - update time, like: yyyy-MM-dd HH:mm:ss.
  * `user_id` - user id.
  * `user_name` - user name.
  * `where_flag` - If add where.


