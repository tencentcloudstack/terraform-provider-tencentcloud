Provides a resource to create a CAM-OIDC-SSO.

Example Usage

```hcl
resource "tencentcloud_cam_oidc_sso" "foo" {
	authorization_endpoint="https://login.microsoftonline.com/.../oauth2/v2.0/authorize"
	client_id="..."
	identity_key="..."
	identity_url="https://login.microsoftonline.com/.../v2.0"
	mapping_filed="name"
	response_mode="form_post"
	response_type="id_token"
	scope=["openid", "email"]
}
```

Import

CAM-OIDC-SSO can be imported using the client_id or any string which can identifier resource, e.g.

```
$ terraform import tencentcloud_cam_oidc_sso.foo xxxxxxxxxxx
```