## ADDED Requirements

### Requirement: Resource exposes computed rules attribute

The `tencentcloud_teo_l7_acc_rule_v2` resource SHALL expose a computed `rules` attribute that contains the full rules list returned by the `DescribeL7AccRules` API response field `response.Rules`.

The `rules` attribute SHALL be a list where each element represents a `RuleEngineItem` with the following fields:
- `status` (string): Rule status, either "enable" or "disable".
- `rule_id` (string): Rule ID, unique identifier of the rule.
- `rule_name` (string): Rule name.
- `description` (list of strings): Rule annotations.
- `branches` (list): Sub-rule branches containing condition, actions, and sub_rules.
- `rule_priority` (int): Rule priority, used for ordering.

The attribute SHALL be `Computed` only (not configurable by users) and SHALL be populated during the Read operation.

#### Scenario: Rules attribute is populated after resource creation

- **WHEN** a `tencentcloud_teo_l7_acc_rule_v2` resource is created and the Read function is called
- **THEN** the `rules` attribute SHALL contain the full rules list returned by `DescribeL7AccRules` for the given zone_id and rule_id, with each rule element containing status, rule_id, rule_name, description, branches, and rule_priority fields

#### Scenario: Rules attribute reflects current API state on refresh

- **WHEN** the resource state is refreshed via `terraform refresh` or `terraform plan`
- **THEN** the `rules` attribute SHALL be updated to reflect the current state of rules as returned by the `DescribeL7AccRules` API

#### Scenario: Rules attribute handles nil response gracefully

- **WHEN** the `DescribeL7AccRules` API returns nil or empty `Rules` in the response
- **THEN** the `rules` attribute SHALL be set to an empty list without causing an error
