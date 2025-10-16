Provides a resource to create a wedata wedata_workflow_folder

Example Usage

```hcl
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = 2905622749543821312
  parent_folder_path = "/"
  folder_name        = "test"
}
```
