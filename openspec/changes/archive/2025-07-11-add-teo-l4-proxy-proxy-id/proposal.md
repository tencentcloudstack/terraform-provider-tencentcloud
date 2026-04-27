## Why

The `tencentcloud_teo_l4_proxy` resource currently does not expose `proxy_id` as a schema attribute. The `proxy_id` is returned by the `CreateL4Proxy` API response and is stored internally as part of the composite resource ID (`zone_id` + `proxy_id`), but users cannot reference `proxy_id` directly in their Terraform configurations. Adding `proxy_id` as a computed attribute allows users to reference the L4 proxy instance ID in other resources or outputs.

## What Changes

- Add `proxy_id` as a computed (`Computed: true`) string attribute to the `tencentcloud_teo_l4_proxy` resource schema
- Set `proxy_id` value in the Create function from the `CreateL4Proxy` API response (`response.Response.ProxyId`)
- Set `proxy_id` value in the Read function from the `DescribeL4Proxy` API response (`respData.ProxyId`)

## Capabilities

### New Capabilities
- `teo-l4-proxy-proxy-id`: Expose the `proxy_id` computed attribute in the `tencentcloud_teo_l4_proxy` resource, allowing users to reference the L4 proxy instance ID returned by the cloud API.

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- Affected file: `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go` (schema definition, create and read functions)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go` (unit tests)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md` (documentation)
- Cloud API: `CreateL4Proxy` response already returns `ProxyId`, no API changes required
- Cloud API: `DescribeL4Proxy` response already returns `ProxyId` in the `L4Proxy` struct, no API changes required
- Backward compatible: Adding a computed attribute does not break existing configurations
