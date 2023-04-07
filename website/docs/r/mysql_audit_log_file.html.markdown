---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_audit_log_file"
sidebar_current: "docs-tencentcloud-resource-mysql_audit_log_file"
description: |-
  Provides a resource to create a mysql audit_log_file
---

# tencentcloud_mysql_audit_log_file

Provides a resource to create a mysql audit_log_file

## Example Usage

```hcl
resource "tencentcloud_mysql_audit_log_file" "audit_log_file" {
  instance_id = "cdb-fitq5t9h"
  start_time  = "2023-03-28 20:14:00"
  end_time    = "2023-03-29 20:14:00"
  order       = "ASC"
  order_by    = "timestamp"
  filter {
    host = ["30.50.207.46"]
    user = ["keep_dbbrain"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) end time.
* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `start_time` - (Required, String, ForceNew) start time.
* `filter` - (Optional, List, ForceNew) Filter condition. Logs can be filtered according to the filter conditions set.
* `order_by` - (Optional, String, ForceNew) Sort field. supported values include:`timestamp` - timestamp; `affectRows` - affected rows; `execTime` - execution time.
* `order` - (Optional, String, ForceNew) Sort by. supported values are: `ASC`- ascending order, `DESC`- descending order.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int) Affects the number of rows. Indicates to filter audit logs whose number of affected rows is greater than this value.
* `db_name` - (Optional, Set) Database name.
* `exec_time` - (Optional, Int) Execution time. The unit is: ms. Indicates to filter audit logs whose execution time is greater than this value.
* `host` - (Optional, Set) Client address.
* `policy_name` - (Optional, Set) The name of policy.
* `sql_type` - (Optional, String) SQL type. Currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.
* `sql_types` - (Optional, Set) SQL type. Supports simultaneous query of multiple types. Currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.
* `sql` - (Optional, String) SQL statement. support fuzzy matching.
* `sqls` - (Optional, Set) SQL statement. Support passing multiple sql statements.
* `table_name` - (Optional, Set) Table name.
* `user` - (Optional, Set) User name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `download_url` - download url.
* `file_size` - size of file(KB).


