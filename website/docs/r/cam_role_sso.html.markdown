---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_sso"
sidebar_current: "docs-tencentcloud-resource-cam_role_sso"
description: |-
  Provides a resource to create a CAM-ROLE-SSO(Only support OIDC).
---

# tencentcloud_cam_role_sso

Provides a resource to create a CAM-ROLE-SSO(Only support OIDC).

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `client_ids` - (Required, Set: [`String`]) Client ids.
* `identity_key` - (Required, String) Sign the public key. Base64 encryption is required.
* `identity_url` - (Required, String) Identity provider URL.
* `name` - (Required, String) The name of resource.
* `auto_rotate_key` - (Optional, Int) OIDC public key auto-rotation switch. Enum values: 0 (disabled), 1 (enabled). Default value: 0.
* `description` - (Optional, String) The description of resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CAM-ROLE-SSO(Only support OIDC) can be imported using the `name`, e.g.

```
terraform import tencentcloud_cam_role_sso.example tf_example
```

