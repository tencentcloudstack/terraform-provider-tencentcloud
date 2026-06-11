---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_provider_runtime"
sidebar_current: "docs-tencentcloud-data_source-provider_runtime"
description: |-
  Use this data source to query the runtime metadata of the TencentCloud
provider, including the active region, client version, mux stack mode and
basic credential presence flags. It does not call any cloud API and can be
used as a lightweight liveness probe for the SDKv2 + framework muxed
provider binary.
---

# tencentcloud_provider_runtime

Use this data source to query the runtime metadata of the TencentCloud
provider, including the active region, client version, mux stack mode and
basic credential presence flags. It does not call any cloud API and can be
used as a lightweight liveness probe for the SDKv2 + framework muxed
provider binary.

## Example Usage

```hcl
data "tencentcloud_provider_runtime" "this" {}

output "stack_mode" {
  value = data.tencentcloud_provider_runtime.this.stack_mode
}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `client_version` - The X-TC-RequestClient version reported by the provider.
* `cos_domain` - The COS root domain the provider currently uses. Empty string when the default public COS domain is used.
* `domain` - The API root domain the provider currently uses. Empty string when the default TencentCloud public domain is used.
* `id` - Synthetic id, equal to the current region.
* `protocol` - The API request protocol the provider currently uses (typically `HTTPS`).
* `region` - The region the provider is currently configured to use.
* `secret_id_present` - Whether a non-empty SecretId is currently configured on the provider. **This field only indicates whether SecretId is configured, NOT whether the credential is valid.** No real credential value is ever exposed by this data source.
* `stack_mode` - Always returns "sdkv2+framework" for this provider, indicating that resources can be implemented in either stack and are served via tf5muxserver.


