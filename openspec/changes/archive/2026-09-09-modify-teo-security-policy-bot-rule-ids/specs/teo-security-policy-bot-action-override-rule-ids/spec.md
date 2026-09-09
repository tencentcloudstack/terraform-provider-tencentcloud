## ADDED Requirements

### Requirement: action_overrides rule_ids schema field
The `tencentcloud_teo_security_policy_config` resource's `bot_management.basic_bot_settings` sub-blocks (`source_idc`, `search_engine_bots`, `known_bot_categories`, `ip_reputation.ip_reputation_group`) SHALL define the `action_overrides` block using a schema where each override element contains:
- `rule_ids` (TypeList, Optional, Elem: TypeString): One or more bot rule IDs (or category IDs) whose action is overridden. Maps to the cloud API field `BotManagementActionOverrides.Ids` (`[]*string`).
- `rule_id` (TypeString, Optional, Deprecated): A single bot rule ID (or category ID) whose action is overridden. Deprecated: use `rule_ids` to support one or more IDs. Retained for backward compatibility.
- `action` (TypeList, Optional, MaxItems: 1): Action override configuration, maps to `BotManagementActionOverrides.Action`.

The `rule_id` field SHALL be marked Deprecated and the `rule_ids` field SHALL be the preferred multi-ID field.

#### Scenario: Schema defines both rule_ids and a deprecated rule_id
- **WHEN** the resource schema for `action_overrides` is defined
- **THEN** it SHALL include `rule_ids` of TypeList with Elem TypeString as an Optional field
- **AND** it SHALL include `rule_id` of TypeString as an Optional field marked Deprecated

#### Scenario: User configures multiple rule IDs
- **WHEN** a user sets `rule_ids = ["rule-a", "rule-b"]` in an `action_overrides` block
- **THEN** the schema SHALL accept the list and mark it as a candidate for create/update

#### Scenario: User configures a single rule ID
- **WHEN** a user sets `rule_ids = ["rule-a"]` in an `action_overrides` block
- **THEN** the schema SHALL accept the single-element list

### Requirement: buildBotManagementActionOverrideFromMap prefers rule_ids with rule_id fallback
The `buildBotManagementActionOverrideFromMap` function SHALL read the `rule_ids` list from the schema map first and convert every non-empty element into a `[]*string` assigned to `BotManagementActionOverrides.Ids`. When `rule_ids` is absent or empty, it SHALL fall back to reading the deprecated `rule_id` string field and construct a single-element `[]*string`.

#### Scenario: Build override with multiple IDs
- **WHEN** the schema map contains `rule_ids = ["rule-a", "rule-b", "rule-c"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL be `["rule-a", "rule-b", "rule-c"]` (as `[]*string`)

#### Scenario: Build override with single ID via rule_ids
- **WHEN** the schema map contains `rule_ids = ["rule-a"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL be `["rule-a"]` (as `[]*string`)

#### Scenario: Build override with empty rule_ids
- **WHEN** the schema map contains an empty or absent `rule_ids` and no `rule_id`
- **THEN** `BotManagementActionOverrides.Ids` SHALL be nil (no IDs assigned)

#### Scenario: Build override with deprecated rule_id fallback
- **WHEN** the schema map contains no `rule_ids` (or an empty list) and `rule_id = "rule-legacy"`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL be `["rule-legacy"]` (as `[]*string`)

#### Scenario: Build override rule_ids takes precedence over rule_id
- **WHEN** the schema map contains both `rule_ids = ["rule-a", "rule-b"]` and `rule_id = "rule-legacy"`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL be `["rule-a", "rule-b"]` (as `[]*string`)

### Requirement: flattenBotManagementActionOverride supports multiple IDs
The `flattenBotManagementActionOverride` function SHALL take the cloud API `BotManagementActionOverrides.Ids` (`[]*string`) and write every element into the schema map's `rule_ids` field as a `[]string`. It SHALL NOT only read the first element. It SHALL NOT write the deprecated `rule_id` field.

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

### Requirement: Unit tests for rule_ids and rule_id fallback
The system SHALL provide unit tests in `resource_tc_teo_security_policy_config_test.go` using gomonkey to mock cloud API calls, covering:
- Multiple `rule_ids` in an `action_overrides` entry flowing through Create/Update (build path)
- Multiple `rule_ids` being flattened from a Read response (flatten path)
- The deprecated `rule_id` fallback path in the build function
- The `rule_ids` precedence over `rule_id` in the build function

#### Scenario: Build path with multiple rule_ids test
- **WHEN** the test exercises the build path with `rule_ids = ["rule-a", "rule-b"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL contain exactly two entries matching the input

#### Scenario: Flatten path with multiple rule_ids test
- **WHEN** the test exercises the flatten path with an API override containing `Ids = ["rule-a", "rule-b"]`
- **THEN** the flattened map SHALL contain `rule_ids` with exactly two entries matching the input

#### Scenario: Build path with deprecated rule_id fallback test
- **WHEN** the test exercises the build path with only `rule_id = "rule-legacy"` (no `rule_ids`)
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL contain exactly one entry matching the input

#### Scenario: Build path rule_ids precedence test
- **WHEN** the test exercises the build path with both `rule_id = "rule-legacy"` and `rule_ids = ["rule-a", "rule-b"]`
- **THEN** the resulting `BotManagementActionOverrides.Ids` SHALL contain exactly `["rule-a", "rule-b"]`
