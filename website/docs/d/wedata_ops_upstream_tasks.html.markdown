---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_upstream_tasks"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_upstream_tasks"
description: |-
  Use this data source to query detailed information of wedata ops upstream task
---

# tencentcloud_wedata_ops_upstream_tasks

Use this data source to query detailed information of wedata ops upstream task

## Example Usage

```hcl
data "tencentcloud_wedata_ops_upstream_tasks" "wedata_ops_upstream_tasks" {
  project_id = "1859317240494305280"
  task_id    = "20250820150144998"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Upstream task details.


