resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "test"
  password     = "test1234"
}

data "tencentcloud_ckafka_users" "foo" {
  instance_id  = tencentcloud_ckafka_user.foo.instance_id
  account_name = tencentcloud_ckafka_user.foo.account_name
}

resource "tencentcloud_ckafka_acl" foo {
  instance_id     = "ckafka-f9ife4zz"
  resource_type   = "TOPIC"
  resource_name   = "topic-tf-test"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "10.10.10.0"
  principal       = tencentcloud_ckafka_user.foo.account_name
}

data "tencentcloud_ckafka_acls" "foo" {
  instance_id   = tencentcloud_ckafka_acl.foo.instance_id
  resource_type = tencentcloud_ckafka_acl.foo.resource_type
  resource_name = tencentcloud_ckafka_acl.foo.resource_name
}