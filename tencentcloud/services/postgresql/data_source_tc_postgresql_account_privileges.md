Use this data source to query detailed information of postgresql account privileges

Example Usage

```hcl
data "tencentcloud_postgresql_account_privileges" "example" {
  db_instance_id = "postgres-3hk6b6tj"
  user_name      = "tf_example"
  database_object_set {
    object_name = "postgres"
    object_type = "database"
  }
}
```
