Provides a resource to create a dcdb account

Example Usage

```hcl
resource "tencentcloud_dcdb_account" "account" {
	instance_id = "tdsqlshard-kkpoxvnv"
	user_name = "mysql"
	host = "127.0.0.1"
	password = "===password==="
	read_only = 0
	description = "this is a test account"
	max_user_connections = 10
}

```
Import

dcdb account can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_account.account account_id
```