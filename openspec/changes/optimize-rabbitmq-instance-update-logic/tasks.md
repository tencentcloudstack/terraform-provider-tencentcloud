## 1. API Research and Analysis

- [x] 1.1 Research TDMQ API capabilities for parameter updates (node_spec, node_num, storage_size, band_width, enable_public_access)
- [x] 1.2 Test if ModifyRabbitMQVipInstance API supports additional parameters beyond cluster_name and tags
- [x] 1.3 Search for alternative APIs (e.g., ModifyCluster, ModifyNetwork) that might support parameter updates
- [x] 1.4 Document which parameters can be updated via API and which require instance recreation
- [x] 1.5 Confirm ForceNew requirements for immutable parameters (zone_ids, vpc_id, subnet_id, cluster_version)

## 2. Code Implementation

- [x] 2.1 Remove unnecessary parameter restrictions from resourceTencentCloudTdmqRabbitmqVipInstanceUpdate function
- [x] 2.2 Add ForceNew property to immutable parameters in schema (zone_ids, vpc_id, subnet_id, cluster_version)
- [x] 2.3 Implement support for updating enable_public_access parameter (ForceNew approach due to API limitation)
- [x] 2.4 Implement support for updating band_width parameter (ForceNew approach due to API limitation)
- [x] 2.5 Implement support for updating node_spec parameter (ForceNew approach due to API limitation)
- [x] 2.6 Implement support for updating node_num parameter (ForceNew approach due to API limitation)
- [x] 2.7 Implement support for updating storage_size parameter (ForceNew approach due to API limitation)
- [x] 2.8 Add comprehensive error handling with helper.Retry() for update operations (already in place)
- [x] 2.9 Add detailed logging using tccommon.LogElapsed() for update operations (already in place)
- [x] 2.10 Implement state consistency checks using tccommon.InconsistentCheck() after update (already in place)
- [x] 2.11 Add support for updating multiple parameters simultaneously in a single apply operation (Terraform handles this automatically)

## 3. Testing

- [ ] 3.1 Write test case for updating enable_public_access (false -> true)
- [ ] 3.2 Write test case for updating enable_public_access (true -> false)
- [ ] 3.3 Write test case for updating band_width
- [ ] 3.4 Write test case for updating node_spec
- [ ] 3.5 Write test case for updating node_num
- [ ] 3.6 Write test case for updating storage_size
- [ ] 3.7 Write test case for updating multiple parameters simultaneously
- [ ] 3.8 Write test case for updating cluster_name (existing functionality)
- [ ] 3.9 Write test case for updating resource_tags (existing functionality)
- [ ] 3.10 Write test case for backward compatibility with existing resources
- [ ] 3.11 Run unit tests to verify update logic
- [ ] 3.12 Run acceptance tests with TF_ACC=1 to verify API integration

## 4. Documentation

- [x] 4.1 Update resource schema descriptions to clarify which parameters are updatable
- [x] 4.2 Update resource schema descriptions to clarify which parameters require recreation (ForceNew)
- [x] 4.3 Update website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown with new update capabilities
- [x] 4.4 Add examples demonstrating parameter updates in documentation
- [x] 4.5 Add examples demonstrating recreation workflow for immutable parameters
- [x] 4.6 Update resource_tc_tdmq_rabbitmq_vip_instance.md example file with new update scenarios

## 5. Code Quality and Validation

- [ ] 5.1 Run go fmt on resource_tc_tdmq_rabbitmq_vip_instance.go to ensure code formatting
- [ ] 5.2 Run go fmt on resource_tc_tdmq_rabbitmq_vip_instance_test.go to ensure code formatting
- [ ] 5.3 Run make doc command to generate updated documentation
- [ ] 5.4 Verify no schema breaking changes (only Optional fields added)
- [ ] 5.5 Verify backward compatibility with existing configurations
- [ ] 5.6 Run make test to ensure all tests pass
- [ ] 5.7 Run make testacc to run acceptance tests if TENCENTCLOUD_SECRET_ID/KEY are set

## 6. Final Verification

- [ ] 6.1 Review code changes against design document decisions
- [ ] 6.2 Verify all spec requirements are implemented
- [ ] 6.3 Verify error handling covers all specified scenarios
- [ ] 6.4 Verify logging covers all specified scenarios
- [ ] 6.5 Verify state consistency is maintained
- [ ] 6.6 Verify backward compatibility is preserved
- [ ] 6.7 Perform end-to-end testing of update operations
- [ ] 6.8 Prepare summary of changes for changelog
