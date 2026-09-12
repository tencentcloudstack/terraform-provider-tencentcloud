## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource needs to add new parameters to support the full CRUD lifecycle of L7 acceleration rules via the cloud API. The current resource already supports `zone_id`, `rule_id`, `status`, `rule_name`, `description`, `branches`, and `rule_priority`, but needs to be updated to properly map all parameters across the Create, Modify, Delete, and Describe API interfaces.

## What Changes

- Add parameter mappings for the `CreateL7AccRules` API interface (`ZoneId`, `Rules`)
- Add parameter mappings for the `ModifyL7AccRule` API interface (`zone_id`, `rule_id`, `status`, `rule_name`, `description`, `branches`)
- Add parameter mappings for the `DeleteL7AccRules` API interface (`zone_id`, `rule_id`)
- Add parameter mappings for the `DescribeL7AccRules` API interface (`rule_id` via `Values` filter)

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-crud-params`: Full CRUD parameter mapping support for the tencentcloud_teo_l7_acc_rule_v2 resource, ensuring all parameters are correctly mapped across CreateL7AccRules, ModifyL7AccRule, DeleteL7AccRules, and DescribeL7AccRules API interfaces.

### Modified Capabilities
- `rule-ids-output`: Update the existing spec to reflect the new CRUD parameter mappings and ensure consistency with the full API parameter set.

## Impact

- Affected files: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`
- Cloud APIs: `CreateL7AccRules`, `ModifyL7AccRule`, `DeleteL7AccRules`, `DescribeL7AccRules` (package: `teo/v20220901`)
- SDK struct: `RuleEngineItem` is the core data structure for L7 Acc Rules
