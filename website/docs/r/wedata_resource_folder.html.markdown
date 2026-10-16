---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource_folder"
sidebar_current: "docs-tencentcloud-resource-wedata_resource_folder"
description: |-
  Provides a resource to create a wedata wedata_resource_folder
---

# tencentcloud_wedata_resource_folder

Provides a resource to create a wedata wedata_resource_folder

## Example Usage

```hcl
resource "tencentcloud_wedata_resource_folder" "wedata_resource_folder" {
  project_id         = 2905622749543821312
  parent_folder_path = "/"
  folder_name        = "folder"
}
```

## Argument Reference

The following arguments are supported:

* `folder_name` - (Required, String) Folder name.
* `parent_folder_path` - (Required, String) Absolute path of parent folder, value example/wedata/test, root directory, please use/.
* `project_id` - (Required, String, ForceNew) Project id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateResourceFolder](https://www.tencentcloud.com/search/?wd=CreateResourceFolder&type=document&lang=en)
* Read: [ListResourceFolders](https://www.tencentcloud.com/search/?wd=ListResourceFolders&type=document&lang=en)
* Update: [UpdateResourceFolder](https://www.tencentcloud.com/search/?wd=UpdateResourceFolder&type=document&lang=en)
* Delete: [DeleteResourceFolder](https://www.tencentcloud.com/search/?wd=DeleteResourceFolder&type=document&lang=en)

