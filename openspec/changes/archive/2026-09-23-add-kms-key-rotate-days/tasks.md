## 1. Service Layer Changes

- [x] 1.1 Update `EnableKeyRotation` function signature in `tencentcloud/services/kms/service_tencentcloud_kms.go` to accept `rotateDays uint64` parameter
- [x] 1.2 Set `request.RotateDays = helper.Uint64(rotateDays)` when `rotateDays > 0` before calling the `EnableKeyRotation` API

## 2. Resource Schema Changes

- [x] 2.1 Add `rotate_days` schema field
- [x] 2.2 Remove `ValidateFunc: tccommon.ValidateIntegerInRange(7, 365)` from `rotate_days` so range validation is delegated to the cloud API

## 3. Create Function Changes

- [x] 3.1 In `resourceTencentCloudKmsKeyCreate`, read `rotate_days` from schema data when `key_usage == ENCRYPT_DECRYPT` and `key_rotation_enabled == true`
- [x] 3.2 Pass `rotate_days` value to `kmsService.EnableKeyRotation(ctx, d.Id(), rotateDays)`

## 4. Update Function Changes

- [x] 4.1 Update `updateKeyRotationStatus` helper function to accept and pass `rotateDays uint64` to `EnableKeyRotation`
- [x] 4.2 In `resourceTencentCloudKmsKeyUpdate`, extend the key-rotation handling block to also trigger when `d.HasChange("rotate_days")` (while `key_rotation_enabled` is true), re-invoking `EnableKeyRotation` with the new `rotate_days` value
- [x] 4.3 Ensure `rotate_days` is NOT added to the `immutableArgs` array (it is updatable)

## 5. Read Function Changes

- [x] 5.1 In `resourceTencentCloudKmsKeyRead`, check `key.RotateDays != nil` and set `d.Set("rotate_days", key.RotateDays)` when non-nil

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/kms/resource_tc_kms_key.md` example file to document `rotate_days` usage

## 7. Tests

- [x] 7.1 Add unit test cases in `tencentcloud/services/kms/resource_tc_kms_key_test.go` covering `rotate_days` in Create/Update/Read flows using gomonkey mocks for cloud APIs

## 8. Validation

- [x] 8.1 Verify the code compiles successfully
- [x] 8.2 Verify no lint errors