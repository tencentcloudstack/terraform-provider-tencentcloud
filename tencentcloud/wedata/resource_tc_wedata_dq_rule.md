Provides a resource to create a wedata dq_rule

Example Usage

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

Import

wedata dq_rule can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_dq_rule.example 1948767646355341312#894
```