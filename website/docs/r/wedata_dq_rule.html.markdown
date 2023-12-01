---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_dq_rule"
sidebar_current: "docs-tencentcloud-resource-wedata_dq_rule"
description: |-
  Provides a resource to create a wedata dq_rule
---

# tencentcloud_wedata_dq_rule

Provides a resource to create a wedata dq_rule

## Example Usage

```hcl
resource "tencentcloud_wedata_dq_rule" "example" {
  project_id                   = "1948767646355341312"
  rule_group_id                = 312
  rule_template_id             = 1
  name                         = "tf_example"
  table_id                     = "N85hbsh5QQ2VLHL2iOUVeQ"
  type                         = 1
  source_object_data_type_name = "table"
  source_object_value          = "表"
  condition_type               = 1
  compare_rule {
    items {
      compare_type = 1
      operator     = "=="
      value_list {
        value_type = 3
        value      = "100"
      }
    }
  }
  alarm_level = 1
  description = "description."
}
```

## Argument Reference

The following arguments are supported:

* `alarm_level` - (Required, Int) Alarm trigger levels 1. Low, 2. Medium, 3. High.
* `compare_rule` - (Required, List) Alarm trigger condition.
* `condition_type` - (Required, Int) Detection scope 1. Full Table 2. Conditional scan.
* `name` - (Required, String) Rule name.
* `project_id` - (Required, String) Project id.
* `rule_template_id` - (Required, Int) Rule template id.
* `source_object_data_type_name` - (Required, String) Source field type. int, string.
* `source_object_value` - (Required, String) Source field name.
* `type` - (Required, Int) Rule Type 1. System Template, 2. Custom Template, 3. Custom SQL.
* `condition_expression` - (Optional, String) Condition scans WHERE condition expressions.
* `custom_sql` - (Optional, String) Custom sql.
* `description` - (Optional, String) Rule description.
* `field_config` - (Optional, List) Custom template sql expression field replacement parameters.
* `quality_dim` - (Optional, Int) Rules belong to quality dimensions (1. accuracy, 2. uniqueness, 3. completeness, 4. consistency, 5. timeliness, 6. effectiveness).
* `rel_condition_expr` - (Optional, String) The source field and the target field are associated with a conditional on expression.
* `rule_group_id` - (Optional, Int) Rule group id.
* `source_engine_types` - (Optional, Set: [`Int`]) List of execution engines supported by this rule.
* `table_id` - (Optional, String) Table id.
* `target_condition_expr` - (Optional, String) Target filter condition expression.
* `target_database_id` - (Optional, String) Target database id.
* `target_object_value` - (Optional, String) Target field name  CITY.
* `target_table_id` - (Optional, String) Target table id.

The `compare_rule` object supports the following:

* `cycle_step` - (Optional, Int) Periodic Indicates the default period of a template, in secondsNote: This field may return null, indicating that a valid value cannot be obtained.
* `items` - (Optional, List) Comparison condition listNote: This field may return null, indicating that a valid value cannot be obtained.

The `field_config` object supports the following:

* `field_data_type` - (Optional, String) Field typeNote: This field may return null, indicating that a valid value cannot be obtained.
* `field_key` - (Optional, String) Field keyNote: This field may return null, indicating that a valid value cannot be obtained.
* `field_value` - (Optional, String) Field valueNote: This field may return null, indicating that a valid value cannot be obtained.

The `field_config` object supports the following:

* `table_config` - (Optional, List) Library table variableNote: This field may return null, indicating that a valid value cannot be obtained.
* `where_config` - (Optional, List) Where variableNote: This field may return null, indicating that a valid value cannot be obtained.

The `items` object supports the following:

* `compare_type` - (Optional, Int) Comparison type 1. Fixed value 2. Fluctuating value 3. Comparison of value range 4. Enumeration range comparison 5. Do not compareNote: This field may return null, indicating that a valid value cannot be obtained.
* `operator` - (Optional, String) Comparison operation type &amp;lt; &amp;lt;= == =&amp;gt; &amp;gt;Note: This field may return null, indicating that a valid value cannot be obtained.
* `value_compute_type` - (Optional, Int) Quality statistics Type 1. Absolute value 2. Increase 3. Decrease 4. C contains 5. N C does not containNote: This field may return null, indicating that a valid value cannot be obtained.
* `value_list` - (Optional, List) Compare the threshold listNote: This field may return null, indicating that a valid value cannot be obtained.

The `table_config` object supports the following:

* `database_id` - (Optional, String) Database idNote: This field may return null, indicating that a valid value cannot be obtained.
* `database_name` - (Optional, String) Database nameNote: This field may return null, indicating that a valid value cannot be obtained.
* `field_config` - (Optional, List) Field variableNote: This field may return null, indicating that a valid value cannot be obtained.
* `table_id` - (Optional, String) Table idNote: This field may return null, indicating that a valid value cannot be obtained.
* `table_key` - (Optional, String) Table keyNote: This field may return null, indicating that a valid value cannot be obtained.
* `table_name` - (Optional, String) Table nameNote: This field may return null, indicating that a valid value cannot be obtained.

The `value_list` object supports the following:

* `value_type` - (Optional, Int) Threshold type 1. Low threshold 2. High threshold 3. Common threshold 4. Enumerated valueNote: This field may return null, indicating that a valid value cannot be obtained.
* `value` - (Optional, String) Threshold valueNote: This field may return null, indicating that a valid value cannot be obtained.

The `where_config` object supports the following:

* `field_data_type` - (Optional, String) Field typeNote: This field may return null, indicating that a valid value cannot be obtained.
* `field_key` - (Optional, String) Field keyNote: This field may return null, indicating that a valid value cannot be obtained.
* `field_value` - (Optional, String) Field valueNote: This field may return null, indicating that a valid value cannot be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

wedata dq_rule can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_dq_rule.example 1948767646355341312#894
```

