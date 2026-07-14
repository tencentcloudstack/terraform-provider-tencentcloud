Provides a resource to create a DLC attach user policyr attachment

Example Usage

```hcl
resource "tencentcloud_dlc_attach_user_policyr_attachment" "example" {
  user_id      = "100032676511"
  account_type = "TencentAccount"
  policy_set {
    database    = "tf_example_db"
    catalog     = "DataLakeCatalog"
    table       = "tf_example_table"
    operation   = "SELECT"
    policy_type = "TABLE"
  }
}
```

Import

DLC attach user policyr attachment can be imported using the composite id, e.g.

```
terraform import tencentcloud_dlc_attach_user_policyr_attachment.example 100032676511#TencentAccount
```
