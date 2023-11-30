Provides a resource to create a mongodb instance_account

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_account" "instance_account" {
  instance_id = "cmgo-lxaz2c9b"
  user_name = "test_account"
  password = "xxxxxxxx"
  mongo_user_password = "xxxxxxxxx"
  user_desc = "test account"
  auth_role {
    mask = 0
    namespace = "*"
  }
}
```