---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_instance_log"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_instance_log"
description: |-
  Use this data source to query detailed information of wedata task instance log
---

# tencentcloud_wedata_task_instance_log

Use this data source to query detailed information of wedata task instance log

## Example Usage

```hcl
data "tencentcloud_wedata_task_instance_log" "wedata_task_instance_log" {
  project_id   = "1859317240494305280"
  instance_key = "20250324192240178_2025-10-13 11:50:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_key` - (Required, String) Unique instance identifier.
* `project_id` - (Required, String) Project ID.
* `life_round_num` - (Optional, Int) Instance lifecycle number, identifying a specific execution of the instance. For example: the first run of a periodic instance is 0, if manually rerun the second execution is 1; defaults to the latest execution.
* `log_level` - (Optional, String) Log level, default All - Info - Debug - Warn - Error - All.
* `next_cursor` - (Optional, String) Pagination cursor for log queries, no business meaning. First query uses null, subsequent queries use NextCursor from previous response.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Scheduled instance details.


