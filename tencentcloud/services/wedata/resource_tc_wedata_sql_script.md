Provides a resource to create a WeData sql script

Example Usage

```hcl
resource "tencentcloud_wedata_sql_script" "example" {
  script_name        = "tf-example"
  project_id         = ""
  parent_folder_path = ""
  script_config {
    datasource_id     = ""
    datasource_env    = ""
    compute_resource  = ""
    executor_group_id = ""
    params            = ""
    advance_config    = ""
  }

  script_content = ""
  access_scope   = ""
}
```

Import

WeData sql script can be imported using the projectId#scriptId, e.g.

```
terraform import tencentcloud_wedata_sql_script.example 2917455276892352512#bf6a325f-ab82-4fba-9eac-1b6ae58f20f6
```
