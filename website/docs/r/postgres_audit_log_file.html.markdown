---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgres_audit_log_file"
sidebar_current: "docs-tencentcloud-resource-postgres_audit_log_file"
description: |-
  Provides a resource to create a PostgreSQL audit log file
---

# tencentcloud_postgres_audit_log_file

Provides a resource to create a PostgreSQL audit log file

## Example Usage

```hcl
resource "tencentcloud_postgres_audit_log_file" "example" {
  instance_id = "postgres-xxxxxxxx"
  start_time  = "2026-03-25 00:00:00"
  end_time    = "2026-03-25 01:00:00"
  product     = "postgres"
}
```

### Create with filter conditions

```hcl
resource "tencentcloud_postgres_audit_log_file" "example_with_filter" {
  instance_id = "postgres-xxxxxxxx"
  start_time  = "2026-03-25 00:00:00"
  end_time    = "2026-03-25 01:00:00"
  product     = "postgres"

  filter {
    affect_rows = 100
    db_name     = ["testdb"]
    exec_time   = 1000
    host        = ["10.0.0.1"]
    sql         = "SELECT"
    user        = ["admin"]
    sql_type    = ["SELECT", "INSERT"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) End time, format: `2026-03-25 01:00:00`.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `product` - (Required, String, ForceNew) Product name, fixed value: `postgres`.
* `start_time` - (Required, String, ForceNew) Start time, format: `2026-03-25 00:00:00`.
* `filter` - (Optional, List, ForceNew) Filter conditions.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int, ForceNew) Affect rows.
* `db_name` - (Optional, List, ForceNew) Database name list.
* `exec_time` - (Optional, Int, ForceNew) Execution time.
* `host` - (Optional, List, ForceNew) Host list.
* `sql_type` - (Optional, List, ForceNew) SQL type list.
* `sql` - (Optional, String, ForceNew) SQL statement.
* `user` - (Optional, List, ForceNew) User name list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `download_url` - Download URL.
* `err_msg` - Error message.
* `file_name` - Audit log file name.
* `file_size` - File size in MB.
* `finish_time` - Finish time.
* `progress` - Download progress.
* `status` - Task status. Values: `success`, `running`, `failed`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `30m`) Used when creating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

PostgreSQL audit log file can be imported using the instanceId#fileName, e.g.

```
terraform import tencentcloud_postgres_audit_log_file.example postgres-xxxxxxxx#audit_log_file_name
```

