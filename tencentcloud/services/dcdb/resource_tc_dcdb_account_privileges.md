Provides a resource to create a dcdb account privileges

Example Usage

```hcl
resource "tencentcloud_dcdb_account_privileges" "example" {
  instance_id = "tdsqlshard-xxxxxxxx"
  account {
    user = "tf_test"
    host = "127.0.0.1"
  }
  global_privileges = ["SHOW DATABASES", "SHOW VIEW"]
  database_privileges {
    privileges = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
    database   = "tf_test_db"
  }

  table_privileges {
    database   = "tf_test_db"
    table      = "tf_test_table"
    privileges = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
  }
}
```

Import

dcdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_account_privileges.example instanceId#userName#host#dbName#tabName#viewName#colName
```