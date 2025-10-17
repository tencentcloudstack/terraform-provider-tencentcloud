---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_code_file"
sidebar_current: "docs-tencentcloud-resource-wedata_code_file"
description: |-
  Provides a resource to create a WeData code file
---

# tencentcloud_wedata_code_file

Provides a resource to create a WeData code file

## Example Usage

```hcl
resource "tencentcloud_wedata_code_folder" "example" {
  project_id         = "2983848457986924544"
  folder_name        = "tf_example"
  parent_folder_path = "/"
}

resource "tencentcloud_wedata_code_file" "example" {
  project_id         = "2983848457986924544"
  code_file_name     = "tf_example_code_file"
  parent_folder_path = tencentcloud_wedata_code_folder.example.path
  code_file_content  = "Hello Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `code_file_name` - (Required, String, ForceNew) Code file name.
* `parent_folder_path` - (Required, String, ForceNew) Parent folder path, for example /aaa/bbb/ccc, path header must start with a slash, root directory pass /.
* `project_id` - (Required, String, ForceNew) Project ID.
* `code_file_config` - (Optional, List) Code file configuration.
* `code_file_content` - (Optional, String) Code file content.

The `code_file_config` object supports the following:

* `notebook_session_info` - (Optional, List) Notebook kernel session information.
* `params` - (Optional, String) Advanced runtime parameters, variable substitution, map-json String,String.

The `notebook_session_info` object of `code_file_config` supports the following:

* `notebook_session_id` - (Optional, String) Session ID.
* `notebook_session_name` - (Optional, String) Session name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_scope` - Permission range: SHARED, PRIVATE.
* `code_file_id` - Code file ID.
* `path` - The full path of the node, /aaa/bbb/ccc.ipynb, consists of the names of each node.


## Import

WeData code file can be imported using the projectId#codeFileId, e.g.

```
terraform import tencentcloud_wedata_code_file.example 1470547050521227264#2bfa8813-344f-4858-a2cc-7a07bd10ac1d
```

