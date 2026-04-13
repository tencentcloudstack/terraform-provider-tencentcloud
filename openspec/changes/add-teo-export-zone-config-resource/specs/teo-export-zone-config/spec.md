## ADDED Requirements

### Requirement: Export zone config
The system SHALL provide a Terraform resource that allows users to export TEO (EdgeOne) zone configuration. The resource SHALL support Create, Read, Update, and Delete operations as defined by the CAPI interface. The resource schema SHALL match the CAPI interface parameters, maintaining the same Required/Optional attributes.

#### Scenario: Create export zone config resource
- **WHEN** user creates a `tencentcloud_teo_export_zone_config` resource with valid zone ID and export parameters
- **THEN** system creates the resource and calls the corresponding CAPI API to export the zone configuration
- **THEN** system stores the exported configuration in the Terraform state

#### Scenario: Read export zone config resource
- **WHEN** user reads an existing `tencentcloud_teo_export_zone_config` resource
- **THEN** system retrieves the current configuration from the CAPI API
- **THEN** system updates the Terraform state with the latest configuration

#### Scenario: Update export zone config resource
- **WHEN** user updates the `tencentcloud_teo_export_zone_config` resource with new export parameters
- **THEN** system calls the CAPI API to export the zone configuration with new parameters
- **THEN** system updates the Terraform state with the new configuration

#### Scenario: Delete export zone config resource
- **WHEN** user deletes the `tencentcloud_teo_export_zone_config` resource
- **THEN** system removes the resource from the Terraform state
- **THEN** system calls the corresponding CAPI API if required

### Requirement: Resource schema validation
The system SHALL validate the resource schema according to the CAPI interface definition. Required parameters MUST be provided by the user, and optional parameters SHALL have default values where applicable.

#### Scenario: Create resource with required parameters
- **WHEN** user creates a resource with all required parameters
- **THEN** system accepts the request and creates the resource

#### Scenario: Create resource missing required parameters
- **WHEN** user creates a resource without all required parameters
- **THEN** system rejects the request with a validation error

#### Scenario: Create resource with optional parameters
- **WHEN** user creates a resource with optional parameters
- **THEN** system accepts the request and applies default values for unspecified optional parameters

### Requirement: Error handling and retry
The system SHALL implement proper error handling and retry mechanisms for asynchronous operations and eventual consistency issues.

#### Scenario: Handle API errors
- **WHEN** CAPI API returns an error
- **THEN** system logs the error and returns a descriptive error message to the user

#### Scenario: Retry on eventual consistency
- **WHEN** system encounters eventual consistency issues
- **THEN** system retries the operation using the helper.Retry() mechanism
- **THEN** system succeeds or fails after max retries

### Requirement: Timeout support
The system SHALL support timeout configuration for resource operations, allowing users to specify custom timeout values for Create, Update, and Delete operations.

#### Scenario: Use default timeouts
- **WHEN** user creates a resource without specifying timeouts
- **THEN** system uses default timeout values

#### Scenario: Use custom timeouts
- **WHEN** user creates a resource with custom timeout values
- **THEN** system uses the specified timeout values for the operations

### Requirement: Comprehensive testing
The system SHALL provide comprehensive unit tests and acceptance tests to ensure the resource functions correctly.

#### Scenario: Run unit tests
- **WHEN** developers run unit tests
- **THEN** all unit tests pass without requiring cloud credentials

#### Scenario: Run acceptance tests
- **WHEN** developers run acceptance tests with TF_ACC=1
- **THEN** system requires TENCENTCLOUD_SECRET_ID and TENCENTCLOUD_SECRET_KEY environment variables
- **THEN** all acceptance tests pass against a real TEO environment
