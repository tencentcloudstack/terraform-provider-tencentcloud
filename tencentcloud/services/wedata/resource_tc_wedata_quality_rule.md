Provides a resource to create a wedata quality_rule

Example Usage

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