## Why

The Tencent Cloud Global Accelerator 2 (GA2) service exposes a full lifecycle of APIs for managing endpoint groups attached to a global accelerator listener — `CreateEndpointGroup`, `DescribeEndpointGroups`, `ModifyEndpointGroup`, and `DeleteEndpointGroups`. The provider does not yet expose any GA2 resources, so users currently cannot manage GA2 endpoint groups as infrastructure-as-code. The newly upgraded SDK (`ga2 v1.3.102`, namespace `v20250115`) now ships the required APIs, unblocking this work.

## What Changes

- Add new general (CRUD) resource `tencentcloud_ga2_endpoint_group` under a new `ga2` service package.
- Resource schema fields mirror the `CreateEndpointGroup` request parameters one-to-one (top-level fields and the nested `endpoint_group_configuration` block).
- Composite resource ID format: `GlobalAcceleratorId#ListenerId#EndpointGroupId` (the three IDs that uniquely identify an endpoint group instance).
- Async API handling: after `CreateEndpointGroup`, `ModifyEndpointGroup`, and `DeleteEndpointGroups` return a `TaskId`, poll `DescribeTaskResult` until `Status == "SUCCESS"`.
- Add service-layer helpers `DescribeGa2EndpointGroupById` and `Ga2DescribeTaskResult` in `service_tencentcloud_ga2.go`.
- Register a new client entry `UseGa2V20250115Client` in `tencentcloud/connectivity/client.go`.
- Register the resource in `provider.go`.
- Generate `.md` documentation and `_test.go` acceptance test file using project naming conventions.

## Capabilities

### New Capabilities
- `ga2-endpoint-group-resource`: Full CRUD lifecycle for a GA2 endpoint group bound to a `(global_accelerator_id, listener_id)` tuple, including async task polling on Create/Update/Delete.

### Modified Capabilities

## Impact

- New files:
  - `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.go`
  - `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.md`
  - `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group_test.go`
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`
- Modified files:
  - `tencentcloud/connectivity/client.go` — add `UseGa2V20250115Client`
  - `tencentcloud/provider.go` — register `tencentcloud_ga2_endpoint_group`
- SDK dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2 v1.3.102` (namespace `v20250115`). Confirmed present in module cache; `go mod vendor` will materialize it once we import it from new code.
