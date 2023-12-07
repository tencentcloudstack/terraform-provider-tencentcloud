Provides a resource to create a postgres backup_plan_config

Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan_config" "backup_plan_config" {
  db_instance_id = local.pgsql_id
  min_backup_start_time = "01:00:00"
  max_backup_start_time = "02:00:00"
  base_backup_retention_period = 7
  backup_period = ["monday","wednesday","friday"]
}
```

Import

postgres backup_plan_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan_config.backup_plan_config backup_plan_config_id
```