---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_sql_folder"
sidebar_current: "docs-tencentcloud-resource-wedata_sql_folder"
description: |-
  Provides a resource to create a WeData sql folder
---

# tencentcloud_wedata_sql_folder

Provides a resource to create a WeData sql folder

## Example Usage

```hcl
resource "tencentcloud_wedata_sql_folder" "example" {
  folder_name        = "tf_example"
  project_id         = "2983848457986924544"
  parent_folder_path = "/"
  access_scope       = "SHARED"
}
```

## Argument Reference

The following arguments are supported:

* `folder_name` - (Required, String) Folder name.
* `parent_folder_path` - (Required, String, ForceNew) The parent folder path is /aaa/bbb/ccc. The path header must have a slash. To query the root directory, pass /.
* `project_id` - (Required, String, ForceNew) Project ID.
* `access_scope` - (Optional, String) Permission range: SHARED, PRIVATE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `folder_id` - Folder ID.
* `path` - Node path.


## Import

WeData sql folder can be imported using the projectId#folderId, e.g.

```
terraform import tencentcloud_wedata_sql_folder.example 2917455276892352512#1c9db971-58c6-43b4-93a0-be526123a1d8
```

