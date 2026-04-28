## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently only exposes a single `rule_id` (TypeString, Computed) from the `CreateL7AccRules` API response, but the API returns `RuleIds` as a list of strings (`[]*string`). Users need access to the full list of created rule IDs returned by the API to properly reference all created rules.

## What Changes

- Add a new computed parameter `rule_ids` (TypeList of TypeString) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema, mapping to `response.Response.RuleIds` from the `CreateL7AccRules` API response.

## Capabilities

### New Capabilities
- `rule-ids-output`: Expose the full list of rule IDs returned by the CreateL7AccRules API as a computed `rule_ids` parameter in the terraform resource schema.

### Modified Capabilities

## Impact

- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` (schema definition, create and read functions)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` (unit tests)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` (documentation)
- Cloud API: `CreateL7AccRules` (response field `RuleIds`)
- Cloud API: `DescribeL7AccRules` (response uses same `RuleEngineItem` which contains `RuleId`)
