Provides a resource to create a mariadb activate_hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_operate_hour_db_instance" "activate_hour_db_instance" {
  instance_id = "tdsql-9vqvls95"
  operate     = "activate"
}
```