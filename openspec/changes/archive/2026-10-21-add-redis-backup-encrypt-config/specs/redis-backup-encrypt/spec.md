## ADDED Requirements

### Requirement: Resource definition for redis backup encrypt config
The provider SHALL register a `tencentcloud_redis_backup_encrypt` resource of kind RESOURCE_KIND_CONFIG that manages the backup encryption configuration of a Redis instance. The resource SHALL expose `instance_id`, `backup_encrypt`, `encrypt_key_type`, `kms_key_id`, and `kms_region` schema fields and SHALL support import using the Redis instance ID.

#### Scenario: Schema fields and required attributes
- **WHEN** the resource schema is defined
- **THEN** `instance_id` SHALL be a required string with `ForceNew=true`; `backup_encrypt` SHALL be a required bool; `encrypt_key_type` SHALL be an optional int; `kms_key_id` SHALL be an optional string; `kms_region` SHALL be an optional string

#### Scenario: Import support
- **WHEN** a user runs `terraform import tencentcloud_redis_backup_encrypt.example <redis-instance-id>`
- **THEN** the resource SHALL import successfully using `schema.ImportStatePassthrough` and treat the instance ID as the resource id

### Requirement: Read backup encrypt config
The resource Read operation SHALL call the cloud API `DescribeBackupEncryptConfig` with `instance_id` from `d.Id()` and populate `backup_encrypt`, `encrypt_key_type`, `kms_key_id`, and `kms_region` into state when the corresponding response fields are non-nil.

#### Scenario: Successful read populates state
- **WHEN** `DescribeBackupEncryptConfig` returns a response with `BackupEncrypt=true`, `EncryptKeyType=1`, `KmsKeyId="xxx"`, `KmsRegion="ap-guangzhou"`
- **THEN** the Read SHALL set `backup_encrypt=true`, `encrypt_key_type=1`, `kms_key_id="xxx"`, `kms_region="ap-guangzhou"` into state and `instance_id` from `d.Id()`

#### Scenario: Instance not found clears state
- **WHEN** `DescribeBackupEncryptConfig` returns an instance-not-found error or an empty configuration
- **THEN** the Read SHALL log `[CRUD] redis_backup_encrypt id=<id>` and call `d.SetId("")` to remove the resource from state

#### Scenario: Read retries on transient errors
- **WHEN** `DescribeBackupEncryptConfig` fails with a transient error
- **THEN** the Read SHALL wrap the error with `tccommon.RetryError` and retry within `tccommon.ReadRetryTimeout`

### Requirement: Create backup encrypt config
The resource Create operation SHALL set the resource id to the `instance_id` value and delegate to the Update operation to write the configuration via `ModifyBackupEncrypt`.

#### Scenario: Create delegates to Update
- **WHEN** a user applies a new `tencentcloud_redis_backup_encrypt` resource
- **THEN** Create SHALL set `d.SetId(instance_id)` and call the Update operation, which SHALL invoke `ModifyBackupEncrypt` with the configured parameters

### Requirement: Update backup encrypt config
The resource Update operation SHALL call the cloud API `ModifyBackupEncrypt` with `instance_id` (from `d.Id()`), `backup_encrypt`, and when provided `encrypt_key_type`, `kms_key_id`, `kms_region`; on success it SHALL refresh state by invoking the Read operation.

#### Scenario: Successful update refreshes state
- **WHEN** Update calls `ModifyBackupEncrypt` with valid parameters and the API succeeds
- **THEN** Update SHALL call the Read operation to refresh state from the cloud

#### Scenario: Update retries on transient errors
- **WHEN** `ModifyBackupEncrypt` fails with a transient error
- **THEN** the Update SHALL wrap the error with `tccommon.RetryError` and retry within `tccommon.WriteRetryTimeout`

### Requirement: Delete is a no-op
The resource Delete operation SHALL be a no-op that returns nil without calling any cloud API, because the backup encryption configuration is bound to the Redis instance lifecycle and has no independent deletion semantics.

#### Scenario: Delete does not modify cloud state
- **WHEN** a user destroys the `tencentcloud_redis_backup_encrypt` resource
- **THEN** Delete SHALL return nil without invoking any cloud API and without clearing the remote encryption configuration
