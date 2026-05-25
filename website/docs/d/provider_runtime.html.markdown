---
subcategory: "Provider Runtime"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_provider_runtime"
sidebar_current: "docs-tencentcloud-datasource-provider_runtime"
description: |-
  Read-only provider runtime metadata served by the framework stack.
---

# tencentcloud_provider_runtime

Read-only provider runtime metadata. Useful for debugging which region the
provider currently targets and which client version is running. This data
source does **not** call any TencentCloud API.

This is the first data source served by the `terraform-plugin-framework`
stack of the TencentCloud provider. The provider serves both
`terraform-plugin-sdk/v2` and `terraform-plugin-framework` resources via a
single `tf5muxserver`, so users can mix legacy SDKv2 resources and new
framework-based resources/data sources in the same configuration without
any extra setup.

## Example Usage

```hcl
data "tencentcloud_provider_runtime" "this" {}

output "region" {
  value = data.tencentcloud_provider_runtime.this.region
}

output "stack_mode" {
  value = data.tencentcloud_provider_runtime.this.stack_mode
}

output "credential_configured" {
  value = data.tencentcloud_provider_runtime.this.secret_id_present
}
```

## Argument Reference

This data source takes no arguments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Synthetic id, equal to the current region.
* `region` - The region the provider is currently configured to use.
* `client_version` - The `X-TC-RequestClient` version reported by the
  provider (matches the value sent in the SDK user-agent header).
* `stack_mode` - Always `sdkv2+framework`, indicating that resources can be
  implemented in either stack and are served via `tf5muxserver`.
* `protocol` - The API request protocol the provider currently uses
  (typically `HTTPS`).
* `domain` - The API root domain the provider currently uses. Empty string
  when the default TencentCloud public domain is used.
* `cos_domain` - The COS root domain the provider currently uses. Empty
  string when the default public COS domain is used.
* `secret_id_present` - Whether a non-empty `SecretId` is currently
  configured on the provider.
  **Note:** This field only indicates whether `SecretId` is configured,
  **NOT** whether the credential is valid. No real credential value is ever
  exposed by this data source.
