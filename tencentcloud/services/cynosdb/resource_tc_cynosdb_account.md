Provides a resource to create a cynosdb account

Example Usage

```hcl
resource "tencentcloud_cynosdb_account" "account" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  account_name         = "terraform_test"
  account_password     = "Password@1234"
  host                 = "%"
  description          = "terraform test"
  max_user_connections = 2
}
```

Import

cynosdb account can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_account.account account_id
```