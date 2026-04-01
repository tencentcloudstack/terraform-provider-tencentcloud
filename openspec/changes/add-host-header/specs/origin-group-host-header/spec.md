## ADDED Requirements

### Requirement: TEO origin group resource supports host_header field
The tencentcloud_teo_origin_group resource SHALL include a `host_header` field that allows users to configure the Host Header for origin requests. This field SHALL be optional and accept string values. When provided, the Host Header SHALL be passed to the CreateOriginGroup API.

#### Scenario: Create origin group with host_header
- **WHEN** user creates a tencentcloud_teo_origin_group resource with host_header set to "www.example.com"
- **THEN** the resource SHALL be created successfully
- **AND** the CreateOriginGroup API SHALL receive the HostHeader parameter with value "www.example.com"

#### Scenario: Create origin group without host_header
- **WHEN** user creates a tencentcloud_teo_origin_group resource without setting host_header
- **THEN** the resource SHALL be created successfully
- **AND** the CreateOriginGroup API SHALL not include the HostHeader parameter in the request

### Requirement: Read operation returns host_header value
The Read operation SHALL retrieve the HostHeader value from the API response and populate the `host_header` field in the Terraform state.

#### Scenario: Read origin group with host_header
- **WHEN** the Read operation is called for an origin group that has host_header set to "www.example.com"
- **THEN** the state SHALL contain host_header with value "www.example.com"

#### Scenario: Read origin group without host_header
- **WHEN** the Read operation is called for an origin group that does not have host_header configured
- **THEN** the state SHALL not contain the host_header field or it SHALL be empty

### Requirement: Update operation supports host_header modification
The Update operation SHALL allow users to modify the `host_header` field and pass the updated value to the CreateOriginGroup API.

#### Scenario: Update origin group host_header
- **WHEN** user updates the host_header field from "www.example.com" to "www.new-example.com"
- **THEN** the Update operation SHALL call the CreateOriginGroup API with the new HostHeader value "www.new-example.com"
- **AND** the state SHALL reflect the updated value

#### Scenario: Add host_header to existing origin group
- **WHEN** user updates an existing origin group by setting host_header to "www.example.com"
- **THEN** the Update operation SHALL call the CreateOriginGroup API with the HostHeader value
- **AND** the state SHALL contain the new host_header field

#### Scenario: Remove host_header from origin group
- **WHEN** user removes the host_header field from an existing origin group
- **THEN** the Update operation SHALL call the CreateOriginGroup API without the HostHeader parameter
- **AND** the state SHALL not contain the host_header field

### Requirement: Delete operation is unaffected by host_header
The Delete operation SHALL not require special handling for the `host_header` field. The deletion process shall proceed normally regardless of whether host_header is configured.

#### Scenario: Delete origin group with host_header
- **WHEN** user deletes a tencentcloud_teo_origin_group resource that has host_header configured
- **THEN** the resource SHALL be deleted successfully
- **AND** the deletion process shall not fail due to the presence of host_header

### Requirement: Schema field definition
The `host_header` field SHALL be defined in the resource Schema with the following properties: Type: String, Optional: true, Computed: false, ForceNew: false.

#### Scenario: Schema definition validation
- **WHEN** the resource Schema is inspected
- **THEN** it SHALL include a field named "host_header"
- **AND** the field type SHALL be schema.TypeString
- **AND** the field SHALL be Optional
- **AND** the field SHALL NOT be Computed
- **AND** the field SHALL NOT be ForceNew