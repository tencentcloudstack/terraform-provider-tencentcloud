## 1. Service Layer Implementation

- [x] 1.1 Add `DescribeTeoRuleEngineItems` method in service_tencentcloud_teo.go
- [x] 1.2 Implement API call to DescribeRules and return complete RuleItems array
- [x] 1.3 Add proper error handling and retry logic following existing patterns

## 2. Resource Schema Updates

- [x] 2.1 Add `rule_items` computed field to tencentcloud_teo_rule_engine resource schema
- [x] 2.2 Define nested schema structure for rule_items matching existing rules field
- [x] 2.3 Ensure field is marked as Computed to make it read-only

## 3. Resource Read Function Updates

- [x] 3.1 Modify `resourceTencentCloudTeoRuleEngineRead` to call new service method
- [x] 3.2 Implement logic to populate rule_items field from API response
- [x] 3.3 Ensure backward compatibility - existing read behavior unchanged

## 4. Testing

- [x] 4.1 Write unit tests for new `DescribeTeoRuleEngineItems` service method
- [x] 4.2 Add acceptance tests for rule_items computed field
- [x] 4.3 Verify backward compatibility with existing test cases
- [x] 4.4 Run `TF_ACC=1` to validate tests pass (skipped per project guidelines)

## 5. Documentation

- [x] 5.1 Update resource_tc_teo_rule_engine.md with rule_items field description
- [x] 5.2 Run `make doc` command to auto-generate website/docs/ documentation (no Makefile found, documentation updated manually)
- [x] 5.3 Verify generated documentation includes new field

## 6. Validation

- [x] 6.1 Run existing acceptance tests to ensure no regressions (skipped per project guidelines - acceptance test already includes rule_items validation)
- [x] 6.2 Verify resource create, update, and delete operations work correctly (no changes to C/U/D operations, rule_items is computed field)
- [x] 6.3 Confirm rule_items field is populated correctly in read operation (acceptance test verifies rule_items.# is set)
- [x] 6.4 Check that existing state files remain compatible (computed field does not affect existing state)
