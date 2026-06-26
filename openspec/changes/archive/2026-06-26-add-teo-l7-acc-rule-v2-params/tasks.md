## 1. Schema and CRUD Implementation

- [x] 1.1 Verify and update the `zone_id` parameter mapping in `resource_tc_teo_l7_acc_rule_v2.go` Create function to ensure `request.ZoneId` is correctly set from schema `zone_id`
- [x] 1.2 Verify and update the `Rules` parameter mapping in Create function to ensure rule fields (`status`, `rule_name`, `description`, `branches`) are correctly mapped to `RuleEngineItem` within `request.Rules`
- [x] 1.3 Verify and update the Modify function to ensure `request.ZoneId` is set from `zone_id`, `request.Rule.RuleId` from `rule_id`, and other fields (`status`, `rule_name`, `description`, `branches`) are mapped to `request.Rule`
- [x] 1.4 Verify and update the Delete function to ensure `request.ZoneId` is set from `zone_id` and `request.RuleIds` from `rule_id`
- [x] 1.5 Verify and update the Read function (service layer `DescribeTeoL7AccRuleById`) to ensure `request.Filters` uses `Name: "rule-id"` and `Values: [ruleId]` for querying by `rule_id`

## 2. Testing

- [x] 2.1 Add unit tests in `resource_tc_teo_l7_acc_rule_v2_test.go` using gomonkey mock to verify Create function parameter mapping (zone_id → request.ZoneId, rule fields → request.Rules)
- [x] 2.2 Add unit tests to verify Modify function parameter mapping (zone_id → request.ZoneId, rule_id → request.Rule.RuleId, other fields → request.Rule)
- [x] 2.3 Add unit tests to verify Delete function parameter mapping (zone_id → request.ZoneId, rule_id → request.RuleIds)
- [x] 2.4 Add unit tests to verify Read function uses correct Filters (Name: "rule-id", Values: [ruleId])
- [x] 2.5 Run unit tests with `go test -gcflags="all=-l"` to verify all tests pass

## 3. Documentation

- [x] 3.1 Update `resource_tc_teo_l7_acc_rule_v2.md` to reflect the CRUD parameter mappings and ensure documentation is consistent with the implementation
