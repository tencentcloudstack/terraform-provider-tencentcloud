## 1. Code Analysis and Preparation

- [x] 1.1 Review current Update function implementation in resource_tc_tdmq_rabbitmq_vip_instance.go
- [x] 1.2 Identify immutable fields that should become updatable (node_spec, node_num, storage_size, auto_renew_flag, band_width, enable_public_access)
- [x] 1.3 Identify truly immutable fields (zone_ids, vpc_id, subnet_id, cluster_version)
- [x] 1.4 Research Tencent Cloud API documentation for ModifyRabbitMQVipInstance to understand supported parameters
- [x] 1.5 Create backup of current implementation before making changes

## 2. Update Function Implementation

- [x] 2.1 Remove updatable fields from immutableArgs list in Update function
- [x] 2.2 Add node_spec update logic with API parameter mapping
- [x] 2.3 Add node_num update logic with API parameter mapping
- [x] 2.4 Add storage_size update logic with API parameter mapping
- [x] 2.5 Add auto_renew_flag update logic with API parameter mapping
- [x] 2.6 Add band_width update logic with API parameter mapping
- [x] 2.7 Add enable_public_access toggle logic with API parameter mapping
- [x] 2.8 Implement instance status validation before update operations
- [x] 2.9 Implement multi-field update logic with proper ordering
- [x] 2.10 Add retry logic for all API calls using resource.Retry

## 3. State Validation and Error Handling

- [x] 3.1 Add state validation after update by calling Read function
- [x] 3.2 Implement field-by-field validation to ensure all updates were applied correctly
- [x] 3.3 Add error handling for partial update failures
- [x] 3.4 Enhance error messages with detailed context (field name, API error, suggested actions)
- [x] 3.5 Add debug logging for all update operations (request/response details)
- [x] 3.6 Implement proper error propagation for update failures

## 4. Testing

- [ ] 4.1 Write unit test for single field update (node_spec)
- [ ] 4.2 Write unit test for single field update (node_num)
- [ ] 4.3 Write unit test for single field update (storage_size)
- [ ] 4.4 Write unit test for single field update (auto_renew_flag)
- [ ] 4.5 Write unit test for single field update (band_width)
- [ ] 4.6 Write unit test for single field update (enable_public_access toggle)
- [ ] 4.7 Write unit test for multi-field update
- [ ] 4.8 Write unit test for immutable field update rejection
- [ ] 4.9 Write unit test for update failure scenarios (instance not in running state)
- [ ] 4.10 Write unit test for retry logic on transient failures
- [ ] 4.11 Run acceptance tests with TF_ACC=1 for RabbitMQ instance updates
- [ ] 4.12 Test backward compatibility with existing resources
- [ ] 4.13 Test idempotent updates (applying same configuration multiple times)

## 5. Documentation Updates

- [x] 5.1 Update resource_tc_tdmq_rabbitmq_vip_instance.md example file with update examples
- [x] 5.2 Run `make doc` to regenerate website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown
- [x] 5.3 Verify generated documentation includes updated field descriptions
- [x] 5.4 Add usage examples for updating each new updatable field
- [x] 5.5 Add usage example for multi-field updates
- [x] 5.6 Document which fields are updatable and which remain immutable
- [x] 5.7 Add notes about instance status requirements for updates

## 6. Code Quality and Formatting

- [ ] 6.1 Run `go fmt` on modified files
- [ ] 6.2 Run `go vet` to check for common issues
- [ ] 6.3 Run `go build` to ensure code compiles successfully
- [ ] 6.4 Check for code style consistency with existing codebase
- [ ] 6.5 Ensure all error messages follow existing patterns
- [ ] 6.6 Verify logging follows existing log format patterns
- [ ] 6.7 Check for memory leaks and resource cleanup

## 7. Verification and Validation

- [ ] 7.1 Verify that existing tests still pass
- [ ] 7.2 Verify that existing resources continue to work without changes
- [ ] 7.3 Verify that terraform plan shows no changes for unchanged resources
- [ ] 7.4 Verify that terraform apply successfully updates updatable fields
- [ ] 7.5 Verify that terraform apply rejects updates to immutable fields with clear error messages
- [ ] 7.6 Verify that state is correctly refreshed after updates
- [ ] 7.7 Verify that documentation is complete and accurate

## 8. Final Review and Cleanup

- [x] 8.1 Review all code changes against design document
- [x] 8.2 Review all code changes against specification requirements
- [x] 8.3 Ensure all TODO comments are resolved
- [x] 8.4 Remove any debugging or temporary code
- [x] 8.5 Update CHANGELOG.md with feature description
- [x] 8.6 Verify backward compatibility with existing Terraform configurations
- [ ] 8.7 Perform final integration test with real Tencent Cloud account
