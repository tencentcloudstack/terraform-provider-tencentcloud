## Context

The `tencentcloud_teo_bind_security_template` resource manages template-to-domain bindings for TEO security policies. Currently, the Read path uses `DescribeTeoBindSecurityTemplateById` which queries `DescribeZones` to enumerate all zones, then paginates through `DescribeWebSecurityTemplates` to find a specific binding. This approach does not expose the `Message` field from `EntityStatus` because `DescribeWebSecurityTemplates` returns `BindDomainInfo` (which lacks `Message`), not `EntityStatus`.

The vendor SDK already includes `DescribeSecurityTemplateBindings` API which returns `SecurityTemplateBinding` containing `TemplateScope`, which in turn contains `EntityStatus` with `Entity`, `Status`, and `Message` fields. This API accepts a specific `ZoneId` and `TemplateId` array, making it a more direct and efficient query.

The existing `EntityStatus` struct already has the `Message` field in the vendor SDK, so no SDK changes are needed.

## Goals / Non-Goals

**Goals:**
- Add a `message` computed attribute to `tencentcloud_teo_bind_security_template` resource
- Switch the Read path from `DescribeWebSecurityTemplates` to `DescribeSecurityTemplateBindings` API
- Preserve all existing behavior (status field, id construction, nil handling, retry logic)

**Non-Goals:**
- Do not change the Create or Delete paths (they use `BindSecurityTemplateToEntity` which is unaffected)
- Do not add `message` as an input parameter (it is read-only/computed)
- Do not change the `EntityStatus` struct in the vendor SDK

## Decisions

### Decision 1: Use `DescribeSecurityTemplateBindings` API instead of `DescribeWebSecurityTemplates`

**Rationale**: `DescribeSecurityTemplateBindings` returns `EntityStatus` structs with the `Message` field, while `DescribeWebSecurityTemplates` returns `BindDomainInfo` which lacks `Message`. The new API also accepts the target `ZoneId` and `TemplateId` directly, eliminating the need to enumerate all zones and iterate through templates.

**Alternatives considered**:
- *Keep `DescribeWebSecurityTemplates` and make a second call to `DescribeSecurityTemplateBindings` for `Message`*: Rejected — adds unnecessary API calls and complexity.
- *Modify `DescribeWebSecurityTemplates` to include `Message`*: Rejected — we cannot change the cloud API.

### Decision 2: Modify `DescribeTeoBindSecurityTemplateById` to accept and return `Message`

**Rationale**: The service function currently constructs an `EntityStatus` with only `Entity` and `Status`. After switching APIs, the function will receive `EntityStatus` from the API response directly, which already contains `Message`. The function signature stays the same (returns `*EntityStatus`), so callers automatically gain access to `Message` without further changes.

### Decision 3: Add `message` as a `Computed` field in the schema

**Rationale**: `Message` is read-only configuration delivery information. It should be `Computed` only (not `Required` or `Optional`), matching the pattern of the existing `status` field.

## Risks / Trade-offs

- **Risk**: `DescribeSecurityTemplateBindings` API behavior differs from `DescribeWebSecurityTemplates` in edge cases (e.g., when template has no bindings).  
  → **Mitigation**: The vendor SDK documentation confirms `DescribeSecurityTemplateBindings` returns empty `TemplateScope` array when no bindings exist. We handle nil/empty responses in the same defensive manner as the current implementation.

- **Risk**: The `DescribeSecurityTemplateBindings` API may have different rate limits or availability.  
  → **Mitigation**: The API is already defined in the same vendor SDK version and is used in the same service. We continue using the existing retry pattern with `tccommon.ReadRetryTimeout`.

- **Risk**: State refresh functions in `_extension.go` depend on the current `DescribeTeoBindSecurityTemplateById` behavior.  
  → **Mitigation**: The function signature and return type (`*EntityStatus`) remain unchanged. The state refresh functions use `resp.Status` which is preserved.
