Provides a resource to create a cynosdb export_instance_error_logs

Example Usage

```hcl
resource "tencentcloud_cynosdb_export_instance_error_logs" "export_instance_error_logs" {
  instance_id   = "cynosdbmysql-ins-afqx1hy0"
  start_time    = "2022-01-01 12:00:00"
  end_time      = "2022-01-01 14:00:00"
  log_levels    = ["note"]
  key_words     = ["content"]
  file_type     = "csv"
  order_by      = "Timestamp"
  order_by_type = "ASC"
}
```