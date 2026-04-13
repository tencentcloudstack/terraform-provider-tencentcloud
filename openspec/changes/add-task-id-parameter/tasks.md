## 1. Verify API Support

- [x] 1.1 Verify ProcessMedia API returns TaskId in vendor SDK models
- [x] 1.2 Verify EditMedia API returns TaskId in vendor SDK models
- [x] 1.3 Verify ProcessLiveStream API returns TaskId in vendor SDK models
- [x] 1.4 Verify StartFlow API returns TaskId in vendor SDK models (if applicable)
- [x] 1.5 Verify WithdrawsWatermark API returns TaskId in vendor SDK models (if applicable)
- [x] 1.6 Identify all MPS operation resources that need task_id field based on API response verification

## 2. Implement task_id in resource_tc_mps_process_media_operation.go

- [x] 2.1 Add task_id computed field to resource schema
- [x] 2.2 Set task_id in Create function after API call succeeds
- [x] 2.3 Verify task_id matches resource ID value

## 3. Implement task_id in resource_tc_mps_edit_media_operation.go

- [x] 3.1 Add task_id computed field to resource schema
- [x] 3.2 Set task_id in Create function after API call succeeds
- [x] 3.3 Verify task_id matches resource ID value

## 4. Implement task_id in resource_tc_mps_process_live_stream_operation.go

- [x] 4.1 Add task_id computed field to resource schema
- [x] 4.2 Set task_id in Create function after API call succeeds
- [x] 4.3 Verify task_id matches resource ID value

## 5. Implement task_id in other MPS operation resources (if applicable)

- [x] 5.1 Add task_id to resource_tc_mps_start_flow_operation.go (if API returns TaskId)
- [x] 5.2 Add task_id to resource_tc_mps_withdraws_watermark_operation.go (if API returns TaskId)
- [x] 5.3 Add task_id to any other identified MPS operation resources (if API returns TaskId)

## 6. Update Documentation Files

- [x] 6.1 Update resource_tc_mps_process_media_operation.md with task_id field
- [x] 6.2 Update resource_tc_mps_edit_media_operation.md with task_id field
- [x] 6.3 Update resource_tc_mps_process_live_stream_operation.md with task_id field
- [x] 6.4 Update documentation for any other modified operation resources
- [x] 6.5 Ensure all documentation includes proper description and usage examples

## 7. Add Test Coverage

- [x] 7.1 Add test assertion for task_id in resource_tc_mps_process_media_operation_test.go
- [x] 7.2 Add test assertion for task_id in resource_tc_mps_edit_media_operation_test.go
- [x] 7.3 Add test assertion for task_id in resource_tc_mps_process_live_stream_operation_test.go
- [x] 7.4 Add test coverage for any other modified operation resources
- [x] 7.5 Verify task_id is not empty in all test cases

## 8. Code Review and Quality Checks

- [x] 8.1 Verify all modified files follow code style guidelines
- [x] 8.2 Verify backward compatibility is maintained
- [x] 8.3 Verify no breaking changes to existing schema
- [x] 8.4 Ensure all computed fields are properly documented
- [x] 8.5 Review that task_id setting follows existing patterns in codebase

## 9. Final Validation

- [x] 9.1 Run acceptance tests for all modified resources
- [x] 9.2 Verify tests pass without errors
- [x] 9.3 Check that task_id field is accessible in Terraform state
- [x] 9.4 Verify task_id can be referenced in other resources or outputs
- [x] 9.5 Confirm all documentation files are updated correctly
