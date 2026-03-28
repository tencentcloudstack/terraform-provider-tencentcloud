## 1. Code Modifications

- [ ] 1.1 Update `immutableArgs` list in `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` function to remove mutable fields: node_spec, node_num, storage_size, band_width, enable_public_access, cluster_version, enable_create_default_ha_mirror_queue
- [ ] 1.2 Add update logic for `node_spec` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.3 Add update logic for `node_num` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.4 Add update logic for `storage_size` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.5 Add update logic for `band_width` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.6 Add update logic for `enable_public_access` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.7 Add update logic for `cluster_version` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.8 Add update logic for `enable_create_default_ha_mirror_queue` field with `d.HasChange()` check and API request parameter setting
- [ ] 1.9 Ensure all mutable field updates are consolidated into a single `ModifyRabbitMQVipInstance` API call when `needUpdate` is true

## 2. Testing

- [ ] 2.1 Write unit tests for `node_spec` update scenario: verify API is called with correct parameters and state is updated
- [ ] 2.2 Write unit tests for `node_num` update scenario: verify API is called with correct parameters for both increase and decrease
- [ ] 2.3 Write unit tests for `storage_size` update scenario: verify API is called with correct parameters
- [ ] 2.4 Write unit tests for `band_width` update scenario: verify API is called with correct parameters
- [ ] 2.5 Write unit tests for `enable_public_access` toggle scenario: verify API is called with correct boolean values
- [ ] 2.6 Write unit tests for `cluster_version` update scenario: verify API is called with correct version string
- [ ] 2.7 Write unit tests for `enable_create_default_ha_mirror_queue` toggle scenario: verify API is called with correct boolean values
- [ ] 2.8 Write unit tests for multiple field updates: verify single API call includes all changed fields
- [ ] 2.9 Write unit tests for immutable field updates: verify error is returned for zone_ids, vpc_id, subnet_id, pay_mode, time_span, auto_renew_flag
- [ ] 2.10 Run acceptance tests with `TF_ACC=1` to verify all update scenarios work correctly against real Tencent Cloud API

## 3. Documentation

- [ ] 3.1 Update field descriptions in `resource_tc_tdmq_rabbitmq_vip_instance.go` schema to indicate which fields are updatable and which are immutable
- [ ] 3.2 Update `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` (generated via `make doc`) to reflect the update capability changes
- [ ] 3.3 Update `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md` example file to demonstrate field update operations

## 4. Code Quality and Validation

- [ ] 4.1 Run `go fmt` on the modified resource file to ensure consistent formatting
- [ ] 4.2 Run `go vet` to check for potential issues
- [ ] 4.3 Run `golangci-lint` to verify code quality standards
- [ ] 4.4 Ensure all existing tests still pass after modifications
- [ ] 4.5 Verify backward compatibility: run `terraform plan` on existing resources to ensure no unexpected changes are shown

## 5. API Verification

- [ ] 5.1 Test `node_spec` update via API to verify Tencent Cloud actually supports this field modification
- [ ] 5.2 Test `node_num` update via API to verify Tencent Cloud actually supports this field modification
- [ ] 5.3 Test `storage_size` update via API to verify Tencent Cloud actually supports this field modification
- [ ] 5.4 Test `band_width` update via API to verify Tencent Cloud actually supports this field modification
- [ ] 5.5 Test `enable_public_access` update via API to verify Tencent Cloud actually supports this field modification
- [ ] 5.6 Test `cluster_version` update via API to verify Tencent Cloud actually supports this field modification (if API returns error, add back to immutable list)
- [ ] 5.7 Test `enable_create_default_ha_mirror_queue` update via API to verify Tencent Cloud actually supports this field modification (if API returns error, add back to immutable list)
- [ ] 5.8 Document any fields that are not actually supported by the API and adjust implementation accordingly

## 6. Error Handling and Edge Cases

- [ ] 6.1 Add error handling for API failures during update operations with clear error messages
- [ ] 6.2 Add validation for invalid field values (e.g., negative node_num, invalid node_spec)
- [ ] 6.3 Test error scenario: API returns validation error for invalid node_spec value
- [ ] 6.4 Test error scenario: API returns quota exceeded error for storage_size increase
- [ ] 6.5 Test error scenario: API returns error for cluster_version downgrade (if supported)
- [ ] 6.6 Test retry logic for transient network errors during update operations

## 7. Final Verification

- [ ] 7.1 Perform full integration test: create instance, update multiple fields, verify all changes are applied
- [ ] 7.2 Test state consistency: after update, verify state matches actual resource by running `terraform refresh`
- [ ] 7.3 Test rollback scenario: if update fails, verify state is not corrupted and instance remains unchanged
- [ ] 7.4 Test concurrent updates: verify Terraform correctly detects and handles state conflicts
- [ ] 7.5 Review all code changes against design document to ensure alignment
- [ ] 7.6 Verify all requirements from specs are satisfied by implementation
