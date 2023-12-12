Provides a resource to create a wedata function

Example Usage

```hcl
resource "tencentcloud_wedata_function" "example" {
  type               = "HIVE"
  kind               = "ANALYSIS"
  name               = "tf_example"
  cluster_identifier = "emr-m6u3qgk0"
  db_name            = "tf_db_example"
  project_id         = "1612982498218618880"
  class_name         = "tf_class_example"
  resource_list {
    path = "/wedata-demo-1314991481/untitled3-1.0-SNAPSHOT.jar"
    name = "untitled3-1.0-SNAPSHOT.jar"
    id   = "5b28bcdf-a0e6-4022-927d-927d399c4593"
    type = "cos"
  }
  description = "description."
  usage       = "usage info."
  param_desc  = "param info."
  return_desc = "return value info."
  example     = "example info."
  comment     = "V1"
}
```