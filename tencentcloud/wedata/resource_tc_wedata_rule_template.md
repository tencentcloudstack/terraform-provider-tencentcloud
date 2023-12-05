Provides a resource to create a wedata rule_template

Example Usage

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

Import

wedata rule_template can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_rule_template.rule_template rule_template_id
```