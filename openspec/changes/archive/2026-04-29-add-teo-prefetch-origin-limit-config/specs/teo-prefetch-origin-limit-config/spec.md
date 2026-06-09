## ADDED Requirements

### Requirement: Resource schema definition
The system SHALL define a Terraform resource `tencentcloud_teo_prefetch_origin_limit` with the following schema:
- `zone_id` (TypeString, Required, ForceNew): 站点 ID
- `domain_name` (TypeString, Required, ForceNew): 加速域名
- `area` (TypeString, Required, ForceNew): 回源限速限制的加速区域，取值为 Overseas 或 MainlandChina
- `bandwidth` (TypeInt, Required): 回源限速带宽，取值范围 100-100000，单位 Mbps
- `enabled` (TypeString, Required): 回源限速限制控制开关，取值为 on 或 off

#### Scenario: Valid resource schema with all required fields
- **WHEN** a user defines a `tencentcloud_teo_prefetch_origin_limit` resource with zone_id, domain_name, area, bandwidth, and enabled
- **THEN** the resource SHALL be accepted and the schema SHALL validate all fields

#### Scenario: zone_id, domain_name, area are ForceNew
- **WHEN** a user changes zone_id, domain_name, or area after creation
- **THEN** Terraform SHALL destroy and recreate the resource

### Requirement: Resource ID composition
The system SHALL compose the resource ID using `zone_id`, `domain_name`, and `area` joined by `tccommon.FILED_SP` separator.

#### Scenario: Resource ID format
- **WHEN** a resource is created with zone_id="zone-1", domain_name="example.com", area="Overseas"
- **THEN** the resource ID SHALL be "zone-1#example.com#Overseas" (using FILED_SP separator)

### Requirement: Create operation
The system SHALL implement the Create operation by calling `ModifyPrefetchOriginLimit` with all schema parameters and Enabled=on.

#### Scenario: Successful creation
- **WHEN** a user creates a `tencentcloud_teo_prefetch_origin_limit` resource
- **THEN** the system SHALL call ModifyPrefetchOriginLimit with ZoneId, DomainName, Area, Bandwidth, and Enabled="on"
- **AND** the resource ID SHALL be set to the composite ID
- **AND** the Read operation SHALL be called after successful creation

#### Scenario: Create API failure
- **WHEN** the ModifyPrefetchOriginLimit API call fails
- **THEN** the system SHALL retry with tccommon.WriteRetryTimeout
- **AND** the error SHALL be wrapped with tccommon.RetryError()

### Requirement: Read operation
The system SHALL implement the Read operation by calling `DescribePrefetchOriginLimit` with ZoneId and Filters (domain-name and area), then matching the returned configuration.

#### Scenario: Successful read
- **WHEN** the Read operation is called
- **THEN** the system SHALL call DescribePrefetchOriginLimit with ZoneId, Limit=100, and Filters for domain-name and area
- **AND** SHALL set domain_name, area, and bandwidth from the matched result
- **AND** SHALL NOT set enabled (as the API does not return it)

#### Scenario: Resource not found
- **WHEN** the DescribePrefetchOriginLimit returns no matching configuration
- **THEN** the system SHALL set d.SetId("") to mark the resource as deleted

### Requirement: Update operation
The system SHALL implement the Update operation by calling `ModifyPrefetchOriginLimit` when bandwidth or enabled changes.

#### Scenario: Successful update
- **WHEN** a user updates bandwidth or enabled
- **THEN** the system SHALL call ModifyPrefetchOriginLimit with all parameters (ZoneId, DomainName, Area, Bandwidth, Enabled)
- **AND** SHALL retry with tccommon.WriteRetryTimeout on failure

### Requirement: Delete operation
The system SHALL implement the Delete operation by calling `ModifyPrefetchOriginLimit` with Enabled=off to remove the limit configuration.

#### Scenario: Successful deletion
- **WHEN** a user deletes a `tencentcloud_teo_prefetch_origin_limit` resource
- **THEN** the system SHALL call ModifyPrefetchOriginLimit with all parameters and Enabled="off"
- **AND** SHALL retry with tccommon.WriteRetryTimeout on failure

### Requirement: Import support
The system SHALL support Terraform import with the composite ID format (zone_id#domain_name#area).

#### Scenario: Import existing resource
- **WHEN** a user imports a resource with ID "zone-1#example.com#Overseas"
- **THEN** the system SHALL parse the ID and call the Read operation to populate the state

### Requirement: Provider registration
The system SHALL register the `tencentcloud_teo_prefetch_origin_limit` resource in `provider.go` and add the resource entry in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_prefetch_origin_limit` SHALL be available for use

### Requirement: Documentation
The system SHALL provide a `.md` documentation file following the project documentation format, including description, example usage, and import section.

#### Scenario: Documentation file exists
- **WHEN** the resource is added
- **THEN** a `resource_tc_teo_prefetch_origin_limit_config.md` file SHALL exist in the teo service directory
