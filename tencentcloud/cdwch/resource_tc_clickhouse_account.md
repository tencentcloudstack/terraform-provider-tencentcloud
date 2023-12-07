Provides a resource to create a clickhouse account

Example Usage

```hcl
resource "tencentcloud_clickhouse_account" "account" {
  instance_id = "cdwch-xxxxxx"
  user_name = "test"
  password = "xxxxxx"
  describe = "xxxxxx"
}
```

Import

clickhouse account can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_account.account ${instance_id}#${user_name}
```