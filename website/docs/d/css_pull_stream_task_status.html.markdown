---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_pull_stream_task_status"
sidebar_current: "docs-tencentcloud-datasource-css_pull_stream_task_status"
description: |-
  Use this data source to query detailed information of css pull_stream_task_status
---

# tencentcloud_css_pull_stream_task_status

Use this data source to query detailed information of css pull_stream_task_status

## Example Usage

```hcl
data "tencentcloud_css_pull_stream_task_status" "pull_stream_task_status" {
  task_id = "63229997"
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_status_info` - Task status info.
  * `file_duration` - The duration of the VOD source file, in seconds.
  * `file_url` - Current use source url.
  * `looped_times` - The number of times a VOD source task is played in a loop.
  * `next_file_url` - The URL of the next progress VOD file.
  * `offset_time` - The playback offset of the VOD source, in seconds.
  * `report_time` - The latest heartbeat reporting time in UTC format, for example: 2022-02-11T10:00:00Z.Note: UTC time is 8 hours ahead of Beijing time.
  * `run_status` - Real run status:active,inactive.


