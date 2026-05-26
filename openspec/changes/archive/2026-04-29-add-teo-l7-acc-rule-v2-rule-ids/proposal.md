## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource had a redundant `rule_ids` computed attribute. Since Terraform manages individual resources (one rule per resource instance), the `rule_id` field is sufficient. Additionally, the `RuleEngineAction` struct in the SDK supports a `Vary` action with `VaryParameters`, but this was not exposed in the Terraform resource.

## What Changes

- Remove the `rule_ids` computed attribute from the schema and its population logic in Create/Read functions
- Remove unit tests related to `rule_ids`
- Add `vary_parameters` action support (schema, Get/flatten, Set/build) in `resource_tc_teo_l7_acc_rule_extension.go`
- Add `origin_authentication_parameters` action support (schema, Get/flatten, Set/build) with nested `request_properties` list (type/name/value)
- Fix `actions.name` field description: correct `SetContentIdentifierParameters` to `SetContentIdentifier`, and add missing action names (`ContentCompression`, `OriginAuthentication`)
- Update documentation with `Vary` and `OriginAuthentication` action examples

## Capabilities

### New Capabilities
- `l7-acc-rule-v2-vary-action`: Add `Vary` action support with `vary_parameters` block (containing `switch` field) to the `tencentcloud_teo_l7_acc_rule_v2` resource's `branches.actions`, mapping to `RuleEngineAction.VaryParameters` in the CreateL7AccRules/ModifyL7AccRule/DescribeL7AccRules APIs.
- `l7-acc-rule-v2-origin-authentication-action`: Add `OriginAuthentication` action support with `origin_authentication_parameters` block (containing `request_properties` list with `type`, `name`, `value` fields) mapping to `RuleEngineAction.OriginAuthenticationParameters` in the APIs.

### Modified Capabilities
- `l7-acc-rule-v2-rule-ids-removed`: Remove the redundant `rule_ids` computed attribute. The `rule_id` field provides the single rule ID managed by this resource.

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`: Remove `rule_ids` schema and related Create/Read logic
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go`: Add `vary_parameters` and `origin_authentication_parameters` to schema, Get function, and Set function
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`: Remove `rule_ids` unit tests
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`: Add `Vary` and `OriginAuthentication` action examples
- API dependency: `CreateL7AccRules`, `ModifyL7AccRule`, `DescribeL7AccRules` (all already used), `VaryParameters`, `OriginAuthenticationParameters` structs in SDK
