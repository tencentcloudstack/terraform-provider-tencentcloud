Provides a resource to create a PostgreSQL audit log file

Example Usage

```hcl
resource "tencentcloud_postgres_audit_log_file" "example" {
  instance_id = "postgres-ckwcgdf1"
  start_time  = "2026-03-25 00:00:00"
  end_time    = "2026-03-25 01:00:00"
}
```

Create with filter conditions

```hcl
resource "tencentcloud_postgres_audit_log_file" "example" {
  instance_id = "postgres-ckwcgdf1"
  start_time  = "2026-03-25 00:00:00"
  end_time    = "2026-03-25 01:00:00"

  filter {
    affect_rows = 100
    db_name     = ["testDB"]
    exec_time   = 1000
    host        = ["10.0.0.1"]
    sql         = "SELECT"
    user        = ["admin"]
    sql_type    = ["SELECT", "INSERT"]
  }
}
```
