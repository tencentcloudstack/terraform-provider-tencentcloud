---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_quality_rule"
sidebar_current: "docs-tencentcloud-resource-wedata_quality_rule"
description: |-
  Provides a resource to create a wedata quality_rule
---

# tencentcloud_wedata_quality_rule

Provides a resource to create a wedata quality_rule

## Example Usage

```hcl
resource "tencentcloud_wedata_quality_rule" "rule" {
  alarm_level                  = 1
  condition_type               = 1
  create_rule_scene            = 1
  database_name                = "default"
  datasource_id                = 65253
  description                  = "tf test rule1"
  name                         = "at_src_mysql2hive_prod_cq_makeup_09db_1_di_表行数_tf_test"
  project_id                   = 3016337760439783424
  quality_dim                  = 1
  rule_group_id                = 949
  rule_template_id             = 6142
  source_engine_types          = [2, 4, 16, 64, 128, 256, 512, 1024]
  source_object_data_type_name = "table"
  source_object_value          = "表"
  table_id                     = 175
  table_name                   = "at_src_mysql2hive_prod_cq_makeup_09db_1_di"
  type                         = 1
  compare_rule {
    compute_expression = "0o1o2"
    cycle_step         = 0
    items {
      compare_type       = 1
      operator           = ">"
      value_compute_type = 0
      value_list {
        value      = 100
        value_type = 3
      }
    }
    items {
      compare_type       = 1
      operator           = "<"
      value_compute_type = 0
      value_list {
        value      = 201
        value_type = 3
      }
    }
    items {
      compare_type       = 1
      operator           = "=="
      value_compute_type = 0
      value_list {
        value      = 102
        value_type = 3
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_level` - (Required, Int) Alarm trigger level. Valid values: `1` (low), `2` (medium), `3` (high).
* `compare_rule` - (Required, List) Alarm trigger condition.
* `create_rule_scene` - (Required, Int) Rule creation scene. Valid values: `1` (single table multiple rules). Other business scenarios are not currently supported.
* `database_name` - (Required, String) Database name.
* `datasource_id` - (Required, String) Data source ID.
* `name` - (Required, String) Rule name.
* `project_id` - (Required, String) Project ID.
* `source_engine_types` - (Required, Set: [`Int`]) Supported execution engine list for this rule. Valid values: `1` (MYSQL), `2` (HIVE), `4` (SPARK), `8` (LIVY), `16` (DLC), `32` (GBASE), `64` (TCHouse-P), `128` (DORIS), `256` (TCHouse-D), `512` (EMR_STARROCKS), `1024` (TCHouse-X).
* `type` - (Required, Int) Rule type. Valid values: `1` (system template), `2` (custom template), `3` (custom SQL).
* `catalog_name` - (Optional, String) Data catalog name, mainly used for DLC data source.
* `condition_expression` - (Optional, String) Conditional scan WHERE condition expression. Required when ConditionType=2 (conditional scan).
* `condition_type` - (Optional, Int) Detection range. Required when Type=1 (system template) or 2 (custom template). Valid values: `1` (full table), `2` (conditional scan). Note: When CompareType is 2 (fluctuation value) or using user-defined template containing filter condition ${FILTER}, detection range must be 2 (conditional scan).
* `custom_sql` - (Optional, String) Custom SQL (Base64 encoded). Required when Type=3 (custom SQL).
* `database_id` - (Optional, String) Database ID.
* `description` - (Optional, String) Rule description.
* `field_config` - (Optional, List) Custom template SQL expression field replacement parameters. Required when Type=2 (custom template).
* `index` - (Optional, String) Index to distinguish different data when adding.
* `quality_dim` - (Optional, Int) Quality dimension of the rule. Required when Type=3 (custom SQL). Valid values: `1` (accuracy), `2` (uniqueness), `3` (completeness), `4` (consistency), `5` (timeliness), `6` (validity).
* `rel_condition_expr` - (Optional, String) Source field and target field association condition ON expression. Required only for field data correlation rules (ruleTemplate qualityDim=4 (consistency) and subQualityDim=3 (field data correlation)). Example: sourceTable.model_id=targetTable.model_id.
* `rule_group_id` - (Optional, Int) Rule group ID.
* `rule_template_id` - (Optional, Int) Rule template ID. Required when Type is not equal to 3 (custom SQL).
* `schema_name` - (Optional, String) Schema name.
* `source_object_data_type_name` - (Optional, String) Source data object (table, field, etc.) detailed type. Required when Type=1 (system template). For table corresponds to fixed value `table` (template is table-level). For field corresponds to field type: int, string, etc. (template is field-level).
* `source_object_value` - (Optional, String) Source data object (table, field, etc.) name. Required when Type=1 (system template).
* `table_id` - (Optional, String) Data table ID. Either TableId or TableName must be provided.
* `table_name` - (Optional, String) Table name. Either TableId or TableName must be provided.
* `target_catalog_name` - (Optional, String) Target data catalog name. Required only for system template field data correlation rules and when data source is DLC (ruleTemplate qualityDim=4 and subQualityDim=3). Used for cross-table data validation and association.
* `target_condition_expr` - (Optional, String) Target filter condition expression.
* `target_database_id` - (Optional, String) Target database ID.
* `target_database_name` - (Optional, String) Target database name. Required only for system template field data correlation rules (ruleTemplate qualityDim=4 and subQualityDim=3). Used for cross-table data validation and association.
* `target_object_value` - (Optional, String) Target field name CITY.
* `target_schema_name` - (Optional, String) Target schema name. Required only for system template field data correlation rules and when data source is TCHouse-P (ruleTemplate qualityDim=4 and subQualityDim=3). Used for cross-table data validation and association.
* `target_table_id` - (Optional, String) Target table ID.
* `target_table_name` - (Optional, String) Target table name. Required only for system template field data correlation rules (ruleTemplate qualityDim=4 and subQualityDim=3). Used for cross-table data validation and association.
* `task_id` - (Optional, String) Task ID.

The `compare_rule` object supports the following:

* `items` - (Required, List) Comparison condition list.
* `compute_expression` - (Optional, String) `o` represents OR, `a` represents AND, numbers represent items index.
* `cycle_step` - (Optional, Int) Periodic template default cycle in seconds.

The `field_config` object of `table_config` supports the following:

* `field_data_type` - (Optional, String) Field data type.
* `field_key` - (Optional, String) Field key.
* `field_value` - (Optional, String) Field value.
* `value_config` - (Optional, List) Field value variable information.

The `field_config` object supports the following:

* `table_config` - (Optional, List) Database and table variables.
* `where_config` - (Optional, List) WHERE variables.

The `items` object of `compare_rule` supports the following:

* `compare_type` - (Optional, Int) Comparison type (required). Valid values: `1` (fixed value), `2` (fluctuation value), `3` (numerical range comparison), `4` (enumeration range comparison), `5` (no comparison), `6` (field data correlation), `7` (fairness).
* `operator` - (Optional, String) Comparison operator type (conditionally required). Required when CompareType belongs to {1,2,6,7}. Valid values: `<`, `<=`, `==`, `=>`, `>`, `!=`, `IRLCRO` (within interval, left closed right open), `IRLORC` (within interval, left open right closed), `IRLCRC` (within interval, left closed right closed), `IRLORO` (within interval, left open right open), `NRLCRO` (not within interval, left closed right open), `NRLORC` (not within interval, left open right closed), `NRLCRC` (not within interval, left closed right closed), `NRLORO` (not within interval, left open right open).
* `value_compute_type` - (Optional, Int) Quality statistics value type (conditionally required). Required when CompareType belongs to {2,3,7}. When compareType = 2 (fluctuation value): `1` = absolute value (ABS), `2` = ascending (ASCEND), `3` = descending (DESCEND). When compareType = 3 (numerical range): `4` = within range (WITH_IN_RANGE), `5` = out of range (OUT_OF_RANGE). When compareType = 7 (fairness): `6` = fairness rate (FAIRNESS_RATE), `7` = fairness gap (FAIRNESS_GAP).
* `value_list` - (Optional, List) Comparison threshold list (required).

The `table_config` object of `field_config` supports the following:

* `database_id` - (Optional, String) Database ID.
* `database_name` - (Optional, String) Database name.
* `field_config` - (Optional, List) Field variables.
* `table_id` - (Optional, String) Table ID.
* `table_key` - (Optional, String) Table key.
* `table_name` - (Optional, String) Table name.

The `value_config` object of `field_config` supports the following:

* `field_data_type` - (Optional, String) Field data type.
* `field_key` - (Optional, String) Field value key.
* `field_value` - (Optional, String) Field value.

The `value_config` object of `where_config` supports the following:

* `field_data_type` - (Optional, String) Field data type.
* `field_key` - (Optional, String) Field value key.
* `field_value` - (Optional, String) Field value.

The `value_list` object of `items` supports the following:

* `value_type` - (Optional, Int) Threshold type (required). Valid values: `1` (low threshold), `2` (high threshold), `3` (normal threshold), `4` (enumeration value).
* `value` - (Optional, String) Threshold value (required).

The `where_config` object of `field_config` supports the following:

* `field_data_type` - (Optional, String) Field data type.
* `field_key` - (Optional, String) Field key.
* `field_value` - (Optional, String) Field value.
* `value_config` - (Optional, List) Field value variable information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


