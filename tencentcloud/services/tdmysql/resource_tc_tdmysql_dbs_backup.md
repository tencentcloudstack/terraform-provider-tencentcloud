Provides a resource to create a TDSQL Boundless(tdmysql) manual backup set.

Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup" "example" {
  instance_id   = "tdsqlc-mysql-abc12345"
  backup_method = "physical"
  backup_type   = "full"
  backup_name   = "tf-example-backup"
}
```

Import

TDSQL Boundless(tdmysql) manual backup set can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup.example 12345
```
