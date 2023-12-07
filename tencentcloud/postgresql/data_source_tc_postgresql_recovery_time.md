Use this data source to query detailed information of postgresql recovery_time

Example Usage

```hcl
data "tencentcloud_postgresql_recovery_time" "recovery_time" {
  db_instance_id = local.pgsql_id
}
```