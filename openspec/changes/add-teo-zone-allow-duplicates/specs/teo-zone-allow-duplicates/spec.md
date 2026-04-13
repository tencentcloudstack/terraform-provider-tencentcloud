## ADDED Requirements

### Requirement: TEO zone resource supports allow_duplicates field
The system SHALL support the `allow_duplicates` field in the `tencentcloud_teo_zone` resource to control whether duplicate rule configurations are allowed in the zone.

#### Scenario: Create zone with allow_duplicates set to true
- **WHEN** user creates a TEO zone with `allow_duplicates = true`
- **THEN** system shall call CreateZone API with `allow_duplicates` parameter set to true
- **THEN** system shall successfully create the zone
- **THEN** the zone shall be configured to allow duplicate rule configurations

#### Scenario: Create zone with allow_duplicates set to false
- **WHEN** user creates a TEO zone with `allow_duplicates = false`
- **THEN** system shall call CreateZone API with `allow_duplicates` parameter set to false
- **THEN** system shall successfully create the zone
- **THEN** the zone shall be configured to not allow duplicate rule configurations

#### Scenario: Create zone without setting allow_duplicates
- **WHEN** user creates a TEO zone without setting `allow_duplicates` field
- **THEN** system shall call CreateZone API without the `allow_duplicates` parameter
- **THEN** system shall successfully create the zone
- **THEN** the zone shall use the API's default value for `allow_duplicates`

### Requirement: TEO zone resource reads allow_duplicates field
The system SHALL read and reflect the `allow_duplicates` field value from the cloud service when reading a TEO zone.

#### Scenario: Read zone returns allow_duplicates value
- **WHEN** user reads an existing TEO zone that has `allow_duplicates` configured
- **THEN** system shall call DescribeZone API to retrieve zone details
- **THEN** system shall read the `allow_duplicates` field from the API response
- **THEN** system shall set the `allow_duplicates` value in the Terraform state

#### Scenario: Read zone with default allow_duplicates
- **WHEN** user reads an existing TEO zone that was created without `allow_duplicates`
- **THEN** system shall call DescribeZone API to retrieve zone details
- **THEN** system shall read the `allow_duplicates` field from the API response
- **THEN** system shall set the actual `allow_duplicates` value (including API default) in the Terraform state

### Requirement: TEO zone resource updates allow_duplicates field
The system SHALL support updating the `allow_duplicates` field in a TEO zone when the field value changes.

#### Scenario: Update allow_duplicates from false to true
- **WHEN** user updates a TEO zone by changing `allow_duplicates` from false to true
- **THEN** system shall detect the change in the `allow_duplicates` field
- **THEN** system shall call ModifyZone API with the new `allow_duplicates` value
- **THEN** system shall successfully update the zone configuration

#### Scenario: Update allow_duplicates from true to false
- **WHEN** user updates a TEO zone by changing `allow_duplicates` from true to false
- **THEN** system shall detect the change in the `allow_duplicates` field
- **THEN** system shall call ModifyZone API with the new `allow_duplicates` value
- **THEN** system shall successfully update the zone configuration

#### Scenario: Update zone without changing allow_duplicates
- **WHEN** user updates other fields of a TEO zone without changing `allow_duplicates`
- **THEN** system shall not call ModifyZone API for the `allow_duplicates` field
- **THEN** system shall update only the changed fields

#### Scenario: Update zone with allow_duplicates not supported by API
- **IF** the ModifyZone API does not support updating `allow_duplicates`
- **WHEN** user attempts to update the `allow_duplicates` field
- **THEN** system shall return an error indicating that this field cannot be updated after creation
- **THEN** system shall provide guidance in the documentation

### Requirement: TEO zone resource delete is unaffected by allow_duplicates
The system SHALL not require special handling for the `allow_duplicates` field when deleting a TEO zone.

#### Scenario: Delete zone with allow_duplicates configured
- **WHEN** user deletes a TEO zone that has `allow_duplicates` configured
- **THEN** system shall call DeleteZone API to delete the zone
- **THEN** the `allow_duplicates` field shall not affect the delete operation
- **THEN** system shall successfully delete the zone

### Requirement: allow_duplicates field is optional and backward compatible
The system SHALL ensure the `allow_duplicates` field is optional and does not break existing configurations.

#### Scenario: Existing zone without allow_duplicates continues to work
- **WHEN** user applies a configuration for an existing TEO zone that does not have `allow_duplicates` set
- **THEN** system shall not require the `allow_duplicates` field to be added
- **THEN** system shall continue to manage the zone without errors
- **THEN** backward compatibility shall be maintained

#### Scenario: Import existing zone
- **WHEN** user imports an existing TEO zone that was created before `allow_duplicates` was added
- **THEN** system shall successfully import the zone
- **THEN** system shall read the actual `allow_duplicates` value from the API
- **THEN** the imported state shall include the correct `allow_duplicates` value
