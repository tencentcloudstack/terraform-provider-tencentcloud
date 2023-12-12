Provides a resource to create a mysql database

Example Usage

```hcl
resource "tencentcloud_mysql_database" "database" {
  instance_id        = "cdb-i9xfdf7z"
  db_name            = "for_tf_test"
  character_set_name = "utf8"
}
```

Import

mysql database can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_database.database instanceId#dbName
```