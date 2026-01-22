Provides a resource to create an identity center user

Example Usage

```hcl
resource "tencentcloud_identity_center_user" "example" {
  zone_id     = "z-1os7c9tyugct"
  user_name   = "tf-example"
  description = "desc."
}
```

Or

```hcl
resource "tencentcloud_identity_center_user" "example" {
  zone_id      = "z-1os7c9tyugct"
  user_name    = "tf-example"
  description  = "desc."
  first_name   = "FirstName"
  last_name    = "LastName"
  display_name = "DisplayName"
  email        = "example@tencent.com"
  user_status  = "Enabled"
}
```

Import

organization identity center user can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_user.example z-1os7c9tyugct#u-rdvm4xdqi8pr
```
