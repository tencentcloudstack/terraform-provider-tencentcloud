---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_code_folder"
sidebar_current: "docs-tencentcloud-resource-wedata_code_folder"
description: |-
  Provides a resource to create a WeData code folder
---

# tencentcloud_wedata_code_folder

Provides a resource to create a WeData code folder

## Example Usage

```hcl
resource "tencentcloud_wedata_code_folder" "example" {
  project_id         = "2983848457986924544"
  folder_name        = "tf_example"
  parent_folder_path = "/"
}
```

## Argument Reference

The following arguments are supported:

* `folder_name` - (Required, String) Folder name.
* `parent_folder_path` - (Required, String, ForceNew) Parent folder path, for example /aaa/bbb/ccc, path header must start with a slash, root directory pass /.
* `project_id` - (Required, String, ForceNew) Project ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_scope` - Permission range: SHARED, PRIVATE.
* `folder_id` - Folder ID.
* `path` - Node path.
* `type` - Type. folder, script.


## Import

WeData code folder can be imported using the projectId#folderId, e.g.

```
terraform import tencentcloud_wedata_code_folder.example 1470547050521227264#2ee111df-5573-4ac4-9f93-cf9e8e438d80
```

