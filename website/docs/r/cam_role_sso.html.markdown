---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_sso"
sidebar_current: "docs-tencentcloud-resource-cam_role_sso"
description: |-
  Provides a resource to create a CAM-ROLE-SSO (Only support OIDC).
---

# tencentcloud_cam_role_sso

Provides a resource to create a CAM-ROLE-SSO (Only support OIDC).

## Example Usage

```hcl
resource "tencentcloud_cam_role_sso" "foo" {
  name         = "tf_cam_role_sso"
  identity_url = "https://login.microsoftonline.com/.../v2.0"
  identity_key = "..."
  client_ids   = ["..."]
  description  = "this is a description"
}
```

## Argument Reference

The following arguments are supported:

* `client_ids` - (Required, Set: [`String`]) Client ids.
* `identity_key` - (Required, String) Sign the public key.
* `identity_url` - (Required, String) Identity provider URL.
* `name` - (Required, String) The name of resource.
* `description` - (Optional, String) The description of resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CAM-ROLE-SSO can be imported using the `name`, e.g.

```
$ terraform import tencentcloud_cam_role_sso.foo "test"
```

