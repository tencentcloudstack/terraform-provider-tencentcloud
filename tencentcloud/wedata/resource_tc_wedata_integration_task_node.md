Provides a resource to create a wedata integration_task_node

Example Usage

```hcl
resource "tencentcloud_wedata_integration_task_node" "example" {
  project_id       = "1612982498218618880"
  task_id          = "20231022181114990"
  name             = "tf_example1"
  node_type        = "INPUT"
  data_source_type = "MYSQL"
  task_type        = 202
  task_mode        = 2
  node_info {
    datasource_id = "5085"
    config {
      name  = "Type"
      value = "MYSQL"
    }
    config {
      name  = "splitPk"
      value = "id"
    }
    config {
      name  = "PrimaryKey"
      value = "id"
    }
    config {
      name  = "isNew"
      value = "true"
    }
    config {
      name  = "PrimaryKey_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "splitPk_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "Database"
      value = "demo_mysql"
    }
    config {
      name  = "TableNames"
      value = "users"
    }
    config {
      name  = "SiblingNodes"
      value = "[]"
    }
    schema {
      id    = "471331072"
      name  = "id"
      type  = "INT"
      alias = "id"
    }
    schema {
      id    = "422052352"
      name  = "username"
      type  = "VARCHAR(50)"
      alias = "username"
    }
  }
}
```