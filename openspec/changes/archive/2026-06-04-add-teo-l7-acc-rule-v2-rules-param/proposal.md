## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently exposes individual fields (status, rule_name, description, branches, rule_id, rule_priority) from a single `RuleEngineItem`, but does not expose the full `Rules` list returned by the `DescribeL7AccRules` API response. Adding a computed `rules` attribute allows users to access the complete rules list information directly from the resource state.

## What Changes

- Add a new `Computed` attribute `rules` (type: list) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, representing the full `Rules` array from the `DescribeL7AccRules` API response.
- Update the Read function to populate the new `rules` attribute from `response.Rules`.
- Update the resource documentation (.md file) to reflect the new attribute.
- Add unit tests covering the new attribute.

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-rules-output`: Add a computed `rules` attribute to `tencentcloud_teo_l7_acc_rule_v2` that exposes the full rules list from the DescribeL7AccRules API response.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`: Add `rules` schema definition and update Read function.
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`: Update documentation with new attribute.
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`: Add unit tests for the new attribute.
- Backward compatible: only adds a new Computed field, no breaking changes to existing configurations.
