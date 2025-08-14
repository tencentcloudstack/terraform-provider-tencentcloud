Provides a resource to create a DLC detach user policy operation

Example Usage

```hcl
resource "tencentcloud_dlc_detach_user_policy_operation" "example" {
  user_id = 100032676511
  policy_set {
    database    = "tf_example_db"
    catalog     = "DataLakeCatalog"
    table       = "tf_example_table"
    operation   = "ASSAYER"
    policy_type = "DATABASE"
    source      = "USER"
    mode        = "COMMON"
  }
}
```
