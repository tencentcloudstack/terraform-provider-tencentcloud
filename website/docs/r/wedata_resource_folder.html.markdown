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

* Create: [CreateResourceFolder](https://cloud.tencent.com/document/api/1267/123641) - 创建资源文件夹
* Read: [ListResourceFolders](https://cloud.tencent.com/document/api/1267/123625) - 查看资源文件夹列表
* Update: [UpdateResourceFolder](https://cloud.tencent.com/document/api/1267/123617) - 更新资源文件夹
* Delete: [DeleteResourceFolder](https://cloud.tencent.com/document/api/1267/123636) - 删除资源文件夹

