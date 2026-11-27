Provides a resource to manage the TDSQL Boundless DBS backup policy configuration.

Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup_policy_config" "example" {
  instance_id = "tdsql3-ibpyogg7tg"
  backup_policy = {
    backup_start_time     = "16:00"
    backup_end_time       = "24:00"
    backup_method         = "physical"
    enable_full           = 1
    enable_log            = 1
    full_retention_period = 7
    log_retention_period  = 7
    period_time           = "0,1,2,4,5,6"
    storage_type          = "COS"
  }
}
```

Import

TDSQL Boundless DBS backup policy config can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup_policy_config.example tdsql3-ibpyogg7tg
```
