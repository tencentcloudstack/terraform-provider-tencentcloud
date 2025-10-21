---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_tree_jobs"
sidebar_current: "docs-tencentcloud-datasource-oceanus_tree_jobs"
description: |-
  Use this data source to query detailed information of oceanus tree_jobs
---

# tencentcloud_oceanus_tree_jobs

Use this data source to query detailed information of oceanus tree_jobs

## Example Usage

```hcl
data "tencentcloud_oceanus_tree_jobs" "example" {
  work_space_id = "space-2idq8wbr"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter rules.
* `result_output_file` - (Optional, String) Used to save results.
* `work_space_id` - (Optional, String) Workspace SerialId.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered. Can only be set `Zone` or `JobType` or `JobStatus`.
* `values` - (Required, Set) Filter values for the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tree_info` - Tree structure information.
  * `children` - Subdirectory Information.
  * `id` - ID.
  * `job_set` - List of jobs.
    * `job_id` - Job ID.
    * `job_type` - Job Type.
    * `name` - Job Name.
    * `running_cu` - Resources occupied by homework.
    * `status` - Job status.
  * `name` - Name.
  * `parent_id` - Parent Id.


