## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource manages EdgeOne (TEO) L7 acceleration rules. The cloud API `RuleEngineItem` struct exposes three new configuration areas that the Terraform resource does not yet surface: rule-level **custom variables** (`CustomVariables`), **custom variable operations** (`CustomVariableOperations`), and the **CustomAction** operation within `branches.actions` (`CustomActionParameters`). Without these parameters, Terraform users cannot define rule-scoped custom variables, perform variable assignment operations, or apply pre-defined custom action configurations — limiting full coverage of the `CreateL7AccRules` / `ModifyL7AccRule` / `DescribeL7AccRules` API capabilities.

## What Changes

- Add a new top-level `custom_variables` parameter (TypeList, Optional) to `tencentcloud_teo_l7_acc_rule_v2`, mapping to `RuleEngineItem.CustomVariables`. Each element exposes `name`, `initial_value`, and `description` (the `CustomVariable` struct).
- Add a new top-level `custom_variable_operations` parameter (TypeList, Optional, MaxItems: 1) mapping to `RuleEngineItem.CustomVariableOperations`. Each element is a `CustomVariableOperation` containing `branches` (a single `CustomVariableOperationRuleBranch` with `condition`, `actions`, `sub_rules`) and `description`. The nested `actions` expose `name` and `set_parameters` (with `name`/`value`, the `SetParameters` struct). The nested `sub_rules` expose `branches` and `description`.
- Add a new `custom_action_parameters` action parameter (TypeList, Optional, MaxItems: 1) under `branches.actions` of the existing `TencentTeoL7RuleBranchBasicInfo` schema, mapping to `RuleEngineAction.CustomActionParameters`. It exposes `custom_actions` (a list of `CustomAction` items, each with `action_id` and `parameters`). Each `parameters` item reuses a single `value` field (TypeString) that the flatten/expand logic maps across all value-type variants (`StringValue`/`IntegerValue`/...); the cloud API `CustomActionParameter` struct's `ValueType` is driven from the `name` field.
- Add the `CustomAction` operation name to the `name` field description in `branches.actions`.
- Wire the new parameters through Create (`CreateL7AccRules` via `request.Rules`), Update (`ModifyL7AccRule` via `request.Rule`), and Read (`DescribeL7AccRules` via `response.Rules`) with nil checks at every nesting level.
- Add unit tests using the gomonkey mock approach for the new parameters.
- Update the resource `.md` documentation (generated via `make doc` in the finalize stage).

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-custom-params`: Adds the `custom_variables`, `custom_variable_operations` top-level parameters and the `custom_action_parameters` action parameter (with the `CustomAction` operation name) to the `tencentcloud_teo_l7_acc_rule_v2` resource, enabling rule-level custom variables, custom variable operations, and custom action configuration through the L7 acceleration rule APIs.

### Modified Capabilities

## Impact

- **Affected files**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go`, `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`
- **Cloud API**: Uses existing `CreateL7AccRules`, `ModifyL7AccRule`, and `DescribeL7AccRules` APIs (no new API dependencies). The `CustomVariables`, `CustomVariableOperations`, `CustomVariableOperation`, `CustomVariableOperationRuleBranch`, `CustomVariableOperationRuleAction`, `CustomVariableOperationSubRule`, `SetParameters`, `CustomActionParameters`, `CustomAction`, and `CustomActionParameter` structs already exist in the vendored SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`).
- **Backward compatibility**: All new parameters are Optional; existing Terraform configurations and state continue to work without changes. The cloud API preserves existing configuration when these fields are not specified.