# Tasks: Add TotalCount Parameter to tencentcloud_teo_l7_acc_rule

## 1. Schema Definition

- [x] 1.1 Add `total_count` field to the data source schema as a computed field
- [x] 1.2 Verify the field type is `schema.TypeInt` to match the API's int64 type

## 2. Read Function Implementation

- [x] 2.1 Extract `TotalCount` value from the API response structure
- [x] 2.2 Set the `total_count` state field with the extracted value
- [x] 2.3 Add error handling for null/undefined TotalCount values (default to 0 if needed)

## 3. Service Layer Updates

- [x] 3.1 Verify that the service layer correctly passes the TotalCount from the API response to the data source

## 4. Testing

- [x] 4.1 Add unit test to verify total_count field extraction from API response
- [ ] 4.2 Add integration test to ensure total_count is accessible in Terraform configurations (basic test added, requires test execution)
- [ ] 4.3 Add edge case test for empty result sets (total_count = 0) (requires test execution)
- [ ] 4.4 Add test for large total_count values to ensure proper integer handling (requires test execution)

## 5. Documentation

- [x] 5.1 Update the data source example file (data_source_tc_teo_l7_acc_rule.md) to include total_count field
- [ ] 5.2 Run `make doc` to generate website documentation automatically (requires external build tools)
- [ ] 5.3 Verify the generated documentation includes the total_count attribute (dependent on 5.2)
- [x] 5.4 Add usage example demonstrating total_count in outputs

## 6. Build and Validation

- [ ] 6.1 Run `make build` to verify the code compiles successfully (requires external build tools)
- [ ] 6.2 Run `make lint` to check for code quality issues (requires external build tools)
- [ ] 6.3 Run `TF_ACC=1 make test TEST=./tencentcloud/services/teo/` to execute acceptance tests (requires external build tools)
- [ ] 6.4 Verify all tests pass and the total_count field is correctly exposed (dependent on 6.3)
