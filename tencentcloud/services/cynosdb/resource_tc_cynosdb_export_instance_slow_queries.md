Provides a resource to create a cynosdb export_instance_slow_queries

Example Usage

```hcl
resource "tencentcloud_cynosdb_export_instance_slow_queries" "export_instance_slow_queries" {
  instance_id = "cynosdbmysql-ins-123"
  start_time  = "2022-01-01 12:00:00"
  end_time    = "2022-01-01 14:00:00"
  username    = "root"
  host        = "10.10.10.10"
  database    = "db1"
  file_type   = "csv"
}
```