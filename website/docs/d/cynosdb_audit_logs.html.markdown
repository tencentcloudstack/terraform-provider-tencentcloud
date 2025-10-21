---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_audit_logs"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_audit_logs"
description: |-
  Use this data source to query detailed information of cynosdb audit_logs
---

# tencentcloud_cynosdb_audit_logs

Use this data source to query detailed information of cynosdb audit_logs

## Example Usage

```hcl
data "tencentcloud_cynosdb_audit_logs" "audit_logs" {
  instance_id = "cynosdbmysql-ins-afqx1hy0"
  start_time  = "2023-06-18 10:00:00"
  end_time    = "2023-06-18 10:00:02"
  order       = "DESC"
  order_by    = "timestamp"
  filter {
    host        = ["30.50.207.176"]
    user        = ["keep_dts"]
    policy_name = ["default_audit"]
    sql_type    = "SELECT"
    sql         = "SELECT @@max_allowed_packet"
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) The end time is in the format of 2017-07-12 10:29:20.
* `instance_id` - (Required, String) Instance ID.
* `start_time` - (Required, String) Start time, format: 2017-07-12 10:29:20.
* `filter` - (Optional, List) Filter conditions. You can filter logs according to the set filtering criteria.
* `order_by` - (Optional, String) Sort fields. The supported values include: timestamp - timestamp; &amp;#39;effectRows&amp;#39; - affects the number of rows; &amp;#39;execTime&amp;#39; - Execution time.
* `order` - (Optional, String) Sort by. The supported values include: ASC - ascending order, DESC - descending order.
* `result_output_file` - (Optional, String) Used to save results.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int) Affects the number of rows. Indicates that filtering affects audit logs with rows greater than this value.
* `db_name` - (Optional, Set) Database name.
* `exec_time` - (Optional, Int) Execution time. Unit: ms. Indicates audit logs with a filter execution time greater than this value.
* `host` - (Optional, Set) Client address.
* `policy_name` - (Optional, Set) Audit policy name.
* `sent_rows` - (Optional, Int) Returns the number of rows.
* `sql_type` - (Optional, String) SQL type. Currently supported: SELECT, Insert, UPDATE, DELETE, CREATE, DROP, ALT, SET, REPLACE, EXECUTE.
* `sql_types` - (Optional, Set) SQL type. Supports simultaneous querying of multiple types. Currently supported: SELECT, Insert, UPDATE, DELETE, CREATE, DROP, ALT, SET, REPLACE, EXECUTE.
* `sql` - (Optional, String) SQL statement. Supports fuzzy matching.
* `sqls` - (Optional, Set) SQL statement. Supports passing multiple SQL statements.
* `table_name` - (Optional, Set) Table name.
* `thread_id` - (Optional, Set) Thread ID.
* `user` - (Optional, Set) User name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Audit log details. Note: This field may return null, indicating that a valid value cannot be obtained.
  * `affect_rows` - Affects the number of rows.
  * `db_name` - Database name.
  * `err_code` - Error code.
  * `exec_time` - Execution time.
  * `host` - Client address.
  * `instance_name` - Instance name.
  * `policy_name` - Audit policy name.
  * `sent_rows` - Number of rows sent.
  * `sql_type` - SQL type.
  * `sql` - SQL statement.
  * `table_name` - Table name.
  * `thread_id` - Execution thread ID.
  * `timestamp` - Timestamp.
  * `user` - User name.


