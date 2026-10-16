## Why

The TencentCloud SSM (Secrets Manager) `CreateSecret` API supports specifying the secret's encryption type via the `EncryptType` parameter (`0` = KMS key encryption, `1` = soft key encryption), and the `DescribeSecret` API returns the actual `EncryptType` of a secret. The `tencentcloud_ssm_secret` Terraform resource currently does not expose this parameter, so users cannot declare the encryption type of their secrets through Terraform — they must rely on the API default. This change adds the `encrypt_type` parameter so users can fully manage secret encryption type as code.

## What Changes

- Add a new optional computed field `encrypt_type` to the `tencentcloud_ssm_secret` resource schema.
  - Type: `schema.TypeInt`.
  - Optional and Computed (so the API default is read back when unset).
  - Not `ForceNew`. Since `EncryptType` is only accepted by `CreateSecret` and is **not** present in `UpdateSecretRequest`/`UpdateSecretDescription`, the secret's encryption type is immutable after creation; immutability is enforced via the `immutableArgs` array in the update method (changing `encrypt_type` returns an error instead of forcing re-creation).
  - Accepted values: `0` (KMS key encryption), `1` (soft key encryption).
- In `resourceTencentCloudSsmSecretCreate`, pass `request.EncryptType` when the user configures the field.
- In `resourceTencentCloudSsmSecretRead`, set `encrypt_type` from the `DescribeSecret` response.
- Add `encrypt_type` to the `immutableArgs` array in `resourceTencentCloudSsmSecretUpdate` so changing it returns an explicit error rather than an error-prone in-place update (no force-new).
- Extend the `SecretInfo` struct and `DescribeSecretByName` service method to carry the `EncryptType` field through to the read operation.
- Add/update the `encrypt_type` argument in the resource documentation (`.md`) and register the field for documentation generation.
- Add unit tests covering the new `encrypt_type` parameter.

No breaking changes: the new field is optional; existing configurations continue to work unchanged.

## Capabilities

### New Capabilities
- `ssm-secret-encrypt-type`: Ability to declare and read back the encryption type (`encrypt_type`) of a `tencentcloud_ssm_secret`, mapping to the SSM `CreateSecret.EncryptType` input and `DescribeSecret.EncryptType` output.

### Modified Capabilities
<!-- None. No existing spec-level requirements are changing. -->

## Impact

- **Code**:
  - `tencentcloud/services/ssm/resource_tc_ssm_secret.go` — schema, create, read, update (immutableArgs).
  - `tencentcloud/services/ssm/service_tencentcloud_ssm.go` — `SecretInfo` struct + `DescribeSecretByName` to capture `EncryptType`.
  - `tencentcloud/services/ssm/resource_tc_ssm_secret_test.go` — unit tests for the new parameter.
- **APIs**: SSM `CreateSecret` (input `EncryptType`), `DescribeSecret` (output `EncryptType`). Verified in `vendor/.../ssm/v20190923/models.go`. `EncryptType` is **not** part of `UpdateSecretRequest`, confirming immutability.
- **Documentation**: `tencentcloud/services/ssm/resource_tc_ssm_secret.md` (updated via `make doc` in finalize phase).
- **Dependencies**: None. TencentCloud SDK Go `ssm/v20190923` already in vendor and already exposes `EncryptType`.