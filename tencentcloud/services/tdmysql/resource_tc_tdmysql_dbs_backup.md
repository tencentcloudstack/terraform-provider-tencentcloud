Provides a resource to create a TDSQL Boundless(tdmysql) manual backup set.

Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup" "example" {
  instance_id   = "tdsql3-ibpyogg7tg"
  backup_method = "physical"
  backup_type   = "full"
  backup_name   = "tf-example"

  timeouts {
    create = "60m"
  }
}
```

Import

TDSQL Boundless(tdmysql) manual backup set can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup.example 94403665169021
```
