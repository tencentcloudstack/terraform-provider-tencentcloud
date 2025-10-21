---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_slow_log_data"
sidebar_current: "docs-tencentcloud-datasource-mysql_slow_log_data"
description: |-
  Use this data source to query detailed information of mysql slow_log_data
---

# tencentcloud_mysql_slow_log_data

Use this data source to query detailed information of mysql slow_log_data

## Example Usage

```hcl
data "tencentcloud_mysql_slow_log_data" "slow_log_data" {
  instance_id = "cdb-fitq5t9h"
  start_time  = 1682664459
  end_time    = 1684392459
  user_hosts  = ["169.254.128.158"]
  user_names  = ["keep_dts"]
  data_bases  = ["tf_ci_test"]
  sort_by     = "Timestamp"
  order_by    = "ASC"
  inst_type   = "slave"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) End timestamp. For example 1585142640.
* `instance_id` - (Required, String) instance id.
* `start_time` - (Required, Int) Start timestamp. For example 1585142640.
* `data_bases` - (Optional, Set: [`String`]) List of databases accessed.
* `inst_type` - (Optional, String) Only valid when the instance is the master instance or disaster recovery instance, the optional value: slave, which means to pull the log of the slave machine.
* `order_by` - (Optional, String) Sort in ascending or descending order. Currently supported: ASC,DESC.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort field. Currently supported: Timestamp, QueryTime, LockTime, RowsExamined, RowsSent.
* `user_hosts` - (Optional, Set: [`String`]) List of client hosts.
* `user_names` - (Optional, Set: [`String`]) A list of client usernames.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Query records.
  * `database` - database name.
  * `lock_time` - Lock duration (seconds).
  * `md5` - The md5 of the Sql statement.
  * `query_time` - Sql execution time (seconds).
  * `rows_examined` - The number of rows to scan.
  * `rows_sent` - The number of rows in the result set.
  * `sql_template` - Sql template.
  * `sql_text` - Sql statement.
  * `timestamp` - Sql execution time.
  * `user_host` - client address.
  * `user_name` - user name.


