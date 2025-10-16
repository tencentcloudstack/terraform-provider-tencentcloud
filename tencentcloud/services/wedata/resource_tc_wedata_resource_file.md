Provides a resource to create a wedata wedata_resource_file

Example Usage

```hcl
resource "tencentcloud_wedata_resource_folder" "wedata_resource_folder" {
  project_id         = 2905622749543821312
  parent_folder_path = "/"
  folder_name        = "folder"
}

resource "tencentcloud_wedata_resource_file" "wedata_resource_file" {
    project_id = 2905622749543821312
    resource_name = "tftest.txt"
    bucket_name = "data-manage-fsi-1315051789"
    cos_region = "ap-beijing-fsi"
    parent_folder_path = "${tencentcloud_wedata_resource_folder.wedata_resource_folder.parent_folder_path}${tencentcloud_wedata_resource_folder.wedata_resource_folder.folder_name}"
    resource_file = "/datastudio/resource/2905622749543821312/${tencentcloud_wedata_resource_folder.wedata_resource_folder.parent_folder_path}${tencentcloud_wedata_resource_folder.wedata_resource_folder.folder_name}/test"
}
```
