# kms-key-rotate-days Specification

## Purpose
TBD - created by archiving change add-kms-key-rotate-days. Update Purpose after archive.
## Requirements
### Requirement: Rotate days parameter on kms key
The `tencentcloud_kms_key` resource SHALL support an optional `rotate_days` parameter (TypeInt, validate range 7–365) that specifies the key rotation period in days. This parameter is passed to the `EnableKeyRotation` API as `RotateDays` when key rotation is enabled. It is only effective when `key_usage` is `ENCRYPT_DECRYPT` and `key_rotation_enabled` is `true`.

#### Scenario: Create kms key with key rotation and rotate_days
- **WHEN** a user specifies `key_rotation_enabled = true` and `rotate_days = 30` in the `tencentcloud_kms_key` resource configuration (with `key_usage = ENCRYPT_DECRYPT`)
- **THEN** the provider SHALL pass `KeyRotationEnabled` by calling `EnableKeyRotation` with `RotateDays=30` in the API request

#### Scenario: Create kms key with key rotation but without rotate_days
- **WHEN** a user specifies `key_rotation_enabled = true` but does NOT specify `rotate_days` in the configuration
- **THEN** the provider SHALL call `EnableKeyRotation` WITHOUT setting `RotateDays` (API defaults to 365 days)

#### Scenario: Create kms key with rotate_days but key rotation disabled
- **WHEN** a user specifies `rotate_days = 30` but `key_rotation_enabled = false` (or unset)
- **THEN** the provider SHALL NOT call `EnableKeyRotation`, and `rotate_days` SHALL have no effect

#### Scenario: Validation of rotate_days range
- **WHEN** a user specifies `rotate_days = 1` or `rotate_days = 366`
- **THEN** the provider SHALL return a validation error indicating the value must be between 7 and 365

#### Scenario: Read existing kms key with rotate_days
- **WHEN** the provider reads an existing `tencentcloud_kms_key` resource where key rotation is enabled
- **THEN** `rotate_days` SHALL be refreshed from the `DescribeKey` API response (`KeyMetadata.RotateDays`)

#### Scenario: Read existing kms key with rotate_days nil
- **WHEN** the provider reads an existing `tencentcloud_kms_key` resource and `KeyMetadata.RotateDays` is nil
- **THEN** the provider SHALL NOT call `d.Set("rotate_days", ...)`

### Requirement: Update rotate_days re-invokes EnableKeyRotation
The `tencentcloud_kms_key` Update function SHALL re-invoke `EnableKeyRotation` with the new `RotateDays` value when `rotate_days` changes and `key_rotation_enabled` is `true`. The `rotate_days` parameter is NOT immutable and does NOT trigger resource recreation.

#### Scenario: Update rotate_days while rotation enabled
- **WHEN** a user changes `rotate_days` from 30 to 60 while `key_rotation_enabled = true`
- **THEN** the provider SHALL call `EnableKeyRotation` with `RotateDays=60`

#### Scenario: Update rotate_days while rotation disabled
- **WHEN** a user changes `rotate_days` but `key_rotation_enabled = false`
- **THEN** the provider SHALL NOT call `EnableKeyRotation` (the change has no effect)

### Requirement: EnableKeyRotation service function accepts rotateDays
The `EnableKeyRotation` service-layer function SHALL accept a `rotateDays uint64` parameter and set `request.RotateDays` when the value is greater than 0.

#### Scenario: EnableKeyRotation with rotateDays
- **WHEN** `EnableKeyRotation` is called with `rotateDays = 30`
- **THEN** the service function SHALL set `request.RotateDays = helper.Uint64(30)` before calling the API

#### Scenario: EnableKeyRotation without rotateDays
- **WHEN** `EnableKeyRotation` is called with `rotateDays = 0`
- **THEN** the service function SHALL NOT set `request.RotateDays` (API defaults to 365)

