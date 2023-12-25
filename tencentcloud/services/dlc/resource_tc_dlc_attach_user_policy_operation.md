Provides a resource to create a dlc attach_user_policy_operation

Example Usage

```hcl
resource "tencentcloud_dlc_attach_user_policy_operation" "attach_user_policy_operation" {
  user_id = "100032676511"
  policy_set {
		database = "test_iac_keep"
		catalog = "DataLakeCatalog"
		table = ""
		operation = "ASSAYER"
		policy_type = "DATABASE"
		function = ""
		view = ""
		column = ""
		source = "USER"
		mode = "COMMON"
  }
}
```
