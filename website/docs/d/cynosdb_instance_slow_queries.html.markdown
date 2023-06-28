---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_instance_slow_queries"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_instance_slow_queries"
description: |-
  Use this data source to query detailed information of cynosdb instance_slow_queries
---

# tencentcloud_cynosdb_instance_slow_queries

Use this data source to query detailed information of cynosdb instance_slow_queries

## Example Usage

Query slow queries of instance

```hcl
variable "cynosdb_cluster_id" {
  default = "default_cynosdb_cluster"
}

data "tencentcloud_cynosdb_instance_slow_queries" "instance_slow_queries" {
  instance_id   = var.cynosdb_cluster_id
  start_time    = "2023-06-20 23:19:03"
  end_time      = "2023-06-30 23:19:03"
  username      = "keep_dts"
  host          = "%%"
  database      = "tf_ci_test"
  order_by      = "QueryTime"
  order_by_type = "desc"
}
```

Query slow queries by time range

```hcl
variable "cynosdb_cluster_id" {
  default = "default_cynosdb_cluster"
}

data "tencentcloud_cynosdb_instance_slow_queries" "instance_slow_queries" {
  instance_id   = var.cynosdb_cluster_id
  start_time    = "2023-06-20 23:19:03"
  end_time      = "2023-06-30 23:19:03"
  order_by      = "QueryTime"
  order_by_type = "desc"
}
```

Query slow queries by user and db name

```hcl
variable "cynosdb_cluster_id" {
  default = "default_cynosdb_cluster"
}

data "tencentcloud_cynosdb_instance_slow_queries" "instance_slow_queries" {
  instance_id   = var.cynosdb_cluster_id
  username      = "keep_dts"
  host          = "%%"
  database      = "tf_ci_test"
  order_by      = "QueryTime"
  order_by_type = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `database` - (Optional, String) Database name.
* `end_time` - (Optional, String) Latest transaction start time.
* `host` - (Optional, String) Client host.
* `order_by_type` - (Optional, String) Sort type, optional values: asc, desc.
* `order_by` - (Optional, String) Sort field, optional values: QueryTime, LockTime, RowsExamined, RowsSent.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, String) Earliest transaction start time.
* `username` - (Optional, String) user name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `slow_queries` - Slow query records.
  * `database` - Database name.
  * `lock_time` - Lock duration in seconds.
  * `query_time` - Execution time in seconds.
  * `rows_examined` - Scan Rows.
  * `rows_sent` - Return the number of rows.
  * `sql_md5` - SQL statement md5.
  * `sql_template` - SQL template.
  * `sql_text` - SQL statement.
  * `timestamp` - Execution timestamp.
  * `user_host` - Client host.
  * `user_name` - user name.


