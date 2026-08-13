Provides a resource to create a PostgreSQL audit log file

Example Usage

```hcl
resource "tencentcloud_postgres_audit_log_file" "example" {
  instance_id = "postgres-xxxxxxxx"
  start_time  = "2026-03-25 00:00:00"
  end_time    = "2026-03-25 01:00:00"
  product     = "postgres"
}
```

Create with filter conditions

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

Import

PostgreSQL audit log file can be imported using the instanceId#fileName, e.g.

```
terraform import tencentcloud_postgres_audit_log_file.example postgres-xxxxxxxx#audit_log_file_name
```
