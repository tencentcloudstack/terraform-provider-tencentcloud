---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_version"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_version"
description: |-
  Use this data source to query detailed information of wedata wedata_task_version
---

# tencentcloud_wedata_task_version

Use this data source to query detailed information of wedata wedata_task_version

## Example Usage

```hcl
data "tencentcloud_wedata_task_version" "wedata_task_version" {
  project_id = "2905622749543821312"
  task_id    = "20251015164958429"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.
* `version_id` - (Optional, String) Submit version ID. If not specified, the latest submit version will be used by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Version detail.


