## ADDED Requirements

### Requirement: Resource supports custom_variables at top level
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional top-level `custom_variables` parameter (TypeList) mapping to `RuleEngineItem.CustomVariables` in the cloud API. Each element SHALL contain:
- `name` (TypeString, Optional): 变量名称，需填写 `user.rule.` 前缀
- `initial_value` (TypeString, Optional): 变量初始值
- `description` (TypeString, Optional): 变量描述

#### Scenario: Create resource with custom_variables
- **WHEN** user defines `custom_variables` blocks (with `name`, `initial_value`, `description`) at the top level of a `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Create function SHALL map each element to a `CustomVariable` struct and set `RuleEngineItem.CustomVariables` in the `CreateL7AccRules` API request

#### Scenario: Read resource returns custom_variables
- **WHEN** the `DescribeL7AccRules` API returns rules with `CustomVariables` populated (non-nil)
- **THEN** the Read function SHALL flatten each `CustomVariable` into the `custom_variables` attribute in the Terraform state

#### Scenario: Update resource custom_variables
- **WHEN** user modifies the `custom_variables` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineItem.CustomVariables` in the `ModifyL7AccRule` API request

#### Scenario: custom_variables not set
- **WHEN** user does not specify `custom_variables` in the resource configuration
- **THEN** the `RuleEngineItem.CustomVariables` field SHALL NOT be set in the request, and resource behavior SHALL be unchanged

### Requirement: Resource supports custom_variable_operations at top level
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional top-level `custom_variable_operations` parameter (TypeList) mapping to `RuleEngineItem.CustomVariableOperations` in the cloud API. Each element SHALL contain:
- `branches` (TypeList, Optional): 子规则分支，含 `condition` (TypeString)、`actions` (TypeList，含 `name` (TypeString) 与 `set_parameters` (TypeList MaxItems=1，含 `name`、`value` (TypeString)))、`sub_rules` (TypeList，含 `branches`（递归同 branches 结构）与 `description` (TypeList))
- `description` (TypeList, Optional): 规则注释

#### Scenario: Create resource with custom_variable_operations
- **WHEN** user defines `custom_variable_operations` blocks (with `branches` containing `condition`/`actions`/`sub_rules` and `description`) at the top level
- **THEN** the Create function SHALL map each element to a `CustomVariableOperation` struct (含 `Branches []*CustomVariableOperationRuleBranch`、`Description []*string`) and set `RuleEngineItem.CustomVariableOperations` in the request

#### Scenario: custom_variable_operations actions set_parameters mapping
- **WHEN** `actions` block contains `name="Set"` and `set_parameters` with `name`/`value`
- **THEN** the Create/Update function SHALL map to `CustomVariableOperationRuleAction.Name` and `SetParameters` (`SetParameters.Name`/`SetParameters.Value`)

#### Scenario: custom_variable_operations sub_rules mapping
- **WHEN** `branches` contains `sub_rules` (含递归 `branches` 与 `description`)
- **THEN** the function SHALL map to `CustomVariableOperationRuleBranch.SubRules` (`[]*CustomVariableOperationSubRule`，每个含 `Branches []*CustomVariableOperationRuleBranch` 与 `Description []*string`)

#### Scenario: Read resource returns custom_variable_operations
- **WHEN** the `DescribeL7AccRules` API returns rules with `CustomVariableOperations` populated (non-nil)
- **THEN** the Read function SHALL flatten each `CustomVariableOperation` (含 branches/actions/set_parameters/sub_rules/description) into the `custom_variable_operations` attribute in state

#### Scenario: Update resource custom_variable_operations
- **WHEN** user modifies the `custom_variable_operations` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineItem.CustomVariableOperations` in the `ModifyL7AccRule` request

#### Scenario: custom_variable_operations not set
- **WHEN** user does not specify `custom_variable_operations`
- **THEN** the `RuleEngineItem.CustomVariableOperations` field SHALL NOT be set, and resource behavior SHALL be unchanged

### Requirement: Resource supports custom_action_parameters in branches.actions
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `custom_action_parameters` parameter within `branches.actions` (TypeList, MaxItems=1). This maps to `RuleEngineAction.CustomActionParameters` in the cloud API, which contains a `custom_actions` list (TypeList). Each `custom_actions` element SHALL contain:
- `action_id` (TypeString, Optional): 定制配置的配置项 Id
- `parameters` (TypeList, Optional): 参数字段值列表，每个参数含 `name` (TypeString)、`value_type` (TypeString)、`string_value` (TypeString)、`integer_value` (TypeInt)、`float_value` (TypeFloat)、`boolean_value` (TypeBool)、`string_array_value` (TypeList of TypeString)、`integer_array_value` (TypeList of TypeInt)、`float_array_value` (TypeList of TypeFloat)

#### Scenario: Create resource with CustomAction action
- **WHEN** user defines an `actions` block with `name = "CustomAction"` and `custom_action_parameters` containing `custom_actions` (with `action_id` and `parameters`)
- **THEN** the Create function SHALL map it to `RuleEngineAction.CustomActionParameters` (`CustomActionParameters.CustomActions []*CustomAction`，每个含 `ActionId` 与 `Parameters []*CustomActionParameter`) in the API request

#### Scenario: Read resource returns custom_action_parameters
- **WHEN** the `DescribeL7AccRules` API returns rules with `Branches[].Actions[].CustomActionParameters` populated (non-nil)
- **THEN** the Read function SHALL set the `custom_action_parameters` attribute (含 custom_actions 与其 action_id/parameters) in the Terraform state

#### Scenario: Update resource custom_action_parameters
- **WHEN** user modifies the `custom_action_parameters` parameter
- **THEN** the Update function SHALL map the new value to `RuleEngineAction.CustomActionParameters` in the API request

#### Scenario: custom_action_parameters not set
- **WHEN** user does not specify `custom_action_parameters` in an actions block
- **THEN** the `RuleEngineAction.CustomActionParameters` field SHALL NOT be set, and resource behavior SHALL be unchanged

### Requirement: Actions name description includes CustomAction
The `name` field description in `branches.actions` SHALL include `CustomAction` as a valid operation name.

#### Scenario: Name field description is updated
- **WHEN** user reads the description of the `name` field in `branches.actions`
- **THEN** the description SHALL list `CustomAction` (定制配置) as a valid operation name

### Requirement: All new parameters are optional
- **WHEN** user does not specify `custom_variables`, `custom_variable_operations`, or `custom_action_parameters` in the resource configuration
- **THEN** the resource SHALL behave exactly as before (no changes to existing behavior)

### Requirement: Unit tests for new parameters
The system SHALL provide unit tests in `resource_tc_teo_l7_acc_rule_v2_test.go` using gomonkey to mock cloud API calls, testing the new `custom_variables`, `custom_variable_operations`, and `custom_action_parameters` parameters across Create, Read, and Update operations.

#### Scenario: Custom variables create test
- **WHEN** a test simulates creating a resource with `custom_variables` set
- **THEN** the mocked `CreateL7AccRules` SHALL be invoked with the `CustomVariables` field populated
- **AND** the test SHALL assert the `RuleEngineItem.CustomVariables` content matches the configured values

#### Scenario: Custom variable operations create test
- **WHEN** a test simulates creating a resource with `custom_variable_operations` set (含 branches/actions/set_parameters/sub_rules/description)
- **THEN** the mocked `CreateL7AccRules` SHALL be invoked with `CustomVariableOperations` populated
- **AND** the test SHALL assert the nested structure is correctly mapped

#### Scenario: Custom action parameters create test
- **WHEN** a test simulates creating a resource with `custom_action_parameters` in an actions block
- **THEN** the mocked `CreateL7AccRules` SHALL be invoked with `RuleEngineAction.CustomActionParameters` populated
- **AND** the test SHALL assert `CustomActions` (含 ActionId/Parameters) is correctly mapped

#### Scenario: Update test for new parameters
- **WHEN** a test simulates updating `custom_variables` or `custom_variable_operations`
- **THEN** the mocked `ModifyL7AccRule` SHALL be invoked with the corresponding fields populated

### Requirement: Resource documentation includes new parameters
The system SHALL update the markdown documentation file `resource_tc_teo_l7_acc_rule_v2.md` to include example usage of `custom_variables`, `custom_variable_operations`, and `custom_action_parameters`.

#### Scenario: Documentation includes custom_variables example
- **WHEN** the resource documentation is updated
- **THEN** the .md file SHALL contain an example HCL block demonstrating `custom_variables` usage

#### Scenario: Documentation includes custom_action_parameters example
- **WHEN** the resource documentation is updated
- **THEN** the .md file SHALL contain an example HCL block showing `custom_action_parameters` within an `actions` block