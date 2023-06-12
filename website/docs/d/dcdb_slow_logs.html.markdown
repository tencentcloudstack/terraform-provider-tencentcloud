---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_slow_logs"
sidebar_current: "docs-tencentcloud-datasource-dcdb_slow_logs"
description: |-
  Use this data source to query detailed information of dcdb slow_logs
---

# tencentcloud_dcdb_slow_logs

Use this data source to query detailed information of dcdb slow_logs

## Example Usage

```hcl
data "tencentcloud_dcdb_slow_logs" "slow_logs" {
  instance_id   = local.dcdb_id
  start_time    = "%s"
  end_time      = "%s"
  shard_id      = "shard-1b5r04az"
  db            = "tf_test_db"
  order_by      = "query_time_sum"
  order_by_type = "desc"
  slave         = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of `tdsqlshard-ow728lmc`.
* `shard_id` - (Required, String) Instance shard ID in the format of `shard-rc754ljk`.
* `start_time` - (Required, String) Query start time in the format of 2016-07-23 14:55:20.
* `db` - (Optional, String) Specific name of the database to be queried.
* `end_time` - (Optional, String) Query end time in the format of 2016-08-22 14:55:20.
* `order_by_type` - (Optional, String) Sorting order. Valid values: desc, asc.
* `order_by` - (Optional, String) Sorting metric. Valid values: query_time_sum, query_count.
* `result_output_file` - (Optional, String) Used to save results.
* `slave` - (Optional, Int) Query slow queries from either the primary or the replica. Valid values: 0 (primary), 1 (replica).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Slow query log data.
  * `check_sum` - Statement checksum for querying details.
  * `db` - Database name.
  * `example_sql` - Sample SQLNote: This field may return null, indicating that no valid values can be obtained.
  * `finger_print` - Abstracted SQL statement.
  * `host` - Host address of account.
  * `lock_time_avg` - Average lock time.
  * `lock_time_max` - Maximum lock time.
  * `lock_time_min` - Minimum lock time.
  * `lock_time_sum` - Total lock time.
  * `query_count` - Number of queries.
  * `query_time_avg` - Average query time.
  * `query_time_max` - Maximum query time.
  * `query_time_min` - Minimum query time.
  * `query_time_sum` - Total query time.
  * `rows_examined_sum` - Number of scanned rows.
  * `rows_sent_sum` - Number of sent rows.
  * `ts_max` - Last execution time.
  * `ts_min` - First execution time.
  * `user` - Account.
* `lock_time_sum` - Total statement lock time.
* `query_count` - Total number of statement queries.
* `query_time_sum` - Total statement query time.


