Provide a resource to configure kubernetes cluster authentication info.

~> **NOTE:** Only available for cluster version >= 1.20

Example Usage

Use TKE default issuer and jwks_uri

```hcl
resource "tencentcloud_kubernetes_auth_attachment" "example" {
  cluster_id                           = "cls-53c7589g"
  use_tke_default                      = true
  auto_create_discovery_anonymous_auth = true
}
```

Use custom issuer and jwks_uri

```hcl
resource "tencentcloud_kubernetes_auth_attachment" "example" {
  cluster_id                           = "cls-53c7589g"
  use_tke_default                      = false
  jwks_uri                             = "https://cls-53c7589g.ccs.tencent-cloud.com/openid/v1/jwks"
  issuer                               = "https://cls-53c7589g.ccs.tencent-cloud.com"
  auto_create_discovery_anonymous_auth = false
}
```

Use OIDC Config

```
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

Import

TKE cluster authentication can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_auth_attachment.example cls-fp5o961e
```
