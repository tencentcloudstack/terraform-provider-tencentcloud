## Context

The `tencentcloud_teo_bind_security_template` resource manages template-to-domain bindings for TEO security policies. The Read path uses `DescribeTeoBindSecurityTemplateById` which returns an `EntityStatus` struct. The vendor SDK's `EntityStatus` struct already includes a `Message` field (instance configuration delivery message), but the resource schema currently does not expose it.

The service layer `DescribeTeoBindSecurityTemplateById` handles special-case logic for locating bindings and is intentionally kept as-is. This change only adds the `message` computed field to the resource schema and sets it in the Read function when the returned `EntityStatus.Message` is non-nil.

## Goals / Non-Goals

**Goals:**
- Add a `message` computed attribute to `tencentcloud_teo_bind_security_template` resource
- Preserve all existing behavior (status field, id construction, nil handling, retry logic)

**Non-Goals:**
- Do not modify the service layer `DescribeTeoBindSecurityTemplateById` function or any other service code
- Do not change the Create or Delete paths (they use `BindSecurityTemplateToEntity` which is unaffected)
- Do not add `message` as an input parameter (it is read-only/computed)
- Do not change the `EntityStatus` struct in the vendor SDK

## Decisions

### Decision 1: Add `message` as a `Computed` field in the schema

**Rationale**: `Message` is read-only configuration delivery information. It should be `Computed` only (not `Required` or `Optional`), matching the pattern of the existing `status` field.

### Decision 2: Set `message` from `respData.Message` only when non-nil

**Rationale**: Following the existing pattern for the `status` field, the Read function only calls `d.Set("message", ...)` when `respData.Message` is not nil. This avoids overwriting the field unnecessarily and keeps the behavior consistent with the rest of the Read function.

## Risks / Trade-offs

- **Risk**: The service layer may not currently populate the `Message` field, so the `message` attribute may remain empty until the service layer is updated separately.
  → **Mitigation**: This is acceptable and intentional — the service layer handles special cases and must not be modified in this change. The schema field is additive and backward compatible. When the service layer returns `Message`, it will automatically flow through to the resource state without further schema changes.
