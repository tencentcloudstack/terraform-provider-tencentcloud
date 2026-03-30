## ADDED Requirements

### Requirement: HostHeader parameter support during resource creation
The system SHALL allow users to specify the `host_header` parameter when creating a `tencentcloud_teo_origin_group` resource. The `host_header` value MUST be passed to the CreateOriginGroup API request.

#### Scenario: Create origin group with HostHeader
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource with `host_header` parameter specified
- **THEN** the system SHALL pass the `host_header` value to the CreateOriginGroup API
- **AND** the system SHALL create the origin group with the specified HostHeader configuration
- **AND** the created resource SHALL reflect the HostHeader value in the state

#### Scenario: Create origin group without HostHeader
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource without specifying `host_header` parameter
- **THEN** the system SHALL create the origin group without HostHeader configuration
- **AND** the system SHALL not pass the `host_header` parameter to the CreateOriginGroup API
- **AND** the created resource SHALL have an empty or default HostHeader value

### Requirement: HostHeader parameter behavior
The `host_header` parameter MUST be treated as optional. When specified, it MUST only take effect when the `type` parameter is set to `HTTP`. The rule engine's Host Header configuration priority SHALL be higher than the origin group's HostHeader.

#### Scenario: HostHeader with HTTP type origin group
- **WHEN** user creates an origin group with `type = "HTTP"` and specifies `host_header`
- **THEN** the system SHALL pass the `host_header` value to the CreateOriginGroup API
- **AND** the HostHeader SHALL take effect for this HTTP-specific origin group

#### Scenario: HostHeader with GENERAL type origin group
- **WHEN** user creates an origin group with `type = "GENERAL"` and specifies `host_header`
- **THEN** the system SHALL pass the `host_header` value to the CreateOriginGroup API
- **AND** the API SHALL handle the parameter according to its own validation rules

### Requirement: Backward compatibility
The implementation MUST maintain backward compatibility with existing configurations that do not specify `host_header` parameter during resource creation.

#### Scenario: Existing configuration without HostHeader
- **WHEN** existing Terraform configuration creates a `tencentcloud_teo_origin_group` resource without `host_header` parameter
- **THEN** the system SHALL create the resource successfully without any changes to behavior
- **AND** no drift or configuration update SHALL be required
