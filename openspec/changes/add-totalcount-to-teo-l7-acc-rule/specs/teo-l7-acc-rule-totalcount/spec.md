## ADDED Requirements

### Requirement: Expose TotalCount parameter in tencentcloud_teo_l7_acc_rule resource
The tencentcloud_teo_l7_acc_rule resource SHALL expose the TotalCount parameter from the DescribeL7AccRules API response. The total_count field SHALL be a Computed field that is automatically populated from the API response during Read operations. The field SHALL return the total number of L7 acceleration rules for the specified zone.

#### Scenario: Read resource with TotalCount populated
- **WHEN** user reads a tencentcloud_teo_l7_acc_rule resource
- **THEN** system SHALL populate total_count field with the TotalCount value from DescribeL7AccRules API response
- **AND** total_count SHALL be a non-negative integer
- **AND** total_count SHALL match the number of rules in the rules list

#### Scenario: TotalCount field is Computed only
- **WHEN** user defines tencentcloud_teo_l7_acc_rule resource configuration
- **THEN** user MUST NOT be able to set total_count field in configuration
- **AND** total_count field SHALL be marked as Computed in schema
- **AND** Terraform SHALL automatically update total_count during refresh and apply operations

#### Scenario: Handle nil TotalCount from API
- **WHEN** DescribeL7AccRules API returns response with TotalCount field as nil
- **THEN** system SHALL handle the nil value without causing panic
- **AND** total_count field SHALL either be omitted from state or set to a default value (0)

#### Scenario: Backward compatibility maintained
- **WHEN** user has existing tencentcloud_teo_l7_acc_rule resources without total_count in state
- **THEN** upgrading provider version SHALL NOT require any configuration changes
- **AND** existing resources SHALL continue to work normally
- **AND** total_count field SHALL be automatically populated during next read operation
