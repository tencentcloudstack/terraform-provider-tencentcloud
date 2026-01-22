---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_task_versions"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_task_versions"
description: |-
  Use this data source to query detailed information of wedata trigger task versions.
---

# tencentcloud_wedata_trigger_task_versions

Use this data source to query detailed information of wedata trigger task versions.

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_task_versions" "trigger_task_versions" {
  project_id        = "1840731342175234"
  task_id           = "20241024174712123456"
  task_version_type = "SAVE"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.
* `task_version_type` - (Optional, String) Saved version: SAVE; Submitted version: SUBMIT. Default is SAVE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Version list.


