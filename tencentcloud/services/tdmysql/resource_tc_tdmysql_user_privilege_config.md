Provides a resource to manage user privileges of a TDSQL Boundless instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_user_privilege_config" "example" {
  instance_id = "tdsql3-5baee8df"
  user_name   = "tf-example-user"
  host        = "localhost"

  global_privileges = ["SELECT", "INSERT"]

  database_privileges = [
    {
      database   = "tf-example-db"
      privileges = ["SELECT"]
    },
  ]

  table_privileges = [
    {
      database   = "tf-example-db"
      table      = "tf-example-table"
      privileges = ["SELECT", "UPDATE"]
    },
  ]
}
```

Import

TDSQL Boundless user privilege config can be imported using the composite id instance_id#user_name#host, e.g.

```
terraform import tencentcloud_tdmysql_user_privilege_config.example tdsql3-5baee8df#tf-example-user#localhost
```
