Provides a resource to manage user privileges of a TDSQL Boundless instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_user_privilege_config_config" "example" {
  instance_id = "tdsql3-5baee8df"
  host        = "localhost"
  user_name   = "tf-example-user"
  db_name     = "*"
  privileges  = ["SELECT", "INSERT"]
}
```

Import

TDSQL Boundless user privilege can be imported using the composite id instance_id#host#user_name#db_name, e.g.

```
terraform import tencentcloud_tdmysql_user_privilege_config_config.example tdsql3-5baee8df#localhost#tf-example-user#*
```
