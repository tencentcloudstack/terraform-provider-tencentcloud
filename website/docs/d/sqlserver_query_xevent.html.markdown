---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_query_xevent"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_query_xevent"
description: |-
  Use this data source to query detailed information of sqlserver query_xevent
---

# tencentcloud_sqlserver_query_xevent

Use this data source to query detailed information of sqlserver query_xevent

## Example Usage

```hcl
data "tencentcloud_sqlserver_query_xevent" "query_xevent" {
  instance_id = "mssql-gyg9xycl"
  event_type  = "blocked"
  start_time  = "2023-06-27 00:00:00"
  end_time    = "2023-07-01 00:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) Generation end time of an extended file.
* `event_type` - (Required, String) Event type. Valid values: slow (Slow SQL event), blocked (blocking event), deadlock` (deadlock event).
* `instance_id` - (Required, String) Instance ID.
* `start_time` - (Required, String) Generation start time of an extended file.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `events` - List of extended events.
  * `end_time` - Generation end time of an extended file.
  * `event_type` - Event type. Valid values: slow (Slow SQL event), blocked (blocking event), deadlock (deadlock event).
  * `external_addr` - Download address on the public network.
  * `file_name` - File name of an extended event.
  * `id` - ID.
  * `internal_addr` - Download address on the private network.
  * `size` - File size of an extended event.
  * `start_time` - Generation start time of an extended file.
  * `status` - Event record status. Valid values: 1 (succeeded), 2 (failed).


