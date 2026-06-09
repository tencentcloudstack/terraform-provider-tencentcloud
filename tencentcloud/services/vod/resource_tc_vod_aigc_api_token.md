Provides a resource to manage a VOD AIGC API Token for a specific sub application.

An AIGC API Token is an independent credential used when calling VOD AIGC-related APIs. One sub application can hold multiple tokens. The token value is returned by the cloud on creation and is treated as a sensitive value.

~> **NOTE:** The cloud side has an approximately 30-second data-sync delay before a newly created or deleted token becomes visible via `DescribeAigcApiTokens`. The provider polls the list internally, so `terraform apply` may take up to a few minutes.

~> **NOTE:** `api_token` is marked `Sensitive` and is still stored in state in plain text. Protect your Terraform state accordingly (remote state encryption, restricted access, etc.).

Example Usage

```hcl
resource "tencentcloud_vod_aigc_api_token" "example" {
  sub_app_id = 251006666
}

output "aigc_api_token" {
  value     = tencentcloud_vod_aigc_api_token.example.api_token
  sensitive = true
}
```

Import

VOD AIGC API Token can be imported using the composite id `sub_app_id#api_token`, e.g.

```
$ terraform import tencentcloud_vod_aigc_api_token.example 251006666#hGjH1dsBbjUbT9Cu
```
