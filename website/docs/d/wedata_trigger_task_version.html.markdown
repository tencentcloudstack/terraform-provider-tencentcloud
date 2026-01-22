---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_task_version"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_task_version"
description: |-
  Use this data source to query detailed information of wedata trigger task version.
---

# tencentcloud_wedata_trigger_task_version

Use this data source to query detailed information of wedata trigger task version.

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_task_version" "trigger_task_version" {
  project_id = "1840731342175234"
  task_id    = "20241024174712123456"
  version_id = "20241024174712123456_1"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.
* `version_id` - (Optional, String) Submitted version ID; if not provided, the latest submitted version is used by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Version details.


