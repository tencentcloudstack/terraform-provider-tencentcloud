## Why

The `tencentcloud_teo_l4_proxy` resource currently does not expose `proxy_id` as a schema attribute. The `ProxyId` is returned by the CreateL4Proxy API response and is also available in the DescribeL4Proxy API response, but users cannot reference it in their Terraform configurations (e.g., for use in other resources or outputs). Adding `proxy_id` as a computed attribute makes the L4 proxy instance ID accessible to users.

## What Changes

- Add a computed attribute `proxy_id` (TypeString) to the `tencentcloud_teo_l4_proxy` resource schema
- Set `proxy_id` from `response.Response.ProxyId` in the create function (already captured internally, just needs schema exposure)
- Set `proxy_id` from `respData.ProxyId` in the read function (data already returned by DescribeL4Proxy API)

## Capabilities

### New Capabilities

- `l4-proxy-proxy-id`: Expose the `proxy_id` computed attribute on the `tencentcloud_teo_l4_proxy` resource, mapping to the `ProxyId` field from the CreateL4Proxy/DescribeL4Proxy cloud API responses.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go`: Add `proxy_id` to schema, set value in read function
- `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go`: Add test coverage for `proxy_id` attribute
- `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md`: Update documentation with `proxy_id` attribute
