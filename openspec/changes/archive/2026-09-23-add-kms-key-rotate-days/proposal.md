## Why

The KMS `EnableKeyRotation` API supports a `RotateDays` parameter to specify the key rotation period (in days, range 7–365, default 365), but the Terraform resource `tencentcloud_kms_key` does not expose it. Users who enable key rotation (`key_rotation_enabled = true`) cannot control the rotation period through Terraform and must use the console or API directly.

## What Changes

- Add `rotate_days` (Optional, TypeInt) parameter to the `tencentcloud_kms_key` resource to specify the key rotation period when key rotation is enabled. Only valid when `key_usage` is `ENCRYPT_DECRYPT` and `key_rotation_enabled` is `true`. The valid range is 7–365 (the API enforces this); no Terraform-side `ValidateFunc` is applied so that API-side validation messages surface directly.
- Pass `RotateDays` to the `EnableKeyRotation` API request when specified during Create and Update flows.
- Read `RotateDays` from the `DescribeKey` API response (`KeyMetadata.RotateDays`) in the Read function to support state refresh and import.
- Treat `rotate_days` as immutable: once a key is created with a rotation period, changing `rotate_days` along with `key_rotation_enabled` updates will re-call `EnableKeyRotation` with the new `RotateDays` value (the `EnableKeyRotation` API accepts `RotateDays` and can be re-invoked to update the period). Changes to `rotate_days` without `key_rotation_enabled = true` are ignored.
- Update the `EnableKeyRotation` service-layer function to accept and pass the `rotateDays` parameter.
- Update `tencentcloud_kms_key` resource documentation (`.md` example file).

## Capabilities

### New Capabilities
- `kms-key-rotate-days`: Enable the `rotate_days` parameter on the `tencentcloud_kms_key` resource to allow users to specify the key rotation period (in days) when key rotation is enabled. Validation of the 7–365 range is delegated to the cloud API (no Terraform-side `ValidateFunc`).

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/kms/resource_tc_kms_key.go` — add `rotate_days` schema field, wire through Create/Update flows, add Read support
  - `tencentcloud/services/kms/service_tencentcloud_kms.go` — update `EnableKeyRotation` function to accept and pass `rotateDays uint64`
  - `tencentcloud/services/kms/resource_tc_kms_key_test.go` — add unit test for the new parameter
  - `tencentcloud/services/kms/resource_tc_kms_key.md` — update documentation with `rotate_days` usage
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118` — `EnableKeyRotationRequest` already includes `RotateDays *uint64` field, and `KeyMetadata` (returned by `DescribeKey`) already includes `RotateDays *uint64`, so no SDK update is required.
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and defaults to not being set (API defaults to 365 days).
- **API constraints:** `RotateDays` is only accepted by `EnableKeyRotation` (used when enabling key rotation). The `DescribeKey` response (`KeyMetadata`) includes `RotateDays`, so Read can refresh the value.