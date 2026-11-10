Provides a resource to manage a VOD AIGC API Token for a specific sub application.

An AIGC API Token is an independent credential used when calling VOD AIGC-related APIs. One sub application can hold multiple tokens. The token value is returned by the cloud on creation and is treated as a sensitive value. The optional `ext_info` field manages the token's ExtInfo (IP whitelist, model whitelist and name) as a JSON string via `UpdateAigcApiToken`.

~> **NOTE:** The cloud side has an approximately 30-second data-sync delay before a newly created or deleted token becomes visible via `DescribeAigcApiTokens`. The provider polls the list internally, so `terraform apply` may take up to a few minutes.

~> **NOTE:** `api_token` is marked `Sensitive` and is still stored in state in plain text. Protect your Terraform state accordingly (remote state encryption, restricted access, etc.).

~> **NOTE:** The `ext_info` value is passed through to the cloud API as a JSON string. The cloud side may normalize the JSON (key order, whitespace), which can cause a persistent diff on the immediately following plan. Use `jsonencode()` to build the value in a stable form.

Example Usage

```hcl
# create sub application
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}

# create api token
resource "tencentcloud_vod_aigc_api_token" "example" {
  sub_app_id = tencentcloud_vod_sub_application.example.sub_app_id
}

output "aigc_api_token" {
  value     = tencentcloud_vod_aigc_api_token.example.api_token
  sensitive = true
}
```

## Update ExtInfo

```hcl
resource "tencentcloud_vod_aigc_api_token" "example" {
  sub_app_id  = tencentcloud_vod_sub_application.example.sub_app_id
  action_type = "Overwrite"
  ext_info = jsonencode({
    name  = "tf-aigc-token"
    ip    = { "11.12.112.90" : true, "15.12.112.90" : true }
    model = { "deepseek-v4-pro" : true }
  })
}
```

Import

VOD AIGC API Token can be imported using the composite id `sub_app_id#api_token`, e.g.

```
terraform import tencentcloud_vod_aigc_api_token.example 1500066380#<YOUR TOKEN>
```
