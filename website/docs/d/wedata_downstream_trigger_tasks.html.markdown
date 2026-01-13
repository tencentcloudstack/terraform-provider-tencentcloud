---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_downstream_trigger_tasks"
sidebar_current: "docs-tencentcloud-datasource-wedata_downstream_trigger_tasks"
description: |-
  Use this data source to query detailed information of wedata downstream trigger tasks.
---

# tencentcloud_wedata_downstream_trigger_tasks

Use this data source to query detailed information of wedata downstream trigger tasks.

## Example Usage

```hcl
data "tencentcloud_wedata_downstream_trigger_tasks" "downstream_trigger_tasks" {
  project_id = "3108707295180644352"
  task_id    = "20241024174712123456"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Downstream dependency details.


