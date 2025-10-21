---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_slow_log_top_sqls"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_slow_log_top_sqls"
description: |-
  Use this data source to query detailed information of dbbrain slow_log_top_sqls
---

# tencentcloud_dbbrain_slow_log_top_sqls

Use this data source to query detailed information of dbbrain slow_log_top_sqls

## Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_top_sqls" "test" {
  instance_id = "%s"
  start_time  = "%s"
  end_time    = "%s"
  sort_by     = "QueryTimeMax"
  order_by    = "ASC"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) The deadline, such as `2019-09-11 10:13:14`, the interval between the deadline and the start time is less than 7 days.
* `instance_id` - (Required, String) instance id.
* `start_time` - (Required, String) Start time, such as `2019-09-10 12:13:14`.
* `order_by` - (Optional, String) The sorting method supports ASC (ascending) and DESC (descending). The default is DESC.
* `product` - (Optional, String) Service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database CynosDB for MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.
* `schema_list` - (Optional, List) Array of database names.
* `sort_by` - (Optional, String) Sort key, currently supports sort keys such as QueryTime, ExecTimes, RowsSent, LockTime and RowsExamined, the default is QueryTime.

The `schema_list` object supports the following:

* `schema` - (Required, String) DB name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rows` - Slow log top sql list.
  * `exec_times` - Execution times.
  * `lock_time_avg` - Average lock waiting time, in seconds.
  * `lock_time_max` - Maximum lock waiting time, in seconds.
  * `lock_time_min` - Minimum lock waiting time, in seconds.
  * `lock_time_ratio` - The ratio of the total lock waiting time of SQL, in %.
  * `lock_time` - SQL total lock waiting time, in seconds.
  * `md5` - MD5 value of SOL template.
  * `query_time_avg` - Average execution time, in seconds.
  * `query_time_max` - The maximum execution time, in seconds.
  * `query_time_min` - The minimum execution time, in seconds.
  * `query_time_ratio` - Total time-consuming ratio, unit %.
  * `query_time` - Total time, in seconds.
  * `rows_examined_avg` - average number of lines scanned.
  * `rows_examined_max` - Maximum number of scan lines.
  * `rows_examined_min` - Minimum number of scan lines.
  * `rows_examined_ratio` - The proportion of the total number of scanned lines, unit %.
  * `rows_examined` - total scan lines.
  * `rows_sent_avg` - average number of rows returned.
  * `rows_sent_max` - Maximum number of rows returned.
  * `rows_sent_min` - Minimum number of rows returned.
  * `rows_sent_ratio` - The proportion of the total number of rows returned, in %.
  * `rows_sent` - total number of rows returned.
  * `schema` - Database name.
  * `sql_template` - sql template.
  * `sql_text` - SQL with parameters (random).


