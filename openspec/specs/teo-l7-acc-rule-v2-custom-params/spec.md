# teo-l7-acc-rule-v2-custom-params Specification

## Purpose
TBD - created by archiving change add-teo-l7-acc-rule-v2-custom-params. Update Purpose after archive.
## Requirements
### Requirement: Resource supports custom_variables top-level parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `custom_variables` parameter (TypeList) at the top level of the resource schema. This parameter defines the rule-level custom variable list and maps to the SDK `RuleEngineItem.CustomVariables` field (a list of `CustomVariable`). Each element SHALL expose `name` (TypeString, Optional), `initial_value` (TypeString, Optional), and `description` (TypeString, Optional), corresponding to the `CustomVariable` struct fields `Name`, `InitialValue`, and `Description`.

#### Scenario: Create resource with custom_variables
- **WHEN** user defines `custom_variables` with `name`, `initial_value`, and `description` in a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Create function SHALL map each element to a `CustomVariable` struct appended to `RuleEngineItem.CustomVariables` in the `CreateL7AccRules` API request (`request.Rules[0].CustomVariables`)

#### Scenario: Read resource returns custom_variables
- **WHEN** the `DescribeL7AccRules` API returns rules with `Rules[0].CustomVariables` populated
- **THEN** the Read function SHALL flatten each `CustomVariable` into the `custom_variables` attribute in the Terraform state, setting `name`, `initial_value`, and `description` only when the corresponding SDK fields are non-nil

#### Scenario: Update resource custom_variables
- **WHEN** user modifies the `custom_variables` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineItem.CustomVariables` in the `ModifyL7AccRule` API request (`request.Rule.CustomVariables`)

#### Scenario: custom_variables omitted
- **WHEN** user does not specify `custom_variables` in the resource configuration
- **THEN** the resource SHALL behave exactly as before and SHALL NOT send the `CustomVariables` field in the API request

### Requirement: Resource supports custom_variable_operations top-level parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `custom_variable_operations` parameter (TypeList, MaxItems: 1) at the top level of the resource schema. This parameter defines rule-level custom variable operations and maps to the SDK `RuleEngineItem.CustomVariableOperations` field (a list of `CustomVariableOperation`). Each element SHALL expose `branches` (TypeList, single `CustomVariableOperationRuleBranch`) and `description` (TypeList of TypeString).

#### Scenario: custom_variable_operations.branches exposes condition, actions, sub_rules
- **WHEN** user defines `custom_variable_operations[0].branches[0]`
- **THEN** the schema SHALL expose `condition` (TypeString, Optional), `actions` (TypeList of `CustomVariableOperationRuleAction`), and `sub_rules` (TypeList of `CustomVariableOperationSubRule`)

#### Scenario: custom_variable_operations.branches.actions exposes name and set_parameters
- **WHEN** user defines `custom_variable_operations[0].branches[0].actions[0]`
- **THEN** the schema SHALL expose `name` (TypeString, Optional) and `set_parameters` (TypeList, MaxItems: 1) where `set_parameters[0]` exposes `name` (TypeString, Optional) and `value` (TypeString, Optional), mapping to the `SetParameters` struct fields `Name` and `Value`

#### Scenario: custom_variable_operations.branches.sub_rules exposes branches and description
- **WHEN** user defines `custom_variable_operations[0].branches[0].sub_rules[0]`
- **THEN** the schema SHALL expose `branches` (TypeList of `CustomVariableOperationRuleBranch`, supporting `condition`, `actions`, and recursive `sub_rules`) and `description` (TypeList of TypeString)

#### Scenario: Create resource with custom_variable_operations
- **WHEN** user defines `custom_variable_operations` in a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Create function SHALL map it to `RuleEngineItem.CustomVariableOperations` in the `CreateL7AccRules` API request

#### Scenario: Read resource returns custom_variable_operations
- **WHEN** the `DescribeL7AccRules` API returns rules with `Rules[0].CustomVariableOperations` populated
- **THEN** the Read function SHALL flatten `CustomVariableOperations` into the `custom_variable_operations` attribute with nil checks at every nesting level

#### Scenario: Update resource custom_variable_operations
- **WHEN** user modifies the `custom_variable_operations` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineItem.CustomVariableOperations` in the `ModifyL7AccRule` API request

### Requirement: Resource supports custom_action_parameters in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `custom_action_parameters` parameter (TypeList, MaxItems: 1) within `branches.actions` of the `TencentTeoL7RuleBranchBasicInfo` schema. This parameter defines the custom action configuration and maps to the SDK `RuleEngineAction.CustomActionParameters` field. It SHALL expose `custom_actions` (TypeList of `CustomAction` items).

#### Scenario: custom_action_parameters.custom_actions exposes action_id and parameters
- **WHEN** user defines `custom_action_parameters[0].custom_actions[0]`
- **THEN** the schema SHALL expose `action_id` (TypeString, Optional, mapping to `CustomAction.ActionId`) and `parameters` (TypeList of `CustomActionParameter` items)

#### Scenario: custom_actions.parameters element exposes name and value
- **WHEN** user defines `custom_action_parameters[0].custom_actions[0].parameters[0]`
- **THEN** the schema SHALL expose `name` (TypeString, Optional, mapping to `CustomActionParameter.Name`) and `value` (TypeString, Optional). The expand logic SHALL inspect the SDK `CustomActionParameter.ValueType` to set the appropriate typed value field (`StringValue`/`IntegerValue`/`FloatValue`/`BooleanValue`/`StringArrayValue`/`IntegerArrayValue`/`FloatArrayValue`), and the flatten logic SHALL serialize the populated typed value back to the single `value` string

#### Scenario: Create resource with CustomAction
- **WHEN** user defines `custom_action_parameters` within an `actions` block (with `name = "CustomAction"`) inside `branches`
- **THEN** the Create function SHALL map it to `RuleEngineAction.CustomActionParameters` in the `CreateL7AccRules` API request

#### Scenario: Read resource returns custom_action_parameters
- **WHEN** the `DescribeL7AccRules` API returns rules with `Branches[0].Actions[0].CustomActionParameters` populated
- **THEN** the Read function SHALL flatten `CustomActionParameters` into the `custom_action_parameters` attribute with nil checks at every nesting level

#### Scenario: Update resource custom_action_parameters
- **WHEN** user modifies the `custom_action_parameters` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineAction.CustomActionParameters` in the `ModifyL7AccRule` API request

### Requirement: Actions name description includes CustomAction
The `name` field description in `branches.actions` SHALL include `CustomAction` as a valid operation name, indicating that `CustomActionParameters` is required when `name = CustomAction`.

#### Scenario: Name field description lists CustomAction
- **WHEN** user reads the description of the `name` field in `branches.actions`
- **THEN** the description SHALL list `CustomAction: Custom configuration` as a valid operation name

### Requirement: All new parameters are optional
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL treat `custom_variables`, `custom_variable_operations`, and `custom_action_parameters` as optional parameters that do not alter existing behavior when omitted.

#### Scenario: All new parameters omitted
- **WHEN** user does not specify `custom_variables`, `custom_variable_operations`, or `custom_action_parameters` in the resource configuration
- **THEN** the resource SHALL behave exactly as before (no changes to existing behavior) and SHALL NOT send the corresponding fields in the API request

