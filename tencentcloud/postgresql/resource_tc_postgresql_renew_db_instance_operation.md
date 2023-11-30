Provides a resource to create a postgresql renew_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_renew_db_instance_operation" "renew_db_instance_operation" {
  db_instance_id = tencentcloud_postgresql_instance.oper_test_PREPAID.id
  period = 1
  auto_voucher = 0
}
```