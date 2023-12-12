Provides a resource to create a postgresql modify_switch_time_period_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_modify_switch_time_period_operation" "modify_switch_time_period_operation" {
  db_instance_id = local.pgsql_id
  switch_tag = 0
}
```