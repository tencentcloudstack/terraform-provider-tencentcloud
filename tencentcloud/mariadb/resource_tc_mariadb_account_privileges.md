Provides a resource to create a mariadb account_privileges

Example Usage

```hcl
resource "tencentcloud_mariadb_account_privileges" "account_privileges" {
  instance_id = "tdsql-9vqvls95"
  accounts {
		user = "keep-modify-privileges"
		host = "127.0.0.1"
  }
  global_privileges = ["ALTER", "CREATE", "DELETE", "SELECT", "UPDATE", "DROP"]
}
```

Import

mariadb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_account_privileges.account_privileges account_privileges_id
```