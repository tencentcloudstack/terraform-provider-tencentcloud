# teo-l7-acc-rule-filters-param Specification

## ADDED Requirements

### Requirement: Filters Parameter Support
The `tencentcloud_teo_l7_acc_rule` data source SHALL support a `filters` parameter that allows users to filter query results using key-value pairs when calling the `DescribeL7AccRules` API.

#### Scenario: Filter by rule name
- **WHEN** user provides a `filters` block with `name = "RuleName"` and `values = ["my-rule"]`
- **THEN** the data source queries the API with the Filters parameter
- **AND** only rules matching the specified name are returned

#### Scenario: Filter by rule status
- **WHEN** user provides a `filters` block with `name = "Status"` and `values = ["Enabled"]`
- **THEN** the data source queries the API with the Filters parameter
- **AND** only enabled rules are returned

#### Scenario: Multiple filters
- **WHEN** user provides multiple `filters` blocks with different keys
- **THEN** all filters are combined with AND logic
- **AND** results match all specified criteria

#### Scenario: Empty filters
- **WHEN** user provides an empty `filters` block or omits it
- **THEN** the data source queries the API without Filters parameter
- **AND** all available rules are returned

#### Scenario: Filter with multiple values
- **WHEN** user provides a `filters` block with `values = ["value1", "value2"]`
- **THEN** the data source queries the API with the Filters parameter
- **AND** results match any of the specified values (OR logic within the filter)
