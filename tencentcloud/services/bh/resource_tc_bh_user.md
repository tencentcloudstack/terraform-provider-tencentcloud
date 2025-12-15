Provides a resource to create a BH user

Example Usage

```hcl
resource "tencentcloud_bh_user" "example" {
  user_name = "tf-example"
  real_name = "Terraform"
  phone     = "+86|18991162528"
  email     = "demo@tencent.com"
  auth_type = 0
}
```

Import

BH user can be imported using the id, e.g.

```
terraform import tencentcloud_bh_user.example 2322
```
