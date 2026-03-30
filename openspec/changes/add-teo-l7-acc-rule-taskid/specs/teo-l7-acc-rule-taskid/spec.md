## ADDED Requirements

### Requirement: Update operation supports TaskId parameter
The tencentcloud_teo_l7_acc_rule resource SHALL support an optional TaskId parameter in the update operation. When provided, the TaskId value SHALL be passed to the ImportZoneConfig API.

#### Scenario: Update with TaskId provided
- **WHEN** user specifies a TaskId in the tencentcloud_teo_l7_acc_rule resource update configuration
- **THEN** the provider shall pass the TaskId value to the ImportZoneConfig API call

#### Scenario: Update without TaskId
- **WHEN** user does not specify a TaskId in the tencentcloud_teo_l7_acc_rule resource update configuration
- **THEN** the provider shall call ImportZoneConfig API without TaskId parameter
