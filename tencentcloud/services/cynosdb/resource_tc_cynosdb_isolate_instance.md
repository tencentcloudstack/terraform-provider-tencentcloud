Provides a resource to create a cynosdb isolate_instance

Example Usage

```hcl
resource "tencentcloud_cynosdb_account" "account" {
  cluster_id           = "cynosdbmysql-bws8h88b"
  account_name         = "terraform_test"
  account_password     = "Password@1234"
  host                 = "%"
  description          = "testx"
  max_user_connections = 2
}
```