Provides a resource to create a cynosdb account_privileges

Example Usage

```hcl
resource "tencentcloud_cynosdb_account_privileges" "account_privileges" {
  cluster_id   = "cynosdbmysql-bws8h88b"
  account_name = "test"
  host         = "%"
  global_privileges = [
    "CREATE",
    "DROP",
    "ALTER",
    "CREATE TEMPORARY TABLES",
    "CREATE VIEW"
  ]
  database_privileges {
    db = "users"
    privileges = [
      "DROP",
      "REFERENCES",
      "INDEX",
      "CREATE VIEW",
      "INSERT",
      "EVENT"
    ]
  }
  table_privileges {
    db         = "users"
    table_name = "tb_user_name"
    privileges = [
      "ALTER",
      "REFERENCES",
      "SHOW VIEW"
    ]
  }
}
```

Import

cynosdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_account_privileges.account_privileges account_privileges_id
```