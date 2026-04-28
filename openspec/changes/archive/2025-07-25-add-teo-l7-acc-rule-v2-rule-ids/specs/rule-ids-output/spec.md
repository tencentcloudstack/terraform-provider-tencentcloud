## ADDED Requirements

### Requirement: Expose rule_ids computed parameter
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL include a computed parameter `rule_ids` of type `TypeList` with element type `TypeString` that exposes the full list of rule IDs returned by the `CreateL7AccRules` API response field `RuleIds`.

#### Scenario: rule_ids populated after resource creation
- **WHEN** the `tencentcloud_teo_l7_acc_rule_v2` resource is created successfully via `CreateL7AccRules` API
- **THEN** the `rule_ids` attribute SHALL be set to the list of rule IDs from `response.Response.RuleIds`

#### Scenario: rule_ids populated during resource read
- **WHEN** the `tencentcloud_teo_l7_acc_rule_v2` resource is read via `DescribeL7AccRules` API
- **THEN** the `rule_ids` attribute SHALL be set to the list of rule IDs extracted from the response rules data

#### Scenario: rule_ids is computed and not user-settable
- **WHEN** a user writes a Terraform configuration for `tencentcloud_teo_l7_acc_rule_v2`
- **THEN** the `rule_ids` parameter SHALL NOT be settable by the user (Computed only)

#### Scenario: backward compatibility with existing rule_id
- **WHEN** the `rule_ids` parameter is added to the schema
- **THEN** the existing `rule_id` (singular, TypeString, Computed) parameter SHALL remain unchanged and continue to hold the first rule ID from the list
