Provides a resource to create a dlc attach_work_group_policy_operation

Example Usage

```hcl
resource "tencentcloud_dlc_attach_work_group_policy_operation" "attach_work_group_policy_operation" {
  work_group_id = 23184
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
