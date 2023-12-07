Provides a resource to create a mariadb log_file_retention_period

Example Usage

```hcl
resource "tencentcloud_mariadb_log_file_retention_period" "log_file_retention_period" {
  instance_id = "tdsql-4pzs5b67"
  days = "8"
}

```
Import

mariadb log_file_retention_period can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_log_file_retention_period.log_file_retention_period tdsql-4pzs5b67
```