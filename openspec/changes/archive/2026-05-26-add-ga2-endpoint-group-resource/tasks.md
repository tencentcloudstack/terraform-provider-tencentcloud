## 1. SDK Vendor Sync

- [x] 1.1 After importing `ga2/v20250115` from new code, run `go mod vendor` to materialize the package under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`.

## 2. Connectivity Client Registration

- [x] 2.1 Add import alias `ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"` to `tencentcloud/connectivity/client.go`.
- [x] 2.2 Add private cached field `ga2v20250115Conn *ga2v20250115.Client` on `TencentCloudClient`.
- [x] 2.3 Add `UseGa2V20250115Client()` method following the `UseIgtmV20231024Client` pattern (lazy-init with `NewClientProfile(300)` and `LogRoundTripper`).

## 3. Service Layer

- [x] 3.1 Create `tencentcloud/services/ga2/service_tencentcloud_ga2.go` declaring `Ga2Service` struct (matches `IgtmService` pattern).
- [x] 3.2 Implement `DescribeGa2EndpointGroupById(ctx, gaId, listenerId, egId)` — paginated `DescribeEndpointGroups` with `Limit=100`, filter `endpoint-group-id`. Returns `*ga2.EndpointGroupConfigurationSet, error`. Treats not-found as `(nil, nil)`. Wraps SDK call in `resource.Retry(tccommon.ReadRetryTimeout)`.
- [x] 3.3 Implement `WaitForGa2TaskFinish(ctx, taskId)` — polls `DescribeTaskResult` with `resource.Retry(tccommon.WriteRetryTimeout*2)` until `Status == "SUCCESS"`.

## 4. Resource Implementation

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.go` with `ResourceTencentCloudGa2EndpointGroup()` schema — top-level: `global_accelerator_id` (Required+ForceNew), `listener_id` (Required+ForceNew), `endpoint_group_type` (Required+ForceNew), `endpoint_group_configuration` (Required, TypeList MaxItems=1 with Resource Elem mirroring `EndpointGroupConfiguration`); computed: `endpoint_group_id`.
- [x] 4.2 Define nested sub-schemas `endpointConfigurationsSchema()` (Elem of inner `endpoint_configurations`) and `portOverridesSchema()` (Elem of inner `port_overrides`) to match the SDK structs `EndpointConfigurations` and `PortOverride`.
- [x] 4.3 Implement `Create`: build `CreateEndpointGroupRequest`, invoke with retry, validate `Response.EndpointGroupId != nil`, then `WaitForGa2TaskFinish(ctx, *Response.TaskId)`, then `d.SetId("<ga>#<listener>#<eg>")`, then call Read.
- [x] 4.4 Implement `Read`: split ID into 3 parts, call `DescribeGa2EndpointGroupById`; if nil set ID empty; otherwise `_ = d.Set(...)` for every present field with nil-guards. Map nested `endpoint_configurations` and `port_overrides`.
- [x] 4.5 Implement `Update`: split ID, build `ModifyEndpointGroupRequest` populated from `endpoint_group_configuration` flat fields when changed; invoke with retry; `WaitForGa2TaskFinish`; call Read. Skip API call entirely if no mutable args changed.
- [x] 4.6 Implement `Delete`: split ID, build `DeleteEndpointGroupsRequest` with single-element `EndpointGroupIds`, invoke with retry, `WaitForGa2TaskFinish`.

## 5. Helper Build/Flatten Functions

- [x] 5.1 `buildEndpointGroupConfiguration(rawList []interface{}) *ga2.EndpointGroupConfiguration`.
- [x] 5.2 `buildEndpointConfigurations(rawList []interface{}) []*ga2.EndpointConfigurations`.
- [x] 5.3 `buildPortOverrides(rawList []interface{}) []*ga2.PortOverride`.
- [x] 5.4 `flattenEndpointConfigurations([]*ga2.EndpointConfigurations) []map[string]interface{}`.
- [x] 5.5 `flattenPortOverrides([]*ga2.PortOverride) []map[string]interface{}`.

## 6. Provider Registration

- [x] 6.1 Add import for the new `ga2` package in `tencentcloud/provider.go`.
- [x] 6.2 Register `"tencentcloud_ga2_endpoint_group": ga2.ResourceTencentCloudGa2EndpointGroup()` in the resources map (alphabetically sorted slot).

## 7. Documentation

- [x] 7.1 Create `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.md` with example HCL (mirroring `resource_tc_config_compliance_pack.md` structure: opening line, Example Usage, Import section).

## 8. Tests

- [x] 8.1 Create `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group_test.go` with `TestAccTencentCloudGa2EndpointGroupResource_basic`. Steps: create + check id/endpoint_group_id; update name + description; import.

## 9. Verification

- [x] 9.1 `go build ./tencentcloud/services/ga2/` — confirm package compiles.
- [x] 9.2 `go build ./tencentcloud/connectivity/` — confirm new client method compiles.
- [x] 9.3 `go build ./tencentcloud/` — confirm full provider compiles.
