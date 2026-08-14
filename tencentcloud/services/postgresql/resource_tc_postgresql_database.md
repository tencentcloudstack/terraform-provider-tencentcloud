Provides a resource to manage a TencentDB for PostgreSQL database within a DB instance.

Example Usage

```hcl
resource "tencentcloud_postgresql_database" "example" {
  db_instance_id  = "postgres-6fego161"
  database_name   = "test_db"
  database_owner  = "tcuser"
  encoding        = "UTF8"
  collate         = "C"
  ctype           = "C"
}
```

Import

PostgreSQL database can be imported using the composite ID `db_instance_id#database_name`, e.g.

```
terraform import tencentcloud_postgresql_database.example postgres-6fego161#test_db
```
