# teo-origin-acl-family Specification

## Purpose
TBD - created by archiving change add-teo-origin-acl-family-param. Update Purpose after archive.
## Requirements
### Requirement: Origin ACL family parameter support
The system SHALL support the `origin_acl_family` parameter in the `tencentcloud_teo_origin_acl` Terraform resource to allow users to configure the control domain for origin ACL configuration.

#### Scenario: Successful creation with origin_acl_family
- **WHEN** user creates a `tencentcloud_teo_origin_acl` resource with `origin_acl_family` parameter set to a valid value (e.g., "gaz", "mlc", "emc", "plat-gaz", "plat-mlc", "plat-emc")
- **THEN** system SHALL set `EnableOriginACLRequest.OriginACLFamily` to the provided value
- **AND** system SHALL create the origin ACL with the specified control domain
- **AND** system SHALL store the `origin_acl_family` value in the Terraform state

#### Scenario: Successful creation without origin_acl_family
- **WHEN** user creates a `tencentcloud_teo_origin_acl` resource without setting the `origin_acl_family` parameter
- **THEN** system SHALL call `EnableOriginACL` API without setting `OriginACLFamily`
- **AND** API SHALL use its default control domain behavior
- **AND** system SHALL read and store the actual control domain value from the API response in the Terraform state

#### Scenario: Successful read operation
- **WHEN** system reads an existing `tencentcloud_teo_origin_acl` resource
- **THEN** system SHALL call `DescribeOriginACL` API
- **AND** system SHALL extract the `OriginACLFamily` value from `OriginACLInfo.OriginACLFamily` in the API response
- **AND** system SHALL store this value in the Terraform state

#### Scenario: Successful update of origin_acl_family
- **WHEN** user updates the `origin_acl_family` parameter of an existing `tencentcloud_teo_origin_acl` resource
- **THEN** system SHALL detect the parameter change
- **AND** system SHALL call `ModifyOriginACL` API with the new `OriginACLFamily` value
- **AND** system SHALL update the Terraform state with the new value

#### Scenario: Successful update without changing origin_acl_family
- **WHEN** user updates other parameters (e.g., `l7_hosts`, `l4_proxy_ids`) of an existing `tencentcloud_teo_origin_acl` resource but does not change `origin_acl_family`
- **THEN** system SHALL NOT call `ModifyOriginACL` API for `OriginACLFamily` parameter
- **AND** system SHALL maintain the existing `origin_acl_family` value in the Terraform state

#### Scenario: Backward compatibility with existing resources
- **WHEN** system applies the new feature to existing `tencentcloud_teo_origin_acl` resources that were created before this parameter existed
- **THEN** system SHALL successfully read the existing resources
- **AND** system SHALL populate `origin_acl_family` in state from the API response (computed field behavior)
- **AND** system SHALL not require any manual migration or configuration changes from users

#### Scenario: Parameter type validation
- **WHEN** user provides a value for `origin_acl_family` parameter
- **THEN** system SHALL accept the value as a string type
- **AND** system SHALL not perform any local validation of the value (validation is performed by the cloud API)

#### Scenario: Computed field behavior
- **WHEN** user creates a `tencentcloud_teo_origin_acl` resource and does not specify `origin_acl_family`
- **THEN** system SHALL mark `origin_acl_family` as computed in the Terraform state
- **AND** system SHALL display the API's default value after the first apply operation
- **AND** subsequent read operations SHALL consistently return this value

#### Scenario: State refresh consistency
- **WHEN** system refreshes the state of an existing `tencentcloud_teo_origin_acl` resource
- **THEN** system SHALL read the current `origin_acl_family` value from the API
- **AND** system SHALL update the Terraform state to match the API's current value
- **AND** system SHALL detect and report any drift between state and API

#### Scenario: Delete operation behavior
- **WHEN** user deletes a `tencentcloud_teo_origin_acl` resource
- **THEN** system SHALL call `DisableOriginACL` API
- **AND** system SHALL remove the resource from Terraform state
- **AND** system SHALL not pass `origin_acl_family` parameter in the delete request (as it's not supported by DisableOriginACL API)

