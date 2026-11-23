# Resource: tencentcloud_tdmysql_dbs_backup

Provides a resource to create a manual backup set for a TDSQL-C (TDSQL Boundless, tdmysql) instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup" "example" {
  instance_id   = "tdsqlc-abc123"
  backup_method = "physical"
  backup_type   = "full"
  backup_name   = "tf-manual-backup"
}
```

Import

A tdmysql backup set can be imported using the `instance_id#backup_set_id` format, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup.example tdsqlc-abc123#10086
```
