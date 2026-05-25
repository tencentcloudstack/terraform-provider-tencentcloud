Use this data source to query the runtime metadata of the TencentCloud
provider, including the active region, client version, mux stack mode and
basic credential presence flags. It does not call any cloud API and can be
used as a lightweight liveness probe for the SDKv2 + framework muxed
provider binary.

Example Usage

```hcl
data "tencentcloud_provider_runtime" "this" {}

output "stack_mode" {
  value = data.tencentcloud_provider_runtime.this.stack_mode
}
```
