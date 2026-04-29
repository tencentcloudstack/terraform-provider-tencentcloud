## ADDED Requirements

### Requirement: ConfirmOriginACLUpdate operation resource
The system SHALL provide a Terraform operation resource `tencentcloud_teo_confirm_origin_acl_update_operation` that calls the `ConfirmOriginACLUpdate` API to confirm that the latest origin IP ACL ranges have been updated on the origin firewall for a given TEO zone.

#### Scenario: Successful confirmation of origin ACL update
- **WHEN** user creates a `tencentcloud_teo_confirm_origin_acl_update_operation` resource with a valid `zone_id`
- **THEN** the system SHALL call `ConfirmOriginACLUpdate` API with the provided `zone_id`
- **AND** the resource ID SHALL be set to a generated token via `helper.BuildToken()`

#### Scenario: Create with retry on transient failure
- **WHEN** the `ConfirmOriginACLUpdate` API call fails with a transient error
- **THEN** the system SHALL retry the API call using `tccommon.WriteRetryTimeout` and `resource.Retry`

### Requirement: Zone ID schema parameter
The resource SHALL expose a `zone_id` parameter of type string that is Required and ForceNew.

#### Scenario: Zone ID is required
- **WHEN** user creates the resource without specifying `zone_id`
- **THEN** the Terraform plan SHALL fail with a required field error

#### Scenario: Zone ID changes trigger recreation
- **WHEN** user changes the `zone_id` value in an existing resource configuration
- **THEN** the system SHALL destroy and recreate the resource

### Requirement: No-op Read handler
The Read handler SHALL be a no-op that returns nil without making any API calls.

#### Scenario: Read operation
- **WHEN** Terraform performs a refresh/read on the resource
- **THEN** the system SHALL return nil without calling any cloud API

### Requirement: No-op Delete handler
The Delete handler SHALL be a no-op that returns nil without making any API calls.

#### Scenario: Delete operation
- **WHEN** user destroys the resource
- **THEN** the system SHALL return nil without calling any cloud API

### Requirement: Provider registration
The resource SHALL be registered in `provider.go` ResourcesMap and documented in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the resource `tencentcloud_teo_confirm_origin_acl_update_operation` SHALL be available for use
