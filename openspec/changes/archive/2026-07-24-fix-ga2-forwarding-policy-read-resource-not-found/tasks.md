## 1. Common Helper Function

- [x] 1.1 Enhance `resource_tc_ga2_common.go` with `HandleGa2ReadNotFound` unified function that handles both SDK `ResourceNotFound` error (via `HandleGa2ResourceNotFoundError`) and nil/empty response from the service layer, with `!d.IsNewResource()` guard in both cases.

## 2. Apply to All GA2 Resource Read Methods

- [x] 2.1 Update `resource_tc_ga2_forwarding_policy.go` Read to use `HandleGa2ReadNotFound`
- [x] 2.2 Update `resource_tc_ga2_global_accelerator.go` Read to use `HandleGa2ReadNotFound`
- [x] 2.3 Update `resource_tc_ga2_listener.go` Read to use `HandleGa2ReadNotFound`
- [x] 2.4 Update `resource_tc_ga2_endpoint_group.go` Read to use `HandleGa2ReadNotFound`
- [x] 2.5 Update `resource_tc_ga2_accelerate_area.go` Read to use `HandleGa2ReadNotFound`
- [x] 2.6 Update `resource_tc_ga2_forwarding_rule.go` Read to use `HandleGa2ReadNotFound`

## 3. Verification

- [x] 3.1 Run `go build ./tencentcloud/services/ga2/` to ensure compilation succeeds.
