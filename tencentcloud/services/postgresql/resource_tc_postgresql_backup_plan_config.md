Provides a resource to create a postgres backup plan config

Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan_config" "example" {
  db_instance_id               = "postgres-ckwcgdf1"
  min_backup_start_time        = "01:00:00"
  max_backup_start_time        = "03:00:00"
  backup_period                = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
  base_backup_retention_period = 7
  log_backup_retention_period  = 7
  backup_method                = "physical"
}
```

Import

postgres backup plan config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan_config.example postgres-ckwcgdf1
```
