---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_dq_rule_template"
sidebar_current: "docs-tencentcloud-resource-wedata_dq_rule_template"
description: |-
  Provides a resource to create a wedata dq_rule_template
---

# tencentcloud_wedata_dq_rule_template

Provides a resource to create a wedata dq_rule_template

## Example Usage

```hcl
resource "tencentcloud_wedata_dq_rule_template" "example" {
  type                = 2
  name                = "tf_example"
  quality_dim         = 1
  source_object_type  = 2
  description         = "description."
  source_engine_types = [2]
  multi_source_flag   = true
  sql_expression      = "c2VsZWN0"
  project_id          = "1948767646355341312"
  where_flag          = true
}
```

## Argument Reference

The following arguments are supported:

* `multi_source_flag` - (Required, Bool) Whether to associate other tables.
* `name` - (Required, String) Template name.
* `project_id` - (Required, String) Project id.
* `quality_dim` - (Required, Int) Quality detection dimension 1. Accuracy 2. Uniqueness 3. Completeness 4. Consistency 5. Timeliness 6. effectiveness.
* `source_engine_types` - (Required, Set: [`Int`]) Type of the engine on the source end.
* `source_object_type` - (Required, Int) Source end data object type 1. Constant 2. Offline table level 2. Offline field level.
* `sql_expression` - (Required, String) SQL expression.
* `type` - (Required, Int) Template Type 1. System Template 2. User-defined template.
* `where_flag` - (Required, Bool) Add where parameter or not.
* `description` - (Optional, String) Template description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `template_id` - Template ID.


## Import

wedata dq_rule_template can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_dq_rule_template.example 1948767646355341312#9480
```

