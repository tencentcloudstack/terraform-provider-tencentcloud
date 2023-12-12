Provides a resource to create a Ckafka user.

Example Usage

Ckafka User

```hcl
resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "tf-test"
  password     = "test1234"
}
```

Import

Ckafka user can be imported using the instance_id#account_name, e.g.

```
$ terraform import tencentcloud_ckafka_user.foo ckafka-f9ife4zz#tf-test
```