## ADDED Requirements

### Requirement: Offset parameter in datasource schema
The tencentcloud_teo_l7_acc_rule datasource schema SHALL include an Offset parameter with the following characteristics:
- Type: Int
- Optional: true
- Description: Offset for pagination, specifies the starting position of the query
- Default value: 0 (when not specified)

#### Scenario: Datasource schema contains Offset parameter
- **WHEN** user queries tencentcloud_teo_l7_acc_rule datasource
- **THEN** the datasource schema exposes an optional Offset parameter of type Int

### Requirement: Offset parameter API integration
The datasource SHALL pass the Offset parameter to the DescribeL7AccRules API when it is specified by the user.

#### Scenario: Offset parameter passed to API
- **WHEN** user provides Offset parameter with value 10 in datasource configuration
- **THEN** the datasource calls DescribeL7AccRules API with Offset parameter set to 10

#### Scenario: Default Offset behavior
- **WHEN** user does not specify Offset parameter in datasource configuration
- **THEN** the datasource calls DescribeL7AccRules API without Offset parameter (default to 0)

### Requirement: Offset parameter validation
The datasource SHALL validate that the Offset parameter is a non-negative integer value.

#### Scenario: Valid Offset parameter
- **WHEN** user specifies Offset parameter with value 100
- **THEN** the datasource accepts the value and passes it to the API

#### Scenario: Invalid Offset parameter
- **WHEN** user specifies Offset parameter with value -5
- **THEN** the datasource returns a validation error indicating Offset must be non-negative
