## 1. Schema Changes

- [x] 1.1 In `tencentcloud/services/dc/resource_tc_dcx.go`, remove `ForceNew: true` from the `bandwidth` schema field in `ResourceTencentCloudDcxInstance()`. Keep the field as `TypeInt`, `Optional`, `Computed`.

## 2. Update Function Changes

- [x] 2.1 In `resourceTencentCloudDcxInstanceUpdate` in `tencentcloud/services/dc/resource_tc_dcx.go`, add a `d.HasChange("bandwidth")` handling block that builds a `dc.NewModifyDirectConnectTunnelAttributeRequest()`, sets `DirectConnectTunnelId` from `d.Id()`, sets `Bandwidth` via `helper.IntInt64()` from `d.Get("bandwidth")`, and calls `ModifyDirectConnectTunnelAttribute` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` error wrapping — consistent with the existing `name` update handling block.

## 3. Unit Test

- [x] 3.1 In `tencentcloud/services/dc/resource_tc_dcx_test.go`, add a unit test using gomonkey mocks for the `ModifyDirectConnectTunnelAttribute` API call covering the `bandwidth` update path (verify the request carries the expected `Bandwidth` value and `DirectConnectTunnelId`).

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/dc/resource_tc_dcx.md` to reflect that `bandwidth` is now updatable (not ForceNew), including an example showing bandwidth update usage.

## 5. Validation

- [x] 5.1 Verify the code compiles successfully (build/lint handled by the downstream CI flow; no `go build`/`go vet` to be run in this phase).
