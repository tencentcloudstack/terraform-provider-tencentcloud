Provides a resource to create a dlc detach_user_policy_operation

Example Usage

```hcl
resource "tencentcloud_dlc_detach_user_policy_operation" "detach_user_policy_operation" {
 user_id = 100032676511
 policy_set {
   database = "test_iac_keep"
   catalog = "DataLakeCatalog"
   table = ""
   operation = "ASSAYER"
   policy_type = "DATABASE"
   re_auth = false
   source = "USER"
   mode = "COMMON"
   operator = "100032669045"
   id = 102533
 }
}
```
