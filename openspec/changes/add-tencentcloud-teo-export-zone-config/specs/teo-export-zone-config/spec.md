# TEO Export Zone Config Specification

## ADDED Requirements

### Requirement: User can export zone configuration
The system SHALL allow users to export TEO zone configuration by providing the zone ID and optional configuration types.

#### Scenario: Export configuration with zone ID only
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource with only `zone_id` specified
- **THEN** system SHALL export all configuration types for the specified zone
- **AND** system SHALL return the configuration content in JSON format

#### Scenario: Export configuration with zone ID and types
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource with `zone_id` and `types` specified
- **THEN** system SHALL export only the specified configuration types for the zone
- **AND** system SHALL return the configuration content in JSON format

#### Scenario: Export configuration fails with invalid zone ID
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource with an invalid `zone_id`
- **THEN** system SHALL return an error indicating the zone ID is invalid

#### Scenario: Export configuration with empty types list
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource with `types` set to an empty list
- **THEN** system SHALL export all configuration types for the specified zone

#### Scenario: Export configuration with unsupported type
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource with an unsupported type in `types` list
- **THEN** system SHALL return an error indicating the type is not supported

### Requirement: User can read exported zone configuration
The system SHALL allow users to read the previously exported zone configuration from the state.

#### Scenario: Read existing exported configuration
- **WHEN** user runs `terraform refresh` on an existing `tencentcloud_teo_export_zone_config` resource
- **THEN** system SHALL re-export the configuration and update the state
- **AND** system SHALL maintain the exported content in the state

#### Scenario: Read configuration after zone modification
- **WHEN** the zone configuration has been modified in the cloud
- **AND** user runs `terraform refresh` on the `tencentcloud_teo_export_zone_config` resource
- **THEN** system SHALL export the latest configuration from the cloud
- **AND** system SHALL update the state with the new content

### Requirement: User can update exported zone configuration parameters
The system SHALL allow users to update the parameters (zone_id, types) of the export operation.

#### Scenario: Update types parameter
- **WHEN** user modifies the `types` parameter in the `tencentcloud_teo_export_zone_config` resource
- **THEN** system SHALL re-export the configuration with the new types
- **AND** system SHALL update the state with the new content

#### Scenario: Update zone_id parameter
- **WHEN** user modifies the `zone_id` parameter in the `tencentcloud_teo_export_zone_config` resource
- **THEN** system SHALL re-export the configuration for the new zone
- **AND** system SHALL update the state with the new content

### Requirement: User can delete exported zone configuration from state
The system SHALL allow users to delete the exported zone configuration from the Terraform state.

#### Scenario: Delete resource from state
- **WHEN** user deletes the `tencentcloud_teo_export_zone_config` resource
- **THEN** system SHALL remove the resource from the state
- **AND** system SHALL NOT delete the configuration from the cloud (export operation is read-only)

### Requirement: System must handle large configuration content
The system SHALL handle large configuration content returned from the export API.

#### Scenario: Export large configuration
- **WHEN** the exported configuration content is large (e.g., multiple configuration types)
- **THEN** system SHALL successfully store the content in the state
- **AND** system SHALL not truncate or modify the content

### Requirement: Zone ID must be required
The system SHALL require the `zone_id` parameter when creating the export resource.

#### Scenario: Create resource without zone ID
- **WHEN** user attempts to create a `tencentcloud_teo_export_zone_config` resource without specifying `zone_id`
- **THEN** system SHALL return a validation error indicating zone_id is required

### Requirement: Types parameter must be optional
The system SHALL make the `types` parameter optional when creating the export resource.

#### Scenario: Create resource without types parameter
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource without specifying `types`
- **THEN** system SHALL not return a validation error
- **AND** system SHALL export all configuration types

### Requirement: Content must be computed
The system SHALL set the `content` parameter as computed (read-only) since it is returned from the API.

#### Scenario: Attempt to set content parameter
- **WHEN** user attempts to set the `content` parameter in the configuration
- **THEN** system SHALL ignore the user-provided value
- **AND** system SHALL use the content returned from the API

### Requirement: System must support retry on transient errors
The system SHALL retry the export operation when encountering transient API errors.

#### Scenario: Transient API error
- **WHEN** the export API returns a transient error (e.g., network timeout)
- **THEN** system SHALL retry the operation with exponential backoff
- **AND** system SHALL eventually succeed or return a final error after maximum retries
