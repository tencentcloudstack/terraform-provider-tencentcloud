---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_downstream_tasks"
sidebar_current: "docs-tencentcloud-datasource-wedata_downstream_tasks"
description: |-
  Use this data source to query detailed information of wedata wedata_downstream_tasks
---

# tencentcloud_wedata_downstream_tasks

Use this data source to query detailed information of wedata wedata_downstream_tasks

## Example Usage

```hcl
data "tencentcloud_wedata_downstream_tasks" "wedata_downstream_tasks" {
  project_id = "2905622749543821312"
  task_id    = "20251015164958429"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Describes the downstream dependency details.


