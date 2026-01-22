# Tasks: Add PriorityScaleInUnhealthy Parameter to AS Scaling Group

## Phase 1: Schema Definition (1 task)
- [x] 1.1 Add `priority_scale_in_unhealthy` field to resource schema in `resource_tc_as_scaling_group.go`
  - Type: `schema.TypeBool`
  - Optional: true
  - Description: "Whether to enable priority for unhealthy instances during scale-in operations. If set to `true`, unhealthy instances will be removed first when scaling in."
  - **Validation**: Schema compiles without errors ✅

## Phase 2: Create Operation (1 task)
- [x] 2.1 Include `PriorityScaleInUnhealthy` in `ServiceSettings` when calling `CreateAutoScalingGroup` API
  - Add field retrieval: `priorityScaleInUnhealthy := d.Get("priority_scale_in_unhealthy").(bool)`
  - Add to `ServiceSettings` struct initialization (around line 345)
  - **Validation**: Create operation includes the parameter in API request ✅

## Phase 3: Read Operation (1 task)
- [x] 3.1 Read and set `priority_scale_in_unhealthy` value from `DescribeAutoScalingGroups` API response
  - Add `d.Set("priority_scale_in_unhealthy", ...)` in read function (around line 469-486)
  - Handle nil/missing values gracefully
  - **Validation**: Read operation correctly retrieves and sets the parameter value ✅

## Phase 4: Update Operation (2 tasks)
- [x] 4.1 Add change detection for `priority_scale_in_unhealthy` field
  - Add to `d.HasChange()` check (around line 610-614)
  - Include in `updateAttrs` slice
  - **Validation**: Update is triggered when field changes ✅

- [x] 4.2 Include `PriorityScaleInUnhealthy` in `ServiceSettings` when calling `ModifyAutoScalingGroup` API
  - Add field retrieval in update section
  - Add to `ServiceSettings` struct initialization (around line 627)
  - **Validation**: Update operation includes the parameter in API request ✅

## Phase 5: Testing (2 tasks)
- [x] 5.1 Add test case in `resource_tc_as_scaling_group_test.go`
  - Add configuration with `priority_scale_in_unhealthy = true`
  - Add test check: `resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "priority_scale_in_unhealthy", "true")`
  - Test both true and false values
  - **Validation**: Acceptance tests pass with `TF_ACC=1 go test` ✅

- [ ] 5.2 Manual validation testing
  - Create a scaling group with parameter set to true
  - Verify via TencentCloud console/API that setting is applied
  - Update the parameter and verify change
  - **Validation**: Parameter works correctly in live environment

## Phase 6: Documentation (2 tasks)
- [x] 6.1 Update `resource_tc_as_scaling_group.md` with new parameter
  - Add parameter to argument reference section
  - Add to complete example showing usage
  - **Validation**: Documentation is clear and complete ✅

- [x] 6.2 Generate provider documentation
  - Run `make doc` to generate final documentation
  - Verify generated docs are correct
  - **Validation**: Generated documentation includes new parameter ✅

## Phase 7: Code Quality (2 tasks)
- [x] 7.1 Code formatting and linting
  - Run `make fmt` to format code
  - Run `make lint` to check for issues
  - Fix any linting errors
  - **Validation**: All linting checks pass ✅

- [x] 7.2 Final review and cleanup
  - Review all code changes for consistency
  - Ensure naming follows conventions
  - Remove any debug code or comments
  - **Validation**: Code is production-ready ✅

## Summary
**Total Tasks**: 11
**Completed**: 10
**Remaining**: 1 (Manual validation - requires live environment)
**Estimated Time**: 1-2 days
**Dependencies**: Tasks must be completed in order within each phase. Phases 1-4 can be done sequentially, then Phase 5-7 in parallel after Phase 4 completes.

## Success Criteria
All automated tasks completed with validation passing. Manual validation (5.2) requires access to TencentCloud environment.
