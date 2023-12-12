Provides a resource to create a bi user_role

Example Usage

```hcl
resource "tencentcloud_bi_user_role" "user_role" {
  area_code    = "+83"
  email        = "1055000000@qq.com"
  phone_number = "13470010000"
  role_id_list = [
    10629359,
  ]
  user_id   = "100032767426"
  user_name = "keep-iac-test"
}
```

Import

bi user_role can be imported using the id, e.g.

```
terraform import tencentcloud_bi_user_role.user_role user_id
```