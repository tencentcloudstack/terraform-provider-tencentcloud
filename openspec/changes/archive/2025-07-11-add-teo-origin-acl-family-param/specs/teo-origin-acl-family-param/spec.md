## ADDED Requirements

### Requirement: origin_acl_family schema parameter for tencentcloud_teo_origin_acl resource
The resource SHALL expose an `origin_acl_family` parameter of type string that is Optional and Computed, describing the origin ACL control domain for source protection.

#### Scenario: Create resource with origin_acl_family specified
- **WHEN** user creates a `tencentcloud_teo_origin_acl` resource with `origin_acl_family` set to a valid value (e.g., "mlc")
- **THEN** the system SHALL set `OriginACLFamily` on the EnableOriginACL API request
- **AND** the resource state SHALL reflect the specified `origin_acl_family` value after read

#### Scenario: Create resource without origin_acl_family specified
- **WHEN** user creates a `tencentcloud_teo_origin_acl` resource without specifying `origin_acl_family`
- **THEN** the system SHALL NOT set `OriginACLFamily` on the EnableOriginACL API request
- **AND** the API SHALL apply the default value (gaz)
- **AND** the Read handler SHALL populate `origin_acl_family` from the API response

#### Scenario: Update origin_acl_family value
- **WHEN** user changes the `origin_acl_family` value in an existing resource configuration
- **THEN** the system SHALL set `OriginACLFamily` on the ModifyOriginACL API request with the new value
- **AND** the resource state SHALL reflect the updated `origin_acl_family` value after read

#### Scenario: Read origin_acl_family from API response
- **WHEN** the Read handler calls DescribeOriginACL
- **AND** the response `OriginACLInfo.OriginACLFamily` is not nil
- **THEN** the system SHALL set `origin_acl_family` in the resource state to the API response value

#### Scenario: Read with nil OriginACLFamily
- **WHEN** the Read handler calls DescribeOriginACL
- **AND** the response `OriginACLInfo.OriginACLFamily` is nil
- **THEN** the system SHALL NOT set `origin_acl_family` in the resource state

### Requirement: origin_acl_family in Create handler
The Create handler SHALL pass `OriginACLFamily` on the EnableOriginACL request when the user specifies `origin_acl_family` in the Terraform configuration.

#### Scenario: OriginACLFamily set on Enable request
- **WHEN** user provides `origin_acl_family` in the resource configuration
- **THEN** the system SHALL set `request.OriginACLFamily` to the provided value before calling EnableOriginACL

### Requirement: origin_acl_family in Update handler
The Update handler SHALL pass `OriginACLFamily` on ModifyOriginACL requests when `origin_acl_family` has changed.

#### Scenario: OriginACLFamily change with entity changes
- **WHEN** `origin_acl_family` has changed AND there are l7_hosts or l4_proxy_ids changes
- **THEN** the system SHALL set `OriginACLFamily` on the first ModifyOriginACL request that processes entity changes

#### Scenario: OriginACLFamily change without entity changes
- **WHEN** `origin_acl_family` has changed AND there are no l7_hosts or l4_proxy_ids changes
- **THEN** the system SHALL make a ModifyOriginACL request with only `ZoneId` and `OriginACLFamily` set

#### Scenario: OriginACLFamily also set in Create batch ModifyOriginACL calls
- **WHEN** the Create handler makes batch ModifyOriginACL calls for overflow L7/L4 entities
- **AND** `origin_acl_family` was specified by the user
- **THEN** the system SHALL set `OriginACLFamily` on each batch ModifyOriginACL request

### Requirement: origin_acl_family in data source
The data source `tencentcloud_teo_origin_acl` SHALL expose `origin_acl_family` as a Computed string field inside the `origin_acl_info` block.

#### Scenario: Read origin_acl_family from data source
- **WHEN** user queries the `tencentcloud_teo_origin_acl` data source
- **AND** the DescribeOriginACL response `OriginACLInfo.OriginACLFamily` is not nil
- **THEN** the system SHALL include `origin_acl_family` in the `origin_acl_info` output with the API response value

#### Scenario: OriginACLFamily is nil in data source
- **WHEN** user queries the `tencentcloud_teo_origin_acl` data source
- **AND** the DescribeOriginACL response `OriginACLInfo.OriginACLFamily` is nil
- **THEN** the system SHALL NOT include `origin_acl_family` in the `origin_acl_info` output

### Requirement: Backward compatibility
The addition of `origin_acl_family` SHALL NOT break existing Terraform configurations or state files.

#### Scenario: Existing configuration without origin_acl_family
- **WHEN** a user applies an existing configuration that does not include `origin_acl_family`
- **THEN** the Terraform plan SHALL show no changes to the resource
- **AND** the `origin_acl_family` value in state SHALL be populated from the API response after the next refresh
