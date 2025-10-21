---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_job_submission_log"
sidebar_current: "docs-tencentcloud-datasource-oceanus_job_submission_log"
description: |-
  Use this data source to query detailed information of oceanus job_submission_log
---

# tencentcloud_oceanus_job_submission_log

Use this data source to query detailed information of oceanus job_submission_log

## Example Usage

```hcl
data "tencentcloud_oceanus_job_submission_log" "example" {
  job_id           = "cql-314rw6w0"
  start_time       = 1696130964345
  end_time         = 1698118169241
  running_order_id = 0
  order_type       = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) End time, unix timestamp, in milliseconds.
* `job_id` - (Required, String) Job ID.
* `start_time` - (Required, Int) Start time, unix timestamp, in milliseconds.
* `cursor` - (Optional, String) Cursor, default empty, first request does not need to pass in.
* `keyword` - (Optional, String) Keyword, default empty.
* `order_type` - (Optional, String) Sorting method, default asc, asc: ascending, desc: descending.
* `result_output_file` - (Optional, String) Used to save results.
* `running_order_id` - (Optional, Int) Job instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `job_instance_list` - Job instance list during the specified time period.
  * `job_instance_start_time` - The startup time of the instance.
  * `running_order_id` - The ID of the instance, starting from 1 in the order of startup time.
  * `starting_millis` - The startup time of the instance in milliseconds.
* `job_request_id` - Request ID of starting job.
* `list_over` - Whether the list is over.
* `log_content_list` - The list of log contents.
  * `container_name` - The name of the container to which the log belongs.
  * `log` - The content of the log.
  * `pkg_id` - The ID of the log group.
  * `pkg_log_id` - The ID of the log, which is unique within the log group.
  * `time` - The timestamp in milliseconds.
* `log_list` - Log list, deprecated.


