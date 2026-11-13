---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_folder"
sidebar_current: "docs-tencentcloud-resource-oceanus_folder"
description: |-
  Provides a resource to create a oceanus folder
---

# tencentcloud_oceanus_folder

Provides a resource to create a oceanus folder

## Example Usage

```hcl
resource "tencentcloud_oceanus_folder" "example" {
  folder_name   = "tf_example"
  parent_id     = "folder-lfqkt11s"
  folder_type   = 0
  work_space_id = "space-125703345ap-shenzhen-fsi"
}
```

## Argument Reference

The following arguments are supported:

* `folder_name` - (Required, String) New file name.
* `parent_id` - (Required, String) Parent folder id.
* `work_space_id` - (Required, String) Workspace SerialId.
* `folder_type` - (Optional, Int) Folder type, 0: job folder, 1: resource folder. Default is 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateFolder](https://www.tencentcloud.com/search/?wd=CreateFolder&type=document&lang=en)
* Read: [DescribeFolder](https://www.tencentcloud.com/search/?wd=DescribeFolder&type=document&lang=en)
* Update: [ModifyFolder](https://www.tencentcloud.com/search/?wd=ModifyFolder&type=document&lang=en)
* Delete: [DeleteFolders](https://www.tencentcloud.com/search/?wd=DeleteFolders&type=document&lang=en)

## Import

oceanus folder can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_folder.example space-125703345ap-shenzhen-fsi#folder-f40fq79g#0
```

