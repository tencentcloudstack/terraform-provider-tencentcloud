## Context

The `tencentcloud_teo_bind_security_template` resource reads its state through `TeoService.DescribeTeoBindSecurityTemplateById`, which historically called `DescribeSecurityTemplateBindings(zoneId, templateId)` and scanned the returned `SecurityTemplate[0].TemplateScope[0].EntityStatus` for the matching `entity`. The binding lookup is unreliable, and the provider already has a supported alternative — `DescribeWebSecurityTemplates(zoneIds)` — which returns `SecurityPolicyTemplateInfo` objects including `BindDomains` (each with `Domain`, `ZoneId`, `Status`). This change swaps the read-path API without altering the resource schema or Terraform-facing behavior.

The existing codebase already contains a paged `DescribeTeoZonesByFilter` helper, but it returns `[]*Zone` and uses `UseTeoClient()`. To keep the new read path self-contained and consistent with the v20220901 client used by the rest of `DescribeTeoBindSecurityTemplateById`, a dedicated `describeTeoAllZoneIds` helper using `UseTeoV20220901Client().DescribeZones` is introduced.

## Goals / Non-Goals

**Goals:**
- Replace `DescribeSecurityTemplateBindings` usage with `DescribeZones` + `DescribeWebSecurityTemplates` in the read path.
- Respect the `DescribeWebSecurityTemplates` constraint of at most 100 zone IDs per request by batching.
- Preserve the resource's user-facing schema, id format (`zoneId#templateId#entity`), and lifecycle (create/read/delete only).
- Keep the create state-refresh polling safe against a nil `Status`.
- Preserve the resource id in logs before clearing state on not-found.

**Non-Goals:**
- Changing the resource schema or adding an update operation.
- Migrating other TEO resources off `DescribeSecurityTemplateBindings`.
- Reusing `DescribeTeoZonesByFilter` (different client wrapper / return type); a focused helper is added instead to avoid cross-coupling.

## Decisions

**Decision 1: Fetch all zone IDs via `DescribeZones`, then batch `DescribeWebSecurityTemplates`.**
- Rationale: `DescribeWebSecurityTemplates` requires zone IDs as input (unlike `DescribeSecurityTemplateBindings` which took a template ID directly). The binding we need is identified by `templateId` + `entity`, and the zone is already part of the resource's composite id — but the template may be returned under any zone the account owns, so we enumerate all zones to be safe.
- Alternative considered: Pass only the resource's own `zoneId` to `DescribeWebSecurityTemplates`. Rejected because the original implementation queried by template ID globally; limiting to a single zone could miss bindings reported under a different zone and cause spurious not-found / state drift.

**Decision 2: Batch size of 100 for `DescribeWebSecurityTemplates`.**
- Rationale: The SDK model doc states "单次查询最多传入 100 个站点" (at most 100 zone IDs per request). A constant `batchSize = 100` enforces this.
- Alternative considered: Smaller batch (e.g., 50). Rejected — would increase request count without benefit.

**Decision 3: `DescribeZones` paging uses `Limit=100` (documented maximum).**
- Rationale: Maximizes pages-per-request efficiency; matches the existing `DescribeTeoZonesByFilter` convention.

**Decision 4: Return a synthesized `EntityStatus{Entity, Status}` instead of reusing the API's `EntityStatus` directly.**
- Rationale: `DescribeWebSecurityTemplates` returns `BindDomainInfo{Domain, ZoneId, Status}`, not `EntityStatus`. To keep the resource Read / state-refresh code unchanged (it expects `*EntityStatus` with `Entity` and `Status`), the helper constructs an `EntityStatus` from the matched `BindDomainInfo`.

**Decision 5: Drop the `DescribeSecurityTemplateBindingsRequest` field from the extension state-refresh function.**
- Rationale: The field was unused after the API swap and referenced a request type no longer constructed. Removing it keeps the extension clean and avoids dead code.

**Decision 6: Guard nil `Status` in the state-refresh function.**
- Rationale: During create polling the freshly-created binding may briefly return a non-nil `EntityStatus` with a nil `Status`. The original code dereferenced `*resp.Status` and could panic. Returning `(resp, "", nil)` keeps the refresh in a pending state until a concrete status arrives.

## Risks / Trade-offs

- [More API calls per read] `DescribeZones` paging + multiple `DescribeWebSecurityTemplates` batches replace a single `DescribeSecurityTemplateBindings` call. → Mitigated: batching caps request count at `ceil(zoneCount/100) + ceil(zoneCount/100)`; accounts rarely have hundreds of zones, and each call is retried with `ReadRetryTimeout`. The trade-off is acceptable for correctness.
- [Read latency increase for large zone counts] → Mitigated: each batch is independent and short; retry is per-batch so a transient failure does not restart all batches.
- [Behavioral parity] The synthesized `EntityStatus.Status` comes from `BindDomainInfo.Status`, whose documented values (`process`/`online`/`fail`) match the original `EntityStatus.Status` values. → No state-machine change for the create `Target: ["online"]`.
- [No migration needed] State id format and schema are unchanged; existing Terraform states remain valid.
