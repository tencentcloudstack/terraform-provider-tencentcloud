---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_workflow_folders"
sidebar_current: "docs-tencentcloud-datasource-wedata_workflow_folders"
description: |-
  Use this data source to query detailed information of wedata wedata_workflow_folders
---

# tencentcloud_wedata_workflow_folders

Use this data source to query detailed information of wedata wedata_workflow_folders

## Example Usage

```hcl
data "tencentcloud_wedata_workflow_folders" "wedata_workflow_folders" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
}
```

## Argument Reference

The following arguments are supported:

* `parent_folder_path` - (Required, String) Parent folder absolute path, for example /abc/de, if it is root directory, pass /.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Paginated folder query result.


