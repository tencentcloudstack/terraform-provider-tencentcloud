## ADDED Requirements

### Requirement: Resource SHALL provide total_count output
The tencentcloud_teo_l7_acc_rule resource SHALL provide a `total_count` computed field that returns the total number of L7 access control rules for the specified zone.

#### Scenario: Read resource returns total_count
- **WHEN** user reads the tencentcloud_teo_l7_acc_rule resource for a zone
- **THEN** the resource SHALL return a `total_count` field with the value from DescribeL7AccRules API response's TotalCount field
- **AND** the `total_count` field SHALL be of type integer

### Requirement: total_count field MUST be computed only
The `total_count` field SHALL be computed only and SHALL NOT allow user input or configuration.

#### Scenario: Attempt to set total_count fails
- **WHEN** user attempts to set `total_count` in the resource configuration
- **THEN** the configuration SHALL be rejected as the field is computed-only
- **AND** the field SHALL be automatically populated during read operations

### Requirement: Backward compatibility MUST be maintained
The addition of `total_count` field SHALL NOT break existing configurations or state files.

#### Scenario: Existing configuration without total_count
- **WHEN** user applies an existing configuration that does not reference `total_count`
- **THEN** the configuration SHALL continue to work without modification
- **AND** the `total_count` field SHALL be automatically added to the state during next refresh

#### Scenario: State refresh includes total_count
- **WHEN** user runs `terraform refresh` on an existing state
- **THEN** the state SHALL be updated with the `total_count` value
- **AND** no changes SHALL be required to the Terraform configuration

### Requirement: total_count reflects API TotalCount
The `total_count` field SHALL accurately reflect the TotalCount value returned by the DescribeL7AccRules API.

#### Scenario: API returns total_count
- **WHEN** DescribeL7AccRules API returns TotalCount = 10
- **THEN** the resource's `total_count` field SHALL be set to 10

#### Scenario: API returns nil TotalCount
- **WHEN** DescribeL7AccRules API returns TotalCount = nil
- **THEN** the resource's `total_count` field SHALL be set to 0 (default integer value)
