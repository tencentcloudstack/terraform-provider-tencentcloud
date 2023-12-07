Provides a resource to create a cynosdb audit_log_file

Example Usage

```hcl
resource "tencentcloud_cynosdb_audit_log_file" "audit_log_file" {
  instance_id = "cynosdbmysql-ins-afqx1hy0"
  start_time  = "2022-07-12 10:29:20"
  end_time    = "2022-08-12 10:29:20"
}
```