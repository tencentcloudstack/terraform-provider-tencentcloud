---
subcategory: "WeData"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_rule_template"
sidebar_current: "docs-tencentcloud-resource-wedata_rule_template"
description: |-
  Provides a resource to create a wedata rule_template
---

# tencentcloud_wedata_rule_template

Provides a resource to create a wedata rule_template

## Example Usage

```hcl
resource "tencentcloud_wedata_rule_template" "rule_template" {
  project_id          = "1840731346428280832"
  type                = 2
  name                = "tf-test"
  quality_dim         = 3
  source_object_type  = 2
  description         = "for tf test"
  source_engine_types = [2, 4, 16]
  multi_source_flag   = false
  sql_expression      = base64encode("select * from db")
  where_flag          = false
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `description` - (Optional, String) Description of Template.
* `multi_source_flag` - (Optional, Bool) Whether to associate other library tables.
* `name` - (Optional, String) Template name.
* `quality_dim` - (Optional, Int) Quality inspection dimensions. `1`: Accuracy, `2`: Uniqueness, `3`: Completeness, `4`: Consistency, `5`: Timeliness, `6`: Effectiveness.
* `source_engine_types` - (Optional, Set: [`Int`]) The engine type corresponding to the source. `2`: hive,`4`: spark, `16`: dlc.
* `source_object_type` - (Optional, Int) Source data object type. `1`: Constant, `2`: Offline table level, `3`: Offline field level.
* `sql_expression` - (Optional, String) SQL Expression.
* `type` - (Optional, Int) Template type. `1` means System template, `2` means Custom template.
* `where_flag` - (Optional, Bool) If add where.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

wedata rule_template can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_rule_template.rule_template rule_template_id
```

