## ADDED Requirements

### Requirement: Export zone configuration
The system SHALL allow users to export TEO (Tencent Edge One) zone configuration using the `tencentcloud_teo_export_zone_config` data source. The data source MUST accept a `zone_id` parameter to specify which zone to export, and MUST return the exported configuration content as a JSON string.

#### Scenario: Export all zone configuration
- **WHEN** user provides only the `zone_id` parameter without specifying `types`
- **THEN** the system MUST export all available configuration types for the zone
- **AND** the system MUST return the configuration content in the `content` output attribute
- **AND** the content MUST be a valid JSON string

#### Scenario: Export specific configuration types
- **WHEN** user provides `zone_id` and `types` parameters
- **THEN** the system MUST export only the specified configuration types
- **AND** the system MUST filter the exported content based on the provided `types` list
- **AND** the system MUST return the filtered configuration content in the `content` output attribute

### Requirement: Validate zone_id parameter
The system SHALL validate that the `zone_id` parameter is provided and is a valid non-empty string. The system MUST return a clear error message if the `zone_id` is missing or invalid.

#### Scenario: Missing zone_id parameter
- **WHEN** user does not provide the `zone_id` parameter
- **THEN** the system MUST return a validation error
- **AND** the error message MUST indicate that `zone_id` is required

#### Scenario: Invalid zone_id format
- **WHEN** user provides an empty string or invalid format for `zone_id`
- **THEN** the system MUST return a validation error
- **AND** the error message MUST indicate that `zone_id` is invalid

### Requirement: Handle types parameter
The system SHALL accept an optional `types` parameter as a list of strings. The parameter MUST specify which configuration types to export. If not provided, the system MUST export all available configuration types.

#### Scenario: Export with empty types list
- **WHEN** user provides an empty `types` list
- **THEN** the system MUST treat it as exporting all configuration types

#### Scenario: Export with single type
- **WHEN** user provides a `types` list with a single configuration type (e.g., `["L7AccelerationConfig"]`)
- **THEN** the system MUST export only that specific configuration type
- **AND** the returned content MUST contain only the requested configuration type

#### Scenario: Export with multiple types
- **WHEN** user provides a `types` list with multiple configuration types
- **THEN** the system MUST export all the specified configuration types
- **AND** the returned content MUST contain all requested configuration types

### Requirement: Handle API errors
The system SHALL properly handle errors from the ExportZoneConfig API and return meaningful error messages to the user. Common errors include invalid zone ID, permission denied, and service unavailable.

#### Scenario: Zone not found
- **WHEN** the provided `zone_id` does not exist
- **THEN** the system MUST return an error indicating the zone was not found
- **AND** the error message MUST include the provided `zone_id`

#### Scenario: Permission denied
- **WHEN** the user does not have permission to access the zone
- **THEN** the system MUST return a permission error
- **AND** the error message MUST indicate insufficient permissions

#### Scenario: Service unavailable
- **WHEN** the TEO service is temporarily unavailable
- **THEN** the system MUST return an error indicating service unavailability
- **AND** the system MAY retry the request with appropriate backoff

### Requirement: Return configuration content
The system MUST return the exported configuration content as a string in the `content` output attribute. The content MUST be in JSON format and UTF-8 encoded.

#### Scenario: Valid configuration content
- **WHEN** the export operation is successful
- **THEN** the system MUST return the configuration content in the `content` attribute
- **AND** the content MUST be a valid JSON string
- **AND** the content MUST be UTF-8 encoded

#### Scenario: Empty configuration
- **WHEN** the zone has no configuration to export
- **THEN** the system MUST return an empty or minimal JSON structure in the `content` attribute
- **AND** the content MUST still be a valid JSON string
