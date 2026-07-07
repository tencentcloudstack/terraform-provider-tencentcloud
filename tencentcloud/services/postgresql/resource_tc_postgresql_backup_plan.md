Provides a resource to create a PostgreSQL backup plan

Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan" "backup_plan" {
  db_instance_id               = "postgres-test123"
  plan_name                    = "tf-example-backup-plan"
  backup_period_type           = "month"
  backup_period                = ["1", "2", "15"]
  min_backup_start_time        = "01:00:00"
  max_backup_start_time        = "02:00:00"
  base_backup_retention_period = 30
}
```

Import

PostgreSQL backup plan can be imported using the db_instance_id#plan_id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan.backup_plan postgres-test123#plan-xxxx
```
