Provides a resource to create a cynosdb account

Example Usage

If host is %

```hcl
resource "tencentcloud_cynosdb_account" "example" {
  cluster_id           = "cynosdbmysql-ddciqx2l"
  account_name         = "tf_example"
  account_password     = "Password@123"
  host                 = "%"
  description          = "remark."
  max_user_connections = 10
}
```

If host is ip

```hcl
resource "tencentcloud_cynosdb_account" "example" {
  cluster_id           = "cynosdbmysql-ddciqx2l"
  account_name         = "tf_example"
  account_password     = "Password@123"
  host                 = "1.1.1.1"
  description          = "remark."
  max_user_connections = 0
}
```

Import

cynosdb account can be imported using the clusterId#accountName#host, e.g.

```
terraform import tencentcloud_cynosdb_account.example cynosdbmysql-ddciqx2l#tf_example#%

or

terraform import tencentcloud_cynosdb_account.example cynosdbmysql-ddciqx2l#tf_example#1.1.1.1
```