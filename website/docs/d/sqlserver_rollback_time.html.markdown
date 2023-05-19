---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_rollback_time"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_rollback_time"
description: |-
  Use this data source to query detailed information of sqlserver rollback_time
---

# tencentcloud_sqlserver_rollback_time

Use this data source to query detailed information of sqlserver rollback_time

## Example Usage

```hcl
data "tencentcloud_sqlserver_rollback_time" "rollback_time" {
  instance_id = "mssql-qelbzgwf"
  dbs         = ["keep_pubsub_db"]
}
```

## Argument Reference

The following arguments are supported:

* `dbs` - (Required, Set: [`String`]) List of databases to be queried.
* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `details` - Information of time range available for database rollback.
  * `db_name` - Database name.
  * `end_time` - End time of time range available for rollback.
  * `start_time` - Start time of time range available for rollback.


