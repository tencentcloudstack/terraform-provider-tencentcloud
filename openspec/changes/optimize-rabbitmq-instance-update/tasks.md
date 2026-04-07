## 1. Code Modification

- [x] 1.1 Modify immutableArgs list in resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function
  - Remove `auto_renew_flag`, `band_width`, and `enable_public_access` from the immutableArgs array
  - Keep 9 truly immutable parameters: `zone_ids`, `vpc_id`, `subnet_id`, `node_spec`, `node_num`, `storage_size`, `enable_create_default_ha_mirror_queue`, `time_span`, `pay_mode`, `cluster_version`

- [x] 1.2 Add update logic for auto_renew_flag
  - Check `d.HasChange("auto_renew_flag")` in update function
  - Set `request.AutoRenewFlag = helper.Bool(v.(bool))` when change detected
  - Set `needUpdate = true` flag when parameter changes

- [x] 1.3 Add update logic for enable_public_access
  - Check `d.HasChange("enable_public_access")` in update function
  - Set `request.EnablePublicAccess = helper.Bool(v.(bool))` when change detected
  - Set `needUpdate = true` flag when parameter changes

- [x] 1.4 Add update logic for band_width
  - Check `d.HasChange("band_width")` in update function
  - Set `request.Bandwidth = helper.IntUint64(v.(int))` when change detected
  - Set `needUpdate = true` flag when parameter changes

- [x] 1.5 Enhance error messages for immutable parameters
  - Update error message format to include parameter name and guidance
  - Change from: `"argument \`%s\` cannot be changed"`
  - Change to: `"argument \`%s\` cannot be changed after instance creation. Please recreate the instance if you need to modify this parameter."`
  - Apply this to all 9 remaining immutable parameters

- [x] 1.6 Add async state waiting logic after update
  - Implement retry mechanism using `resource.Retry()` after successful API call
  - Poll instance status using `service.DescribeTdmqRabbitmqVipInstanceByFilter()` with timeout `tccommon.ReadRetryTimeout*10`
  - Wait for status to change from "Updating" to "Running" or "Success"
  - Handle retryable errors with `tccommon.RetryError()`
  - Return non-retryable error for unexpected states (e.g., "Failed", "Rollback")

- [x] 1.7 Verify update parameter passing to API
  - Ensure all changed parameters (`cluster_name`, `resource_tags`, `auto_renew_flag`, `enable_public_access`, `band_width`) are passed to `ModifyRabbitMQVipInstance` API
  - Ensure only one API call is made for multiple parameter changes (not separate calls for each field)
  - Maintain backward compatibility with existing `cluster_name` and `resource_tags` update logic

## 2. Testing

- [x] 2.1 Update acceptance test for enable_public_access update
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateEnablePublicAccess`
  - Create instance with public access disabled, then enable it via update
  - Verify instance is not recreated
  - Verify `enable_public_access` value is updated in state

- [x] 2.2 Update acceptance test for band_width update
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateBandwidth`
  - Create instance with initial bandwidth, then modify it via update
  - Verify instance is not recreated
  - Verify `band_width` value is updated in state

- [x] 2.3 Update acceptance test for auto_renew_flag update
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateAutoRenewFlag`
  - Create instance with auto renew disabled, then enable it via update
  - Verify instance is not recreated
  - Verify `auto_renew_flag` value is updated in state

- [x] 2.4 Update acceptance test for multiple parameter updates
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateMultipleFields`
  - Create instance, then update `enable_public_access`, `band_width`, `cluster_name`, and `resource_tags` simultaneously
  - Verify instance is not recreated
  - Verify all parameters are updated in a single operation

- [x] 2.5 Update acceptance test for async state waiting
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateAsyncWait`
  - Update instance and verify state polling occurs
  - Verify operation waits for status to stabilize
  - Verify no premature return before status update

- [x] 2.6 Update acceptance test for immutable parameter errors
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateImmutableParams`
  - Attempt to update each immutable parameter (`zone_ids`, `vpc_id`, `subnet_id`, `node_spec`, `node_num`, `storage_size`, `enable_create_default_ha_mirror_queue`, `time_span`, `pay_mode`, `cluster_version`)
  - Verify error message includes parameter name and guidance
  - Verify instance is not modified

- [x] 2.7 Update acceptance test for idempotency
  - Add test case: `TestAccTdmqRabbitmqVipInstance_updateIdempotent`
  - Run apply multiple times with same configuration
  - Verify subsequent applies detect no changes
  - Verify no API calls are made for operations with no changes

- [x] 2.8 Run acceptance tests locally
  - Execute `TF_ACC=1 go test -v -run TestAccTdmqRabbitmqVipInstance_update`
  - Verify all new test cases pass
  - Verify existing test cases still pass
  - Fix any failures before proceeding

## 3. Documentation

- [x] 3.1 Update resource example file
  - Update `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md`
  - Add example showing update of `enable_public_access`, `band_width`, and `auto_renew_flag`
  - Document that these parameters are now updateable after creation
  - Update notes about remaining immutable parameters

- [x] 3.2 Generate website documentation
  - Run `make doc` command to auto-generate documentation
  - Verify `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` is updated
  - Check that `enable_public_access`, `band_width`, and `auto_renew_flag` are documented as updateable
  - Check that remaining immutable parameters are correctly marked

- [x] 3.3 Verify documentation completeness
  - Review generated documentation for accuracy
  - Ensure all updateable parameters are clearly labeled
  - Ensure immutable parameters have clear error messages documented
  - Verify examples demonstrate update capabilities

## 4. Code Quality

- [x] 4.1 Format code with go fmt
  - Run `go fmt ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - Verify all code adheres to Go formatting standards
  - Commit formatting changes if needed

- [x] 4.2 Verify backward compatibility
  - Review all changes to ensure no breaking changes
  - Ensure existing functionality is preserved
  - Verify state schema is unchanged (no schema modifications required)

- [x] 4.3 Add comments for complex logic
  - Add inline comments explaining async state waiting logic
  - Document retry mechanism and timeout values
  - Clarify immutability decisions for remaining parameters

- [x] 4.4 Review error handling
  - Ensure all API errors are properly handled
  - Verify error messages are clear and actionable
  - Check that retry logic handles transient errors correctly

## 5. Verification

- [x] 5.1 Manual verification of update functionality
  - Manually create an instance via Terraform
  - Update `enable_public_access` and verify changes are applied
  - Update `band_width` and verify changes are applied
  - Update `auto_renew_flag` and verify changes are applied
  - Verify instance is not recreated during updates

- [x] 5.2 Manual verification of immutable parameter enforcement
  - Attempt to update an immutable parameter (e.g., `zone_ids`)
  - Verify error message is clear and actionable
  - Verify instance is not modified

- [x] 5.3 Manual verification of async waiting
  - Update an instance and monitor Terraform output
  - Verify operation waits for status to stabilize
  - Verify no premature return occurs

- [x] 5.4 Final review of all changes
  - Review all modified files for correctness
  - Ensure all tasks are completed
  - Verify no outstanding TODO comments remain
  - Prepare for submission
