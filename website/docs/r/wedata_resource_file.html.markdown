---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource_file"
sidebar_current: "docs-tencentcloud-resource-wedata_resource_file"
description: |-
  Provides a resource to create a wedata wedata_resource_file
---

# tencentcloud_wedata_resource_file

Provides a resource to create a wedata wedata_resource_file

## Example Usage

```hcl
resource "tencentcloud_wedata_resource_folder" "wedata_resource_folder" {
  project_id         = 2905622749543821312
  parent_folder_path = "/"
  folder_name        = "folder"
}

resource "tencentcloud_wedata_resource_file" "wedata_resource_file" {
  project_id         = 2905622749543821312
  resource_name      = "tftest.txt"
  bucket_name        = "data-manage-fsi-1315051789"
  cos_region         = "ap-beijing-fsi"
  parent_folder_path = "${tencentcloud_wedata_resource_folder.wedata_resource_folder.parent_folder_path}${tencentcloud_wedata_resource_folder.wedata_resource_folder.folder_name}"
  resource_file      = "/datastudio/resource/2905622749543821312/${tencentcloud_wedata_resource_folder.wedata_resource_folder.parent_folder_path}${tencentcloud_wedata_resource_folder.wedata_resource_folder.folder_name}/test"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, String) cos bucket name, which can be obtained from the GetResourceCosPath interface.
* `cos_region` - (Required, String) The cos bucket area corresponding to the BucketName bucket.
* `parent_folder_path` - (Required, String) The path to upload resource files in the project, example value: /wedata/qxxxm/, root directory, please use/.
* `project_id` - (Required, String, ForceNew) Project id.
* `resource_file` - (Required, String) - You can only choose one of the two methods of uploading a file and manually filling. If both are provided, the order of values is file> manual filling value
-the manual filling value must be the existing cos path, /datastudio/resource/is a fixed prefix, projectId is the project ID, and a specific value needs to be passed in, parentFolderPath is the parent folder path, name is the file name, and examples of manual filling value values are: /datastudio/resource/projectId/parentFolderPath/name 
.
* `resource_name` - (Required, String) The resource file name should be consistent with the uploaded file name as much as possible.
* `bundle_id` - (Optional, String) bundle client ID.
* `bundle_info` - (Optional, String) bundle client information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



