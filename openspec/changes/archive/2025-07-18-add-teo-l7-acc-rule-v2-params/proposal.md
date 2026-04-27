## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource currently lacks a `rule_ids` computed output parameter to expose all rule IDs returned by the `CreateL7AccRules` API response. While the existing code internally reads `response.Response.RuleIds[0]` to set the composite resource ID, the full list of rule IDs is not available to Terraform users as a schema field. Adding `rule_ids` as a computed attribute will allow users to reference all created rule IDs in their Terraform configurations.

## What Changes

- Add a new computed schema field `rule_ids` (TypeList of TypeString) to the `tencentcloud_teo_l7_acc_rule_v2` resource
- Populate `rule_ids` from `response.Response.RuleIds` in the Create function
- Populate `rule_ids` from the DescribeL7AccRules response in the Read function

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-rule-ids`: Expose the `rule_ids` computed attribute in the `tencentcloud_teo_l7_acc_rule_v2` resource, allowing Terraform users to access all rule IDs returned by the CreateL7AccRules API

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- Affected code: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`
- Affected tests: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`
- Affected docs: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`
- API dependency: `CreateL7AccRules` response (`RuleIds` field) and `DescribeL7AccRules` response
- Backward compatible: Only adding a new computed field, no breaking changes
