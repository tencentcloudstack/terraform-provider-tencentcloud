## ADDED Requirements

### Requirement: Create Origin Group with Host Header
The system SHALL support creating an origin group with the `host_header` parameter set during resource creation. The `host_header` parameter is optional and, when provided, MUST be included in the `CreateOriginGroup` API call.

#### Scenario: Create origin group with host header
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource with `host_header` set to a valid string value
- **THEN** the `CreateOriginGroup` API is called with `HostHeader` parameter set to the provided value
- **THEN** the origin group is created with the specified host header configuration
- **THEN** the resource state includes the `host_header` value

#### Scenario: Create origin group without host header
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource without setting `host_header`
- **THEN** the `CreateOriginGroup` API is called without the `HostHeader` parameter
- **THEN** the origin group is created successfully
- **THEN** the resource state does not include a `host_header` value or it matches the API default

#### Scenario: Create and update origin group with host header
- **WHEN** user creates a `tencentcloud_teo_origin_group` resource with `host_header` set to "example.com"
- **THEN** the origin group is created with `host_header` = "example.com"
- **WHEN** user updates the resource to change `host_header` to "new-example.com"
- **THEN** the `ModifyOriginGroup` API is called with `HostHeader` = "new-example.com"
- **THEN** the resource state reflects the updated `host_header` value
