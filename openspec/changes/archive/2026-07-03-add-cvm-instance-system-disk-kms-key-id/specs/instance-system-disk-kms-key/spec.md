## ADDED Requirements

### Requirement: System disk KMS key ID configuration

The `tencentcloud_instance` resource SHALL support a top-level `kms_key_id` parameter (TypeString, Optional, ForceNew) that specifies the custom KMS key ID used to encrypt the system disk.

#### Scenario: Create instance with system disk KMS key

- **WHEN** a user specifies `kms_key_id` in the `tencentcloud_instance` resource configuration
- **THEN** the provider SHALL pass the value as `SystemDisk.KmsKeyId` in the `RunInstances` API request
- **THEN** the instance SHALL be created with the system disk encrypted using the specified KMS key

#### Scenario: Create instance without system disk KMS key

- **WHEN** a user does NOT specify `kms_key_id` in the `tencentcloud_instance` resource configuration
- **THEN** the provider SHALL NOT set `SystemDisk.KmsKeyId` in the `RunInstances` API request
- **THEN** the instance SHALL be created with default encryption behavior (no custom KMS key)

#### Scenario: Reinstall instance preserves system disk KMS key

- **WHEN** an instance has `kms_key_id` configured and the `image_id` changes (triggering a `ResetInstance` call)
- **THEN** the provider SHALL pass the `kms_key_id` value as `SystemDisk.KmsKeyId` in the `ResetInstance` API request
- **THEN** the reinstalled system disk SHALL be encrypted using the same KMS key

#### Scenario: Changing kms_key_id forces instance recreation

- **WHEN** a user modifies the `kms_key_id` value for an existing instance
- **THEN** the provider SHALL force recreation of the instance (ForceNew)