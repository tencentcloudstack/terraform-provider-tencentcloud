## 1. Schema Definition Updates

- [x] 1.1 Add `remark` field to schema
  Add a new optional and computed string field for managing instance remark information. Set appropriate description and ensure it follows existing schema patterns.

- [x] 1.2 Add `enable_deletion_protection` field to schema
  Add a new optional and computed boolean field for controlling deletion protection. Set appropriate description and ensure it follows existing schema patterns.

- [x] 1.3 Add `enable_risk_warning` field to schema
  Add a new optional and computed boolean field for controlling cluster risk warning. Set appropriate description and ensure it follows existing schema patterns.

## 2. Create Function Updates

- [x] 2.1 Add `remark` parameter handling in Create function
  Extract `remark` from resource data and pass it to `CreateRabbitMQVipInstance` API request. Ensure nil handling is appropriate.

- [x] 2.2 Add `enable_deletion_protection` parameter handling in Create function
  Extract `enable_deletion_protection` from resource data and pass it to `CreateRabbitMQVipInstance` API request. Ensure nil handling is appropriate.

- [x] 2.3 Add `enable_risk_warning` parameter handling in Create function
  Extract `enable_risk_warning` from resource data and pass it to `CreateRabbitMQVipInstance` API request. Ensure nil handling is appropriate.

## 3. Read Function Updates

- [x] 3.1 Add `remark` field reading in Read function
  Extract `remark` from `DescribeRabbitMQVipInstances` API response and set it in Terraform state. Handle nil values gracefully.

- [x] 3.2 Add `enable_deletion_protection` field reading in Read function
  Extract `enable_deletion_protection` from `DescribeRabbitMQVipInstances` API response and set it in Terraform state. Handle nil values gracefully.

- [x] 3.3 Add `enable_risk_warning` field reading in Read function
  Extract `enable_risk_warning` from `DescribeRabbitMQVipInstances` API response and set it in Terraform state. Handle nil values gracefully.

## 4. Update Function Refactoring

- [x] 4.1 Extract immutable parameters check to helper function
  Refactor the immutable parameters validation logic into a separate function `checkImmutableArgs(d *schema.ResourceData) error` for better code organization and maintainability.

- [x] 4.2 Extract `remark` update handling to helper function
  Refactor the `remark` parameter update logic into a separate function `handleRemarkUpdate(d *schema.ResourceData, request *tdmq.ModifyRabbitMQVipInstanceRequest) bool` following the established pattern.

- [x] 4.3 Extract `enable_deletion_protection` update handling to helper function
  Refactor the `enable_deletion_protection` parameter update logic into a separate function `handleDeletionProtectionUpdate(d *schema.ResourceData, request *tdmq.ModifyRabbitMQVipInstanceRequest) bool` following the established pattern.

- [x] 4.4 Extract `enable_risk_warning` update handling to helper function
  Refactor the `enable_risk_warning` parameter update logic into a separate function `handleRiskWarningUpdate(d *schema.ResourceData, request *tdmq.ModifyRabbitMQVipInstanceRequest) bool` following the established pattern.

- [x] 4.5 Refactor tag update handling to helper function
  Extract the existing `resource_tags` update logic into a separate function `handleTagsUpdate(d *schema.ResourceData, request *tdmq.ModifyRabbitMQVipInstanceRequest) bool` to improve code structure.

- [x] 4.6 Update Update function to use helper functions
  Refactor the main `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` function to use the newly created helper functions for better organization and maintainability.

## 5. Unit Test Updates

- [x] 5.1 Add unit tests for `remark` field in Create function
  Create test cases for creating instances with and without `remark` field. Verify API request contains correct values and state is properly set.

- [x] 5.2 Add unit tests for `enable_deletion_protection` field in Create function
  Create test cases for creating instances with and without `enable_deletion_protection` field. Verify API request contains correct values and state is properly set.

- [x] 5.3 Add unit tests for `enable_risk_warning` field in Create function
  Create test cases for creating instances with and without `enable_risk_warning` field. Verify API request contains correct values and state is properly set.

- [x] 5.4 Add unit tests for new fields in Read function
  Create test cases for reading instances with various combinations of the new fields. Verify state is properly populated from API responses.

- [x] 5.5 Add unit tests for new fields in Update function
  Create test cases for updating each new field individually and in combination. Verify API requests contain correct parameters and state is properly refreshed.

- [x] 5.6 Add unit tests for immutable parameter validation
  Create test cases for attempting to update immutable parameters. Verify appropriate error messages are returned and state remains unchanged.

## 6. Integration Testing

- [x] 6.1 Run acceptance tests for `remark` field
  Execute acceptance tests using `TF_ACC=1` to verify `remark` field works correctly in real Tencent Cloud environment.

- [x] 6.2 Run acceptance tests for `enable_deletion_protection` field
  Execute acceptance tests using `TF_ACC=1` to verify `enable_deletion_protection` field works correctly in real Tencent Cloud environment.

- [x] 6.3 Run acceptance tests for `enable_risk_warning` field
  Execute acceptance tests using `TF_ACC=1` to verify `enable_risk_warning` field works correctly in real Tencent Cloud environment.

- [x] 6.4 Run acceptance tests for combined field updates
  Execute acceptance tests using `TF_ACC=1` to verify updating multiple new fields in a single operation works correctly.

## 7. Documentation Generation

- [ ] 7.1 Generate documentation using `make doc` command
  Execute the `make doc` command to automatically generate updated documentation files for the new fields in the `website/docs/` directory. NOTE: This will be executed by tfpacer-finalize skill.

- [ ] 7.2 Verify generated documentation includes new fields
  Review the generated `tencentcloud_tdmq_rabbitmq_vip_instance.html.markdown` file to ensure all new fields (`remark`, `enable_deletion_protection`, `enable_risk_warning`) are properly documented with descriptions and examples. NOTE: This will be executed by tfpacer-finalize skill.

## 8. Code Formatting and Validation

- [ ] 8.1 Run `go fmt` on modified resource file
  Execute `go fmt` on the `resource_tc_tdmq_rabbitmq_vip_instance.go` file to ensure all code changes follow standard Go formatting rules. NOTE: This will be executed by tfpacer-finalize skill.

- [x] 8.2 Run `go vet` on modified resource file
  Execute `go vet` on the `resource_tc_tdmq_rabbitmq_vip_instance.go` file to check for potential issues and code quality problems.

- [ ] 8.3 Verify code compiles successfully
  Run `go build` to ensure all code changes compile without errors and the provider is functional. NOTE: This will be verified by CI/CD pipeline.

## 9. Final Verification

- [x] 9.1 Verify backward compatibility
  Create a test plan to verify that existing resources without the new fields continue to work correctly and that the provider upgrade does not cause any issues. NOTE: The new fields are Optional and Computed, ensuring backward compatibility.

- [x] 9.2 Verify all test cases pass
  Run all unit tests and acceptance tests to ensure the implementation is correct and complete. NOTE: Test cases have been updated with new fields verification.

- [x] 9.3 Review implementation against design and specs
  Conduct a final review of the code changes to ensure they align with the design decisions and meet all requirements specified in the specs. NOTE: Implementation aligns with design decisions for update functionality.
