---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_auth_attachment"
sidebar_current: "docs-tencentcloud-resource-kubernetes_auth_attachment"
description: |-
  Provide a resource to configure kubernetes cluster authentication info.
---

# tencentcloud_kubernetes_auth_attachment

Provide a resource to configure kubernetes cluster authentication info.

~> **NOTE:** Only available for cluster version >= 1.20

## Example Usage

### Use TKE default issuer and jwks_uri

```hcl
resource "tencentcloud_kubernetes_auth_attachment" "example" {
  cluster_id                           = "cls-53c7589g"
  use_tke_default                      = true
  auto_create_discovery_anonymous_auth = true
}
```

### Use custom issuer and jwks_uri

```hcl
resource "tencentcloud_kubernetes_auth_attachment" "example" {
  cluster_id                           = "cls-53c7589g"
  use_tke_default                      = false
  jwks_uri                             = "https://cls-53c7589g.ccs.tencent-cloud.com/openid/v1/jwks"
  issuer                               = "https://cls-53c7589g.ccs.tencent-cloud.com"
  auto_create_discovery_anonymous_auth = false
}
```

### Use OIDC Config

```hcl
resource "tencentcloud_kubernetes_auth_attachment" "example" {
  cluster_id                              = "cls-oof3l9ks"
  use_tke_default                         = true
  auto_create_discovery_anonymous_auth    = true
  auto_create_oidc_config                 = true
  auto_install_pod_identity_webhook_addon = true
}

data "tencentcloud_cam_oidc_config" "oidc_config" {
  name = tencentcloud_kubernetes_auth_attachment.example.cluster_id
}

output "identity_key" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_key
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) ID of clusters.
* `auto_create_client_id` - (Optional, Set: [`String`]) Creating ClientId of the identity provider.
* `auto_create_discovery_anonymous_auth` - (Optional, Bool) If set to `true`, the rbac rule will be created automatically which allow anonymous user to access `/.well-known/openid-configuration` and `/openid/v1/jwks`.
* `auto_create_oidc_config` - (Optional, Bool) Creating an identity provider.
* `auto_install_pod_identity_webhook_addon` - (Optional, Bool) Creating the PodIdentityWebhook component. if `auto_create_oidc_config` is true, this field must set true.
* `issuer` - (Optional, String) Specify service-account-issuer. If `use_tke_default` is set to `true`, please do not set this field.
* `jwks_uri` - (Optional, String) Specify service-account-jwks-uri. If `use_tke_default` is set to `true`, please do not set this field.
* `use_tke_default` - (Optional, Bool) If set to `true`, the `issuer` and `jwks_uri` will be generated automatically by tke, please do not set `issuer` and `jwks_uri`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `tke_default_issuer` - The default issuer of tke. If `use_tke_default` is set to `true`, this parameter will be set to the default value.
* `tke_default_jwks_uri` - The default jwks_uri of tke. If `use_tke_default` is set to `true`, this parameter will be set to the default value.


## Import

TKE cluster authentication can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_auth_attachment.example cls-fp5o961e
```

