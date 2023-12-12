Provides a resource to create a dcdb account_privileges

Example Usage

```hcl
resource "tencentcloud_dcdb_account_privileges" "account_privileges" {
  instance_id = "%s"
  account {
		user = "tf_test"
		host = "%s"
  }
  global_privileges = ["SHOW DATABASES","SHOW VIEW"]
  database_privileges {
		privileges = ["SELECT","INSERT","UPDATE","DELETE","CREATE"]
		database = "tf_test_db"
  }

  table_privileges {
		database = "tf_test_db"
		table = "tf_test_table"
		privileges = ["SELECT","INSERT","UPDATE","DELETE","CREATE"]

  }
```

Import

dcdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_account_privileges.account_privileges instanceId#userName#host#dbName#tabName#viewName#colName
```