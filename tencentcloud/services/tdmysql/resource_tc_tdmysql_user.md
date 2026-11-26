Provides a resource to manage a TDSQL-C for MySQL user.

Example Usage

```hcl
resource "tencentcloud_tdmysql_user" "example" {
  instance_id = "tdsqlc-c1q7xi8l"
  user_name   = "tf-example-user"
  host        = "10.0.0.1"
  password    = "TF@example-12345"
  description = "created by terraform"

  timeouts {
    create = "15m"
    delete = "15m"
  }
}
```

Import

TDSQL-C for MySQL user can be imported using the composite id `instance_id#user_name#host`, e.g.

```
terraform import tencentcloud_tdmysql_user.example tdsqlc-c1q7xi8l#tf-example-user#10.0.0.1
```
