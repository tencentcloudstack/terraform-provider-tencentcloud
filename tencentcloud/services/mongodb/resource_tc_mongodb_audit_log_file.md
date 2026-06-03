Provides a resource to create a MongoDB audit log file

Example Usage

```hcl
resource "tencentcloud_mongodb_audit_log_file" "example" {
  instance_id = "cmgo-xfts1234"
  start_time  = "2021-07-12 10:29:20"
  end_time    = "2021-07-12 10:39:20"
  order       = "ASC"
  order_by    = "timestamp"

  filter {
    host        = ["10.0.0.1"]
    user        = ["admin"]
    exec_time   = 100
    affect_rows = 10
    atype       = ["insert", "update"]
    result      = ["ok"]
    param       = ["test"]
  }
}
```

Import

mongodb audit_log_file can be imported using the composite id, e.g.

```
terraform import tencentcloud_mongodb_audit_log_file.example cmgo-xfts1234#audit_log_20210712.log
```
