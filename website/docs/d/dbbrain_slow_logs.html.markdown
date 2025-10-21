---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_slow_logs"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_slow_logs"
description: |-
  Use this data source to query detailed information of dbbrain slow_logs
---

# tencentcloud_dbbrain_slow_logs

Use this data source to query detailed information of dbbrain slow_logs

## Example Usage

```hcl
data "tencentcloud_dbbrain_slow_logs" "slow_logs" {
  product     = "mysql"
  instance_id = "%s"
  md5         = "4961208426639258265"
  start_time  = "%s"
  end_time    = "%s"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) The deadline, such as 2019-09-11 10:13:14, the interval between the deadline and the start time is less than 7 days.
* `instance_id` - (Required, String) instance Id.
* `md5` - (Required, String) md5 value of sql template.
* `product` - (Required, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.
* `start_time` - (Required, String) Start time, such as 2019-09-10 12:13:14.
* `db` - (Optional, Set: [`String`]) database list.
* `ip` - (Optional, Set: [`String`]) ip.
* `key` - (Optional, Set: [`String`]) keywords.
* `result_output_file` - (Optional, String) Used to save results.
* `time` - (Optional, Set: [`Int`]) Time-consuming interval, the left and right boundaries of the time-consuming interval correspond to the 0th element and the first element of the array respectively.
* `user` - (Optional, Set: [`String`]) user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rows` - Slow log details.
  * `database` - database.
  * `lock_time` - lock time, in secondsNote: This field may return null, indicating that no valid value can be obtained.
  * `query_time` - Execution time, in seconds.
  * `rows_examined` - scan linesNote: This field may return null, indicating that no valid value can be obtained.
  * `rows_sent` - Return the number of rowsNote: This field may return null, indicating that no valid value can be obtained.
  * `sql_text` - sql statement.
  * `timestamp` - Slow log start time.
  * `user_host` - Ip sourceNote: This field may return null, indicating that no valid value can be obtained.
  * `user_name` - User sourceNote: This field may return null, indicating that no valid value can be obtained.


