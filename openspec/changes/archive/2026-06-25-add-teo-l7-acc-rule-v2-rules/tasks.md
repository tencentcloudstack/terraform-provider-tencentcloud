## 1. Schema & CRUD Implementation

- [x] 1.1 Add `rules` computed field to `ResourceTencentCloudTeoL7AccRuleV2()` schema in `resource_tc_teo_l7_acc_rule_v2.go`, using `TypeList` with an `Elem` that mirrors the `RuleEngineItem` fields (`rule_id`, `status`, `rule_name`, `description`, `branches`, `rule_priority`) all as Computed
- [x] 1.2 Create helper function `flattenTeoL7AccRuleV2Rules` to convert `[]*teov20220901.RuleEngineItem` to `[]map[string]interface{}` for use with `d.Set("rules", ...)`
- [x] 1.3 In `ResourceTencentCloudTeoL7AccRuleV2Read`, after `resourceTencentCloudTeoL7AccRuleSetBranchs`, set the `rules` field using `flattenTeoL7AccRuleV2Rules(respData.Rules)`

## 2. Unit Tests

- [x] 2.1 Add unit tests for `flattenTeoL7AccRuleV2Rules` helper function and verify `rules` field is correctly populated in Read path using gomonkey mocks
- [x] 2.2 Run unit tests with `go test -gcflags=all=-l` and ensure all pass

## 3. Documentation

- [x] 3.1 Update `resource_tc_teo_l7_acc_rule_v2.md` with `rules` field usage example in the HCL configuration block

## 4. Validation

- [x] 4.1 Ensure the change compiles cleanly and existing unit tests still pass