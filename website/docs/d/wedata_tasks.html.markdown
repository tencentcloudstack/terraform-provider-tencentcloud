---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_tasks"
sidebar_current: "docs-tencentcloud-datasource-wedata_tasks"
description: |-
  Use this data source to query detailed information of wedata wedata_tasks
---

# tencentcloud_wedata_tasks

Use this data source to query detailed information of wedata wedata_tasks

## Example Usage

```hcl
data "tencentcloud_wedata_tasks" "wedata_tasks" {
  project_id = 2905622749543821312
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `bundle_id` - (Optional, String) Bundle id.
* `create_time` - (Optional, Set: [`String`]) Creation time range (yyyy-MM-dd HH:MM:ss). Two time values must be provided in the array.
* `create_user_uin` - (Optional, String) Creator ID.
* `modify_time` - (Optional, Set: [`String`]) Modification time range (yyyy-MM-dd HH:mm:ss). Two time values must be provided in the array.
* `owner_uin` - (Optional, String) Owner ID.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, String) Task Status:
* N: New
* Y: Scheduling
* F: Offline
* O: Paused
* T: Offlining
* INVALID: Invalid.
* `submit` - (Optional, Bool) Submission status.
* `task_name` - (Optional, String) Task name.
* `task_type_id` - (Optional, Int) Task type.
* `workflow_id` - (Optional, String) Workflow ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Describes the task pagination information.


