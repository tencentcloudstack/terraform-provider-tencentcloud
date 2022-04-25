---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_oidc_sso"
sidebar_current: "docs-tencentcloud-resource-cam_oidc_sso"
description: |-
  Provides a resource to create a CAM-OIDC-SSO.
---

# tencentcloud_cam_oidc_sso

Provides a resource to create a CAM-OIDC-SSO.

## Example Usage

```hcl
resource "tencentcloud_cam_oidc_sso" "foo" {
  authorization_endpoint = "https://login.microsoftonline.com/.../oauth2/v2.0/authorize"
  client_id              = "..."
  identity_key           = "..."
  identity_url           = "https://login.microsoftonline.com/.../v2.0"
  mapping_filed          = "name"
  response_mode          = "form_post"
  response_type          = "id_token"
  scope                  = ["openid", "email"]
}
```

## Argument Reference

The following arguments are supported:

* `authorization_endpoint` - (Required) Authorization request Endpoint, OpenID Connect identity provider authorization address. Corresponds to the value of the `authorization_endpoint` field in the Openid-configuration provided by the Enterprise IdP.
* `client_id` - (Required) Client ID, the client ID registered with the OpenID Connect identity provider.
* `identity_key` - (Required) The signature public key requires base64_encode. Verify the public key signed by the OpenID Connect identity provider ID Token. For the security of your account, we recommend that you rotate the signed public key regularly.
* `identity_url` - (Required) Identity provider URL. OpenID Connect identity provider identity.Corresponds to the value of the `issuer` field in the Openid-configuration provided by the Enterprise IdP.
* `mapping_filed` - (Required) Map field names. Which field in the IdP's id_token maps to the user name of the subuser, usually the sub or name field.
* `response_mode` - (Required) Authorize the request Forsonse mode. Authorization request return mode, form_post and frogment two optional modes, recommended to select form_post mode.
* `response_type` - (Required) Authorization requests The Response type, with a fixed value id_token.
* `scope` - (Optional) Authorize the request Scope. openid; email; profile; Authorization request information scope. The default is required openid.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CAM-OIDC-SSO can be imported using the client_id or any string which can identifier resource, e.g.

```
$ terraform import tencentcloud_cam_oidc_sso.foo xxxxxxxxxxx
```

