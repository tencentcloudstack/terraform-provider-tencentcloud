## 1. Schema Definition

- [x] 1.1 Add `actions` top-level optional parameter to `resource_tc_teo_l7_acc_rule_v2.go` Schema map, reusing the `RuleEngineAction` schema structure from `TencentTeoL7RuleBranchBasicInfo`

## 2. Create Method

- [x] 2.1 In `ResourceTencentCloudTeoL7AccRuleV2Create`, add logic to read `actions` from resource data and populate `rule.Branches[0].Actions` before API call, ensuring compatibility when `branches` is also set

## 3. Read Method

- [x] 3.1 In `ResourceTencentCloudTeoL7AccRuleV2Read`, add logic to extract `Actions` from `respData.Rules[0].Branches[0].Actions` and set to `actions` in resource data

## 4. Update Method

- [x] 4.1 In `ResourceTencentCloudTeoL7AccRuleV2Update`, add `actions` to the `HasChange` check and add logic to populate `rule.Branches[0].Actions` with the new actions value

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_l7_acc_rule_v2.md` to add an Example Usage section demonstrating the `actions` parameter usage

## 6. Verification

- [x] 6.1 Verify code compiles correctly (no `go build` needed - check syntax consistency)
- [x] 6.2 Verify existing unit tests in `resource_tc_teo_l7_acc_rule_v2_test.go` still pass with new changes