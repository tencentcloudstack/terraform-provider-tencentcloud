Provides a resource to create a postgresql disisolate_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_disisolate_db_instance_operation" "disisolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
  period = 1
  auto_voucher = false
}
```