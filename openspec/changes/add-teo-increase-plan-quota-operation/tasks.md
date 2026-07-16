## 1. Resource Implementation

- [ ] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go` with schema definition (plan_id, quota_type, quota_number as Required+ForceNew; deal_name as Computed), Create function calling IncreasePlanQuota API with retry, and empty Read/Delete functions

## 2. Unit Tests

- [ ] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation_test.go` with unit tests using gomonkey to mock the IncreasePlanQuota API call, covering successful creation, API error, and nil response scenarios

## 3. Provider Registration

- [ ] 3.1 Register the new resource `tencentcloud_teo_increase_plan_quota` in `tencentcloud/provider.go` (add to ResourcesMap)
- [ ] 3.2 Register the new resource in `tencentcloud/provider.md` (add to the TEO resources list)

## 4. Documentation

- [ ] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md` with description, Example Usage, and Import sections following the existing TEO resource doc format

## 5. Validation

- [ ] 5.1 Run `go test -gcflags=all=-l` on the unit test file to verify tests pass
- [ ] 5.2 Verify the code compiles successfully (via `go build` check in later steps)