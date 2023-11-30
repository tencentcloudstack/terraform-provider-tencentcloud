Use this data source to query detailed information of cam oidc_config

Example Usage

```hcl
data "tencentcloud_cam_oidc_config" "oidc_config" {
  name = "cls-kzilgv5m"
}

output "identity_key" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_key
}

output "identity_url" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_url
}

```