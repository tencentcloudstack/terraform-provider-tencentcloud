# teo-security-policy-bot-action-override-rule-ids Specification

## Purpose
TBD - created by archiving change modify-teo-security-policy-bot-rule-ids. Update Purpose after archive.
## Requirements
### Requirement: action_overrides rule_ids schema field
The `tencentcloud_teo_security_policy_config` resource's `bot_management.basic_bot_settings` sub-blocks (`source_idc`, `search_engine_bots`, `known_bot_categories`, `ip_reputation.ip_reputation_group`) SHALL define the `action_overrides` block using a schema where each override element contains:
- `rule_ids` (TypeList, Required, Elem: TypeString): One or more bot rule IDs (or category IDs) whose action is overridden. Maps to the cloud API field `BotManagementActionOverrides.Ids` (`[]*string`).
- `action` (TypeList, Optional, MaxItems: 1): Action override configuration, maps to `BotManagementActionOverrides.Action`.

The legacy `rule_id` (TypeString) field SHALL be removed and replaced by `rule_ids`.

#### Scenario: Schema defines rule_ids as a list
- **WHEN** the resource schema for `action_overrides` is defined
- **THEN** it SHALL include `rule_ids` of TypeList with Elem TypeString as a Required field
- **AND** it SHALL NOT include a `rule_id` field

#### Scenario: User configures multiple rule IDs
- **WHEN** a user sets `rule_ids = ["rule-a", "rule-b"]` in an `action_overrides` block
- **THEN** the schema SHALL accept the list and mark it as a candidate for create/update

#### Scenario: User configures a single rule ID
- **WHEN** a user sets `rule_ids = ["rule-a"]` in an `action_overrides` block
- **THEN** the schema SHALL accept the single-element list

### Requirement: buildBotManagementActionOverrideFromMap supports multiple IDs
The `buildBotManagementActionOverrideFromMap` function SHALL read the `rule_ids` list from the schema map and convert every non-empty element into a `[]*string` assigned to `BotManagementActionOverrides.Ids`. It SHALL NOT read a `rule_id` string field.

#### Scenario: Build override with multiple IDs
- **WHEN** the schema map contains `rule_ids = ["rule-a", "rule-b", "rule-c"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL be `["rule-a", "rule-b", "rule-c"]` (as `[]*string`)

#### Scenario: Build override with single ID
- **WHEN** the schema map contains `rule_ids = ["rule-a"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL be `["rule-a"]` (as `[]*string`)

#### Scenario: Build override with empty rule_ids
- **WHEN** the schema map contains an empty or absent `rule_ids`
- **THEN** `BotManagementActionOverrides.Ids` SHALL be nil (no IDs assigned)

### Requirement: flattenBotManagementActionOverride supports multiple IDs
The `flattenBotManagementActionOverride` function SHALL take the cloud API `BotManagementActionOverrides.Ids` (`[]*string`) and write every element into the schema map's `rule_ids` field as a `[]string`. It SHALL NOT only read the first element.

#### Scenario: Flatten override with multiple IDs
- **WHEN** the API response contains `BotManagementActionOverrides` with `Ids = ["rule-a", "rule-b"]`
- **THEN** the flattened schema map SHALL contain `rule_ids = ["rule-a", "rule-b"]`

#### Scenario: Flatten override with single ID
- **WHEN** the API response contains `BotManagementActionOverrides` with `Ids = ["rule-a"]`
- **THEN** the flattened schema map SHALL contain `rule_ids = ["rule-a"]`

#### Scenario: Flatten override with nil Ids
- **WHEN** the API response contains `BotManagementActionOverrides` with `Ids = nil`
- **THEN** the flattened schema map SHALL NOT set `rule_ids`

### Requirement: Multiple action_overrides entries are preserved
The Read/Create/Update operations SHALL preserve the ability to have multiple `action_overrides` entries in each `basic_bot_settings` sub-block, with each entry independently carrying its own `rule_ids` list and `action`.

#### Scenario: Multiple override entries each with multiple rule IDs
- **WHEN** a user configures two `action_overrides` blocks, one with `rule_ids = ["a", "b"]` and another with `rule_ids = ["c"]`
- **THEN** the Create/Update operation SHALL produce two `BotManagementActionOverrides` entries with the corresponding `Ids` arrays

#### Scenario: Read reflects multiple override entries
- **WHEN** the DescribeSecurityPolicy API returns multiple `BotManagementActionOverrides` entries each with multiple `Ids`
- **THEN** the Read function SHALL flatten each entry into a separate `action_overrides` list element with its own `rule_ids`

### Requirement: Unit tests for rule_ids
The system SHALL provide unit tests in `resource_tc_teo_security_policy_config_test.go` using gomonkey to mock cloud API calls, covering:
- Multiple `rule_ids` in an `action_overrides` entry flowing through Create/Update (build path)
- Multiple `rule_ids` being flattened from a Read response (flatten path)

#### Scenario: Build path with multiple rule_ids test
- **WHEN** the test exercises the build path with `rule_ids = ["rule-a", "rule-b"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL contain exactly two entries matching the input

#### Scenario: Flatten path with multiple rule_ids test
- **WHEN** the test exercises the flatten path with an API override containing `Ids = ["rule-a", "rule-b"]`
- **THEN** the flattened map SHALL contain `rule_ids` with exactly two entries matching the input

