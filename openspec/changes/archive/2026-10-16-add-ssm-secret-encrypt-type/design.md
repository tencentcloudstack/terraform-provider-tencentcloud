## Context

The `tencentcloud_ssm_secret` resource (RESOURCE_KIND_GENERAL) in `tencentcloud/services/ssm/resource_tc_ssm_secret.go` manages the full lifecycle of a TencentCloud SSM (Secrets Manager) secret. It currently supports `secret_name`, `description`, `kms_key_id`, `secret_type`, `additional_config`, `tags`, `is_enabled`, `recovery_window_in_days`, and the computed `status`.

The TencentCloud SSM SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923`, already vendored) now exposes an `EncryptType` field on the `CreateSecret` request (input) and on the `DescribeSecret` response (output). This field lets users choose between KMS key encryption (`0`) and soft key encryption (`1`) when creating a secret.

This change adds a single new parameter `encrypt_type` to the resource. It is the only new parameter in this change.

### Current state of the resource (relevant parts)

Schema (lines 30-83 of `resource_tc_ssm_secret.go`) — `encrypt_type` does **not** exist yet. Tags (`tags`) are already implemented via the shared tag service (`svctag.ModifyTags`) in Create/Read/Update, so the `ModifyTags` API entries in the requirement are the existing tag infrastructure, **not** new parameters.

Create (`resourceTencentCloudSsmSecretCreate`): builds `ssm.NewCreateSecretRequest()`, sets `SecretName`, `Description`, `KmsKeyId`, `SecretType`, `AdditionalConfig`, and (for `secret_type==0`) default `VersionId`/`SecretString`. Sends within `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

Read (`resourceTencentCloudSsmSecretRead`): calls `ssmService.DescribeSecretByName(ctx, secretName)` returning a `*SecretInfo`, then sets schema fields from it. `SecretInfo` (see below) **does not** currently carry `EncryptType`.

Update (`resourceTencentCloudSsmSecretUpdate`): has an `immutableArgs` slice `["secret_type", "additional_config"]` that rejects in-place changes; handles `description`, `is_enabled`, `tags`.

### Vendor API struct excerpts (from vendor/.../ssm/v20190923/models.go)

**`CreateSecretRequest`** (lines 348-383) — `EncryptType` is a request input field:
```go
type CreateSecretRequest struct {
	*tchttp.BaseRequest
	SecretName      *string  `json:"SecretName,omitnil,omitempty" name:"SecretName"`
	VersionId       *string  `json:"VersionId,omitnil,omitempty" name:"VersionId"`
	Description     *string  `json:"Description,omitnil,omitempty" name:"Description"`
	KmsKeyId        *string  `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
	SecretType      *uint64  `json:"SecretType,omitnil,omitempty" name:"SecretType"`
	SecretBinary    *string  `json:"SecretBinary,omitnil,omitempty" name:"SecretBinary"`
	SecretString    *string  `json:"SecretString,omitnil,omitempty" name:"SecretString"`
	AdditionalConfig *string `json:"AdditionalConfig,omitnil,omitempty" name:"AdditionalConfig"`
	Tags            []*Tag   `json:"Tags,omitnil,omitempty" name:"Tags"`
	KmsHsmClusterId *string  `json:"KmsHsmClusterId,omitnil,omitempty" name:"KmsHsmClusterId"`
	// 凭据加密类型 枚举值: 0 KMS密钥加密; 1 软密钥加密
	EncryptType     *uint64  `json:"EncryptType,omitnil,omitempty" name:"EncryptType"`
}
```

**`DescribeSecretResponseParams`** (lines 821-893) — `EncryptType` is a response output field:
```go
type DescribeSecretResponseParams struct {
	SecretName        *string  `json:"SecretName,omitnil,omitempty" name:"SecretName"`
	Description       *string  `json:"Description,omitnil,omitempty" name:"Description"`
	KmsKeyId          *string  `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
	CreateUin         *uint64  `json:"CreateUin,omitnil,omitempty" name:"CreateUin"`
	Status            *string  `json:"Status,omitnil,omitempty" name:"Status"`
	DeleteTime        *uint64  `json:"DeleteTime,omitnil,omitempty" name:"DeleteTime"`
	CreateTime        *uint64  `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`
	SecretType        *int64   `json:"SecretType,omitnil,omitempty" name:"SecretType"`
	ProductName      *string  `json:"ProductName,omitnil,omitempty" name:"ProductName"`
	ResourceID        *string  `json:"ResourceID,omitnil,omitempty" name:"ResourceID"`
	// ... rotation/ssh fields ...
	AdditionalConfig    *string `json:"AdditionalConfig,omitnil,omitempty" name:"AdditionalConfig"`
	// 凭据加密类型 枚举值: 0 KMS密钥加密; 1 软密钥加密 默认值0
	EncryptType         *uint64 `json:"EncryptType,omitnil,omitempty" name:"EncryptType"`
	EncryptSwitching    *bool   `json:"EncryptSwitching,omitnil,omitempty" name:"EncryptSwitching"`
	// ...
}
```

**`UpdateSecretRequest`** (lines 2107-2122) — `EncryptType` is **NOT** present (only `SecretName`, `VersionId`, `SecretBinary`, `SecretString`), confirming the encryption type is immutable after creation:
```go
type UpdateSecretRequest struct {
	*tchttp.BaseRequest
	SecretName    *string `json:"SecretName,omitnil,omitempty" name:"SecretName"`
	VersionId     *string `json:"VersionId,omitnil,omitempty" name:"VersionId"`
	SecretBinary  *string `json:"SecretBinary,omitnil,omitempty" name:"SecretBinary"`
	SecretString  *string `json:"SecretString,omitnil,omitempty" name:"SecretString"`
}
```

### Service layer (`service_tencentcloud_ssm.go`)

`SecretInfo` struct (lines 26-37) — needs a new `encryptType` field:
```go
type SecretInfo struct {
	secretName       string
	description       string
	kmsKeyId          string
	secretType        int64
	additionalConfig string
	createUin         uint64
	status           string
	deleteTime       uint64
	createTime       uint64
	resourceId       string
}
```

`DescribeSecretByName` (lines 103-130) returns `*SecretInfo`, populating it by dereferencing `*response.Response.Xxx`. `EncryptType` (`*uint64`) must be captured here; note `SecretInfo.secretType` is `int64` while the API returns `*int64`, and `EncryptType` is `*uint64` — the new field should be `uint64` to match the API type.

## Goals / Non-Goals

**Goals:**
- Expose `encrypt_type` on `tencentcloud_ssm_secret` so users can set KMS (`0`) or soft-key (`1`) encryption at creation time.
- Read back the effective `EncryptType` from `DescribeSecret` so state reflects reality.
- Treat `encrypt_type` as immutable (force re-creation when changed), matching the SDK contract that `EncryptType` is not part of `UpdateSecret`.
- Preserve backward compatibility: existing configs without `encrypt_type` keep working (API default `0` applies and is read back via Computed).

**Non-Goals:**
- Not exposing `EncryptSwitching` (rotation/switching state) — out of scope for this single-parameter change.
- Not adding encryption-type switching workflows (the API does not support updating `EncryptType` via `UpdateSecret`).
- Not modifying the `tags` implementation — it is already complete.

## Decisions

### Decision 1: `encrypt_type` schema type = `schema.TypeInt`, Optional + Computed + ForceNew
- The SDK field is `*uint64`. The provider conventionally represents such numeric enums with `schema.TypeInt` (e.g., the existing `secret_type` is `schema.TypeInt` while the API uses `*int64`/`*uint64`). Using `TypeInt` keeps consistency with the resource's existing style.
- Optional + Computed so that when the user omits it, the API default is applied and the value is read back into state.
- ForceNew because `EncryptType` is only on `CreateSecret`, never on `UpdateSecret`. Terraform's `ForceNew: true` cleanly handles this by destroying+recreating rather than erroring. We additionally list it in `immutableArgs` for defense-in-depth, but `ForceNew` is the primary mechanism.

**Alternatives considered:**
- `schema.TypeString` ("0"/"1"): rejected; the resource already uses `TypeInt` for `secret_type`, so `TypeInt` is more consistent.
- Plain `Optional` without `Computed`: rejected — then `d.Set("encrypt_type", ...)` on read would warn/err for unset computed values; Computed is required to round-trip the API default.
- Making it updatable: rejected — the SDK has no update path for `EncryptType`.

### Decision 2: Carry `EncryptType` through `SecretInfo`
- `DescribeSecretByName` is the single place that calls `DescribeSecret` and maps the response into the in-memory `SecretInfo`. Adding `encryptType uint64` to `SecretInfo` and populating it from `*response.Response.EncryptType` (with a nil-guard) keeps the read path uniform with the other fields.
- In `resourceTencentCloudSsmSecretRead`, set `encrypt_type` via `d.Set("encrypt_type", secretInfo.encryptType)`.

**Alternatives considered:**
- Calling `DescribeSecret` directly in the resource read: rejected; would bypass the existing service abstraction and duplicate logic.

### Decision 3: Create path — set `request.EncryptType` only when configured
- Use `d.GetOk("encrypt_type")` to detect a user-supplied value; when present, set `request.EncryptType = helper.IntUint64(v.(int))` (matching the existing `secret_type` pattern that uses `helper.IntUint64`). When absent, leave it nil so the API applies its default (`0`).

### Decision 4: Type mapping — `helper.IntUint64`
- The existing code uses `helper.IntUint64(v.(int))` for `secret_type` (`*uint64`). We follow the same helper for `EncryptType` to stay consistent.

## Risks / Trade-offs

- **[Risk] ForceNew on `encrypt_type` change destroys the secret** → Mitigation: This matches the API contract (encryption type is immutable). Document clearly that changing `encrypt_type` triggers re-creation. Since secrets may hold sensitive data, users must use `recovery_window_in_days` or re-seed versions after recreation; this is inherent to the API limitation, not a provider bug.
- **[Risk] SDK response `EncryptType` could be nil for very old secrets** → Mitigation: nil-guard in `DescribeSecretByName` (only populate when non-nil); in read, `d.Set` with the zero value (`uint64(0)`) is acceptable since `0` (KMS) is the API default. Alternatively skip `d.Set` when nil — but Computed fields should be set; we set `0`.
- **[Trade-off] `TypeInt` vs `TypeBool`** → `EncryptType` is an enum (`0`/`1`), but future values are possible, so `TypeInt` is more future-proof than a boolean.
- **[Risk] `DescribeSecretByName` dereferences pointers without nil-guards today** → Mitigation: This pre-existing pattern assumes all fields are populated. We add a nil-guard specifically for `EncryptType` to be safe, while not refactoring the existing fields (out of scope).

## Migration Plan

- No migration required. The new field is optional+computed. On `terraform apply`, existing resources are unaffected; `terraform refresh` will populate `encrypt_type` from the API for existing state.
- Rollback: removing the field reverts to current behavior; existing state with the key is ignored by older provider versions.

## Open Questions

None. The SDK contract is clear: `EncryptType` is create-only input and read output.