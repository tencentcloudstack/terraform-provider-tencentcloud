Provides a resource to manage a TDSQL Boundless(tdmysql) user.

Example Usage

```hcl
resource "tencentcloud_tdmysql_user" "example" {
  instance_id = "tdsql3-ibpyogg7tg"
  user_name   = "tf-example"
  host        = "10.0.0.1"
  password    = "Password@2026"
  description = "created by terraform"

  timeouts {
    create = "3m"
    delete = "3m"
  }
}
```

Import

TDSQL Boundless(tdmysql) user can be imported using the composite id `instance_id#user_name#host`, e.g.

```
terraform import tencentcloud_tdmysql_user.example tdsql3-ibpyogg7tg#tf-example#10.0.0.1
```
