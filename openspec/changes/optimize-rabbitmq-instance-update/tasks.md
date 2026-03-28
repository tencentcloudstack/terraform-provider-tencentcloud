## 1. Code Analysis and Preparation

- [x] 1.1 Review current update logic in `resource_tc_tdmq_rabbitmq_vip_instance.go` (lines 450-523)
- [x] 1.2 Identify immutable fields that should remain unchanged (zone_ids, vpc_id, subnet_id, node_spec, node_num, storage_size, enable_create_default_ha_mirror_queue, band_width, enable_public_access)
- [x] 1.3 Verify API capabilities for auto_renew_flag and time_span updates via API documentation or testing
- [x] 1.4 Create backup branch for the changes

## 2. Update Function Implementation

- [x] 2.1 Remove auto_renew_flag, time_span from immutableArgs array in Update function
- [x] 2.2 Add validation logic for auto_renew_flag updates to check pay_mode = 1 (prepaid only)
- [x] 2.3 Add validation logic for time_span updates to check pay_mode = 1 (prepaid only)
- [x] 2.4 Implement auto_renew_flag update logic with d.HasChange() and API parameter mapping
- [x] 2.5 Implement time_span update logic with d.HasChange() and API parameter mapping
- [x] 2.6 Add needUpdate flag to only call API when changes are detected
- [x] 2.7 Wrap ModifyRabbitMQVipInstance API call in resource.Retry() with WriteRetryTimeout
- [x] 2.8 Add detailed error logging for API request/response bodies
- [x] 2.9 Improve error messages with context about which field caused the failure
- [x] 2.10 Add call to resourceTencentCloudTdmqRabbitmqVipInstanceRead() after successful update to verify state synchronization

## 3. Post-Update State Verification

- [ ] 3.1 Add retry logic wrapper around Read operation after update with ReadRetryTimeout*10
- [ ] 3.2 Implement state comparison logic to verify updated fields reflect new values
- [ ] 3.3 Add logging for state verification process (attempts, success, timeout)
- [ ] 3.4 Handle eventual consistency scenarios where Read returns stale data

## 4. Update Documentation

- [ ] 4.1 Review current documentation in `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
- [ ] 4.2 Update field descriptions to indicate auto_renew_flag and time_span are now updatable
- [ ] 4.3 Add notes about update constraints (prepaid instances only for auto_renew_flag and time_span)
- [ ] 4.4 Add example usage for update operations showing auto_renew_flag and time_span updates
- [ ] 4.5 Document the retry and state verification behavior
- [ ] 4.6 Run `make doc` to regenerate documentation from code comments

## 5. Unit Tests

- [ ] 5.1 Add unit test for auto_renew_flag update on prepaid instance (positive case)
- [ ] 5.2 Add unit test for auto_renew_flag update validation failure on postpaid instance (negative case)
- [ ] 5.3 Add unit test for time_span update on prepaid instance (positive case)
- [ ] 5.4 Add unit test for time_span update validation failure on postpaid instance (negative case)
- [ ] 5.5 Add unit test for differential updates (only changed fields sent to API)
- [ ] 5.6 Add unit test for retry logic on transient API errors
- [ ] 5.7 Add unit test for post-update state verification
- [ ] 5.8 Add unit test for error handling with detailed error messages
- [ ] 5.9 Run `go test ./tencentcloud/services/trabbit/` to verify unit tests pass

## 6. Acceptance Tests

- [ ] 6.1 Add acceptance test for auto_renew_flag update with TF_ACC=1
- [ ] 6.2 Add acceptance test for time_span update with TF_ACC=1
- [ ] 6.3 Add acceptance test for concurrent updates (multiple fields) with TF_ACC=1
- [ ] 6.4 Add acceptance test for update validation errors with TF_ACC=1
- [ ] 6.5 Run full acceptance test suite: `TF_ACC=1 go test ./tencentcloud/services/trabbit/ -v`
- [ ] 6.6 Verify all acceptance tests pass with real Tencent Cloud API

## 7. Backward Compatibility Verification

- [ ] 7.1 Run existing unit tests to ensure no regressions
- [ ] 7.2 Run existing acceptance tests to ensure no regressions
- [ ] 7.3 Test with existing Terraform configurations (without new fields) to ensure they still work
- [ ] 7.4 Verify state migration scenarios (existing resources with no changes applied)

## 8. Code Quality and Linting

- [ ] 8.1 Run `go fmt ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` to format code
- [ ] 8.2 Run `go vet ./tencentcloud/services/trabbit/` to check for code issues
- [ ] 8.3 Run `golangci-lint run` if available in the project
- [ ] 8.4 Ensure all error messages are user-friendly and provide actionable guidance

## 9. Final Verification and Cleanup

- [ ] 9.1 Verify all files are properly formatted with `go fmt`
- [ ] 9.2 Ensure no TODO or FIXME comments left in production code
- [ ] 9.3 Verify logging is at appropriate levels (DEBUG for detailed info, INFO for important events, ERROR for failures)
- [ ] 9.4 Review code comments for clarity and completeness
- [ ] 9.5 Verify documentation matches implementation behavior
- [ ] 9.6 Perform final end-to-end testing with Terraform plan/apply workflow

## 10. Deployment Preparation

- [ ] 10.1 Review all code changes one final time
- [ ] 10.2 Update CHANGELOG.md with details about the new update capabilities
- [ ] 10.3 Create summary of changes for PR description
- [ ] 10.4 Verify all tests pass: `go test ./tencentcloud/services/trabbit/` and `TF_ACC=1 go test ./tencentcloud/services/trabbit/`
- [ ] 10.5 Prepare for code review with clear documentation of what was changed and why
