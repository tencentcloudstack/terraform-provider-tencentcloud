## Why

The `DeleteApplicationProxy` cloud API requires `ZoneId` and `ProxyId` as input parameters. Currently, the `tencentcloud_teo_application_proxy` resource's delete function retrieves these values by parsing the composite ID (`d.Id()`) instead of using `d.Get()`. Per the coding guidelines, when using composite IDs with `FILED_SP` as separator, the read, update, and delete methods should obtain ID component fields from `d.Get()` rather than parsing `d.Id()`. This change aligns the delete operation with the established pattern.

## What Changes

- Modify `resourceTencentCloudTeoApplicationProxyDelete` to read `zone_id` and `proxy_id` from `d.Get()` instead of parsing from `d.Id()`, and pass them as explicit input parameters to the `DeleteApplicationProxy` API request.

## Capabilities

### New Capabilities
- `teo-application-proxy-delete-params`: Adds `zone_id` and `proxy_id` as explicit input parameters for the `DeleteApplicationProxy` API in the `tencentcloud_teo_application_proxy` resource, following the coding guideline of using `d.Get()` over `d.Id()` for composite ID fields.

### Modified Capabilities

## Impact

- Affected file: `tencentcloud/services/teo/resource_tc_teo_application_proxy.go` (delete function)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_application_proxy_test.go` (unit tests)
- Cloud API: `DeleteApplicationProxy` (parameters already exist in SDK, no vendor change needed)
