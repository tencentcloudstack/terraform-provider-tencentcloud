Provides a resource to create a Ckafka Acl.

Example Usage

Ckafka Acl

```hcl
resource "tencentcloud_ckafka_acl" "foo" {
  instance_id     = "ckafka-f9ife4zz"
  resource_type   = "TOPIC"
  resource_name   = "topic-tf-test"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "*"
  principal       = tencentcloud_ckafka_user.foo.account_name
}
```

Import

Ckafka acl can be imported using the instance_id#permission_type#principal#host#operation_type#resource_type#resource_name, e.g.

```
$ terraform import tencentcloud_ckafka_acl.foo ckafka-f9ife4zz#ALLOW#test#*#WRITE#TOPIC#topic-tf-test
```