Provides a resource to manage the TDSQL Boundless DBS backup policy configuration.

Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup_policy_config" "example" {
  instance_id = "tdsql-12345"
  backup_policy = {
    backup_method         = "physical"
    backup_start_time     = "02:00"
    backup_end_time       = "04:00"
    enable_full           = 1
    enable_log            = 1
    full_retention_period = 7
    log_retention_period  = 7
    period_time           = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"
    storage_type          = "COS"
  }
}
```

Import

TDSQL Boundless DBS backup policy config can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup_policy_config.example tdsql-12345
```
