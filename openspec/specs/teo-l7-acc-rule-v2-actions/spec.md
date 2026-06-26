## ADDED Requirements

### Requirement: Resource supports top-level actions parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL support an optional `actions` parameter at the top level of the resource schema. This parameter defines the list of acceleration rule actions (e.g., Cache, CacheKey, AccessURLRedirect) and maps to the SDK `RuleEngineItem.Branches[0].Actions` field.

#### Scenario: Create resource with actions
- **WHEN** user defines `actions` in a `tencentcloud_teo_l7_acc_rule_v2` resource configuration
- **THEN** the Create function SHALL map actions to `RuleEngineItem.Branches[0].Actions` in the `CreateL7AccRules` API request

#### Scenario: Read resource returns actions
- **WHEN** the `DescribeL7AccRules` API returns rules with `Branches[0].Actions` populated
- **THEN** the Read function SHALL set the `actions` attribute in the Terraform state from the API response

#### Scenario: Update resource actions
- **WHEN** user modifies the `actions` parameter in an existing `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the Update function SHALL map the new actions to `RuleEngineItem.Branches[0].Actions` in the `ModifyL7AccRule` API request

#### Scenario: Actions parameter is optional
- **WHEN** user does not specify `actions` in the resource configuration
- **THEN** the resource SHALL behave exactly as before (no changes to existing behavior)

#### Scenario: Actions parameter schema matches existing RuleEngineAction structure
- **WHEN** user configures `actions` with fields like `name`, `cache_parameters`, `cache_key_parameters`, etc.
- **THEN** the schema SHALL accept the same structure as `branches.actions` (i.e., the `RuleEngineAction` schema from `TencentTeoL7RuleBranchBasicInfo`)
