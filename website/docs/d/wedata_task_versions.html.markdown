---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_versions"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_versions"
description: |-
  Use this data source to query detailed information of wedata wedata_task_versions
---

# tencentcloud_wedata_task_versions

Use this data source to query detailed information of wedata wedata_task_versions

## Example Usage

```hcl
data "tencentcloud_wedata_task_versions" "wedata_task_versions" {
  project_id = "2905622749543821312"
  task_id    = "20251015164958429"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.
* `task_version_type` - (Optional, String) SAVE version.
SUBMIT version.
Defaults to SAVE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Task version list.


