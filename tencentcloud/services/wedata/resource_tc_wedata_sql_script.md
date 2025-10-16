Provides a resource to create a WeData sql script

Example Usage

```hcl
resource "tencentcloud_wedata_sql_script" "example" {
  script_name        = "tf_example_script"
  project_id         = "2983848457986924544"
  parent_folder_path = "/tf_example"
  script_config {
    datasource_id    = "108826"
    compute_resource = "svmgao_stability"
  }

  script_content = "SHOW DATABASES;"
  access_scope   = "SHARED"
}
```

Import

WeData sql script can be imported using the projectId#scriptId, e.g.

```
terraform import tencentcloud_wedata_sql_script.example 2983848457986924544#cccc3170-6334-46c3-adce-c5776ad2280c
```
