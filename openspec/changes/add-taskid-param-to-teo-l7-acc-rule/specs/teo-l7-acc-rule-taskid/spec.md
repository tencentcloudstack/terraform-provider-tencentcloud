## ADDED Requirements

### Requirement: TaskId parameter support
The tencentcloud_teo_l7_acc_rule resource SHALL support an optional TaskId parameter of type String. The TaskId parameter SHALL be used to identify asynchronous tasks when calling the ImportZoneConfig API during update operations. The TaskId parameter MUST be optional to maintain backward compatibility with existing configurations.

#### Scenario: Add TaskId parameter to resource schema
- **WHEN** user specifies the TaskId parameter in tencentcloud_teo_l7_acc_rule resource configuration
- **THEN** the resource schema SHALL accept and store the TaskId value as a String type

#### Scenario: Update resource with TaskId parameter
- **WHEN** user updates tencentcloud_teo_l7_acc_rule resource with TaskId parameter specified
- **THEN** the provider SHALL pass the TaskId value to ImportZoneConfig API call
- **AND** the API call SHALL include the TaskId in the request parameters

#### Scenario: Backward compatibility without TaskId
- **WHEN** user configures tencentcloud_teo_l7_acc_rule resource without specifying TaskId parameter
- **THEN** the resource SHALL function normally without TaskId
- **AND** the ImportZoneConfig API call SHALL proceed without TaskId parameter

#### Scenario: TaskId validation
- **WHEN** user provides an invalid TaskId value (e.g., empty string if specified, special characters)
- **THEN** the provider SHALL validate the TaskId according to API constraints
- **AND** appropriate validation errors SHALL be returned to the user
