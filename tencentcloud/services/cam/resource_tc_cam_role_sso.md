Provides a resource to create a CAM-ROLE-SSO(Only support OIDC).

Example Usage

```hcl
resource "tencentcloud_cam_role_sso" "example" {
  name            = "tf_example"
  identity_url    = "https://login.microsoftonline.com/.../v2.0"
  identity_key    = "baz****"
  client_ids      = ["61adcf00620c31e3ddbc9546"]
  description     = "this is a description"
  auto_rotate_key = 1
}
```

Import

CAM-ROLE-SSO(Only support OIDC) can be imported using the `name`, e.g.

```
terraform import tencentcloud_cam_role_sso.example tf_example
```
