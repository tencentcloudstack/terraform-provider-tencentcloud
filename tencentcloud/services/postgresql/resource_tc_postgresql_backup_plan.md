Provides a resource to create a PostgreSQL backup plan

Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan" "example" {
  db_instance_id               = "postgres-ckwcgdf1"
  plan_name                    = "tf-example"
  backup_period_type           = "month"
  backup_period                = ["1", "2", "15"]
  min_backup_start_time        = "01:00:00"
  max_backup_start_time        = "03:00:00"
  base_backup_retention_period = 30
  log_backup_retention_period  = 7
  backup_method                = "physical"
}
```

Import

PostgreSQL backup plan can be imported using the db_instance_id#plan_id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan.example postgres-ckwcgdf1#0ba458c5-2468-5804-ab50-ca556e88c6a3
```
