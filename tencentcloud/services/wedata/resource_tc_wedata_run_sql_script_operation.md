Provides a resource to create a WeData run sql script operation

Example Usage

```hcl
resource "tencentcloud_wedata_run_sql_script_operation" "example" {
  script_id  = "195a5f49-8e43-4e05-8b42-cecdfb6349f8"
  project_id = "2983848457986924544"
}
```

Or

```hcl
resource "tencentcloud_wedata_run_sql_script_operation" "example" {
  script_id      = "195a5f49-8e43-4e05-8b42-cecdfb6349f8"
  project_id     = "2983848457986924544"
  script_content = "SHOW DATABASES;"
}
```
