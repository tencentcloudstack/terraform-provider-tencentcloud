## Context

The `tencentcloud_kms_key` resource currently supports creating KMS keys with parameters such as `alias`, `description`, `key_usage`, `key_rotation_enabled`, `hsm_cluster_id`, `is_enabled`, `is_archived`, `pending_delete_window_in_days`, and `tags`. The KMS `EnableKeyRotation` API also accepts a `RotateDays` parameter to specify the key rotation period (in days, range 7–365, default 365), but the Terraform resource does not expose it.

**Current state:**
- Resource file: `tencentcloud/services/kms/resource_tc_kms_key.go`
- Service layer: `tencentcloud/services/kms/service_tencentcloud_kms.go` (`EnableKeyRotation` function at line 154)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118`

**API behavior analysis:**

| API | RotateDays in Request | RotateDays in Response |
|-----|----------------------|------------------------|
| `EnableKeyRotation` | Yes (`RotateDays *uint64`, range 7–365, default 365) | N/A (response only has RequestId) |
| `DescribeKey` (returns `KeyMetadata`) | N/A | Yes (`KeyMetadata.RotateDays *uint64`) |
| `DisableKeyRotation` | No | N/A |
| `CreateKey` | No | N/A |

**Key constraint:** `RotateDays` is only accepted by `EnableKeyRotation` (used when enabling key rotation). It is NOT in `CreateKey`. The `DescribeKey` response (`KeyMetadata`) includes `RotateDays`, so it can be refreshed on Read.

## Cloud API struct definitions (from vendor)

### EnableKeyRotationRequest

File: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118/models.go` (lines 2154–2165)

```go
type EnableKeyRotationRequest struct {
	*tchttp.BaseRequest

	// CMK唯一标识符
	KeyId *string `json:"KeyId,omitnil,omitempty" name:"KeyId"`

	// 密钥轮转周期，单位天，允许范围 7 ~ 365，默认值 365。
	RotateDays *uint64 `json:"RotateDays,omitnil,omitempty" name:"RotateDays"`

	// 可信服务成员账号信息,当前账号时管理员或者委派管理员时有效。
	MemberAccount *MemberAccount `json:"MemberAccount,omitnil,omitempty" name:"MemberAccount"`
}
```

### EnableKeyRotationResponse

File: `vendor/.../kms/v20190118/models.go` (lines 2194–2197)

```go
type EnableKeyRotationResponse struct {
	*tchttp.BaseResponse
	Response *EnableKeyRotationResponseParams `json:"Response"`
}
```

`EnableKeyRotationResponseParams` only contains `RequestId *string` (line 2191).

### KeyMetadata (read path — DescribeKey response)

File: `vendor/.../kms/v20190118/models.go` (lines 3441–3510). Relevant fields:

```go
type KeyMetadata struct {
	// ... other fields ...

	// <p>是否开启了密钥轮换功能</p>
	KeyRotationEnabled *bool `json:"KeyRotationEnabled,omitnil,omitempty" name:"KeyRotationEnabled"`  // line 3467

	// <p>在密钥轮换开启状态下，下次轮换的时间</p>
	NextRotateTime *uint64 `json:"NextRotateTime,omitnil,omitempty" name:"NextRotateTime"`  // line 3473

	// <p>HSM 集群 ID（仅对 KMS 独占版/托管版服务实例有效）</p>
	HsmClusterId *string `json:"HsmClusterId,omitnil,omitempty" name:"HsmClusterId"`  // line 3488

	// <p>密钥轮转周期（天）</p>
	RotateDays *uint64 `json:"RotateDays,omitnil,omitempty" name:"RotateDays"`  // line 3491

	// <p>上次轮转时间（Unix timestamp）</p>
	LastRotateTime *uint64 `json:"LastRotateTime,omitnil,omitempty" name:"LastRotateTime"`  // line 3494

	// ... other fields ...
}
```

## Existing resource schema (current state)

File: `tencentcloud/services/kms/resource_tc_kms_key.go`

The schema is built from two maps merged in `ResourceTencentCloudKmsKey()`:
- `TencentKmsBasicInfo()` (lines 20–62): `alias`, `description`, `is_enabled`, `is_archived`, `pending_delete_window_in_days`, `tags`, `key_state`
- `specialInfo` (lines 65–84): `key_usage` (ForceNew, default `ENCRYPT_DECRYPT`), `key_rotation_enabled` (Optional, default false), `hsm_cluster_id`

**Key rotation in Create (lines 165–181):** When `key_usage == ENCRYPT_DECRYPT` and `key_rotation_enabled == true`, calls `kmsService.EnableKeyRotation(ctx, d.Id())`. Currently does NOT pass `RotateDays`.

**Key rotation in Update (lines 342–352):** When `key_usage == ENCRYPT_DECRYPT` and `d.HasChange("key_rotation_enabled")`, calls `updateKeyRotationStatus(ctx, kmsService, keyId, keyRotationEnabled)` which calls `EnableKeyRotation` or `DisableKeyRotation`. Currently does NOT pass `RotateDays`.

**Read (lines 248–255):** Sets `alias`, `description`, `key_state`, `key_usage`, `key_rotation_enabled`, `hsm_cluster_id`. Currently does NOT read `RotateDays`.

## Existing service-layer function (current state)

File: `tencentcloud/services/kms/service_tencentcloud_kms.go` (lines 154–169)

```go
func (me *KmsService) EnableKeyRotation(ctx context.Context, keyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := kms.NewEnableKeyRotationRequest()
	request.KeyId = helper.String(keyId)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseKmsClient().EnableKeyRotation(request)
	if err != nil {
		errRet = errors.WithStack(err)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}
```

Currently the function only accepts `keyId` and does not pass `RotateDays`.

## Goals / Non-Goals

**Goals:**
- Add `rotate_days` (Optional, TypeInt, validate range 7–365) parameter to `tencentcloud_kms_key` to specify the key rotation period when key rotation is enabled.
- Pass `RotateDays` to the `EnableKeyRotation` API request when specified, during both Create and Update flows.
- Read `RotateDays` from the `DescribeKey` API response (`KeyMetadata.RotateDays`) in the Read function to support state refresh and import.
- Update the `EnableKeyRotation` service-layer function to accept and pass the `rotateDays` parameter.
- Maintain full backward compatibility — existing configurations continue to work unchanged (API defaults to 365 days when `RotateDays` is not set).

**Non-Goals:**
- Adding `rotate_days` to KMS datasources (`tencentcloud_kms_keys`, `tencentcloud_kms_describe_keys`) — out of scope.
- Making `rotate_days` independently updatable separate from `key_rotation_enabled` — `rotate_days` is only meaningful when `key_rotation_enabled` is `true`, and changes are handled via re-invoking `EnableKeyRotation` with the new value.

## Decisions

### Decision 1: `rotate_days` is Optional (not ForceNew, not immutable)

**Rationale:** The `EnableKeyRotation` API accepts `RotateDays` and can be re-invoked to update the rotation period. When `key_rotation_enabled` is true and the user changes `rotate_days`, the Update flow can re-call `EnableKeyRotation` with the new `RotateDays` value. This avoids unnecessary resource recreation. The parameter is NOT added to `immutableArgs`.

**Alternatives considered:**
- `ForceNew: true` — rejected, because it forces resource destruction/recreation unnecessarily; the API supports updating the rotation period.
- Immutable args pattern — rejected for the same reason; the value can be updated.

### Decision 2: Update `EnableKeyRotation` service function signature to accept `rotateDays`

**Rationale:** Add `rotateDays uint64` parameter to `EnableKeyRotation(ctx, keyId, rotateDays)`. When `rotateDays > 0`, set `request.RotateDays = helper.Uint64(rotateDays)`. When `rotateDays == 0`, do not set it (API defaults to 365). This is the minimal change needed and keeps the function backward-compatible within the package.

### Decision 3: Wire `rotate_days` through Create and Update

**Rationale:**
- **Create:** In the existing block that calls `EnableKeyRotation` (when `key_usage == ENCRYPT_DECRYPT` and `key_rotation_enabled == true`), read `rotate_days` from schema and pass it to `EnableKeyRotation`.
- **Update:** The existing `updateKeyRotationStatus` helper function needs to be updated to also accept and pass `rotateDays`. When `key_rotation_enabled` changes to `true` or when `rotate_days` changes (while rotation is enabled), re-call `EnableKeyRotation` with the new `rotate_days` value. When `key_rotation_enabled` changes to `false`, call `DisableKeyRotation` as before.

### Decision 4: Read `RotateDays` from `KeyMetadata`

**Rationale:** The `DescribeKey` API returns `KeyMetadata` which includes `RotateDays *uint64`. In the Read function, after checking `key.RotateDays != nil`, set `d.Set("rotate_days", key.RotateDays)`. This enables proper state refresh and import support.

### Decision 5: `rotate_days` only valid with `ENCRYPT_DECRYPT` key usage

**Rationale:** The `EnableKeyRotation` API and the existing `key_rotation_enabled` parameter are only valid when `key_usage` is `ENCRYPT_DECRYPT`. The `rotate_days` parameter follows the same constraint. The Create/Update flows already guard the `EnableKeyRotation` call with `keyUsage == KMS_KEY_USAGE_ENCRYPT_DECRYPT`, so `rotate_days` will only be passed when rotation is enabled.

## Risks / Trade-offs

- **[Risk] Changing `rotate_days` without `key_rotation_enabled` change in Update**: If the user changes only `rotate_days` (not `key_rotation_enabled`), the existing Update code only triggers rotation status change when `d.HasChange("key_rotation_enabled")`.
  - **Mitigation:** Extend the Update condition to also trigger when `d.HasChange("rotate_days")` (while `key_rotation_enabled` is true). Re-call `EnableKeyRotation` with the new `rotate_days` value.

- **[Risk] `rotate_days` set but `key_rotation_enabled` is false**: The user specifies `rotate_days` but does not enable key rotation.
  - **Mitigation:** `rotate_days` is simply ignored if `key_rotation_enabled` is false, since `EnableKeyRotation` is not called. This is consistent with the existing behavior where `key_rotation_enabled` only takes effect for `ENCRYPT_DECRYPT` keys. Documentation should clarify this.

- **[Risk] Imported resources have `rotate_days` populated**: `DescribeKey` returns `RotateDays`, so imported resources will have correct values.
  - **Mitigation:** No special handling needed; value is read from API response.