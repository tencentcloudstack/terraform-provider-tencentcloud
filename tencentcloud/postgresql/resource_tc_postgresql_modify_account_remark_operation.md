Provides a resource to create a postgresql modify_account_remark_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_modify_account_remark_operation" "modify_account_remark_operation" {
  db_instance_id = local.pgsql_id
  user_name = "root"
  remark = "hello_world"
}
```