Provides a resource to create a sqlserver restore_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_restore_instance" "restore_instance" {
  instance_id = "mssql-qelbzgwf"
  backup_id   = 3482091273
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "restore_keep_pubsub_db2"
  }
}
```

Import

sqlserver restore_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restore_instance.restore_instance mssql-qelbzgwf#3482091273#keep_pubsub_db2#restore_keep_pubsub_db2
```