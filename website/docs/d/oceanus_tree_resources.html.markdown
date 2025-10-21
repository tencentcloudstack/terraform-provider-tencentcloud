---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_tree_resources"
sidebar_current: "docs-tencentcloud-datasource-oceanus_tree_resources"
description: |-
  Use this data source to query detailed information of oceanus tree_resources
---

# tencentcloud_oceanus_tree_resources

Use this data source to query detailed information of oceanus tree_resources

## Example Usage

```hcl
data "tencentcloud_oceanus_tree_resources" "example" {
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `work_space_id` - (Required, String) Workspace SerialId.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tree_info` - Tree structure information.
  * `children` - Subdirectory Information.
  * `id` - ID.
  * `items` - List of items.
    * `file_name` - File name.
    * `folder_id` - Folder id.
    * `name` - Name.
    * `ref_job_status_count_set` - Counting the number of associated tasks by state.
      * `count` - Job count.
      * `job_status` - Job status.
    * `remark` - Remark.
    * `resource_id` - Resource Id.
    * `resource_type` - Resource Type.
  * `name` - Name.
  * `parent_id` - Parent Id.


