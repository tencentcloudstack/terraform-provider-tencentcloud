## ADDED Requirements

### Requirement: Support auto-mount instance ID parameter

The `tencentcloud_cbs_storage` resource SHALL support specifying a CVM instance ID to auto-mount and auto-initialize a CBS data disk during creation via the `instance_id` schema field, which maps to the `AutoMountConfiguration.InstanceId` parameter of the `CreateDisks` API.

#### Scenario: Create CBS storage with instance_id for auto-mount

- **WHEN** user specifies `instance_id` in the `tencentcloud_cbs_storage` configuration
- **THEN** the provider SHALL pass `AutoMountConfiguration.InstanceId` (as a single-element `[]*string` containing the instance ID) to the `CreateDisks` API request
- **AND** the disk SHALL be auto-mounted and auto-initialized to the specified CVM instance during creation

#### Scenario: Create CBS storage without instance_id (backward compatibility)

- **WHEN** user creates a `tencentcloud_cbs_storage` without specifying `instance_id`
- **THEN** the provider SHALL NOT set `AutoMountConfiguration` in the `CreateDisks` request
- **AND** the disk SHALL be created without auto-mounting, identical to previous behavior

### Requirement: Read instance_id from API response

The `tencentcloud_cbs_storage` resource SHALL populate the `instance_id` field in Terraform state from the `Disk.InstanceId` field returned by the `DescribeDisks` API response, so that state reflects the actual CVM instance the disk is currently mounted to.

#### Scenario: Read instance_id when disk is mounted

- **WHEN** the provider reads a CBS storage whose `Disk.InstanceId` is non-nil
- **THEN** the provider SHALL set `instance_id` in Terraform state to the value of `Disk.InstanceId`
- **AND** the state SHALL reflect the CVM instance the disk is currently mounted to

#### Scenario: Read instance_id when disk is not mounted

- **WHEN** the provider reads a CBS storage whose `Disk.InstanceId` is nil or empty
- **THEN** the provider SHALL NOT call `d.Set("instance_id", ...)` to avoid overwriting state with an empty value
- **AND** no spurious diffs SHALL be generated

### Requirement: Schema field definition for instance_id

The resource schema SHALL include a new optional `instance_id` field for specifying the auto-mount CVM instance.

#### Scenario: instance_id field specification

- **WHEN** defining the `instance_id` schema field
- **THEN** it SHALL be of type String
- **AND** it SHALL be Optional
- **AND** it SHALL have a description explaining it specifies the CVM instance ID for auto-mounting the data disk during creation

### Requirement: Backward compatibility

The addition of the `instance_id` parameter SHALL NOT break existing CBS storage configurations.

#### Scenario: Existing configurations without instance_id

- **WHEN** user applies an existing `tencentcloud_cbs_storage` configuration that does not specify `instance_id`
- **THEN** the storage SHALL be created/updated successfully without auto-mounting
- **AND** no changes SHALL be required to existing Terraform configurations

#### Scenario: State file compatibility

- **WHEN** provider reads state from a CBS storage created before this feature was added
- **THEN** the provider SHALL handle the missing `instance_id` field gracefully
- **AND** no spurious diffs SHALL be generated

### Requirement: Documentation

The resource documentation SHALL include examples and guidance for the `instance_id` parameter.

#### Scenario: Usage example with auto-mount

- **WHEN** user views the resource documentation
- **THEN** it SHALL include an example showing how to create a CBS storage with `instance_id` for auto-mounting
- **AND** the example SHALL demonstrate correct usage of the `instance_id` field

#### Scenario: Parameter description

- **WHEN** user views the Argument Reference section
- **THEN** the `instance_id` field SHALL be documented with its type, optionality, and purpose
