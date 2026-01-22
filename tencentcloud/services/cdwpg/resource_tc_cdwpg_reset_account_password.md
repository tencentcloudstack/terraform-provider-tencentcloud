Provides a resource to reset cdwpg account password

Example Usage

```hcl
resource "tencentcloud_cdwpg_reset_account_password" "cdwpg_reset_account_password" {
	instance_id = "cdwpg-zpiemnyd"
	user_name = "dbadmin"
	new_password = "testpassword"
}
```

Import

cdwpg reset account password can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_reset_account_password.cdwpg_account cdwpg_reset_account_password_id
```
