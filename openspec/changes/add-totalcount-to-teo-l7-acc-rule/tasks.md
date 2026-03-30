## 1. Schema Modification

- [x] 1.1 Add `total_count` field to tencentcloud_teo_l7_acc_rule resource schema with Type: TypeInt, Computed: true, Description: "Total count of L7 acceleration rules"
- [x] 1.2 Add `total_count` field to ResourceTencentCloudTeoL7AccRule() function in resource_tc_teo_l7_acc_rule.go

## 2. Read Function Implementation

- [x] 2.1 In resourceTencentCloudTeoL7AccRuleRead(), after retrieving respData, check if TotalCount is not nil
- [x] 2.2 If respData.TotalCount is not nil, set the value to state using `d.Set("total_count", *respData.TotalCount)`
- [x] 2.3 Handle nil TotalCount case to avoid panic (either skip setting or set to 0)

## 3. Documentation Update

- [x] 3.1 Add `total_count` field documentation to resource_tc_teo_l7_acc_rule.md with description and type information
- [x] 3.2 Run `make doc` command to generate website documentation from the updated .md file
- [x] 3.3 Verify the generated website/docs/r/teo_l7_acc_rule.html.markdown contains the new total_count field

## 4. Testing

- [x] 4.1 Add test case in resource_tc_teo_l7_acc_rule_test.go to verify total_count field is correctly populated during read
- [x] 4.2 Add test to verify total_count matches the length of rules list (optional, for data consistency)
- [x] 4.3 Add test case to verify existing resources without total_count in state work correctly after provider upgrade
- [x] 4.4 Run acceptance tests with `TF_ACC=1 go test ./tencentcloud/services/teo -run TestAccTencentCloudTeoL7AccRule -v`

## 5. Build and Verification

- [x] 5.1 Run `go build` to ensure the code compiles without errors
- [x] 5.2 Run `go fmt` to ensure code formatting is correct
- [x] 5.3 Run `go vet` to check for potential issues
- [x] 5.4 Run `go test ./tencentcloud/services/teo -run TestAccTencentCloudTeoL7AccRule` to verify unit tests pass
- [x] 5.5 Manual test: Create a test Terraform configuration with tencentcloud_teo_l7_acc_rule resource, run `terraform apply`, then `terraform show` to verify total_count field is populated

## 6. Code Review

- [x] 6.1 Review changes for backward compatibility (no breaking changes to existing schema)
- [x] 6.2 Review error handling for nil TotalCount case
- [x] 6.3 Review documentation accuracy and completeness
- [x] 6.4 Verify all code follows project coding standards and patterns
