Use this data source to query detailed user information of Ckafka

Example Usage

```hcl
data "tencentcloud_ckafka_users" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "test"
}
```