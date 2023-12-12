Provides a resource to create a mariadb account

Example Usage

```hcl
resource "tencentcloud_mariadb_account" "account" {
	instance_id = "tdsql-4pzs5b67"
	user_name   = "account-test"
	host        = "10.101.202.22"
	password    = "Password123."
	read_only   = 0
	description = "desc"
}

```
Import

mariadb account can be imported using the instance_id#user_name#host, e.g.
```
$ terraform import tencentcloud_mariadb_account.account tdsql-4pzs5b67#account-test#10.101.202.22
```