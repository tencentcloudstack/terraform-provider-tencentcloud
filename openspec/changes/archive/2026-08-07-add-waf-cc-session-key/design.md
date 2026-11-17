## Context

The `tencentcloud_waf_cc_session` resource (in `tencentcloud/services/waf/resource_tc_waf_cc_session.go`) manages WAF session definitions via the TencentCloud WAF `v20180125` SDK. It currently uses the `UpsertSession` API for both create and update operations and `DescribeSession` for read operations.

The `UpsertSessionRequest` struct already contains a `Key *string` field (annotated as "精准匹配时配置的key" — the key configured for precise matching). The read-side `SessionItem` struct returned by `DescribeSession` also contains the same `Key *string` field. Today the Terraform schema does not expose this parameter, so users cannot configure precise-match keys from Terraform.

The resource uses a composite ID (`domain#edition#sessionID`) and an immutable-args pattern in update (`immutableArgs` slice guards `domain`, `edition`, `session_name`).

## Goals / Non-Goals

**Goals:**
- Expose the `Key` parameter in the `tencentcloud_waf_cc_session` resource schema so users can configure precise-match session keys.
- Wire `Key` into the create, read, and update code paths using the exact same pattern as the other string fields (`source`, `category`, `key_or_start_mat`, etc.).
- Keep the change backward compatible (optional field, no state migration, no ForceNew).

**Non-Goals:**
- Do not modify any existing schema fields or their behavior.
- Do not change the composite ID structure.
- Do not introduce a `_extension.go` file.
- Do not add retry logic beyond what the existing `resource.Retry` blocks already provide.

## Decisions

### 1. Parameter Naming: `key`

**Decision**: Use `key` (lowercase, matching the SDK field name `Key`).

**Rationale**:
- The existing resource schema maps SDK fields with straightforward snake_case where names are single words. `key` is a single word, so no transformation is needed.
- Matches the SDK request field `request.Key` directly, keeping the mapping obvious.

**Alternatives considered**:
- `session_key`: More explicit, but the resource is already named `waf_cc_session`, so `session_key` would be redundant. `key` is consistent with the API naming.

### 2. Parameter Type: `schema.TypeString`

**Decision**: `schema.TypeString`.

**Rationale**: The SDK field is `Key *string`, so a string schema type is the direct match.

### 3. Optional vs Required

**Decision**: Optional (`Optional: true`), no default, not ForceNew.

**Rationale**:
- The parameter is not always needed — it only applies to precise-matching scenarios. Making it optional avoids breaking existing configurations and lets the API default apply when omitted.
- Not ForceNew so updates can modify it in-place via `UpsertSession` (which is an upsert by design, using `SessionID` to identify the target).

### 4. Update Path

**Decision**: Treat `key` as a mutable argument in the update operation — do NOT add it to the `immutableArgs` slice.

**Rationale**: The `UpsertSession` API is an upsert; it accepts `Key` on both create and update. Since the API supports updating `Key`, it should be mutable. The existing mutable fields (`source`, `category`, `key_or_start_mat`, `end_mat`, `start_offset`, `end_offset`) follow the same pattern of populating the request from `d.GetOk`.

### 5. Read Path

**Decision**: In `resourceTencentCloudWafCcSessionRead`, after fetching `ccSession` via `DescribeWafCcSessionById`, set `d.Set("key", ccSession.Key)` only when `ccSession.Key != nil`, following the existing nil-guard pattern used for every other field.

**Rationale**: Consistent with the existing read pattern and the project rule that requires nil-checking response fields before calling `d.Set`.

## Risks / Trade-offs

- **[Risk] API response omits `Key` when not set** → Mitigation: nil-guard in read (`if ccSession.Key != nil`) ensures we don't overwrite state with an empty value when the API returns nil.
- **[Risk] `Key` only meaningful for precise-match category** → Mitigation: This is an API-level concern; Terraform simply forwards the value. No schema-level validation is added to stay consistent with other similarly constrained fields (e.g., `end_mat` which is only meaningful when `category == match`).
- **[Trade-off] No ForceNew on `key`** → Acceptable because `UpsertSession` supports in-place update; aligns with existing mutable fields.

## Migration Plan

No migration required. The field is purely additive and optional. Existing configurations without `key` continue to work unchanged. Rollback is as simple as removing the parameter from the configuration; the provider remains backward compatible.

## Open Questions

None. The SDK supports `Key` on both the request and response structs, and the API semantics are clear.
