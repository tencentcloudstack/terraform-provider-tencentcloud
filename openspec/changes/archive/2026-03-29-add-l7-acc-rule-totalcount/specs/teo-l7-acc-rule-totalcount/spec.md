## ADDED Requirements

### Requirement: TotalCount field must be returned from DescribeL7AccRules API
The system SHALL return the TotalCount field from the DescribeL7AccRules API response when querying L7 acceleration rules. This field represents the total number of rules available for the specified zone.

#### Scenario: Successful API response with TotalCount
- **WHEN** DescribeL7AccRules API is called with a valid zone_id
- **THEN** the response MUST include TotalCount field
- **AND** TotalCount MUST be a non-negative integer representing the total number of rules
- **AND** TotalCount value MUST match or be greater than the number of rules in the Rules array

#### Scenario: Empty result set
- **WHEN** no rules exist for the specified zone
- **THEN** TotalCount MUST be 0
- **AND** Rules array MUST be empty or null

### Requirement: TotalCount must be exposed in data source schema
If a data source exists for teo_l7_acc_rule, the system SHALL expose the TotalCount field in the data source schema as a computed field named "total_count".

#### Scenario: Data source returns TotalCount
- **WHEN** user queries the teo_l7_acc_rule data source
- **THEN** the data source MUST return a "total_count" field
- **AND** the "total_count" field MUST be of type integer
- **AND** the value MUST match the TotalCount from the API response

#### Scenario: TotalCount field is computed
- **WHEN** data source schema is defined
- **THEN** "total_count" MUST be marked as Computed
- **AND** "total_count" MUST NOT be marked as Required or Optional (user cannot set it)

### Requirement: Backward compatibility must be maintained
The system MUST maintain backward compatibility when adding the TotalCount field. Existing resources and configurations must continue to work without modification.

#### Scenario: Existing resource read operation
- **WHEN** reading an existing tencentcloud_teo_l7_acc_rule resource
- **THEN** the resource MUST read successfully without requiring TotalCount field
- **AND** the resource state MUST not be affected by the TotalCount field

#### Scenario: Existing data source queries
- **WHEN** existing data source queries are executed
- **THEN** all queries MUST return successfully
- **AND** the TotalCount field MUST be added as an additional output field
- **AND** existing output fields MUST continue to work as before

## REMOVED Requirements

None

## MODIFIED Requirements

None
