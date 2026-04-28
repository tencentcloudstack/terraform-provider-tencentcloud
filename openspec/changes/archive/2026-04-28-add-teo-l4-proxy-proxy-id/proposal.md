## Why

The `tencentcloud_teo_l4_proxy` resource currently does not expose `proxy_id` as a schema attribute. The `ProxyId` is returned by the CreateL4Proxy API response and is also available in the DescribeL4Proxy API response, but users cannot reference it in their Terraform configurations (e.g., for use in other resources or outputs). Adding `proxy_id` as a computed attribute makes the L4 proxy instance ID accessible to users.

## What Changes

- Add a computed attribute `proxy_id` (TypeString) to the `tencentcloud_teo_l4_proxy` resource schema
- Set `proxy_id` from `response.Response.ProxyId` in the create function (already captured internally, just needs schema exposure)
- Set `proxy_id` from `respData.ProxyId` in the read function (data already returned by DescribeL4Proxy API)
- Mark `ddos_protection_config` and its sub-fields as deprecated (adding `Deprecated` and `Computed: true` attributes)

## Capabilities

### New Capabilities

- `l4-proxy-proxy-id`: Expose the `proxy_id` computed attribute on the `tencentcloud_teo_l4_proxy` resource, mapping to the `ProxyId` field from the CreateL4Proxy/DescribeL4Proxy cloud API responses.

### Modified Capabilities

- `l4-proxy-ddos-protection-config-deprecated`: Mark `ddos_protection_config` and its sub-fields (`level_mainland`, `max_bandwidth_mainland`, `level_overseas`) as deprecated from version 1.82.90. The field remains functional for backward compatibility but users are advised to stop using it.

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go`: Add `proxy_id` to schema, set value in both create and read functions; mark `ddos_protection_config` as deprecated with `Deprecated` and `Computed: true`
- `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go`: Add test coverage for `proxy_id` attribute (Create, Read, NotFound, Schema tests)
- `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md`: Update documentation with `proxy_id` attribute and output example
