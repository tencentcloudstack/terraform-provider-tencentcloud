Provides a resource to create a WeData sql folder

Example Usage

```hcl
resource "tencentcloud_wedata_sql_folder" "example" {
  folder_name        = "tf_example"
  project_id         = "2983848457986924544"
  parent_folder_path = "/"
  access_scope       = "SHARED"
}
```

Import

WeData sql folder can be imported using the projectId#folderId, e.g.

```
terraform import tencentcloud_wedata_sql_folder.example 2917455276892352512#1c9db971-58c6-43b4-93a0-be526123a1d8
```
