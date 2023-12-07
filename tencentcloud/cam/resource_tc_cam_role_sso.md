Provides a resource to create a CAM-ROLE-SSO (Only support OIDC).

Example Usage

```hcl
resource "tencentcloud_cam_role_sso" "foo" {
	name="tf_cam_role_sso"
	identity_url="https://login.microsoftonline.com/.../v2.0"
	identity_key="..."
	client_ids=["..."]
	description="this is a description"
}
```

Import

CAM-ROLE-SSO can be imported using the `name`, e.g.

```
$ terraform import tencentcloud_cam_role_sso.foo "test"
```