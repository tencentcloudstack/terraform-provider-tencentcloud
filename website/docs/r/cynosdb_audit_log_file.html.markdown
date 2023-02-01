---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_audit_log_file"
sidebar_current: "docs-tencentcloud-resource-cynosdb_audit_log_file"
description: |-
  Provides a resource to create a cynosdb audit_log_file
---

# tencentcloud_cynosdb_audit_log_file

Provides a resource to create a cynosdb audit_log_file

## Example Usage

```hcl
resource "tencentcloud_cynosdb_audit_log_file" "audit_log_file" {
  instance_id = "xxxxxxx"
  start_time  = "2022-07-12 10:29:20"
  end_time    = "2022-08-12 10:29:20"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) End time.
* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `start_time` - (Required, String, ForceNew) Start time.
* `filter` - (Optional, List, ForceNew) Filter condition. Logs can be filtered according to the filter conditions set.
* `order_by` - (Optional, String, ForceNew) Sort field. supported values are:
`timestamp` - timestamp
`affectRows` - affected rows
`execTime` - execution time.
* `order` - (Optional, String, ForceNew) Sort by. Supported values are: `ASC` - ascending, `DESC` - descending.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int) Affects the number of rows. Indicates that the audit log whose number of affected rows is greater than this value is filtered.
* `db_name` - (Optional, Set) The name of database.
* `exec_time` - (Optional, Int) Execution time. The unit is: ms. Indicates to filter audit logs whose execution time is greater than this value.
* `host` - (Optional, Set) Client host.
* `policy_name` - (Optional, Set) The name of audit policy.
* `sent_rows` - (Optional, Int) Return the number of rows.
* `sql_type` - (Optional, String) SQL type. currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.
* `sql_types` - (Optional, Set) SQL type. Supports simultaneous query of multiple types. currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.
* `sql` - (Optional, String) SQL statement. Support fuzzy matching.
* `sqls` - (Optional, Set) SQL statement. Support passing multiple sql statements.
* `table_name` - (Optional, Set) The name of table.
* `thread_id` - (Optional, Set) The ID of thread.
* `user` - (Optional, Set) User name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Audit log file creation time. The format is 2019-03-20 17:09:13.
* `download_url` - The download address of the audit logs.
* `err_msg` - Error message.
* `file_name` - Audit log file name.
* `file_size` - File size, The unit is KB.


