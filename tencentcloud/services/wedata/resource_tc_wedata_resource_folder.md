Provides a resource to create a wedata wedata_resource_folder

Example Usage

```hcl
resource "tencentcloud_wedata_resource_folder" "wedata_resource_folder" {
  project_id         = 2905622749543821312
  parent_folder_path = "/"
  folder_name        = "folder"
}
```

Import

wedata wedata_resource_folder can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_resource_folder.wedata_resource_folder wedata_resource_folder_id
```
