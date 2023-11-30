Provides a resource to create a bi project_user_role

~> **NOTE:** You cannot use `tencentcloud_bi_user_role` and `tencentcloud_bi_project_user_role` at the same time to modify the `phone_number` and `email` of the same user.

Example Usage

```hcl
resource "tencentcloud_bi_project_user_role" "project_user_role" {
  area_code    = "+86"
  project_id   = 11015030
  role_id_list = [10629453]
  email        = "123456@qq.com"
  phone_number = "13130001000"
  user_id      = "100024664626"
  user_name    = "keep-cam-user"
}
```

Import

bi project_user_role can be imported using the id, e.g.

```
terraform import tencentcloud_bi_project_user_role.project_user_role projectId#userId
```