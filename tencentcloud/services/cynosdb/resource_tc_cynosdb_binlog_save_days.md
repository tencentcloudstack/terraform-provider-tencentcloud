Provides a resource to create a cynosdb binlog_save_days

Example Usage

```hcl
resource "tencentcloud_cynosdb_binlog_save_days" "binlog_save_days" {
  cluster_id       = "cynosdbmysql-123"
  binlog_save_days = 7
}
```

Import

cynosdb binlog_save_days can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_binlog_save_days.binlog_save_days binlog_save_days_id
```