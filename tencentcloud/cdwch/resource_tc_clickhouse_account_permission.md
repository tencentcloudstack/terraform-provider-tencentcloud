Provides a resource to create a clickhouse account_permission

Example Usage

```hcl
resource "tencentcloud_clickhouse_account_permission" "account_permission_all_database" {
  instance_id = "cdwch-xxxxxx"
  cluster = "default_cluster"
  user_name = "user1"
  all_database = true
  global_privileges = ["SELECT", "ALTER"]
}

resource "tencentcloud_clickhouse_account_permission" "account_permission_not_all_database" {
	instance_id = "cdwch-xxxxxx"
  cluster = "default_cluster"
  user_name = "user2"
  all_database = false
  database_privilege_list {
    database_name = "xxxxxx"
    database_privileges = ["SELECT", "ALTER"]
  }
}
```

Import

clickhouse account_permission can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_account_permission.account_permission ${instanceId}#${cluster}#${userName}
```