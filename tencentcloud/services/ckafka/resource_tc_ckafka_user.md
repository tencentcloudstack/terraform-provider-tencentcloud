Provides a resource to create a Ckafka user.

Example Usage

Ckafka User

```hcl
resource "tencentcloud_ckafka_user" "example" {
  instance_id  = "ckafka-7k5nbnem"
  account_name = "tf-example"
  password     = "Password@123"
}
```

Import

Ckafka user can be imported using the instance_id#account_name, e.g.

```
$ terraform import tencentcloud_ckafka_user.example ckafka-7k5nbnem#tf-example
```
