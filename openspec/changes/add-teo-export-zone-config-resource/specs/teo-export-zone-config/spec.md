## ADDED Requirements

### Requirement: Export zone configuration
The Terraform Provider SHALL support exporting TEO zone configuration through the `tencentcloud_teo_export_zone_config` resource. The resource SHALL be created, read, updated, and deleted using standard Terraform CRUD operations, with all parameters matching the CAPI interface definition exactly.

#### Scenario: Create export zone config resource
- **WHEN** user applies Terraform configuration with `tencentcloud_teo_export_zone_config` resource
- **THEN** provider SHALL call CAPI interface to create the export configuration
- **AND** provider SHALL store the resource ID in Terraform state
- **AND** provider SHALL return all exported configuration attributes to the user

#### Scenario: Read export zone config resource
- **WHEN** user queries an existing `tencentcloud_teo_export_zone_config` resource
- **THEN** provider SHALL call CAPI interface to retrieve the current configuration
- **AND** provider SHALL refresh the resource state with latest values from the API

#### Scenario: Update export zone config resource
- **WHEN** user modifies any optional parameter in `tencentcloud_teo_export_zone_config` resource
- **THEN** provider SHALL call CAPI interface to update the configuration
- **AND** provider SHALL update the resource state with new values

#### Scenario: Delete export zone config resource
- **WHEN** user removes `tencentcloud_teo_export_zone_config` resource from Terraform configuration
- **THEN** provider SHALL call CAPI interface to delete the export configuration
- **AND** provider SHALL remove the resource from Terraform state

### Requirement: Schema parameters match CAPI interface
The resource schema SHALL include all parameters defined in the CAPI interface `iacpres-i2etM5NTBN` (version `iacpresv-k90qSqRTp9`), with Required/Optional attributes exactly matching the interface definition.

#### Scenario: Required parameters validation
- **WHEN** user creates `tencentcloud_teo_export_zone_config` resource without required parameters
- **THEN** provider SHALL reject the configuration with a clear error message listing missing required fields

#### Scenario: Optional parameters handling
- **WHEN** user omits optional parameters in `tencentcloud_teo_export_zone_config` resource
- **THEN** provider SHALL create the resource successfully using default values from the CAPI interface

#### Scenario: Schema compatibility
- **WHEN** CAPI interface version `iacpresv-k90qSqRTp9` defines a parameter as Required
- **THEN** resource schema SHALL mark that parameter as Required
- **AND** **WHEN** CAPI interface defines a parameter as Optional
- **THEN** resource schema SHALL mark that parameter as Optional

### Requirement: Timeout support for async operations
The resource SHALL support configurable timeouts for Create, Read, Update, and Delete operations to handle asynchronous CAPI calls, with appropriate default values for TEO service characteristics.

#### Scenario: Custom timeout configuration
- **WHEN** user specifies custom timeout in `tencentcloud_teo_export_zone_config` resource
- **THEN** provider SHALL use the custom timeout value for the operation

#### Scenario: Default timeout usage
- **WHEN** user does not specify timeout in `tencentcloud_teo_export_zone_config` resource
- **THEN** provider SHALL use the default timeout values defined in the resource schema

#### Scenario: Timeout error handling
- **WHEN** CAPI operation exceeds the configured timeout
- **THEN** provider SHALL return a clear timeout error message to the user

### Requirement: Retry mechanism for eventual consistency
The resource SHALL implement retry logic using the project's standard `helper.Retry()` function to handle eventual consistency of TEO service operations.

#### Scenario: Retry on transient errors
- **WHEN** CAPI call returns a transient error (e.g., network timeout, rate limit)
- **THEN** provider SHALL retry the operation according to the retry configuration
- **AND** provider SHALL log retry attempts for debugging

#### Scenario: Final consistency handling
- **WHEN** Create or Update operation completes but state is not immediately consistent
- **THEN** provider SHALL retry the Read operation until state is consistent or timeout is reached

### Requirement: Comprehensive test coverage
The resource SHALL include both unit tests and acceptance tests to ensure correct behavior of all CRUD operations and edge cases.

#### Scenario: Unit test coverage
- **WHEN** running unit tests for `resource_tencentcloud_teo_export_zone_config`
- **THEN** all CRUD functions SHALL be tested with mock data
- **AND** edge cases SHALL be covered (missing required fields, invalid parameters, etc.)

#### Scenario: Acceptance test coverage
- **WHEN** running acceptance tests with `TF_ACC=1`
- **THEN** provider SHALL execute tests against real TEO API
- **AND** tests SHALL verify resource creation, reading, updating, and deletion
- **AND** tests SHALL clean up all created resources after completion

### Requirement: Documentation and examples
The resource SHALL include comprehensive documentation and usage examples to help users understand how to use the `tencentcloud_teo_export_zone_config` resource effectively.

#### Scenario: Resource documentation
- **WHEN** users access the provider documentation
- **THEN** documentation SHALL describe all parameters and their types
- **AND** documentation SHALL include examples of common usage patterns
- **AND** documentation SHALL explain any specific behaviors or constraints

#### Scenario: Example configuration
- **WHEN** users reference the resource examples
- **THEN** example file SHALL demonstrate a complete working configuration
- **AND** example SHALL include all required parameters
- **AND** example SHALL show usage of common optional parameters
