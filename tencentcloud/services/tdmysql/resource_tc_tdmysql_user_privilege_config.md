Provides a resource to manage user privileges of a TDSQL Boundless instance.

Example Usage

```hcl
resource "tencentcloud_tdmysql_user_privilege_config" "example" {
  instance_id = "tdsql3-ibpyogg7tg"
  user_name   = "tf-example"
  host        = "10.0.0.1"

  global_privileges = [
    "INSERT",
    "SELECT",
    "ALTER",
    "CREATE",
    "DELETE"
  ]
}
```

Import

TDSQL Boundless user privilege config can be imported using the composite id instance_id#user_name#host, e.g.

```
terraform import tencentcloud_tdmysql_user_privilege_config.example tdsql3-ibpyogg7tg#tf-example#10.0.0.1
```
