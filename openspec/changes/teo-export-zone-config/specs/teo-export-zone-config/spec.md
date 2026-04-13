## ADDED Requirements

### Requirement: Export zone configuration
The system SHALL provide a Terraform resource `tencentcloud_teo_export_zone_config` that allows users to export EdgeOne zone configuration.

#### Scenario: Successful export zone configuration
- **WHEN** user provides valid `zone_id` and `export_type` parameters
- **THEN** system shall successfully create the export resource
- **AND** the resource ID shall be in format `zoneId#exportType`
- **AND** the exported configuration shall be available in the resource state

#### Scenario: Export with invalid zone ID
- **WHEN** user provides invalid `zone_id` parameter
- **THEN** system shall return error indicating invalid zone ID
- **AND** no resource shall be created

### Requirement: Read exported zone configuration
The system SHALL support reading the exported zone configuration through Terraform refresh operation.

#### Scenario: Successful read of existing export
- **WHEN** user performs Terraform refresh on existing `tencentcloud_teo_export_zone_config` resource
- **THEN** system shall retrieve the latest exported configuration from CAPI
- **AND** the resource state shall be updated with the latest data

#### Scenario: Read non-existent export
- **WHEN** user attempts to read a non-existent export resource
- **THEN** system shall return error indicating resource not found
- **AND** the resource shall be marked as removed in Terraform state

### Requirement: Update export configuration
The system SHALL support updating the export configuration parameters.

#### Scenario: Update export type
- **WHEN** user modifies `export_type` parameter
- **THEN** system shall update the export configuration with new type
- **AND** the resource ID shall remain the same if zone_id is unchanged

### Requirement: Delete export resource
The system SHALL support deleting the export resource from Terraform state.

#### Scenario: Successful delete
- **WHEN** user deletes the `tencentcloud_teo_export_zone_config` resource
- **THEN** system shall remove the resource from Terraform state
- **AND** the exported configuration file on cloud side shall remain intact (as it's a read-only export operation)

### Requirement: Schema definition
The system SHALL define a schema that matches the CAPI interface parameters.

#### Scenario: Required parameters
- **WHEN** user creates the resource
- **THEN** `zone_id` parameter SHALL be required
- **AND** `export_type` parameter SHALL be required

#### Scenario: Optional parameters
- **WHEN** user creates the resource
- **THEN** optional parameters SHALL have appropriate default values
- **AND** all parameters SHALL match CAPI interface definition

### Requirement: Asynchronous operation handling
The system SHALL handle asynchronous operations with proper retry and timeout mechanisms.

#### Scenario: Export operation in progress
- **WHEN** export operation is in progress
- **THEN** system shall retry reading the operation status
- **AND** system shall respect the timeout configuration
- **AND** operation shall complete successfully or return appropriate error

### Requirement: Error handling
The system SHALL provide clear error messages for all failure scenarios.

#### Scenario: API error
- **WHEN** CAPI API returns error
- **THEN** system shall propagate the error message to user
- **AND** error message shall include relevant context (zone_id, operation type)

#### Scenario: Network timeout
- **WHEN** API call times out
- **THEN** system shall retry according to retry policy
- **AND** if all retries fail, system shall return timeout error
