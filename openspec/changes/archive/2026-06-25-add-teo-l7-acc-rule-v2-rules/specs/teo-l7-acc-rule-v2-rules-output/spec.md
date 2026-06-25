## ADDED Requirements

### Requirement: Rules output field
The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL expose a computed `rules` field that mirrors the `response.Rules` output from the `DescribeL7AccRules` API. This field SHALL be a TypeList containing `RuleEngineItem` structured blocks, providing users with a unified view of the complete rule data as returned by the cloud API.

#### Scenario: Rules field populated after create
- **WHEN** a user creates a `tencentcloud_teo_l7_acc_rule_v2` resource and Terraform reads the resource state
- **THEN** the `rules` field in the Terraform state SHALL contain the rule data as returned by `DescribeL7AccRules`, with each `RuleEngineItem` block containing `rule_id`, `status`, `rule_name`, `description`, `branches`, and `rule_priority`

#### Scenario: Rules field populated after read
- **WHEN** a user runs `terraform refresh` or `terraform plan` on an existing `tencentcloud_teo_l7_acc_rule_v2` resource
- **THEN** the `rules` field SHALL reflect the current state of the rule from the cloud API

#### Scenario: Rules field empty when resource not found
- **WHEN** the `DescribeL7AccRules` API returns an empty or nil `Rules` list
- **THEN** the `rules` field SHALL be empty and the resource SHALL be removed from state (id set to `""`)
