## ADDED Requirements

### Requirement: Total count field must be available in resource state
The tencentcloud_teo_l7_acc_rule resource SHALL expose the total number of L7 access rules as a computed field named `total_count` in the resource schema.

#### Scenario: Total count is populated from API response
- **WHEN** the provider reads the tencentcloud_teo_l7_acc_rule resource
- **THEN** the `total_count` field SHALL be populated with the value from the DescribeL7AccRules API response's TotalCount field

#### Scenario: Total count is read-only
- **WHEN** a user attempts to set the `total_count` field in their Terraform configuration
- **THEN** the provider SHALL ignore the user-provided value and use the API response value instead

#### Scenario: Total count type is integer
- **WHEN** the provider exposes the `total_count` field
- **THEN** the field SHALL be of type Int to match the API's int64 response type

### Requirement: Total count must handle nil API responses gracefully
The provider SHALL handle cases where the API returns a nil value for TotalCount without causing errors.

#### Scenario: Total count is nil in API response
- **WHEN** the DescribeL7AccRules API returns a nil value for TotalCount
- **THEN** the provider SHALL NOT set the total_count field in the resource state (leave it unset)

### Requirement: Total count must be documented
The provider SHALL include documentation for the total_count field in the resource documentation.

#### Scenario: Documentation includes total_count
- **WHEN** the resource documentation is generated or updated
- **THEN** the documentation SHALL include a description of the total_count field, indicating it is computed and represents the total number of L7 access rules

### Requirement: Total count must be testable
The provider SHALL include acceptance tests to verify the total_count field is correctly populated.

#### Scenario: Acceptance test verifies total_count
- **WHEN** acceptance tests run for the tencentcloud_teo_l7_acc_rule resource
- **THEN** tests SHALL verify that the total_count field is present in the resource state
- **AND** tests SHALL verify that the value is greater than or equal to the number of rules in the state
