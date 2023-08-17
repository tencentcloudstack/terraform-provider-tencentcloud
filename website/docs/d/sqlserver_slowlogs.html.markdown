---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_slowlogs"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_slowlogs"
description: |-
  Use this data source to query detailed information of sqlserver slowlogs
---

# tencentcloud_sqlserver_slowlogs

Use this data source to query detailed information of sqlserver slowlogs

## Example Usage

```hcl
data "tencentcloud_sqlserver_slowlogs" "example" {
  instance_id = "mssql-qelbzgwf"
  start_time  = "2023-08-01 00:00:00"
  end_time    = "2023-08-07 00:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) Query end time.
* `instance_id` - (Required, String) Instance ID.
* `start_time` - (Required, String) Query start time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `slowlogs` - Information list of slow query logs.
  * `count` - Number of logs in file.
  * `end_time` - File generation end time.
  * `external_addr` - Download address for public network.
  * `id` - Unique ID of slow query log file.
  * `internal_addr` - Download address for private network.
  * `size` - File size in KB.
  * `start_time` - File generation start time.
  * `status` - Status (1: success, 2: failure) Note: this field may return null, indicating that no valid values can be obtained.


