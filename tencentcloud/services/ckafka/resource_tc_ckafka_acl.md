Provides a resource to create a Ckafka Acl.

-> **Note:** When creating the ACL, if the cloud API `CreateAcl` returns a `FailedOperation` error (e.g. the instance is being modified or the cluster is busy), the provider will retry up to 3 times with a fixed 5-second interval between each attempt to improve the eventual success rate of ACL creation.

Example Usage

Ckafka Acl

```hcl
resource "tencentcloud_ckafka_user" "example" {
  instance_id  = "ckafka-7k5nbnem"
  account_name = "tf-example"
  password     = "Password@123"
}

resource "tencentcloud_ckafka_acl" "example" {
  instance_id     = "ckafka-7k5nbnem"
  resource_type   = "TOPIC"
  resource_name   = "tf-example-resource"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "*"
  principal       = tencentcloud_ckafka_user.example.account_name
}
```

Import

Ckafka Acl can be imported using the instance_id#permission_type#principal#host#operation_type#resource_type#resource_name, e.g.

```
$ terraform import tencentcloud_ckafka_acl.example ckafka-7k5nbnem#ALLOW#tf-example#*#WRITE#TOPIC#tf-example-resource
```
