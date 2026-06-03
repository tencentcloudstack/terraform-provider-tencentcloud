Provides a resource to create a MongoDB audit log file

Example Usage

```hcl
resource "tencentcloud_mongodb_audit_log_file" "example" {
  instance_id = "cmgo-5aqo4yf7"
  start_time  = "2026-06-01 10:29:20"
  end_time    = "2026-06-01 10:39:20"
  order       = "ASC"
  order_by    = "timestamp"

  filter {
    host        = ["10.0.0.1"]
    user        = ["admin"]
    exec_time   = 100
    affect_rows = 10
    atype       = ["insert", "update"]
    result      = ["ok"]
    param       = ["keyword"]
  }
}
```

Import

mongodb audit_log_file can be imported using the composite instance_id#file_name, e.g.

```
terraform import tencentcloud_mongodb_audit_log_file.example cmgo-5aqo4yf7#1309118522_cmgo-5aqo4yf7_1780474413_109642711.csv
```
