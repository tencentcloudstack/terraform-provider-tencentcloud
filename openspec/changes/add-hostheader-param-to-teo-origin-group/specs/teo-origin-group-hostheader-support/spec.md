## ADDED Requirements

### Requirement: Create origin group with host_header
The system SHALL support passing the `host_header` parameter to the CreateOriginGroup API when creating a tencentcloud_teo_origin_group resource. This parameter shall be correctly transmitted to the API's HostHeader field when provided by the user.

#### Scenario: Create origin group with host_header specified
- **WHEN** user creates a tencentcloud_teo_origin_group resource with `host_header` parameter set to "example.com"
- **THEN** the CreateOriginGroup API request shall include HostHeader field with value "example.com"
- **AND** the resource creation shall complete successfully

#### Scenario: Create origin group without host_header
- **WHEN** user creates a tencentcloud_teo_origin_group resource without specifying `host_header` parameter
- **THEN** the CreateOriginGroup API request shall not include HostHeader field (or include it as nil)
- **AND** the resource creation shall complete successfully

### Requirement: Read host_header after creation
The system SHALL correctly read and return the `host_header` value from the API response after creating a tencentcloud_teo_origin_group resource with host_header specified.

#### Scenario: Read host_header after creation
- **WHEN** user creates a tencentcloud_teo_origin_group resource with `host_header` set to "example.com"
- **AND** the API returns HostHeader with value "example.com"
- **THEN** the resource state shall have `host_header` set to "example.com"
- **AND** subsequent Read operations shall correctly retrieve this value

### Requirement: Consistency with Update operation
The system shall ensure that the `host_header` parameter handling in the Create operation is consistent with the Update operation, following the same pattern of parameter retrieval, validation, and API transmission.

#### Scenario: Parameter retrieval consistency
- **WHEN** `host_header` parameter is provided in either Create or Update operations
- **THEN** both operations shall use the same method to retrieve the parameter value (`d.GetOk("host_header")`)
- **AND** both operations shall use the same approach to assign the value to the API request (`helper.String(v.(string))`)

#### Scenario: API field transmission consistency
- **WHEN** `host_header` is set in Create or Update operations
- **THEN** both operations shall assign the value to the corresponding API field (`request.HostHeader`) in the same manner
- **AND** the type conversion shall be consistent between Create and Update

### Requirement: Backward compatibility
The system shall maintain full backward compatibility with existing configurations and state. The addition of host_header parameter transmission in Create shall not break any existing resources or configurations.

#### Scenario: Existing resources without host_header
- **WHEN** existing resources are updated that do not have `host_header` in their state
- **THEN** these resources shall continue to function normally
- **AND** no schema changes shall be required
- **AND** existing tests shall continue to pass

#### Scenario: New resources with optional host_header
- **WHEN** users create new resources without specifying `host_header`
- **THEN** the resource shall be created successfully
- **AND** the behavior shall be identical to the current implementation
- **AND** no breaking changes shall be introduced
