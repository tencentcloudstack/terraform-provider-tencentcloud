---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_workflows"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_workflows"
description: |-
  Use this data source to query detailed information of wedata ops workflows
---

# tencentcloud_wedata_ops_workflows

Use this data source to query detailed information of wedata ops workflows

## Example Usage

```hcl
data "tencentcloud_wedata_ops_workflows" "wedata_ops_workflows" {
  project_id    = "2905622749543821312"
  folder_id     = "720ecbfb-7e5a-11f0-ba36-b8cef6a5af5c"
  status        = "ALL_RUNNING"
  owner_uin     = "100044349576"
  workflow_type = "Cycle"
  sort_type     = "ASC"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `create_time` - (Optional, String) Creation time, format yyyy-MM-dd HH:mm:ss.
* `create_user_uin` - (Optional, String) Creator ID.
* `folder_id` - (Optional, String) File ID.
* `key_word` - (Optional, String) Workflow keyword filter, supports fuzzy matching by workflow ID/name.
* `modify_time` - (Optional, String) Update time, format yyyy-MM-dd HH:mm:ss.
* `owner_uin` - (Optional, String) Responsible person ID.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_item` - (Optional, String) Sorting field, optional values: `CreateTime`, `TaskCount`.
* `sort_type` - (Optional, String) Sorting order, `DESC` or `ASC`, uppercase.
* `status` - (Optional, String) Workflow status filter: `ALL_RUNNING`: All scheduled, `ALL_FREEZED`: All paused, `ALL_STOPPTED`: All offline, `PART_RUNNING`: Partially scheduled, `ALL_NO_RUNNING`: All unscheduled, `ALL_INVALID`: All invalid.
* `workflow_type` - (Optional, String) Workflow type filter, supported values: `Cycle` or `Manual`. By default, only `Cycle` is queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Record list.


