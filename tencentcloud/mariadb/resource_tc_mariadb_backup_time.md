Provides a resource to create a mariadb backup_time

Example Usage

```hcl
resource "tencentcloud_mariadb_backup_time" "backup_time" {
  instance_id       = "tdsql-9vqvls95"
  start_backup_time = "01:00"
  end_backup_time   = "04:00"
}
```

Import

mariadb backup_time can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_backup_time.backup_time backup_time_id
```