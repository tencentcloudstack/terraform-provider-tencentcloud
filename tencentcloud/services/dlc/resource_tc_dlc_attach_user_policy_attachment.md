Provides a resource to create a DLC attach user policy attachment

~> **NOTE:** `policy_set` only supports attaching exactly one policy per resource.

Example Usage

```hcl
resource "tencentcloud_dlc_attach_user_policy_attachment" "example" {
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

DLC attach user policy attachment can be imported using the composite id (`user_id#policy_id`), e.g.

```
terraform import tencentcloud_dlc_attach_user_policy_attachment.example 100032676511#1400538864-1721880558-fh1lTgYD
```
