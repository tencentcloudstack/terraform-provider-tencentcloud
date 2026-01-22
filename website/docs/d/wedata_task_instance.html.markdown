---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_instance"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_instance"
description: |-
  Use this data source to query detailed information of wedata task instance
---

# tencentcloud_wedata_task_instance

Use this data source to query detailed information of wedata task instance

## Example Usage

```hcl
data "tencentcloud_wedata_task_instance" "wedata_task_instance" {
  project_id   = "1859317240494305280"
  instance_key = "20250324192240178_2025-10-13 11:50:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_key` - (Required, String) Unique instance identifier, can be obtained via ListInstances.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `time_zone` - (Optional, String) Time zone, the time zone of the input time string, default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Instance details.


