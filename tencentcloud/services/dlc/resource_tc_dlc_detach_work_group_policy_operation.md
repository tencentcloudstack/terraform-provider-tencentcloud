Provides a resource to create a dlc detach_work_group_policy_operation

Example Usage

```hcl
resource "tencentcloud_dlc_detach_work_group_policy_operation" "detach_work_group_policy_operation" {
  work_group_id = 23184
  policy_set {
    database = "test_iac_keep"
    catalog = "DataLakeCatalog"
    table = ""
    operation = "ASSAYER"
    policy_type = "DATABASE"
    re_auth = false
    source = "WORKGROUP"
    mode = "COMMON"
    operator = "100032669045"
    id = 102535
  }
}
```
