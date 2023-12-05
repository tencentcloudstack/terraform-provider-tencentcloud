Provides a resource to create a sqlserver rollback_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-qelbzgwf"
  time        = "2023-05-23 01:00:00"
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "rollback_pubsub_db3"
  }
}
```

Import

sqlserver rollback_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_rollback_instance.rollback_instance mssql-qelbzgwf#2023-05-23 01:00:00#keep_pubsub_db2#rollback_pubsub_db3
```