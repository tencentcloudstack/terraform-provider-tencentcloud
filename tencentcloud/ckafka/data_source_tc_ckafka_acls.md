Use this data source to query detailed acl information of Ckafka

Example Usage

```hcl
data "tencentcloud_ckafka_acls" "foo" {
  instance_id   = "ckafka-f9ife4zz"
  resource_type = "TOPIC"
  resource_name = "topic-tf-test"
  host          = "2"
}
```