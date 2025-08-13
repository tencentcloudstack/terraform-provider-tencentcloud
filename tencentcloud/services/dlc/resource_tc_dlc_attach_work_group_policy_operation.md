Provides a resource to create a DLC attach work group policy operation

Example Usage

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_operation" "example" {
  work_group_id = 23184
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
