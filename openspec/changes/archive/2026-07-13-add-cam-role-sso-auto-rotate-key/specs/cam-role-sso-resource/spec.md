## ADDED Requirements

### Requirement: Auto Rotate Key Parameter Management
The `tencentcloud_cam_role_sso` resource SHALL support managing the OIDC public key auto-rotation switch through the `auto_rotate_key` parameter.

#### Scenario: Define auto_rotate_key schema field
- **WHEN** defining the resource schema for `tencentcloud_cam_role_sso`
- **THEN** the schema SHALL include an `auto_rotate_key` field of type `schema.TypeInt`
- **AND** the field SHALL be Optional (not Required)
- **AND** the field SHALL NOT be ForceNew (changes SHALL trigger in-place update)
- **AND** the field description SHALL document the enum values 0 (disabled) and 1 (enabled) with default value 0

#### Scenario: Create resource with auto_rotate_key specified
- **GIVEN** a user creates a `tencentcloud_cam_role_sso` resource with `auto_rotate_key = 1`
- **WHEN** the create operation is invoked
- **THEN** the system SHALL call `CreateOIDCConfig` API with `request.AutoRotateKey` set to the provided value (converted to `*uint64`)
- **AND** after a successful API call, the system SHALL set the resource id to the resource `name`
- **AND** invoke the read operation to refresh state

#### Scenario: Create resource without auto_rotate_key specified
- **GIVEN** a user creates a `tencentcloud_cam_role_sso` resource without specifying `auto_rotate_key`
- **WHEN** the create operation is invoked
- **THEN** the system SHALL NOT send `AutoRotateKey` (or send the zero value) to `CreateOIDCConfig`, letting the API apply its default value (0)

#### Scenario: Update auto_rotate_key
- **GIVEN** an existing `tencentcloud_cam_role_sso` resource
- **WHEN** the user updates the `auto_rotate_key` field
- **THEN** the system SHALL detect the change via `d.HasChange("auto_rotate_key")`
- **AND** the system SHALL call `UpdateOIDCConfig` API with `request.AutoRotateKey` set to the new value (converted to `*uint64`)
- **AND** the update SHALL be performed in-place without recreating the resource

#### Scenario: Read auto_rotate_key from API
- **GIVEN** a `tencentcloud_cam_role_sso` resource exists and is being read or refreshed
- **WHEN** the `DescribeOIDCConfig` API returns a response
- **THEN** the system SHALL check whether `response.Response.AutoRotateKey` is nil
- **AND** if it is NOT nil, the system SHALL set `auto_rotate_key` in Terraform state to the returned value (converted from `*uint64` to `int`)
- **AND** if it IS nil, the system SHALL skip setting `auto_rotate_key` to avoid a nil pointer dereference

#### Scenario: Delete resource with auto_rotate_key
- **GIVEN** a `tencentcloud_cam_role_sso` resource is being destroyed
- **WHEN** the delete operation is invoked
- **THEN** the system SHALL call `DeleteOIDCConfig` API with only the `Name` parameter (the resource id)
- **AND** the system SHALL NOT include `AutoRotateKey` in the delete request, because `DeleteOIDCConfigRequest` does not support this field

### Requirement: Backward Compatibility
The addition of `auto_rotate_key` SHALL maintain full backward compatibility with existing `tencentcloud_cam_role_sso` configurations and state.

#### Scenario: Existing configurations without auto_rotate_key
- **GIVEN** an existing Terraform configuration for `tencentcloud_cam_role_sso` that does not specify `auto_rotate_key`
- **WHEN** the user applies the configuration
- **THEN** the resource SHALL continue to work normally without errors
- **AND** the API default value (0, disabled) SHALL apply

### Requirement: Documentation
The resource documentation SHALL reflect the new `auto_rotate_key` parameter.

#### Scenario: Update resource documentation
- **WHEN** the resource documentation `resource_tc_cam_role_sso.md` is updated
- **THEN** the example usage SHALL include the `auto_rotate_key` parameter
- **AND** the documentation SHALL be regenerated via `make doc` (handled in the finalize phase) to update `website/docs/r/cam_role_sso.html.markdown`
