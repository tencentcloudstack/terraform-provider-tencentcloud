---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_xlogs"
sidebar_current: "docs-tencentcloud-datasource-postgresql_xlogs"
description: |-
  Provide a datasource to query PostgreSQL Xlogs.
---

# tencentcloud_postgresql_xlogs

Provide a datasource to query PostgreSQL Xlogs.

## Example Usage

```hcl
data "tencentcloud_postgresql_xlogs" "foo" {
  instance_id = "postgres-xxxxxxxx"
  start_time  = "2022-01-01 00:00:00"
  end_time    = "2022-01-07 01:02:03"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) PostgreSQL instance id.
* `end_time` - (Optional, String) Xlog end time, format `yyyy-MM-dd hh:mm:ss`.
* `result_output_file` - (Optional, String) Used for save results.
* `start_time` - (Optional, String) Xlog start time, format `yyyy-MM-dd hh:mm:ss`, start time cannot before 7 days ago.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of Xlog query result.
  * `end_time` - Xlog file created end time.
  * `external_addr` - Xlog external download address.
  * `id` - Xlog id.
  * `internal_addr` - Xlog internal download address.
  * `size` - Xlog file size.
  * `start_time` - Xlog file created start time.


