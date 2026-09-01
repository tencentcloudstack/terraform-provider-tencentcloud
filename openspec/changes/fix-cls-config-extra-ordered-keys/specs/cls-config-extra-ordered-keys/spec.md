## ADDED Requirements

### Requirement: ConfigExtra extraction keys SHALL preserve positional order

The `tencentcloud_cls_config_extra.extract_rule.keys` attribute SHALL be an
ordered list of strings. Create and Update SHALL send keys to CLS in
configuration order, and Read SHALL store keys in CLS API response order.

#### Scenario: Create preserves declared key order

- **WHEN** a user creates a ConfigExtra with `keys = ["first", "second"]`
- **THEN** the Create request sends `ExtractRule.Keys` as `first`, then
  `second`
- **AND** the refreshed Terraform state stores `first` at index 0 and `second`
  at index 1

#### Scenario: Reordering keys updates in place

- **GIVEN** an existing ConfigExtra has `keys = ["first", "second"]`
- **WHEN** the user changes the configuration to
  `keys = ["second", "first"]`
- **THEN** Terraform detects a meaningful in-place change
- **AND** the Update request sends `second`, then `first`
- **AND** a subsequent refresh and plan converge without replacement or a
  perpetual diff

#### Scenario: Import preserves API key order

- **GIVEN** CLS returns extraction keys in a defined order
- **WHEN** the ConfigExtra is imported or refreshed
- **THEN** Terraform state stores the keys in the same order

### Requirement: Version 0 ConfigExtra state SHALL upgrade without replacement

The resource SHALL declare schema version 1 and a version 0 state upgrader whose
legacy type models `extract_rule.keys` as a set. The upgrade SHALL accept
supported JSON and legacy flatmap state, retain all stored keys, and encode the
result using the current list schema without error or resource replacement.

#### Scenario: JSON state upgrades through the provider protocol

- **GIVEN** version 0 JSON state contains extraction keys represented as an
  array
- **WHEN** SDKv2 executes `UpgradeResourceState`
- **THEN** the response contains no error diagnostics
- **AND** the upgraded msgpack decodes using the version 1 schema
- **AND** the decoded key list contains the same values in the JSON array order

#### Scenario: Legacy flatmap state upgrades through the provider protocol

- **GIVEN** version 0 flatmap state contains extraction keys represented by set
  hashes
- **WHEN** SDKv2 executes `UpgradeResourceState`
- **THEN** the response contains no error diagnostics
- **AND** the upgraded msgpack decodes using the version 1 schema
- **AND** every legacy set value is present in the decoded key list

#### Scenario: Lost version 0 declaration order is not fabricated

- **GIVEN** version 0 set state no longer contains the user's original HCL
  order
- **WHEN** the state is upgraded and refreshed
- **THEN** the provider uses the order available from decoded state and the CLS
  API
- **AND** Terraform may plan one in-place correction when current configuration
  order differs
- **AND** the provider SHALL NOT suppress that corrective diff or replace the
  resource
