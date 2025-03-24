Provides a resource to create a cdwpg cdwpg_account

Example Usage

```hcl
resource "tencentcloud_cdwpg_account" "cdwpg_account" {
	instance_id = "cdwpg-zpiemnyd"
	user_name = "dbadmin"
	new_password = "testpassword"
}
```

Import

cdwpg cdwpg_account can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_account.cdwpg_account cdwpg_account_id
```
