## ADDED Requirements

### Requirement: tencentcloud_teo_origin_group resource supports host_header parameter
The tencentcloud_teo_origin_group resource SHALL provide a `host_header` parameter of type String that allows users to configure a custom Host header when making requests to origin servers. This parameter SHALL map to the HostHeader field in the CreateOriginGroup and ModifyOriginGroup APIs.

#### Scenario: Create origin group with host_header
- **WHEN** user creates a tencentcloud_teo_origin_group resource with `host_header` set to "www.example.com"
- **THEN** system SHALL create the origin group successfully
- **AND** system SHALL pass the `host_header` value to CreateOriginGroup API as the HostHeader field
- **AND** system SHALL store the `host_header` value in the Terraform state

#### Scenario: Create origin group without host_header
- **WHEN** user creates a tencentcloud_teo_origin_group resource without setting `host_header`
- **THEN** system SHALL create the origin group successfully
- **AND** system SHALL NOT pass the HostHeader field to CreateOriginGroup API (or pass nil/empty based on API requirements)
- **AND** system SHALL NOT store any value for `host_header` in the Terraform state

#### Scenario: Update origin group host_header
- **WHEN** user updates an existing tencentcloud_teo_origin_group resource by changing `host_header` from "www.example.com" to "cdn.example.com"
- **THEN** system SHALL call ModifyOriginGroup API with the new `host_header` value
- **AND** system SHALL update the `host_header` value in the Terraform state
- **AND** all other resource attributes SHALL remain unchanged

#### Scenario: Add host_header to existing origin group
- **WHEN** user updates an existing tencentcloud_teo_origin_group resource that did not have `host_header` set by adding `host_header` with value "www.example.com"
- **THEN** system SHALL call ModifyOriginGroup API with the `host_header` value
- **AND** system SHALL store the `host_header` value in the Terraform state

#### Scenario: Remove host_header from existing origin group
- **WHEN** user updates an existing tencentcloud_teo_origin_group resource by removing the `host_header` parameter (setting it to null or removing from config)
- **THEN** system SHALL call ModifyOriginGroup API without the HostHeader field (or with empty value based on API requirements)
- **AND** system SHALL remove the `host_header` value from the Terraform state

#### Scenario: Read origin group with host_header
- **WHEN** system reads an existing tencentcloud_teo_origin_group resource that has `host_header` configured
- **THEN** system SHALL retrieve the HostHeader value from ReadOriginGroup or DescribeOriginGroup API response
- **AND** system SHALL update the `host_header` value in the Terraform state to match the API response

#### Scenario: Read origin group without host_header
- **WHEN** system reads an existing tencentcloud_teo_origin_group resource that does not have HostHeader configured
- **THEN** system SHALL NOT set any value for `host_header` in the Terraform state
- **AND** the `host_header` attribute SHALL remain unset

### Requirement: host_header parameter is optional and backward compatible
The `host_header` parameter SHALL be optional and SHALL NOT affect existing resource configurations or states. Existing resources without `host_header` SHALL continue to work without any changes.

#### Scenario: Existing resource configuration remains valid
- **WHEN** user applies Terraform configuration for existing tencentcloud_teo_origin_group resources without `host_header`
- **THEN** system SHALL apply the configuration without errors
- **AND** system SHALL NOT require any migration or state updates

#### Scenario: New parameter does not force resource recreation
- **WHEN** user adds `host_header` to an existing tencentcloud_teo_origin_group resource configuration
- **THEN** system SHALL update the resource in-place without destroying and recreating it
- **AND** system SHALL use ModifyOriginGroup API to apply the change

### Requirement: host_header parameter validation
The `host_header` parameter SHALL accept valid Host header values and SHALL provide appropriate error messages for invalid inputs.

#### Scenario: Valid host_header value
- **WHEN** user sets `host_header` to a valid Host header value (e.g., "www.example.com", "api.example.com:8080")
- **THEN** system SHALL accept the value without validation errors
- **AND** system SHALL pass the value to the underlying API

#### Scenario: Empty string handling
- **WHEN** user sets `host_header` to an empty string ""
- **THEN** system SHALL handle the empty string according to API requirements (either pass as-is or omit from API call)
- **AND** system SHALL provide consistent behavior between Create and Update operations

### Requirement: Documentation and examples
The tencentcloud_teo_origin_group resource documentation SHALL include `host_header` parameter description, usage examples, and notes about its purpose and behavior.

#### Scenario: Documentation includes parameter description
- **WHEN** user views the tencentcloud_teo_origin_group resource documentation
- **THEN** documentation SHALL include a clear description of the `host_header` parameter
- **AND** documentation SHALL specify the parameter type (String)
- **AND** documentation SHALL indicate that the parameter is optional

#### Scenario: Documentation includes usage examples
- **WHEN** user views the tencentcloud_teo_origin_group resource documentation
- **THEN** documentation SHALL include at least one example showing how to use `host_header`
- **AND** example SHALL demonstrate the parameter syntax and typical use case
