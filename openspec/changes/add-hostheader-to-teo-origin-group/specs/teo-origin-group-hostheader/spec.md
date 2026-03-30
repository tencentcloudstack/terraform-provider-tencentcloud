## ADDED Requirements

### Requirement: Create origin group with HostHeader
The system SHALL allow users to specify the `host_header` parameter when creating a `tencentcloud_teo_origin_group` resource. The HostHeader value MUST be passed to the CreateOriginGroup API to configure the back-to-origin Host Header during resource creation.

#### Scenario: Create origin group with HostHeader specified
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource with `host_header` parameter set to a valid value (e.g., "example.com")
- **THEN** system MUST pass the `host_header` value to the CreateOriginGroup API's `HostHeader` parameter
- **AND** the origin group MUST be created with the specified HostHeader configuration

#### Scenario: Create origin group without HostHeader
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource without specifying the `host_header` parameter
- **THEN** system MUST NOT pass the HostHeader parameter to the CreateOriginGroup API
- **AND** the origin group MUST be created successfully without HostHeader configuration

### Requirement: HostHeader parameter handling in Create operation
The system MUST process the `host_header` parameter in the Create function using the same pattern as the Update function, ensuring consistency across resource operations.

#### Scenario: HostHeader parameter processing matches Update pattern
- **WHEN** processing the `host_header` parameter in the Create function
- **THEN** system MUST use `d.GetOk("host_header")` to check if the parameter is set
- **AND** system MUST use `helper.String()` to convert the value to a string pointer
- **AND** system MUST assign the value to `request.HostHeader` in the CreateOriginGroupRequest
