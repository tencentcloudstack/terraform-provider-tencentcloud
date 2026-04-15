## 1. Resource Schema Definition

- [x] 1.1 Create resource schema definition for `tencentcloud_teo_function_v2` with all required fields (zone_id, name, content) and optional fields (remark)
- [x] 1.2 Add computed fields to schema (function_id, domain, create_time, update_time)
- [x] 1.3 Configure zone_id and name as ForceNew (immutable fields)
- [x] 1.4 Set proper descriptions for all schema fields
- [x] 1.5 Configure resource importer with StatePassthrough for composite ID support

## 2. Service Layer Implementation

- [x] 2.1 Add `DescribeTeoFunctionV2ById` method to TeoService in service_tencentcloud_teo.go
- [x] 2.2 Implement DescribeFunctions API call with zone_id and function_id parameters
- [x] 2.3 Add retry logic with tccommon.ReadRetryTimeout for the describe operation
- [x] 2.4 Implement error handling and logging for the describe operation
- [x] 2.5 Handle case where function is not found (return nil without error)

## 3. Create Operation Implementation

- [x] 3.1 Implement `resourceTencentCloudTeoFunctionV2Create` function with proper defer handlers (LogElapsed, InconsistentCheck)
- [x] 3.2 Implement CreateFunction API call with all required parameters (zone_id, name, content, remark)
- [x] 3.3 Add retry logic with tccommon.WriteRetryTimeout for the create operation
- [x] 3.4 Implement state refresh function to poll DescribeFunctions until domain is assigned
- [x] 3.5 Configure state refresh with 10s delay, 3s min timeout, 600s max timeout
- [x] 3.6 Set composite resource ID using `zone_id#function_id` format with tccommon.FILED_SP separator
- [x] 3.7 Call Read function after successful creation to populate state

## 4. Read Operation Implementation

- [x] 4.1 Implement `resourceTencentCloudTeoFunctionV2Read` function with proper defer handlers
- [x] 4.2 Parse composite resource ID to extract zone_id and function_id
- [x] 4.3 Call TeoService.DescribeTeoFunctionV2ById to fetch function details
- [x] 4.4 Handle function not found case by clearing resource state
- [x] 4.5 Set all resource state fields from API response (function_id, name, remark, content, domain, create_time, update_time)
- [x] 4.6 Set zone_id from parsed ID

## 5. Update Operation Implementation

- [x] 5.1 Implement `resourceTencentCloudTeoFunctionV2Update` function with proper defer handlers
- [x] 5.2 Validate that immutable fields (name, zone_id) are not changed
- [x] 5.3 Detect which mutable fields (remark, content) have changed
- [x] 5.4 Call ModifyFunction API with changed fields only
- [x] 5.5 Add retry logic with tccommon.WriteRetryTimeout for the update operation
- [x] 5.6 Call Read function after successful update to refresh state

## 6. Delete Operation Implementation

- [x] 6.1 Implement `resourceTencentCloudTeoFunctionV2Delete` function with proper defer handlers
- [x] 6.2 Parse composite resource ID to extract zone_id and function_id
- [x] 6.3 Call DeleteFunction API with zone_id and function_id parameters
- [x] 6.4 Add retry logic with tccommon.WriteRetryTimeout for the delete operation
- [x] 6.5 Handle delete response appropriately (ignore response as per existing pattern)

## 7. Resource Registration

- [x] 7.1 Register `resourceTencentCloudTeoFunctionV2` in the provider's resource map
- [x] 7.2 Verify resource name follows naming convention (tencentcloud_teo_function_v2)

## 8. Test Implementation

- [x] 8.1 Create unit test file `resource_tc_teo_function_v2_test.go`
- [x] 8.2 Implement test case for resource creation with valid parameters
- [x] 8.3 Implement test case for resource read after creation
- [x] 8.4 Implement test case for resource update (remark and content fields)
- [x] 8.5 Implement test case for resource deletion
- [x] 8.6 Implement test case for import using composite ID format
- [x] 8.7 Implement test case for validation (invalid name format)
- [x] 8.8 Implement test case for update attempt on immutable name field

## 9. Documentation and Examples

- [x] 9.1 Create example file `resource_tc_teo_function_v2.md` with Terraform configuration example
- [x] 9.2 Include example for resource creation with all parameters
- [x] 9.3 Include example for resource update
- [x] 9.4 Include import example with composite ID format
- [ ] 9.5 Run `make doc` to generate website documentation
- [ ] 9.6 Verify generated documentation file in `website/docs/r/teo_function_v2.md`

## 10. Code Quality Verification

- [x] 10.1 Run `go build` to verify the code compiles without errors
- [x] 10.2 Run `gofmt` to ensure code formatting follows Go standards
- [x] 10.3 Run `go vet` to check for common code issues
- [x] 10.4 Verify all imports are correct and follow project conventions
- [x] 10.5 Check for any TODO comments and address if necessary

## 11. Integration Testing

- [x] 11.1 Set up TENCENTCLOUD_SECRET_ID and TENCENTCLOUD_SECRET_KEY environment variables
- [x] 11.2 Run acceptance tests with `TF_ACC=1 go test -v -run TestAccTencentCloudTeoFunctionV2`
- [x] 11.3 Verify all test cases pass successfully
- [x] 11.4 Test resource creation, update, and deletion in a real environment
- [x] 11.5 Test resource import functionality
- [x] 11.6 Verify error handling with invalid inputs

## 12. Final Review and Cleanup

- [x] 12.1 Review all generated code for consistency with existing patterns
- [x] 12.2 Verify all error messages are clear and helpful
- [x] 12.3 Check that all logging statements include appropriate log IDs
- [x] 12.4 Ensure all defer statements are properly placed
- [x] 12.5 Verify code follows the project's naming conventions
- [x] 12.6 Add appropriate code comments for complex logic
- [x] 12.7 Remove any debugging or temporary code
- [x] 12.8 Final verification that all artifacts are complete and correct
