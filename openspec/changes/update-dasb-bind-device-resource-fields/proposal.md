## Why

The `tencentcloud_dasb_bind_device_resource` resource was created against an older version of the DASB API (`v20191018`). The `BindDeviceResource` API has since added new K8S-related fields (`ManageDimension`, `ManageAccountId`, `ManageAccount`, `ManageKubeconfig`, `Namespace`, `Workload`) and the `DescribeDevices` API returns additional fields on the `Device` struct (`DomainName`). However, the current Go SDK vendored in this provider does **not yet include** these new fields in `BindDeviceResourceRequest`. The resource must be updated once the SDK is upgraded.

## What Changes

- **SDK gap identified**: `BindDeviceResourceRequest` in the vendored SDK (`v20191018`) is missing `ManageDimension`, `ManageAccountId`, `ManageAccount`, `ManageKubeconfig`, `Namespace`, and `Workload` fields — implementation is **blocked until SDK is upgraded**.
- Add new optional schema fields to `tencentcloud_dasb_bind_device_resource` for the K8S cluster managed account dimension parameters (post SDK upgrade).
- Update Create and Update functions to pass the new fields to `BindDeviceResource`.
- Update Read function to map `DomainName` from `DescribeDevices` response into state.
- Fix Read function to correctly dereference `*uint64` pointer for `device_id_set` items.
- Fix Read function: `domain_id` should only be set once (not inside the per-device loop).

## Capabilities

### New Capabilities
- `dasb-bind-device-resource-k8s-fields`: Extend `tencentcloud_dasb_bind_device_resource` with K8S managed account dimension fields and map additional read-back fields from `DescribeDevices`.

### Modified Capabilities

## Impact

- File: `tencentcloud/services/bh/resource_tc_dasb_bind_device_resource.go`
- File: `tencentcloud/services/bh/service_tencentcloud_dasb.go` (potentially, for `DescribeDasbDeviceByResourceId` query improvement)
- Requires SDK upgrade: `github.com/tencentcloud/tencentcloud-sdk-go` must include the new `BindDeviceResourceRequest` fields before full implementation.
