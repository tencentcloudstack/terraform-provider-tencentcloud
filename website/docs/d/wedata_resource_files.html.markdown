---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource_files"
sidebar_current: "docs-tencentcloud-datasource-wedata_resource_files"
description: |-
  Use this data source to query detailed information of wedata wedata_resource_files
---

# tencentcloud_wedata_resource_files

Use this data source to query detailed information of wedata wedata_resource_files

## Example Usage

```hcl
data "tencentcloud_wedata_resource_files" "wedata_resource_files" {
  project_id    = 2905622749543821312
  resource_name = "tftest.txt"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `create_time_end` - (Optional, String) Create time range. specifies the termination time in yyyy-MM-dd HH:MM:ss format.
* `create_time_start` - (Optional, String) Create time range. specifies the start time in yyyy-MM-dd HH:MM:ss format.
* `create_user_uin` - (Optional, String) Creator ID. obtain through the DescribeCurrentUserInfo API.
* `modify_time_end` - (Optional, String) Update time range. specifies the end time in yyyy-MM-dd HH:MM:ss format.
* `modify_time_start` - (Optional, String) Update time range. specifies the start time in yyyy-MM-dd HH:MM:ss format.
* `parent_folder_path` - (Optional, String) Specifies the path of the file's parent folder (for example /a/b/c, querying resource files under the folder c).
* `resource_name` - (Optional, String) Resource file name (fuzzy search keyword).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Retrieve the resource file list.


