Provides a resource to create a WeData code folder

Example Usage

```hcl
resource "tencentcloud_wedata_code_folder" "example" {
  project_id         = "2983848457986924544"
  folder_name        = "tf_example"
  parent_folder_path = "/"
}
```

Import

WeData code folder can be imported using the projectId#folderId, e.g.

```
terraform import tencentcloud_wedata_code_folder.example 1470547050521227264#2ee111df-5573-4ac4-9f93-cf9e8e438d80
```
